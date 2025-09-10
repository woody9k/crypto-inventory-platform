# API Documentation

*Last Updated: 2025-01-09*
*Platform: Crypto Inventory Management System*

## 🎯 Overview

This document provides comprehensive API documentation for all services in the crypto inventory management platform, including the newly implemented SaaS admin service.

## 📋 Table of Contents

1. [Authentication Service API](#authentication-service-api)
2. [Inventory Service API](#inventory-service-api)
3. [SaaS Admin Service API](#saas-admin-service-api)
4. [Common Response Formats](#common-response-formats)
5. [Error Handling](#error-handling)
6. [Authentication & Authorization](#authentication--authorization)

---

## Authentication Service API

**Base URL**: `http://localhost:8081`
**Service**: Tenant authentication and user management

### Authentication Endpoints

#### POST /api/v1/auth/register
Register a new tenant user.

**Request Body**:
```json
{
  "email": "user@example.com",
  "password": "SecurePassword123!",
  "first_name": "John",
  "last_name": "Doe",
  "tenant_name": "Example Corp",
  "subscription_tier": "professional"
}
```

**Response** (201 Created):
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "tenant_id": "uuid",
    "role": "tenant_owner",
    "created_at": "2025-01-09T10:00:00Z"
  }
}
```

#### POST /api/v1/auth/login
Authenticate a tenant user.

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
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "tenant_id": "uuid",
    "role": "tenant_owner",
    "is_active": true,
    "email_verified": false,
    "last_login_at": "2025-01-09T10:00:00Z",
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:00:00Z"
  },
  "access_token": "jwt_token",
  "refresh_token": "jwt_refresh_token",
  "expires_in": 86400
}
```

#### POST /api/v1/auth/refresh
Refresh an access token.

**Request Body**:
```json
{
  "refresh_token": "jwt_refresh_token"
}
```

**Response** (200 OK):
```json
{
  "access_token": "new_jwt_token",
  "expires_in": 86400
}
```

### User Management Endpoints

#### GET /api/v1/users
Get current user information.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "user": {
    "id": "uuid",
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "tenant_id": "uuid",
    "role": "tenant_owner",
    "is_active": true,
    "email_verified": false,
    "last_login_at": "2025-01-09T10:00:00Z",
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:00:00Z"
  }
}
```

---

## Inventory Service API

**Base URL**: `http://localhost:8082`
**Service**: Asset and sensor management

### Asset Endpoints

#### GET /api/v1/assets
Get list of network assets.

**Headers**: `Authorization: Bearer <token>`
**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)
- `search` (optional): Search term
- `type` (optional): Asset type filter

**Response** (200 OK):
```json
{
  "assets": [
    {
      "id": "uuid",
      "name": "Web Server 01",
      "type": "server",
      "ip_address": "192.168.1.100",
      "status": "active",
      "last_seen": "2025-01-09T10:00:00Z",
      "created_at": "2025-01-09T10:00:00Z",
      "updated_at": "2025-01-09T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "pages": 5
  }
}
```

#### GET /api/v1/assets/:id
Get specific asset details.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "asset": {
    "id": "uuid",
    "name": "Web Server 01",
    "type": "server",
    "ip_address": "192.168.1.100",
    "status": "active",
    "last_seen": "2025-01-09T10:00:00Z",
    "metadata": {
      "os": "Ubuntu 20.04",
      "cpu": "Intel Xeon",
      "memory": "16GB"
    },
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:00:00Z"
  }
}
```

### Sensor Endpoints

#### GET /api/v1/sensors
Get list of sensors.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "sensors": [
    {
      "id": "uuid",
      "name": "Sensor 01",
      "type": "network_monitor",
      "status": "active",
      "last_heartbeat": "2025-01-09T10:00:00Z",
      "created_at": "2025-01-09T10:00:00Z"
    }
  ]
}
```

---

## SaaS Admin Service API

**Base URL**: `http://localhost:8084`
**Service**: Platform administration and tenant management

### Authentication Endpoints

#### POST /api/v1/auth/login
Authenticate a platform administrator.

**Request Body**:
```json
{
  "email": "admin@crypto-inventory.com",
  "password": "admin123"
}
```

**Response** (200 OK):
```json
{
  "user": {
    "id": "uuid",
    "email": "admin@crypto-inventory.com",
    "first_name": "Platform",
    "last_name": "Administrator",
    "role": "super_admin",
    "is_active": true,
    "email_verified": true,
    "last_login_at": "2025-01-09T10:00:00Z",
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:00:00Z"
  },
  "access_token": "jwt_token",
  "refresh_token": "jwt_refresh_token",
  "expires_in": 86400
}
```

### Tenant Management Endpoints

#### GET /api/v1/admin/tenants
Get list of all tenants.

**Headers**: `Authorization: Bearer <token>`
**Query Parameters**:
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)

