package auth

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"context"

	"github.com/democorp/crypto-inventory/services/auth-service/internal/models"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrEmailNotVerified   = errors.New("email not verified")
)

// AuthService handles authentication operations
type AuthService struct {
	db       *sql.DB
	redis    *redis.Client
	jwt      *JWTService
	password *PasswordService
}

// NewAuthService creates a new authentication service
func NewAuthService(db *sql.DB, redis *redis.Client, jwt *JWTService) *AuthService {
	return &AuthService{
		db:       db,
		redis:    redis,
		jwt:      jwt,
		password: NewPasswordService(),
	}
}

// Register creates a new user account
func (a *AuthService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Validate password strength
	if err := ValidatePasswordStrength(req.Password); err != nil {
		return nil, err
	}

	// Check if email already exists
	existingUser, err := a.GetUserByEmail(req.Email)
	if err != nil && err != ErrUserNotFound {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	passwordHash, err := a.password.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create tenant if specified
	var tenantID uuid.UUID
	if req.TenantName != "" {
		tenant, err := a.createTenant(req.TenantName)
		if err != nil {
			return nil, fmt.Errorf("failed to create tenant: %w", err)
		}
		tenantID = tenant.ID
	} else {
		// TODO: Handle case where user joins existing tenant
		return nil, errors.New("tenant selection not implemented")
	}

	// Create user
	userID := uuid.New()
	user := &models.User{
		ID:            userID,
		TenantID:      tenantID,
		Email:         strings.ToLower(req.Email),
		PasswordHash:  passwordHash,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		Role:          "admin", // First user in tenant is admin
		IsActive:      true,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err = a.createUser(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Login authenticates a user and returns tokens
func (a *AuthService) Login(req *models.LoginRequest) (*models.AuthResponse, error) {
	// Get user by email
	user, err := a.GetUserByEmail(req.Email)
	if err != nil {
		if err == ErrUserNotFound {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserInactive
	}

	// Verify password
	valid, err := a.password.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessToken, refreshToken, err := a.jwt.GenerateTokens(
		user.ID, user.TenantID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	err = a.updateUserLastLogin(user.ID, now)
	if err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to update last login time: %v\n", err)
	}

	// Store refresh token in Redis
	err = a.storeRefreshToken(user.ID, refreshToken, a.jwt.GetRefreshExpiry())
	if err != nil {
		// Log error but don't fail login
		fmt.Printf("Failed to store refresh token: %v\n", err)
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(a.jwt.GetAccessExpiry().Seconds()),
	}, nil
}

// RefreshToken generates new tokens using a refresh token
func (a *AuthService) RefreshToken(refreshToken string) (*models.AuthResponse, error) {
	// Validate refresh token
	claims, err := a.jwt.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if claims.Type != "refresh" {
		return nil, ErrInvalidToken
	}

	// Check if refresh token exists in Redis
	exists, err := a.redis.Exists(context.Background(), fmt.Sprintf("refresh_token:%s", claims.UserID.String())).Result()
	if err != nil {
		return nil, err
	}
	if exists == 0 {
		return nil, ErrInvalidToken
	}

	// Get user
	user, err := a.GetUserByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := a.jwt.GenerateTokens(
		user.ID, user.TenantID, user.Email, user.Role)
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Store new refresh token and remove old one
	err = a.storeRefreshToken(user.ID, newRefreshToken, a.jwt.GetRefreshExpiry())
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    int64(a.jwt.GetAccessExpiry().Seconds()),
	}, nil
}

// Logout invalidates the refresh token
func (a *AuthService) Logout(userID uuid.UUID) error {
	return a.redis.Del(context.Background(), fmt.Sprintf("refresh_token:%s", userID.String())).Err()
}

// GetUserByEmail retrieves a user by email
func (a *AuthService) GetUserByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, role, 
		       active, email_verified, last_login_at, created_at, updated_at, deleted_at
		FROM users 
		WHERE email = $1 AND deleted_at IS NULL`

	user := &models.User{}
	err := a.db.QueryRow(query, strings.ToLower(email)).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Role, &user.IsActive,
		&user.EmailVerified, &user.LastLoginAt, &user.CreatedAt,
		&user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (a *AuthService) GetUserByID(userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, tenant_id, email, password_hash, first_name, last_name, role, 
		       active, email_verified, last_login_at, created_at, updated_at, deleted_at
		FROM users 
		WHERE id = $1 AND deleted_at IS NULL`

	user := &models.User{}
	err := a.db.QueryRow(query, userID).Scan(
		&user.ID, &user.TenantID, &user.Email, &user.PasswordHash,
		&user.FirstName, &user.LastName, &user.Role, &user.IsActive,
		&user.EmailVerified, &user.LastLoginAt, &user.CreatedAt,
		&user.UpdatedAt, &user.DeletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// createUser inserts a new user into the database
func (a *AuthService) createUser(user *models.User) error {
	query := `
		INSERT INTO users (id, tenant_id, email, password_hash, first_name, last_name, role, active, email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	_, err := a.db.Exec(query,
		user.ID, user.TenantID, user.Email, user.PasswordHash,
		user.FirstName, user.LastName, user.Role, user.IsActive,
		user.EmailVerified, user.CreatedAt, user.UpdatedAt,
	)

	return err
}

// createTenant creates a new tenant
func (a *AuthService) createTenant(name string) (*models.Tenant, error) {
	tenantID := uuid.New()
	slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))

	tenant := &models.Tenant{
		ID:               tenantID,
		Name:             name,
		Slug:             slug,
		SubscriptionTier: "basic",
		MaxEndpoints:     100,
		MaxUsers:         10,
		Settings:         make(map[string]interface{}),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	query := `
		INSERT INTO tenants (id, name, slug, subscription_tier, max_endpoints, max_users, settings, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := a.db.Exec(query,
		tenant.ID, tenant.Name, tenant.Slug, tenant.SubscriptionTier,
		tenant.MaxEndpoints, tenant.MaxUsers, "{}", tenant.CreatedAt, tenant.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return tenant, nil
}

// updateUserLastLogin updates the user's last login time
func (a *AuthService) updateUserLastLogin(userID uuid.UUID, lastLogin time.Time) error {
	query := `UPDATE users SET last_login_at = $1, updated_at = $2 WHERE id = $3`
	_, err := a.db.Exec(query, lastLogin, time.Now(), userID)
	return err
}

// storeRefreshToken stores a refresh token in Redis
func (a *AuthService) storeRefreshToken(userID uuid.UUID, token string, expiry time.Duration) error {
	key := fmt.Sprintf("refresh_token:%s", userID.String())
	return a.redis.Set(context.Background(), key, token, expiry).Err()
}
