# Fuzzy Search Tests
# Test intelligent search functionality across multiple entity types

# Setup authentication for search tests
POST {{host}}/api/v1/auth/login
Content-Type: application/json
{
  "username": "testuser",
  "password": "newpass123"
}

HTTP 200
[Captures]
auth_token: jsonpath "$.data.token"
user_org_id: jsonpath "$.data.user.organization_id"

# Test basic search across all entities
GET {{host}}/api/v1/search?q=test&limit=10
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.query" == "test"
jsonpath "$.data.total_results" isInteger
jsonpath "$.data.results" isCollection

# Test quick search for autocomplete functionality
GET {{host}}/api/v1/search/quick?q=te&limit=5
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.query" == "te"
jsonpath "$.data.results" isCollection
jsonpath "$.data.results" count <= 5

# Test search with specific entity types
GET {{host}}/api/v1/search?q=user&entity_types=user,customer&limit=15
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.query" == "user"
jsonpath "$.data.results" isCollection

# Test advanced search with filters
POST {{host}}/api/v1/search/advanced
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "test",
  "entity_types": ["user", "customer", "contact"],
  "filters": {
    "date_range": {
      "field": "created_at",
      "start": "2024-01-01",
      "end": "2024-12-31"
    },
    "boolean_filters": {
      "active": true
    }
  },
  "sort_by": "relevance",
  "limit": 20
}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.query" == "test"
jsonpath "$.data.filters" exists
jsonpath "$.data.results" isCollection

# Test advanced search with numeric range filters
POST {{host}}/api/v1/search/advanced
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "customer",
  "entity_types": ["customer"],
  "filters": {
    "numeric_range": {
      "field": "id",
      "min": 1,
      "max": 100
    }
  },
  "limit": 10
}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.results" isCollection

# Test entity-specific search endpoints
GET {{host}}/api/v1/search/users?q=test&limit=5
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.entity_type" == "user"
jsonpath "$.data.results" isCollection

# Test customer-specific search
GET {{host}}/api/v1/search/customers?q=inc&status=active&limit=10
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.entity_type" == "customer"
jsonpath "$.data.results" isCollection

# Test contact-specific search
GET {{host}}/api/v1/search/contacts?q=support&limit=8
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.entity_type" == "contact"
jsonpath "$.data.results" isCollection

# Test plan-specific search
GET {{host}}/api/v1/search/plans?q=premium&limit=5
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.entity_type" == "plan"
jsonpath "$.data.results" isCollection

# Test email-specific search
GET {{host}}/api/v1/search/emails?q=welcome&limit=10
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.entity_type" == "email"
jsonpath "$.data.results" isCollection

# Test getting user search preferences
GET {{host}}/api/v1/search/preferences
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.user_id" exists
jsonpath "$.data.preferred_entities" isCollection
jsonpath "$.data.results_per_page" isInteger
jsonpath "$.data.highlight_enabled" isBoolean

# Test updating user search preferences
POST {{host}}/api/v1/search/preferences
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "preferred_entities": ["user", "customer", "contact"],
  "results_per_page": 25,
  "highlight_enabled": true,
  "save_history": true,
  "default_sort": "relevance"
}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.preferred_entities" contains "user"
jsonpath "$.data.preferred_entities" contains "customer"
jsonpath "$.data.results_per_page" == 25

# Test getting saved searches
GET {{host}}/api/v1/search/saved
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data" isCollection

# Test saving a search query
POST {{host}}/api/v1/search/saved
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "name": "Test Search",
  "query": "test user",
  "entity_types": ["user"],
  "filters": {
    "boolean_filters": {
      "active": true
    }
  },
  "description": "Search for active test users"
}

HTTP 201
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.name" == "Test Search"
jsonpath "$.data.query" == "test user"

[Captures]
saved_search_id: jsonpath "$.data.id"

# Test getting saved searches after creating one
GET {{host}}/api/v1/search/saved
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data" count >= 1

# Test deleting a saved search
DELETE {{host}}/api/v1/search/saved/{{saved_search_id}}
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.message" contains "deleted"

# Test search with minimum query length validation
GET {{host}}/api/v1/search?q=a&limit=10
Authorization: Bearer {{auth_token}}

HTTP 400
[Asserts]
jsonpath "$.success" == false
jsonpath "$.error" contains "minimum"

# Test search with empty query
GET {{host}}/api/v1/search?q=&limit=10
Authorization: Bearer {{auth_token}}

