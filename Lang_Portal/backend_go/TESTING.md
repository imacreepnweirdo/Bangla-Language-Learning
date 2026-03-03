# API Testing Guide

This document explains how to run and write tests for the Bangla Language Learning Portal API.

## Test Structure

### Unit Tests
Located in `internal/handlers/handlers_test.go` - Tests individual handler functions with in-memory database.

### Integration Tests  
Located in `tests/integration_test.go` - Tests full API endpoints (placeholder for future implementation).

## Running Tests

### Prerequisites
Install test dependencies:
```bash
go mod tidy
```

### Run All Tests
```bash
go test ./...
```

### Run Specific Test File
```bash
go test ./internal/handlers
```

### Run Tests with Verbose Output
```bash
go test -v ./internal/handlers
```

### Run Tests with Coverage
```bash
go test -cover ./internal/handlers
```

### Generate Coverage Report
```bash
go test -coverprofile=coverage.out ./internal/handlers
go tool cover -html=coverage.out -o coverage.html
```

## Test Coverage Areas

### ✅ Currently Tested
- **Dashboard Handlers**: Last study session, study progress, quick stats
- **Word Handlers**: Get words list, get word details, not found cases
- **Group Handlers**: Get groups list, get group details
- **Study Activity Handlers**: Get activity details, create study activity
- **Study Session Handlers**: Review word functionality
- **Reset Handlers**: Reset history, full reset

### 🔄 Test Features
- **In-memory SQLite database** for isolated testing
- **Test data seeding** with realistic Bengali content
- **JSON response validation** against model structures
- **HTTP status code verification**
- **Error case testing** (404, 500, etc.)
- **Pagination testing** where applicable

## Writing New Tests

### Test Structure Template
```go
func TestHandlerName_FunctionName(t *testing.T) {
    // 1. Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    // 2. Seed test data if needed
    seedTestData(t, db)
    
    // 3. Setup Gin test mode
    gin.SetMode(gin.TestMode)
    handler := NewHandlerName(db)
    
    // 4. Create test request
    w := httptest.NewRecorder()
    c, _ := gin.CreateTestContext(w)
    c.Request, _ = http.NewRequest("GET", "/api/endpoint", nil)
    
    // 5. Execute handler
    handler.FunctionName(c)
    
    // 6. Assert results
    assert.Equal(t, http.StatusOK, w.Code)
    
    // 7. Validate response structure
    var response ExpectedResponseType
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, expectedValue, response.Field)
}
```

### Testing POST Requests
```go
func TestPOSTEndpoint(t *testing.T) {
    // ... setup code ...
    
    requestBody := models.RequestType{
        Field1: "value1",
        Field2: 123,
    }
    jsonBody, _ := json.Marshal(requestBody)
    
    c.Request, _ = http.NewRequest("POST", "/api/endpoint", bytes.NewBuffer(jsonBody))
    c.Request.Header.Set("Content-Type", "application/json")
    
    // ... execute and assert ...
}
```

### Testing Error Cases
```go
func TestErrorCase(t *testing.T) {
    // ... setup code ...
    
    // Test invalid ID
    c.Params = gin.Params{{Key: "id", Value: "invalid"}}
    handler.GetByID(c)
    
    assert.Equal(t, http.StatusBadRequest, w.Code)
    
    // Test not found
    c.Params = gin.Params{{Key: "id", Value: "999"}}
    handler.GetByID(c)
    
    assert.Equal(t, http.StatusNotFound, w.Code)
}
```

## Test Data Management

### Database Schema
Tests use the same schema as production via `setupTestDB()` which creates:
- All required tables (words, groups, study_sessions, etc.)
- Foreign key constraints
- Proper indexes

### Test Data
`seedTestData()` creates realistic test data:
- Bengali words with parts of speech
- Group relationships
- Study sessions and activities
- Review items for testing statistics

## Best Practices

### 1. Isolation
- Each test gets a fresh in-memory database
- Tests don't share state
- Use `defer db.Close()` for cleanup

### 2. Realistic Data
- Use actual Bengali words and phrases
- Test with realistic data volumes
- Include edge cases (empty results, invalid data)

### 3. Comprehensive Coverage
- Test success cases
- Test error cases (invalid input, not found)
- Test response structure validation
- Test HTTP status codes

### 4. Performance
- Use in-memory database for speed
- Keep tests focused and small
- Use table-driven tests for multiple cases

## Future Testing Plans

### Integration Tests
- Full HTTP server testing
- Database integration with real SQLite file
- End-to-end API workflows

### Load Testing
- Concurrent request handling
- Database connection pooling
- Memory usage under load

### API Contract Testing
- OpenAPI/Swagger validation
- Response time testing
- Backward compatibility testing

## Troubleshooting

### Common Issues

**Database connection errors**
- Ensure SQLite driver is imported: `_ "github.com/mattn/go-sqlite3"`
- Check in-memory database path: `:memory:`

**JSON parsing errors**
- Verify response models match actual JSON structure
- Check for missing JSON tags in models

**Test isolation problems**
- Ensure each test creates its own database
- Don't share test data between tests

**Import errors**
- Run `go mod tidy` to update dependencies
- Check package imports in test files

### Debugging Tests

Run specific test with verbose output:
```bash
go test -v -run TestSpecificFunction ./internal/handlers
```

Add debugging output:
```go
t.Logf("Request: %+v", request)
t.Logf("Response: %s", w.Body.String())
```

Use testify's `require` for fatal errors:
```go
require.NoError(t, err)  // Stops test on error
assert.NoError(t, err)  // Continues test on error
```
