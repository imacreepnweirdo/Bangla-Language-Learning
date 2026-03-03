# Ruby + RSpec API Testing Guide

This guide explains how to use Ruby with RSpec to test the Bangla Language Learning Portal API.

## Why Ruby + RSpec?

**Advantages over Go testing:**
- **More readable syntax** - Natural language descriptions
- **Better error messages** - Clear failure descriptions
- **Rich ecosystem** - Many testing libraries available
- **Flexible matchers** - Easy to write custom assertions
- **Less boilerplate** - More concise test code

## Setup

### Prerequisites
- Ruby 2.7+ installed
- Your Go server running on `http://localhost:8080`

### Installation
```bash
# Install Ruby gems
bundle install

# Or without bundler
gem install rspec httparty json
```

## Running Tests

### Start Your Go Server First
```bash
# In one terminal
cd backend_go
go run cmd/server/main.go
```

### Run All Tests
```bash
# In another terminal
bundle exec rspec
```

### Run Specific Test Groups
```bash
# Test only dashboard endpoints
bundle exec rspec -e "Dashboard Endpoints"

# Test only words endpoints
bundle exec rspec -e "Words Endpoints"

# Run with line numbers
bundle exec rspec spec/api_spec.rb:25
```

### Run with Different Output Formats
```bash
# Documentation format (default)
bundle exec rspec --format documentation

# Progress format
bundle exec rspec --format progress

# JUnit XML for CI/CD
bundle exec rspec --format RspecJunitFormatter --out rspec.xml
```

## Test Structure

### File Organization
```
spec/
├── spec_helper.rb          # Global configuration
├── api_spec.rb             # Main API tests
└── support/                # Additional helpers
    ├── api_helpers.rb
    └── matchers.rb
```

### Test Organization
```ruby
describe 'API Feature' do
  describe 'GET /endpoint' do
    context 'when valid request' do
      it 'returns expected response' do
        # Test implementation
      end
    end
    
    context 'when invalid request' do
      it 'returns error response' do
        # Error testing
      end
    end
  end
end
```

## Test Examples

### Basic GET Request Test
```ruby
describe 'GET /api/words' do
  it 'returns paginated words list' do
    response = HTTParty.get("#{base_url}/api/words?page=1")
    
    expect_success(response)
    expect_pagination(response)
    
    data = response.parsed_response
    expect(data['items']).to be_an(Array)
  end
end
```

### POST Request Test
```ruby
describe 'POST /api/study_activities' do
  it 'creates new study activity' do
    payload = {
      group_id: 1,
      study_activity_id: 1
    }
    
    response = HTTParty.post(
      "#{base_url}/api/study_activities",
      headers: json_headers,
      body: payload.to_json
    )
    
    expect_success(response)
    
    data = response.parsed_response
    expect(data).to have_key('id')
    expect(data['group_id']).to eq(1)
  end
end
```

### Error Testing
```ruby
describe 'Error Handling' do
  it 'returns 404 for non-existent word' do
    response = HTTParty.get("#{base_url}/api/words/99999")
    expect_error(response, 404)
  end
  
  it 'validates required fields' do
    payload = { group_id: 1 } # missing study_activity_id
    
    response = HTTParty.post(
      "#{base_url}/api/study_activities",
      headers: json_headers,
      body: payload.to_json
    )
    
    expect_error(response, 400)
  end
end
```

## Custom Matchers

### Built-in Matchers
```ruby
# Check response structure
expect(data).to have_api_structure(['id', 'name', 'created_at'])

# Check pagination
expect_pagination(response)

# Check success/error
expect_success(response)
expect_error(response, 404)
```

### Custom Matcher Example
```ruby
RSpec::Matchers.define :be_bengali_word do
  match do |actual|
    actual.match?(/[\u0980-\u09FF]/) # Bengali Unicode range
  end
  
  failure_message do |actual|
    "Expected '#{actual}' to contain Bengali characters"
  end
end

# Usage
expect(word['bengali']).to be_bengali_word
```

## Advanced Testing

### Shared Examples
```ruby
# Define reusable test patterns
RSpec.shared_examples 'paginated response' do |endpoint|
  it 'returns paginated structure' do
    response = HTTParty.get("#{base_url}#{endpoint}?page=1")
    expect_pagination(response)
  end
end

# Use in multiple tests
describe 'Words API' do
  it_behaves_like 'paginated response', '/api/words'
end

describe 'Groups API' do
  it_behaves_like 'paginated response', '/api/groups'
end
```

### Data Setup and Cleanup
```ruby
describe 'Study Sessions' do
  before(:all) do
    # Setup test data once
    @test_session_id = create_test_session
  end
  
  after(:all) do
    # Cleanup test data
    cleanup_test_data
  end
  
  it 'works with test session' do
    # Test using @test_session_id
  end
end
```

### Performance Testing
```ruby
describe 'Performance' do
  it 'responds quickly' do
    start_time = Time.now
    response = HTTParty.get("#{base_url}/api/words")
    end_time = Time.now
    
    expect_success(response)
    expect(end_time - start_time).to be < 0.5 # 500ms limit
  end
end
```

## Integration with CI/CD

### GitHub Actions Example
```yaml
name: API Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      app:
        image: golang:1.21
        command: ./server
        ports: [8080:8080]
    
    steps:
      - uses: actions/checkout@v3
      - uses: ruby/setup-ruby@v1
        with:
          ruby-version: '3.2'
      - run: bundle install
      - run: bundle exec rspec --format RspecJunitFormatter --out rspec.xml
```

## Debugging

### Pry Integration
```ruby
# Add to Gemfile
gem 'pry'

# Use in tests
it 'debugs response' do
  response = HTTParty.get("#{base_url}/api/words")
  
  binding.pry # Debug here
  
  expect_success(response)
end
```

### Verbose Output
```bash
# Show all test output
bundle exec rspec --format documentation --backtrace

# Show HTTP requests/responses
bundle exec rspec --format documentation --require spec_helper
```

## Best Practices

### 1. Test Organization
- Group related tests in `describe` blocks
- Use `context` for different scenarios
- Write descriptive test names

### 2. Data Management
- Use test data factories
- Clean up after tests
- Don't rely on specific data ordering

### 3. Error Testing
- Test both success and failure cases
- Verify HTTP status codes
- Check error message formats

### 4. Performance
- Keep tests fast
- Use shared setup where appropriate
- Avoid unnecessary HTTP calls

### 5. Maintenance
- Keep tests DRY (Don't Repeat Yourself)
- Use meaningful variable names
- Add comments for complex logic

## Troubleshooting

### Common Issues

**Server not running**
```
Error: Connection refused
```
Solution: Start your Go server first

**JSON parsing errors**
```
Error: JSON::ParserError
```
Solution: Check response body and content-type

**Timeout errors**
```
Error: Net::ReadTimeout
```
Solution: Increase timeout or check server performance

### Debug Commands
```bash
# Run specific failing test
bundle exec rspec spec/api_spec.rb:123

# Run with backtrace
bundle exec rspec --backtrace

# Run with debugging
bundle exec rspec --require pry
```

## Next Steps

1. **Add more comprehensive tests** for edge cases
2. **Set up CI/CD integration** for automated testing
3. **Add performance benchmarks** for regression testing
4. **Create test data factories** for complex scenarios
5. **Add API contract testing** for backward compatibility