HTTP 400
[Asserts]
jsonpath "$.success" == false
jsonpath "$.error" contains "query"

# Test search with invalid entity type
GET {{host}}/api/v1/search?q=test&entity_types=invalid_entity
Authorization: Bearer {{auth_token}}

HTTP 400
[Asserts]
jsonpath "$.success" == false
jsonpath "$.error" contains "entity"

# Test search without authentication
GET {{host}}/api/v1/search?q=test

HTTP 401
[Asserts]
jsonpath "$.success" == false
jsonpath "$.error" contains "authorization"

# Test advanced search with invalid filters
POST {{host}}/api/v1/search/advanced
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "test",
  "filters": {
    "invalid_filter": {
      "field": "nonexistent",
      "value": "invalid"
    }
  }
}

HTTP 400
[Asserts]
jsonpath "$.success" == false

# Test search with extremely large limit
GET {{host}}/api/v1/search?q=test&limit=10000
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
# Should be capped at max allowed limit

# Test search with negative limit
GET {{host}}/api/v1/search?q=test&limit=-5
Authorization: Bearer {{auth_token}}

HTTP 400
[Asserts]
jsonpath "$.success" == false

# Test search with special characters in query
GET {{host}}/api/v1/search?q=test%40email.com&limit=10
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true

# Test search with SQL injection attempt
GET {{host}}/api/v1/search?q='; DROP TABLE users; --&limit=10
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
# Should handle safely without SQL injection

# Test updating preferences with invalid data
POST {{host}}/api/v1/search/preferences
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "results_per_page": -10,
  "preferred_entities": ["invalid_entity"]
}

HTTP 400
[Asserts]
jsonpath "$.success" == false

# Test saving search with missing required fields
POST {{host}}/api/v1/search/saved
Authorization: Bearer {{auth_token}}
Content-Type: application/json
{
  "query": "test"
}

HTTP 400
[Asserts]
jsonpath "$.success" == false
jsonpath "$.error" contains "name"

# Setup for admin tests - create admin user
POST {{host}}/api/v1/auth/register
Content-Type: application/json
{
  "username": "testadmin",
  "email": "admin@example.com",
  "password": "adminpass123",
  "first_name": "Test",
  "last_name": "Admin",
  "organization_name": "Admin Organization",
  "role": "admin"
}

HTTP 201
[Captures]
admin_token: jsonpath "$.data.token"

# Test admin endpoint - get entity types
GET {{host}}/api/v1/admin/search/entities
Authorization: Bearer {{admin_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data" isCollection
jsonpath "$.data" count >= 5

# Test admin endpoint - get search analytics
GET {{host}}/api/v1/admin/search/analytics
Authorization: Bearer {{admin_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.total_searches" isInteger
jsonpath "$.data.popular_queries" isCollection
jsonpath "$.data.entity_popularity" isCollection

# Test admin endpoint - update search configuration
PUT {{host}}/api/v1/admin/search/config
Authorization: Bearer {{admin_token}}
Content-Type: application/json
{
  "min_search_length": 3,
  "max_results": 100,
  "score_threshold": 0.4,
  "enable_highlight": true,
  "cache_results": true,
  "cache_timeout_min": 60
}

HTTP 200
[Asserts]
jsonpath "$.success" == true
jsonpath "$.data.min_search_length" == 3
jsonpath "$.data.max_results" == 100

# Test admin endpoint access with regular user token
GET {{host}}/api/v1/admin/search/entities
Authorization: Bearer {{auth_token}}

HTTP 403
[Asserts]
jsonpath "$.success" == false
jsonpath "$.error" contains "admin"

# Test admin endpoint without authentication
GET {{host}}/api/v1/admin/search/entities

HTTP 401
[Asserts]
jsonpath "$.success" == false
jsonpath "$.error" contains "authorization"

# Test search performance with common queries
GET {{host}}/api/v1/search?q=user&limit=50
Authorization: Bearer {{auth_token}}

HTTP 200
[Asserts]
jsonpath "$.success" == true
# Response should be reasonably fast

# Test concurrent search requests
GET {{host}}/api/v1/search?q=test1&limit=10
Authorization: Bearer {{auth_token}}

HTTP 200

GET {{host}}/api/v1/search?q=test2&limit=10
Authorization: Bearer {{auth_token}}

HTTP 200

GET {{host}}/api/v1/search?q=test3&limit=10
Authorization: Bearer {{auth_token}}

HTTP 200