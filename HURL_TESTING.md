# HURL Testing Documentation

This document describes the HURL test suite for AE SaaS Basic, covering PDF generation, static assets, and fuzzy search functionality.

## Overview

The HURL test suite provides comprehensive API testing for all endpoints in the AE SaaS Basic system. Tests are organized into logical groups and can be run individually or as a complete suite.

## Test Files

### Core System Tests
- `01_health.hurl` - Health check and basic connectivity tests
- `02_auth.hurl` - Authentication, registration, and user management tests

### Business Logic Tests  
- `03_plans.hurl` - Subscription plans and pricing tests
- `04_customers.hurl` - Customer management and billing tests
- `05_contacts.hurl` - Contact management and communication tests
- `06_emails.hurl` - Email system and template tests
- `07_user_settings.hurl` - User preferences and configuration tests

### New Feature Tests
- `08_pdf_generation.hurl` - **PDF generation system tests**
- `09_static_assets.hurl` - **Static file serving tests**

## Prerequisites

### System Requirements
- [HURL](https://hurl.dev/) installed
- Running AE SaaS Basic server
- PostgreSQL database running (for full functionality)
- wkhtmltopdf installed (for PDF generation tests)

### Installation
```bash
# Install HURL (macOS)
brew install hurl

# Install HURL (Ubuntu)
curl -LO https://github.com/Orange-OpenSource/hurl/releases/latest/download/hurl_x.x.x_amd64.deb
sudo dpkg -i hurl_x.x.x_amd64.deb

# Install wkhtmltopdf (for PDF tests)
# macOS
brew install wkhtmltopdf

# Ubuntu
sudo apt-get install wkhtmltopdf
```

## Running Tests

### Complete Test Suite
```bash
# Run all tests
./run-hurl-tests.sh

# Run with custom host
TEST_HOST=http://localhost:3000 ./run-hurl-tests.sh
```

### Individual Test Files
```bash
# Run specific test file
hurl --variable host=http://localhost:8080 tests/hurl/08_pdf_generation.hurl

# Run with verbose output
hurl --variable host=http://localhost:8080 --very-verbose tests/hurl/08_pdf_generation.hurl

# Save results to file
hurl --variable host=http://localhost:8080 --json --output results.json tests/hurl/09_static_assets.hurl
```

### Environment Variables
- `TEST_HOST` - Server host (default: http://localhost:8080)
- `HURL_VERBOSE` - Enable verbose output (default: false)

## Test Details

### PDF Generation Tests (`08_pdf_generation.hurl`)

Tests the comprehensive PDF generation system including:

#### Configuration Tests
- Get PDF configuration settings
- Validate template directories and settings

#### Template Management
- List available PDF templates
- Get specific template information
- Preview templates with sample data

#### PDF Generation
- Generate PDF from template with data
- Generate PDF from raw HTML
- Stream PDF generation for large documents
- Test various page sizes and orientations

#### Error Handling
- Invalid template names
- Missing required data
- Malformed HTML content
- Authentication failures

#### Sample Test Commands
```bash
# Test PDF generation specifically
hurl --variable host=http://localhost:8080 tests/hurl/08_pdf_generation.hurl

# Test only template listing
hurl --variable host=http://localhost:8080 \
     --include-only "Test listing available PDF templates" \
     tests/hurl/08_pdf_generation.hurl
```

### Static Assets Tests (`09_static_assets.hurl`)

Tests static file serving capabilities including:

#### Asset Types
- CSS stylesheets (`/statics/css/`)
- JavaScript files (`/statics/js/`)
- Images (`/statics/images/`)
- Email templates (`/statics/email_templates/`)
- PDF templates (`/statics/templates/pdf/`)
- Fonts (`/statics/fonts/`)

#### Security Tests
- Directory traversal protection
- Access control for sensitive files
- Case sensitivity handling
- Special character handling

#### Performance Tests
- Caching headers validation
- Compression support
- Multiple concurrent requests

#### Sample Test Commands
```bash
# Test static assets
hurl --variable host=http://localhost:8080 tests/hurl/09_static_assets.hurl

# Test only security aspects
hurl --variable host=http://localhost:8080 \
     --include-only "directory traversal" \
     tests/hurl/09_static_assets.hurl
```



## Test Data Setup

### Authentication
Tests automatically create test users and obtain authentication tokens:

```hurl
# Standard user for most tests
POST {{host}}/api/v1/auth/login
Content-Type: application/json
{
  "username": "testuser",
  "password": "newpass123"
}

# Admin user for administrative tests
POST {{host}}/api/v1/auth/register
Content-Type: application/json
{
  "username": "testadmin",
  "email": "admin@example.com",
  "password": "adminpass123",
  "role": "admin"
}
```

### Sample Data
Tests create their own sample data as needed:

- Customers with various attributes
- Contacts with different types
- Plans with different features
- Email templates
- Search queries and preferences

## Expected Results

### Success Criteria
All tests should pass when:
- Server is running and accessible
- Database is connected and migrated
- Required dependencies (wkhtmltopdf) are installed
- Static files are properly served

### Common Issues

#### PDF Generation Failures
```
Error: wkhtmltopdf not found
Solution: Install wkhtmltopdf system dependency
```

#### Static Asset 404 Errors
```
Error: Static files not found
Solution: Ensure statics directory exists with proper structure
```

#### Authentication Failures
```
Error: JWT secret not configured
Solution: Set JWT_SECRET environment variable
```

#### Database Connection Issues
```
Error: Database connection failed
Solution: Start PostgreSQL and verify connection settings
```

## Test Results

### Output Format
Test results are saved in JSON format in the `test_results/` directory:

```json
{
  "entries": [
    {
      "request": {
        "method": "GET",
        "url": "http://localhost:8080/api/v1/search"
      },
      "response": {
        "status": 200,
        "headers": [...],
        "body": {...}
      }
    }
  ]
}
```

### Analyzing Results
```bash
# View test summary
cat test_results/09_static_assets.json | jq '.entries[].response.status'

# Check for failures
grep -r "error" test_results/

# View specific test details
cat test_results/08_pdf_generation.json | jq '.entries[-1]'
```

## Continuous Integration

### GitHub Actions Integration
```yaml
name: HURL Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Start services
        run: docker-compose up -d
      - name: Install HURL
        run: curl -LO https://github.com/Orange-OpenSource/hurl/releases/latest/download/hurl_x.x.x_amd64.deb && sudo dpkg -i hurl_x.x.x_amd64.deb
      - name: Run tests
        run: ./run-hurl-tests.sh
```

### Local Development
```bash
# Watch mode for development
while inotifywait -e modify tests/hurl/; do
  ./run-hurl-tests.sh
done

# Run specific tests during development
hurl --variable host=http://localhost:8080 \
     --very-verbose \
     tests/hurl/09_static_assets.hurl
```

## Extending Tests

### Adding New Test Files
1. Create new `.hurl` file in `tests/hurl/`
2. Add to `TEST_FILES` array in `run-hurl-tests.sh`
3. Follow existing patterns for authentication and assertions

### Test Best Practices
- Always test error conditions
- Include authentication setup in each file
- Use descriptive test names
- Validate both success and failure scenarios
- Test boundary conditions (empty queries, large limits)
- Include security tests (SQL injection, directory traversal)

### Sample New Test Structure
```hurl
# Feature Tests
# Description of what this feature does

# Setup authentication
POST {{host}}/api/v1/auth/login
Content-Type: application/json
{
  "username": "testuser",
  "password": "newpass123"
}

HTTP 200
[Captures]
auth_token: jsonpath "$.data.token"

# Test positive case
GET {{host}}/api/v1/new-feature
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true

# Test error case
GET {{host}}/api/v1/new-feature/invalid

HTTP 404
[Asserts]
jsonpath "$.success" == false
```

This comprehensive test suite ensures that all features (PDF generation and static assets) are thoroughly tested and working correctly in the AE SaaS Basic system.

**Note**: Fuzzy search functionality has been moved to the backend module. See the backend test suite for fuzzy search and general search API tests.