**Response** (200 OK):
```json
{
  "tenants": [
    {
      "id": "uuid",
      "name": "Demo Corporation",
      "slug": "demo-corp",
      "domain": "demo.example.com",
      "subscription_tier": "professional",
      "trial_ends_at": "2025-02-09T10:00:00Z",
      "billing_email": "billing@demo.com",
      "payment_status": "active",
      "stripe_customer_id": "cus_xxx",
      "sso_enabled": false,
      "is_active": true,
      "created_at": "2025-01-09T10:00:00Z",
      "updated_at": "2025-01-09T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 50
  }
}
```

#### GET /api/v1/admin/tenants/:id
Get specific tenant details.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "tenant": {
    "id": "uuid",
    "name": "Demo Corporation",
    "slug": "demo-corp",
    "domain": "demo.example.com",
    "subscription_tier": "professional",
    "trial_ends_at": "2025-02-09T10:00:00Z",
    "billing_email": "billing@demo.com",
    "payment_status": "active",
    "stripe_customer_id": "cus_xxx",
    "sso_enabled": false,
    "is_active": true,
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:00:00Z"
  }
}
```

#### POST /api/v1/admin/tenants
Create a new tenant.

**Headers**: `Authorization: Bearer <token>`
**Request Body**:
```json
{
  "name": "New Corporation",
  "slug": "new-corp",
  "domain": "new.example.com",
  "subscription_tier": "professional",
  "billing_email": "billing@new.com"
}
```

**Response** (201 Created):
```json
{
  "message": "Tenant created successfully",
  "tenant_id": "uuid"
}
```

#### PUT /api/v1/admin/tenants/:id
Update tenant information.

**Headers**: `Authorization: Bearer <token>`
**Request Body**:
```json
{
  "name": "Updated Corporation",
  "domain": "updated.example.com",
  "billing_email": "billing@updated.com",
  "payment_status": "active"
}
```

**Response** (200 OK):
```json
{
  "message": "Tenant updated successfully"
}
```

#### DELETE /api/v1/admin/tenants/:id
Delete a tenant (soft delete).

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "message": "Tenant deleted successfully"
}
```

#### POST /api/v1/admin/tenants/:id/suspend
Suspend a tenant.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "message": "Tenant suspended successfully"
}
```

#### POST /api/v1/admin/tenants/:id/activate
Activate a tenant.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "message": "Tenant activated successfully"
}
```

#### GET /api/v1/admin/tenants/:id/stats
Get tenant statistics.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "stats": {
    "tenant_id": "uuid",
    "tenant_name": "Demo Corporation",
    "user_count": 25,
    "asset_count": 150,
    "sensor_count": 5,
    "last_activity": "2025-01-09T10:00:00Z",
    "storage_used": 1024000,
    "api_requests": 5000
  }
}
```

### Platform User Management Endpoints

#### GET /api/v1/admin/users
Get list of platform users.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "users": [
    {
      "id": "uuid",
      "email": "admin@crypto-inventory.com",
      "first_name": "Platform",
      "last_name": "Administrator",
      "role": "super_admin",
      "is_active": true,
      "email_verified": true,
      "last_login_at": "2025-01-09T10:00:00Z",
      "created_at": "2025-01-09T10:00:00Z",
      "updated_at": "2025-01-09T10:00:00Z"
    }
  ]
}
```

