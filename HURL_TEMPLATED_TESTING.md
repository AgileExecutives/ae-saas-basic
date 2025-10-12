# HURL Templated Testing System

## Overview

The AE SaaS Basic project now includes a sophisticated templated testing system that automatically generates unique user credentials and identifiers for each test run, preventing conflicts and ensuring test isolation.

## How It Works

### Unique Identifier Generation

Each test run generates:
- `{{UNIQUE_ID}}`: Timestamp + hash (e.g., `1760208981_e157c10b`)
- `{{UNIQUE_USERNAME}}`: `testuser_1760208981_e157c10b`
- `{{UNIQUE_EMAIL}}`: `test_1760208981_e157c10b@example.com`
- `{{UNIQUE_CUSTOMER}}`: `customer_1760208981_e157c10b`
- `{{UNIQUE_ORG}}`: `org_1760208981_e157c10b`
- `{{UNIQUE_PASSWORD}}`: `Pass123_e157c10b`
- `{{HOST}}`: Server URL (e.g., `http://localhost:8080`)

### Directory Structure

```
tests/hurl/
├── templates/           # Template files with {{VARIABLES}}
│   ├── 00_quick_integration.hurl
│   ├── 01_health.hurl
│   ├── 02_auth.hurl
│   └── 04_customers.hurl
├── processed/           # Auto-generated files with substituted values
└── hurl.config         # Configuration file
```

## Template Variables

Use these variables in your HURL template files:

### User Authentication
- `{{UNIQUE_USERNAME}}` - Unique username for test user
- `{{UNIQUE_EMAIL}}` - Unique email address
- `{{UNIQUE_PASSWORD}}` - Unique password

### Identifiers
- `{{UNIQUE_ID}}` - Base unique identifier
- `{{UNIQUE_CUSTOMER}}` - Unique customer identifier
- `{{UNIQUE_ORG}}` - Unique organization identifier

### Environment
- `{{HOST}}` or `{{host}}` - Server host URL

## Example Template

```hurl
# Authentication Test Template
POST {{HOST}}/api/v1/auth/register
Content-Type: application/json
{
  "username": "{{UNIQUE_USERNAME}}",
  "email": "{{UNIQUE_EMAIL}}",
  "password": "{{UNIQUE_PASSWORD}}",
  "first_name": "Test",
  "last_name": "User",
  "organization_id": 1
}

HTTP 201
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.user.username" == "{{UNIQUE_USERNAME}}"
jsonpath "$.data.user.email" == "{{UNIQUE_EMAIL}}"

[Captures]
auth_token: jsonpath "$.data.token"

# Use the captured token in subsequent requests
GET {{HOST}}/api/v1/user-settings
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
```

## Running Tests

### Basic Usage
```bash
./run-hurl-tests.sh
```

### What Happens
1. **Unique ID Generation**: Creates timestamp-based unique identifiers
2. **Template Processing**: Converts template files to executable HURL files
3. **Test Execution**: Runs processed tests with unique data
4. **Results**: Saves JSON results and provides colorized output

### Example Output
```
🚀 Starting AE SaaS Basic HURL Tests with Templating
Host: http://localhost:8080
Results Directory: test_results
Unique Test ID: 1760208981_e157c10b
Test Username: testuser_1760208981_e157c10b
Test Email: test_1760208981_e157c10b@example.com

🔍 Checking server availability...
✅ Server is running

🧪 Running 00_quick_integration.hurl...
✅ 00_quick_integration.hurl passed
🧪 Running 01_health.hurl...
✅ 01_health.hurl passed
```

## Benefits

### Test Isolation
- Each test run uses completely unique user credentials
- No conflicts between concurrent or sequential test runs
- Clean database state for each test execution

### Realistic Testing
- Tests actual user registration and authentication flows
- Validates unique constraint handling
- Tests real API behavior with dynamic data

### Debugging
- Processed files saved in `tests/hurl/processed/` for inspection
- JSON results saved for detailed error analysis
- Clear colorized output shows pass/fail status

## Creating New Templates

1. Create a new template file in `tests/hurl/templates/`
2. Use template variables for dynamic data
3. Follow existing naming convention (e.g., `05_new_feature.hurl`)
4. Run tests - the system auto-processes templates

## Troubleshooting

### Common Issues
- **User conflicts**: Templates should prevent this, but check if using hardcoded values
- **Email validation**: Ensure email format is valid (no special characters)
- **Server not running**: The script checks server availability first

### Debugging
- Check `tests/hurl/processed/` for actual values being used
- Examine `test_results/*.json` for detailed error information
- Run individual tests: `hurl tests/hurl/processed/02_auth.hurl --test`

## Migration from Static Tests

To convert existing static HURL tests to templates:

1. Copy test file to `tests/hurl/templates/`
2. Replace hardcoded usernames with `{{UNIQUE_USERNAME}}`
3. Replace hardcoded emails with `{{UNIQUE_EMAIL}}`
4. Replace hardcoded passwords with `{{UNIQUE_PASSWORD}}`
5. Replace host URLs with `{{HOST}}`
6. Update assertions to use template variables where needed

This system ensures robust, reliable API testing with proper test isolation and realistic data scenarios.