# 🔐 Crypto Inventory Platform API Documentation

*Version: 1.0*  
*Last Updated: 2025-01-09*

## 📋 Overview

The Crypto Inventory Platform provides a comprehensive REST API for managing cryptographic assets, user authentication, and compliance reporting. The API follows RESTful principles and uses JWT-based authentication with multi-tenant support.

### **Base URLs**
- **Auth Service**: `http://localhost:8081`
- **Inventory Service**: `http://localhost:8082`
- **Compliance Engine**: `http://localhost:8083` (Coming Soon)
- **Report Generator**: `http://localhost:8084` (Coming Soon)
- **Sensor Manager**: `http://localhost:8085` (Coming Soon)

### **Authentication**
All API endpoints (except health checks) require JWT authentication via the `Authorization` header:
```
Authorization: Bearer <access_token>
```

### **Response Format**
All responses follow this structure:
```json
{
  "data": {}, // Response data (varies by endpoint)
  "message": "Success", // Human-readable message
  "error": null, // Error details if applicable
  "timestamp": "2025-01-09T10:30:00Z"
}
```

---

## 🔑 Authentication Service API

### **Health Check**
```http
GET /health
```
**Description**: Check if the auth service is running  
**Authentication**: None required  
**Response**:
```json
{
  "status": "healthy",
  "service": "auth-service",
  "timestamp": "2025-01-09T10:30:00Z"
}
```

### **User Registration**
```http
POST /api/v1/auth/register
```
**Description**: Create a new user account and tenant  
**Authentication**: None required  
**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Doe",
  "tenant_name": "Acme Corp"
}
```
**Validation Rules**:
- `email`: Valid email format, globally unique
- `password`: 8+ characters, mixed case, numbers, special characters
- `first_name`: Required, 1+ characters
- `last_name`: Required, 1+ characters
- `tenant_name`: Required, 2+ characters

**Response** (201 Created):
```json
{
  "data": {
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "admin",
      "is_active": true,
      "email_verified": false,
      "created_at": "2025-01-09T10:30:00Z"
    },
    "tenant": {
      "id": "123e4567-e89b-12d3-a456-426614174001",
      "name": "Acme Corp",
      "slug": "acme-corp",
      "subscription_tier": "trial",
      "trial_ends_at": "2025-02-09T10:30:00Z"
    }
  },
  "message": "User registered successfully",
  "error": null
}
```

### **User Login**
```http
POST /api/v1/auth/login
```
**Description**: Authenticate user and return JWT tokens  
**Authentication**: None required  
**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!"
}
```
**Response** (200 OK):
```json
{
  "data": {
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "admin",
      "tenant_id": "123e4567-e89b-12d3-a456-426614174001"
    },
    "tokens": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
      "expires_in": 900,
      "token_type": "Bearer"
    }
  },
  "message": "Login successful",
  "error": null
}
```

### **Get Current User**
```http
GET /api/v1/auth/me
```
**Description**: Get current authenticated user information  
**Authentication**: Required (JWT)  
**Response** (200 OK):
```json
{
  "data": {
    "user": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "email": "user@example.com",
      "first_name": "John",
      "last_name": "Doe",
      "role": "admin",
      "is_active": true,
      "email_verified": true,
      "last_login_at": "2025-01-09T10:30:00Z",
      "created_at": "2025-01-09T10:30:00Z"
    },
    "tenant": {
      "id": "123e4567-e89b-12d3-a456-426614174001",
      "name": "Acme Corp",
      "subscription_tier": "trial",
      "trial_ends_at": "2025-02-09T10:30:00Z"
    }
  },
  "message": "User information retrieved",
  "error": null
}
```

### **Refresh Token**
```http
POST /api/v1/auth/refresh
```
**Description**: Refresh access token using refresh token  
**Authentication**: None required  
**Request Body**:
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
**Response** (200 OK):
```json
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 900,
    "token_type": "Bearer"
  },
  "message": "Token refreshed successfully",
  "error": null
}
```

