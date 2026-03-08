# Apartment System - E2E Test Guide

## Overview

This guide provides comprehensive end-to-end tests for the Apartment System's authentication and user management endpoints. Tests are organized by scenario and can be executed using VS Code REST Client extension or similar HTTP testing tools.

## Setup Requirements

### Prerequisites

1. Go server running on `http://localhost:8080` (default)
2. VS Code REST Client extension installed (or similar HTTP client)
3. Database configured and running
4. `.env.local` file with `JWT_SECRET` configured

### Environment Variables

All tests use the following base configuration:

- **Base URL**: `http://localhost:8080`
- **Admin Email**: `admin@test.com`
- **Admin Password**: `AdminPassword123`
- **Tenant 1 Email**: `tenant@test.com`
- **Tenant 1 Password**: `TenantPassword123`
- **Tenant 2 Email**: `tenant2@test.com`
- **Tenant 2 Password**: `TenantPassword456`

## Test Scenarios

### 1. Health Check

**File**: `auth-and-user-management.http` - TEST 1
**Purpose**: Verify server is running
**Endpoint**: `GET /health`
**Expected Response**:

```json
{
  "status": 200,
  "version": "0.0.1"
}
```

### 2. User Registration Tests

#### 2.1 Register Tenant - Success

**File**: `auth-and-user-management.http` - TEST 2
**Purpose**: Register a new tenant user (normal use case)
**Endpoint**: `POST /auth/register`
**Request Body**:

```json
{
  "name": "John Tenant",
  "phone": "0812345678",
  "email": "tenant@test.com",
  "password": "TenantPassword123",
  "role": "TENANT"
}
```

**Expected Response**: HTTP 200 with user data
**Key Points**:

- Registration sets role to TENANT regardless of request
- Email must be unique
- All fields are required

#### 2.2 Register Admin - FAILURE

**File**: `auth-and-user-management.http` - TEST 4
**Purpose**: Demonstrate that ADMIN role cannot be self-assigned during registration
**Endpoint**: `POST /auth/register`
**Request Body**:

```json
{
  "name": "Admin User",
  "phone": "0811111111",
  "email": "admin@test.com",
  "password": "AdminPassword123",
  "role": "ADMIN"
}
```

**Expected Response**: HTTP 400
**Error Message**: "Cannot set ADMIN role during tenant registration"
**Key Points**:

- System prevents privilege escalation
- Only existing ADMINs can create other ADMINs via API
- This is a security control

#### 2.3 Register with Duplicate Email - FAILURE

**File**: `auth-and-user-management.http` - TEST 16
**Purpose**: Ensure email uniqueness constraint
**Endpoint**: `POST /auth/register`
**Expected Response**: HTTP 400
**Error Message**: "email already exists"

#### 2.4 Register with Missing Fields - FAILURE

**File**: `auth-and-user-management.http` - TEST 12
**Purpose**: Validate request body requirements
**Endpoint**: `POST /auth/register`
**Expected Response**: HTTP 400
**Error Message**: "Invalid request body"

### 3. Authentication Tests

#### 3.1 Login Success

**File**: `auth-and-user-management.http` - TEST 5
**Purpose**: Login with valid credentials and receive JWT token
**Endpoint**: `POST /auth/login`
**Request Body**:

```json
{
  "email": "tenant@test.com",
  "password": "TenantPassword123"
}
```

**Expected Response**: HTTP 200
**Response Body**:

```json
{
  "status": 200,
  "message": "Access token generated successfully",
  "data": {
    "access_token": "eyJhbGc..."
  }
}
```

**Key Points**:

- Token is stored in `@tenantToken` variable for reuse
- Token expires based on JWT configuration
- Token contains email and role information

#### 3.2 Login with Invalid Password - FAILURE

**File**: `auth-and-user-management.http` - TEST 10
**Purpose**: Reject login with wrong password
**Endpoint**: `POST /auth/login`
**Expected Response**: HTTP 401
**Error Message**: "Invalid username or password"

#### 3.3 Login with Non-existent Email - FAILURE

**File**: `auth-and-user-management.http` - TEST 11
**Purpose**: Reject login for non-existent user
**Endpoint**: `POST /auth/login`
**Expected Response**: HTTP 401
**Error Message**: "Invalid username or password"

### 4. User Management Tests

#### 4.1 Create User - ADMIN Only

**File**: `auth-and-user-management.http` - TEST 8
**Purpose**: Create a new user (restricted to ADMIN role)
**Endpoint**: `POST /user/create`
**Headers Required**: `Authorization: Bearer {admin_token}`
**Request Body**:

```json
{
  "name": "New User",
  "phone": "0899999999",
  "email": "newuser@test.com",
  "password": "Password123",
  "role": "TENANT"
}
```

**Expected Response**:

- HTTP 202 (if ADMIN token is used)
- HTTP 401/403 (if non-ADMIN token or no token)

#### 4.2 Delete User - ADMIN Only

**File**: `auth-and-user-management.http` - TEST 14
**Purpose**: Delete a user (restricted to ADMIN role)
**Endpoint**: `DELETE /user/{userId}`
**Headers Required**: `Authorization: Bearer {admin_token}`
**Expected Response**:

