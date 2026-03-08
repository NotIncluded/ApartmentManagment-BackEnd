# E2E Test Environment Setup & Configuration

## Overview

This document describes how to set up the environment for running E2E tests, including database setup, environment variables, and test data preparation.

## Prerequisites

### 1. Server Requirements

- Go 1.16+ installed
- PostgreSQL or MySQL database
- `.env.local` file in project root

### 2. Testing Tools

- VS Code with REST Client extension (or similar HTTP client)
- Database GUI client (optional, for manual verification)
- Postman/Insomnia (optional, as alternative to REST Client)

## Environment Setup

### Step 1: Create `.env.local` file

```bash
# Create file in project root: .env.local
PORT=8080
JWT_SECRET=your-secret-key-change-this-in-production
DATABASE_URL=postgres://user:password@localhost:5432/apartment_system
```

**Important**:

- Change `JWT_SECRET` to a secure random string in production
- Ensure database credentials are correct
- PORT must match `@baseUrl` in test files (default: 8080)

### Step 2: Start the Server

```bash
cd c:\Users\punna\OneDrive\Desktop\WORK\OOAD-Apartment
go run ./main.go
```

**Expected Output**:

```
Starting server on :8080
Database connected successfully
Server running...
```

### Step 3: Verify Server is Running

```bash
# Quick health check
curl http://localhost:8080/health

# Expected response:
# {"status":200,"version":"0.0.1"}
```

## Test Data Management

### Database State Management

**Before Each Test Run:**

1. Start with clean database (optional, depends on test goals)
2. Pre-load default data if needed
3. Clear session state/tokens

**Database Reset** (if needed):

```sql
-- Run against your database to reset test data
DELETE FROM users;
DELETE FROM contracts;
DELETE FROM bills;
DELETE FROM utility_usages;
DELETE FROM payments;
```

### Pre-loaded Test Users (Optional)

For quick testing, you can pre-load an ADMIN user:

```sql
INSERT INTO users (id, name, phone, email, password, role, created_at, updated_at)
VALUES (
  'admin-uuid-here',
  'System Admin',
  '0811111111',
  'admin@system.com',
  'AdminPassword123',  -- Note: should be hashed in production
  'ADMIN',
  NOW(),
  NOW()
);
```

## Test Execution Workflows

### Workflow 1: Basic Authentication Test

**Time**: ~5 minutes

```
1. Health Check (TEST 1)
2. Register Tenant (TEST 2)
3. Login (TEST 5)
4. Login with Wrong Password (TEST 10)
5. Verify token-based access (TEST 6)
```

### Workflow 2: Authorization Testing

**Time**: ~10 minutes

```
1. Register Tenant (TEST 2)
2. Login Tenant (TEST 5)
3. Try Admin-only Operations (TEST 7, 8, 9)
4. Test Role-based Rejection (TEST 14)
5. Verify error responses match expectations
```

### Workflow 3: Error Handling

**Time**: ~8 minutes

```
1. Invalid Credentials (TEST 10, 11)
2. Missing Fields (TEST 12)
3. Duplicate Email (TEST 16)
4. No Token (TEST 6)
5. Invalid Token (TEST 15)
```

### Workflow 4: Full End-to-End

**Time**: ~15 minutes

Run all tests in sequence (TEST 1-16) to verify complete flow.

## VS Code REST Client Configuration

### Installation

1. Open VS Code
2. Extensions → Search "REST Client"
3. Install by Huachao Mao
4. Reload VS Code

### Settings (Optional)

Add to `.vscode/settings.json`:

```json
{
  "rest-client.timeout": 10000,
  "rest-client.showResponseInDifferentTab": true,
  "rest-client.previewResponseInUnicodeFormat": true,
  "rest-client.excludeHostsForProxy": ["localhost", "127.0.0.1"]
}
```

### Usage

1. Open `auth-and-user-management.http` file
2. Click "Send Request" above any test
3. View response in side panel
4. Variables (like `@tenantToken`) auto-populate