### **Logout**
```http
POST /api/v1/auth/logout
```
**Description**: Invalidate refresh token and logout user  
**Authentication**: Required (JWT)  
**Response** (200 OK):
```json
{
  "data": null,
  "message": "Logged out successfully",
  "error": null
}
```

---

## 📊 Inventory Service API

### **Health Check**
```http
GET /health
```
**Description**: Check if the inventory service is running  
**Authentication**: None required  
**Response**:
```json
{
  "status": "healthy",
  "service": "inventory-service",
  "timestamp": "2025-01-09T10:30:00Z"
}
```

### **Get Assets**
```http
GET /api/v1/assets
```
**Description**: Retrieve paginated list of assets with filtering and search  
**Authentication**: Required (JWT)  
**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)
- `search` (optional): Search term for hostname, IP, or description
- `asset_type` (optional): Filter by asset type (server, endpoint, service, appliance)
- `environment` (optional): Filter by environment (production, staging, development)
- `protocol` (optional): Filter by protocol (TLS, SSH, etc.)
- `risk_level` (optional): Filter by risk level (high, medium, low, unknown)
- `sort_by` (optional): Sort field (risk_score, hostname, first_discovered_at)
- `sort_order` (optional): Sort direction (asc, desc)

**Example Request**:
```http
GET /api/v1/assets?page=1&limit=20&search=web&asset_type=server&environment=production&risk_level=high&sort_by=risk_score&sort_order=desc
```

**Response** (200 OK):
```json
{
  "data": {
    "assets": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "hostname": "web-server-01",
        "ip_address": "192.168.1.100",
        "port": 443,
        "asset_type": "server",
        "operating_system": "Ubuntu 20.04",
        "environment": "production",
        "business_unit": "Engineering",
        "owner_email": "admin@acme.com",
        "description": "Main web server",
        "tags": {
          "critical": true,
          "department": "engineering"
        },
        "risk_score": 85,
        "risk_level": "high",
        "first_discovered_at": "2025-01-01T00:00:00Z",
        "last_seen_at": "2025-01-09T10:30:00Z",
        "created_at": "2025-01-01T00:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8,
      "has_next": true,
      "has_prev": false
    }
  },
  "message": "Assets retrieved successfully",
  "error": null
}
```

### **Get Asset by ID**
```http
GET /api/v1/assets/{id}
```
**Description**: Get detailed information about a specific asset  
**Authentication**: Required (JWT)  
**Path Parameters**:
- `id`: Asset UUID

**Response** (200 OK):
```json
{
  "data": {
    "asset": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "hostname": "web-server-01",
      "ip_address": "192.168.1.100",
      "port": 443,
      "asset_type": "server",
      "operating_system": "Ubuntu 20.04",
      "environment": "production",
      "business_unit": "Engineering",
      "owner_email": "admin@acme.com",
      "description": "Main web server",
      "tags": {
        "critical": true,
        "department": "engineering"
      },
      "metadata": {
        "os_version": "20.04.3 LTS",
        "kernel": "5.4.0-89-generic"
      },
      "risk_score": 85,
      "risk_level": "high",
      "first_discovered_at": "2025-01-01T00:00:00Z",
      "last_seen_at": "2025-01-09T10:30:00Z",
      "created_at": "2025-01-01T00:00:00Z",
      "crypto_implementations": [
        {
          "id": "123e4567-e89b-12d3-a456-426614174002",
          "protocol": "TLS",
          "protocol_version": "1.3",
          "cipher_suite": "TLS_AES_256_GCM_SHA384",
          "key_exchange_algorithm": "ECDHE",
          "key_size": 256,
          "certificate_issuer": "Let's Encrypt",
          "certificate_subject": "web-server-01.acme.com",
          "certificate_valid_from": "2025-01-01T00:00:00Z",
          "certificate_valid_to": "2025-04-01T00:00:00Z",
          "risk_score": 85,
          "risk_level": "high",
          "created_at": "2025-01-01T00:00:00Z"
        }
      ]
    }
  },
  "message": "Asset retrieved successfully",
  "error": null
}
```