#### GET /api/v1/admin/users/:id
Get specific platform user details.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "user": {
    "id": "uuid",
    "email": "admin@crypto-inventory.com",
    "first_name": "Platform",
    "last_name": "Administrator",
    "role": "super_admin",
    "is_active": true,
    "email_verified": true,
    "last_login_at": "2025-01-09T10:00:00Z",
    "created_at": "2025-01-09T10:00:00Z",
    "updated_at": "2025-01-09T10:00:00Z"
  }
}
```

#### POST /api/v1/admin/users
Create a new platform user.

**Headers**: `Authorization: Bearer <token>`
**Request Body**:
```json
{
  "email": "newadmin@crypto-inventory.com",
  "password": "SecurePassword123!",
  "first_name": "New",
  "last_name": "Admin",
  "role": "platform_admin"
}
```

**Response** (201 Created):
```json
{
  "message": "Platform user created successfully",
  "user_id": "uuid"
}
```

### Platform Statistics Endpoints

#### GET /api/v1/admin/stats/platform
Get platform-wide statistics.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "stats": {
    "total_tenants": 50,
    "active_tenants": 45,
    "total_users": 1250,
    "total_assets": 5000,
    "total_sensors": 200,
    "total_api_requests": 100000,
    "storage_used": 50000000,
    "revenue": 250000
  }
}
```

#### GET /api/v1/admin/stats/tenants
Get statistics for all tenants.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "tenants_stats": [
    {
      "tenant_id": "uuid",
      "tenant_name": "Demo Corporation",
      "user_count": 25,
      "asset_count": 150,
      "sensor_count": 5,
      "last_activity": "2025-01-09T10:00:00Z",
      "storage_used": 1024000,
      "api_requests": 5000
    }
  ]
}
```

### System Monitoring Endpoints

#### GET /api/v1/admin/monitoring/health
Get system health status.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "database": "healthy",
  "timestamp": "2025-01-09T10:00:00Z"
}
```

#### GET /api/v1/admin/monitoring/logs
Get system logs.

**Headers**: `Authorization: Bearer <token>`

**Response** (200 OK):
```json
{
  "logs": [],
  "message": "System logs endpoint - to be implemented"
}
```

---

## Common Response Formats

### Success Response
```json
{
  "data": { ... },
  "message": "Operation successful",
  "timestamp": "2025-01-09T10:00:00Z"
}
```

### Error Response
```json
{
  "error": "Error message",
  "code": "ERROR_CODE",
  "details": { ... },
  "timestamp": "2025-01-09T10:00:00Z"
}
```

### Pagination Response
```json
{
  "data": [ ... ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "pages": 5
  }
}
```

---

## Error Handling

### HTTP Status Codes

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required
- `403 Forbidden` - Insufficient permissions
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., duplicate email)
- `422 Unprocessable Entity` - Validation error
- `500 Internal Server Error` - Server error

### Error Response Format

```json
{
  "error": "Validation failed",
  "code": "VALIDATION_ERROR",
  "details": {
    "field": "email",
    "message": "Invalid email format"
  },
  "timestamp": "2025-01-09T10:00:00Z"
}
```

---

## Authentication & Authorization

### JWT Token Format

```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "tenant_id": "uuid",  // For tenant users
  "role": "tenant_owner", // or platform role
  "type": "access",
  "exp": 1641234567,
  "iat": 1641148167
}
```

### Authorization Headers

```
Authorization: Bearer <jwt_token>
```

### Role-Based Access Control

**Tenant Roles**:
- `tenant_owner` - Full tenant access
- `tenant_admin` - Tenant management
- `security_admin` - Security settings
- `analyst` - Data analysis

**Platform Roles**:
- `super_admin` - Full platform access
- `platform_admin` - Platform management
- `support_admin` - Support and monitoring

---

*This API documentation should be updated as new endpoints are added or existing ones are modified. Last updated: 2025-01-09*