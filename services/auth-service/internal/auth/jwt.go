package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	UserID   uuid.UUID `json:"user_id"`
	TenantID uuid.UUID `json:"tenant_id"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	Type     string    `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateTokens generates both access and refresh tokens
func (j *JWTService) GenerateTokens(userID, tenantID uuid.UUID, email, role string) (string, string, error) {
	// Generate access token
	accessToken, err := j.generateToken(userID, tenantID, email, role, "access", j.accessExpiry)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err := j.generateToken(userID, tenantID, email, role, "refresh", j.refreshExpiry)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// generateToken generates a JWT token with the given parameters
func (j *JWTService) generateToken(userID, tenantID uuid.UUID, email, role, tokenType string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID:   userID,
		TenantID: tenantID,
		Email:    email,
		Role:     role,
		Type:     tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "crypto-inventory-auth",
			Audience:  jwt.ClaimStrings{"crypto-inventory"},
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

// ValidateToken validates and parses a JWT token
func (j *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return j.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GetAccessExpiry returns the access token expiry duration
func (j *JWTService) GetAccessExpiry() time.Duration {
	return j.accessExpiry
}

// GetRefreshExpiry returns the refresh token expiry duration
func (j *JWTService) GetRefreshExpiry() time.Duration {
	return j.refreshExpiry
}