### **Search Assets**
```http
GET /api/v1/assets/search
```
**Description**: Advanced search across assets with multiple criteria  
**Authentication**: Required (JWT)  
**Query Parameters**: Same as Get Assets, plus:
- `q` (required): Search query
- `fields` (optional): Comma-separated fields to search (hostname,ip_address,description)

**Example Request**:
```http
GET /api/v1/assets/search?q=web&fields=hostname,description&asset_type=server&environment=production
```

**Response**: Same format as Get Assets

### **Get Asset Crypto Implementations**
```http
GET /api/v1/assets/{id}/crypto
```
**Description**: Get all cryptographic implementations for a specific asset  
**Authentication**: Required (JWT)  
**Path Parameters**:
- `id`: Asset UUID

**Response** (200 OK):
```json
{
  "data": {
    "crypto_implementations": [
      {
        "id": "123e4567-e89b-12d3-a456-426614174002",
        "protocol": "TLS",
        "protocol_version": "1.3",
        "cipher_suite": "TLS_AES_256_GCM_SHA384",
        "key_exchange_algorithm": "ECDHE",
        "key_size": 256,
        "certificate_issuer": "Let's Encrypt",
        "certificate_subject": "web-server-01.acme.com",
        "certificate_valid_from": "2025-01-01T00:00:00Z",
        "certificate_valid_to": "2025-04-01T00:00:00Z",
        "risk_score": 85,
        "risk_level": "high",
        "created_at": "2025-01-01T00:00:00Z"
      }
    ]
  },
  "message": "Crypto implementations retrieved successfully",
  "error": null
}
```

### **Get Risk Summary**
```http
GET /api/v1/risk/summary
```
**Description**: Get risk analysis summary across all assets  
**Authentication**: Required (JWT)  
**Response** (200 OK):
```json
{
  "data": {
    "risk_summary": {
      "total_assets": 150,
      "high_risk_assets": 25,
      "medium_risk_assets": 80,
      "low_risk_assets": 40,
      "unknown_risk_assets": 5,
      "average_risk_score": 65.5,
      "risk_distribution": {
        "high": 16.7,
        "medium": 53.3,
        "low": 26.7,
        "unknown": 3.3
      },
      "top_risks": [
        {
          "risk_type": "Weak Cipher Suites",
          "count": 45,
          "percentage": 30.0
        },
        {
          "risk_type": "Expired Certificates",
          "count": 12,
          "percentage": 8.0
        }
      ],
      "compliance_status": {
        "nist_compliant": 120,
        "fips_compliant": 95,
        "pci_compliant": 110
      }
    }
  },
  "message": "Risk summary retrieved successfully",
  "error": null
}
```

---

## ❌ Error Responses

### **Common Error Codes**

| Code | Description | Example |
|------|-------------|---------|
| 400 | Bad Request | Invalid request body or parameters |
| 401 | Unauthorized | Missing or invalid JWT token |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource not found |
| 409 | Conflict | Resource already exists (e.g., email in use) |
| 422 | Unprocessable Entity | Validation errors |
| 500 | Internal Server Error | Server-side error |

### **Error Response Format**
```json
{
  "data": null,
  "message": "Validation failed",
  "error": {
    "code": "VALIDATION_ERROR",
    "details": {
      "email": ["Email is required"],
      "password": ["Password must be at least 8 characters"]
    }
  },
  "timestamp": "2025-01-09T10:30:00Z"
}
```

### **Authentication Errors**
```json
{
  "data": null,
  "message": "Authentication failed",
  "error": {
    "code": "INVALID_CREDENTIALS",
    "details": "Invalid email or password"
  },
  "timestamp": "2025-01-09T10:30:00Z"
}
```

---

## 🔒 Security Considerations