- HTTP 200 (if ADMIN token is used)
- HTTP 401/403 (if non-ADMIN token or no token)
- HTTP 400 (if user ID is invalid)

### 5. Middleware & Role-Based Access Control Tests

#### 5.1 Protected Route - No Token

**File**: `auth-and-user-management.http` - TEST 6
**Purpose**: Verify protected routes reject unauthenticated requests
**Endpoint**: `POST /user/create` (without Authorization header)
**Expected Response**: HTTP 401 (Unauthorized)
**Key Points**:

- All endpoints under `/user` group require ADMIN token
- Missing token is rejected at middleware level

#### 5.2 Protected Route - Invalid Token

**File**: `auth-and-user-management.http` - TEST 15
**Purpose**: Verify malformed/invalid tokens are rejected
**Endpoint**: `POST /user/create`
**Headers**: `Authorization: Bearer invalid.token.here`
**Expected Response**: HTTP 401 (Unauthorized)

#### 5.3 Protected Route - Wrong Role

**File**: `auth-and-user-management.http` - TEST 7, 9
**Purpose**: Verify role-based access control (TENANT cannot access ADMIN endpoints)
**Endpoint**: `POST /user/create`
**Headers**: `Authorization: Bearer {tenant_token}`
**Expected Response**: HTTP 401 or 403 (Forbidden)
**Key Points**:

- Even with valid token, non-ADMIN users are rejected
- Middleware checks both token validity and role

## Running Tests Sequentially

### Recommended Test Order

1. **TEST 1** - Health Check (verify server is up)
2. **TEST 2** - Register Tenant #1 (creates main test user)
3. **TEST 3** - Register Tenant #2 (for delete testing)
4. **TEST 4** - Register as ADMIN (should fail)
5. **TEST 5** - Login Tenant (stores token for later use)
6. **TEST 10** - Login with Wrong Password (should fail)
7. **TEST 11** - Login Non-existent User (should fail)
8. **TEST 6** - Protected Route without Token (should fail)
9. **TEST 15** - Protected Route with Invalid Token (should fail)
10. **TEST 7** - Protected Route with Tenant Token (should fail)
11. **TEST 12** - Register with Missing Fields (should fail)
12. **TEST 16** - Register Duplicate Email (should fail)
13. **TEST 13** - Delete without Token (should fail)
14. **TEST 14** - Delete with Tenant Token (should fail)

### Using VS Code REST Client

1. Open `auth-and-user-management.http` in VS Code
2. Click "Send Request" above each test case
3. View response in the side panel
4. Variables like `@tenantToken` are automatically captured from responses

## Response Format

All responses follow this format:

```json
{
  "status": 200,
  "message": "User registered successfully",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Tenant",
    "email": "tenant@test.com",
    "role": "TENANT",
    "created_at": "2026-03-08T10:30:00Z"
  }
}
```

## Security Considerations

### 1. Admin Privilege Escalation Prevention

- Users cannot set ADMIN role during registration
- Only existing ADMINs can create other ADMINs
- System rejects any attempt to bypass this

### 2. Role-Based Access Control

- `/user/*` endpoints require ADMIN token
- Middleware validates both authentication and authorization
- Non-ADMIN users receive 401/403 errors

### 3. Token Security

- Tokens contain user email and role
- Tokens are signed with JWT_SECRET
- Invalid/expired tokens are rejected
- Token validation happens at middleware level

### 4. Password Security

- Passwords are stored (note: should be hashed in production)
- Login validates exact password match
- Failed logins don't reveal user existence

## Troubleshooting

### Issue: "Connection refused"

**Solution**: Ensure server is running on localhost:8080

```bash
go run ./main.go
```

### Issue: "Invalid JWT_SECRET"

**Solution**: Verify `.env.local` has valid `JWT_SECRET`

```bash
echo "JWT_SECRET=your-secret-key" >> .env.local
```

### Issue: "email already exists"

**Solution**: Registration uses unique emails

- Change email in test request
- Or clear database and restart tests

### Issue: Token not captured in variable

**Solution**:

1. Verify login returns 200 status
2. Check response has `data.access_token` field
3. Manually copy token if needed

### Issue: "jwt token missing or invalid"

**Solution**:

1. Ensure token is correctly formatted in header
2. Check token is not expired
3. Verify bearer keyword is present: `Bearer {token}`

## Integration with Postman/Insomnia

If using Postman/Insomnia instead of VS Code REST Client:

1. Import the `.http` file or create requests manually
2. Set up collection-level variables for tokens
3. Use pre-request scripts to manage token assignment
4. Create test suites for sequential execution

## Next Steps

- Implement actual ADMIN user creation endpoint
- Add password hashing (bcrypt)
- Implement token expiration
- Add rate limiting for login attempts
- Implement token refresh mechanism
- Add request validation layer
- Consider two-factor authentication

## Contact & Support

For issues or questions about these tests, refer to:

- System architecture documentation
- API specification in code comments
- Git commit history for recent changes
