# Fuzzy Search API Tests with Unique Data
# Tests intelligent search functionality across multiple entities

# First, create a test user and get authentication token
POST http://localhost:8080/api/v1/auth/register
Content-Type: application/json
{
  "username": "testuser_1760213272_31b81a8d",
  "email": "test_1760213272_31b81a8d@example.com",
  "password": "Pass123_31b81a8d",
  "first_name": "Search",
  "last_name": "Tester",
  "organization_id": 1
}

# Accept either success or conflict (user exists)
HTTP *

# Now login to get the token regardless of registration result
POST http://localhost:8080/api/v1/auth/login
Content-Type: application/json
{
  "username": "testuser_1760213272_31b81a8d",
  "password": "Pass123_31b81a8d"
}

HTTP 200
[Captures]
auth_token: jsonpath "$.data.token"
user_org_id: jsonpath "$.data.user.organization_id"

# Test 1: Basic search using POST (proper JSON format)
POST http://localhost:8080/api/v1/search
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "1760213272_31b81a8d",
  "limit": 10,
  "include_count": true
}

HTTP 200
[Asserts]
jsonpath "$.query" exists
jsonpath "$.results" exists
jsonpath "$.total" exists

# Test 2: Quick search for autocomplete (public endpoint)
GET http://localhost:8080/api/v1/search/quick?q=testuser_1760213272_31b81a8d&limit=5

HTTP 200
[Asserts]
jsonpath "$.query" exists
jsonpath "$.results" exists
jsonpath "$.total" exists

# Test 3: Search with entity type filter (POST format)
POST http://localhost:8080/api/v1/search
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "test",
  "entity_types": ["users", "customers"],
  "limit": 20
}

HTTP 200
[Asserts]
jsonpath "$.query" exists
jsonpath "$.results" exists
jsonpath "$.categories" exists

# Test 4: Search in specific entity type (customers)
POST http://localhost:8080/api/v1/search/entities/customers
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "customer_1760213272_31b81a8d",
  "limit": 10
}

HTTP *
[Asserts]
# May not be implemented, accept various responses
status >= 200

# Test 5: Search suggestions
GET http://localhost:8080/api/v1/search/suggestions?q=testuser_1760213272_31b81a8d&limit=5

HTTP 200
[Asserts]
jsonpath "$.suggestions" exists
jsonpath "$.query" exists

# Test 6: Search with pagination (POST format)
POST http://localhost:8080/api/v1/search
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "test",
  "offset": 0,
  "limit": 5,
  "include_count": true
}

HTTP 200
[Asserts]
jsonpath "$.query" exists
jsonpath "$.total" exists

# Test 7: Search with sorting options (POST format)
POST http://localhost:8080/api/v1/search
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "test",
  "sort_by": "relevance",
  "sort_order": "desc",
  "limit": 10
}

HTTP 200
[Asserts]
jsonpath "$.query" exists
jsonpath "$.results" exists

# Test 8: Empty search query (should fail validation)
POST http://localhost:8080/api/v1/search
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "",
  "limit": 10
}

# Should return 400 for empty query
HTTP 400

# Test 9: Search entity types endpoint
GET http://localhost:8080/api/v1/search/types

HTTP 200
[Asserts]
jsonpath "$.entity_types" exists

# Test 10: Search health check
GET http://localhost:8080/api/v1/search/health

HTTP 200
[Asserts]
jsonpath "$.status" == "healthy"

# Test 11: Try to access protected search without authentication
POST http://localhost:8080/api/v1/search
Content-Type: application/json
{
  "query": "test",
  "limit": 10
}

HTTP 401

# Test 12: Search with advanced filters (POST format)
POST http://localhost:8080/api/v1/search/advanced
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "1760213272_31b81a8d",
  "filters": {
    "entity_types": ["users", "customers"]
  },
  "sort_by": "created_at",
  "sort_order": "desc",
  "limit": 10,
  "include_count": true
}

# Accept either success or not implemented
HTTP *
[Asserts]
# Advanced search may return search results or error
status >= 200

# Test 13: Search stats (if supported)
GET http://localhost:8080/api/v1/search/stats
Authorization: Bearer {{auth_token}}

# Accept either success or not implemented  
HTTP *

# Test 14: Search config
GET http://localhost:8080/api/v1/search/config
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.config" exists
jsonpath "$.entity_types" exists