### **JWT Token Security**
- Access tokens expire in 15 minutes
- Refresh tokens expire in 7 days
- Tokens are stored securely in Redis
- All tokens are invalidated on logout

### **Password Security**
- Passwords are hashed using Argon2id
- Minimum 8 characters with complexity requirements
- Password strength validation on registration

### **Multi-Tenant Isolation**
- All data is isolated by tenant ID
- Users can only access their tenant's data
- JWT tokens include tenant context

### **Rate Limiting**
- 100 requests per minute per IP
- 10 login attempts per minute per IP
- Exponential backoff on failed attempts

---

## 📝 Data Models

### **Asset Model**
```typescript
interface Asset {
  id: string;                    // UUID
  tenant_id: string;             // UUID
  hostname?: string;             // Optional hostname
  ip_address?: string;           // Optional IP address
  port?: number;                 // Optional port number
  asset_type: string;            // server, endpoint, service, appliance
  operating_system?: string;     // Optional OS information
  environment?: string;          // production, staging, development
  business_unit?: string;        // Optional business unit
  owner_email?: string;          // Optional owner email
  description?: string;          // Optional description
  tags: Record<string, any>;     // Key-value tags
  metadata: Record<string, any>; // Additional metadata
  first_discovered_at: string;   // ISO 8601 timestamp
  last_seen_at: string;          // ISO 8601 timestamp
  created_at: string;            // ISO 8601 timestamp
  updated_at: string;            // ISO 8601 timestamp
  deleted_at?: string;           // ISO 8601 timestamp (soft delete)
  
  // Calculated fields
  risk_score: number;            // 0-100 risk score
  risk_level: string;            // high, medium, low, unknown
  crypto_implementations?: CryptoImplementation[];
}
```

### **CryptoImplementation Model**
```typescript
interface CryptoImplementation {
  id: string;                    // UUID
  tenant_id: string;             // UUID
  asset_id: string;              // UUID
  protocol: string;              // TLS, SSH, etc.
  protocol_version?: string;     // 1.3, 2.0, etc.
  cipher_suite?: string;         // TLS_AES_256_GCM_SHA384
  key_exchange_algorithm?: string; // ECDHE, RSA, etc.
  key_size?: number;             // 256, 2048, etc.
  certificate_issuer?: string;   // Certificate authority
  certificate_subject?: string;  // Certificate subject
  certificate_valid_from?: string; // ISO 8601 timestamp
  certificate_valid_to?: string;   // ISO 8601 timestamp
  risk_score: number;            // 0-100 risk score
  risk_level: string;            // high, medium, low, unknown
  created_at: string;            // ISO 8601 timestamp
  updated_at: string;            // ISO 8601 timestamp
}
```

---

## 🚀 Getting Started

### **1. Authentication Flow**
```bash
# 1. Register a new user
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePassword123!",
    "first_name": "Admin",
    "last_name": "User",
    "tenant_name": "My Company"
  }'

# 2. Login to get tokens
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "SecurePassword123!"
  }'

# 3. Use access token for authenticated requests
curl -X GET http://localhost:8082/api/v1/assets \
  -H "Authorization: Bearer <access_token>"
```

### **2. Asset Management**
```bash
# Get all assets with filtering
curl -X GET "http://localhost:8082/api/v1/assets?page=1&limit=20&environment=production&risk_level=high" \
  -H "Authorization: Bearer <access_token>"

# Get specific asset details
curl -X GET http://localhost:8082/api/v1/assets/{asset_id} \
  -H "Authorization: Bearer <access_token>"

# Search assets
curl -X GET "http://localhost:8082/api/v1/assets/search?q=web&asset_type=server" \
  -H "Authorization: Bearer <access_token>"
```

---

## 📞 Support

For API support and questions:
- **Documentation**: This file and inline code comments
- **Health Checks**: Use `/health` endpoints to verify service status
- **Error Handling**: All errors include detailed messages and codes

---

*This API documentation is automatically generated and maintained. Last updated: 2025-01-09*