## Alternative Test Tools

### Using Postman

**Import Steps:**

1. Open Postman → File → Import
2. Select `auth-and-user-management.http`
3. Create environment with variables:
   - `baseUrl` = `http://localhost:8080`
   - `tenantToken` = (auto-populate from login)
4. Run requests individually or in collection

**Collection Setup:**

1. Create collection "Apartment API Tests"
2. Add folder "Authentication"
3. Add folder "User Management"
4. Organize requests accordingly

### Using cURL

**Example - Register:**

```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "phone": "0812345678",
    "email": "test@example.com",
    "password": "TestPassword123",
    "role": "TENANT"
  }'
```

**Example - Login:**

```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "TestPassword123"
  }'
```

**Example - Protected Route:**

```bash
curl -X POST http://localhost:8080/user/create \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "name": "New User",
    "phone": "0899999999",
    "email": "newuser@example.com",
    "password": "Password123",
    "role": "TENANT"
  }'
```

## Continuous Integration Options

### GitHub Actions Example

Create `.github/workflows/e2e-tests.yml`:

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: postgres
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.16

      - name: Start Server
        run: go run ./main.go &
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost/test_db
          JWT_SECRET: test-secret
          PORT: 8080

      - name: Wait for Server
        run: sleep 5

      - name: Run Tests
        run: |
          # Run HTTP tests using REST client or similar
          # This would import and execute the test files
```

## Debugging Tests

### Enable Verbose Logging

Add to test requests:

```http
X-Debug: true
X-Log-Level: debug
```

### Common Issues & Solutions

**Issue: "Address already in use"**

```bash
# Find process using port 8080
lsof -i :8080
# Kill process
kill -9 <PID>
```

**Issue: "Database connection failed"**

```bash
# Verify database is running
# Check credentials in .env.local
# Ensure database exists
createdb apartment_system
```

**Issue: Tests pass locally but fail in CI**

- Ensure all environment variables are set in CI
- Check database initialization script
- Verify dependencies are installed
- Check Go version compatibility

## Test Maintenance

### Weekly Review

- Check for deprecated endpoints
- Update test data if schema changes
- Review error messages for clarity
- Update documentation if needed

### Monthly Review

- Analyze test coverage
- Add new test scenarios
- Remove obsolete tests
- Performance optimization

## Best Practices

### ✅ DO:

- Run tests on clean database state
- Use meaningful variable names
- Add clear comments for each test
- Verify both success AND failure cases
- Keep tests isolated and independent
- Document expected vs actual responses
- Version control test files

### ❌ DON'T:

- Hardcode credentials in test files
- Use production data in testing
- Skip error scenario testing
- Leave commented-out tests in code
- Run tests against production database
- Ignore test failures
- Mix different test types in one file

## Performance Considerations

### Response Time Expectations

- Registration: < 200ms
- Login: < 200ms
- Create User: < 200ms
- Delete User: < 100ms
- Protected Route Access: < 50ms

### Load Testing (Future)

Consider using tools like:

- Apache JMeter
- Locust
- k6 (Grafana)
- Artillery

## Security Testing Checklist

- [ ] SQL Injection attempts blocked
- [ ] Token expiration enforced
- [ ] Invalid tokens rejected
- [ ] Role-based access enforced
- [ ] Privilege escalation prevented
- [ ] Password requirements validated
- [ ] Duplicate email prevention
- [ ] Rate limiting (if implemented)
- [ ] CORS properly configured
- [ ] HTTPS enforced (in production)

## Next Steps

1. Set up environment as described above
2. Start running tests from `auth-and-user-management.http`
3. Refer to `TEST_GUIDE.md` for detailed test documentation
4. Use `quick-reference.http` for quick testing
5. Integrate with CI/CD pipeline for automated testing
6. Add more test scenarios as system evolves

## Support & Contact

For issues or improvements:

- Check test output for error messages
- Review API implementation code
- Consult system documentation
- Create issue tracking entry if bug found
