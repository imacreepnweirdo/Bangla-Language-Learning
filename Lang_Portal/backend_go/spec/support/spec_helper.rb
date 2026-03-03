require 'rspec'
require 'httparty'
require 'json'

RSpec.configure do |config|
  config.color = true
  config.formatter = :documentation
  
  # Enable test ordering for better isolation
  config.order = :random
  
  # Include HTTP helpers
  config.include HTTParty
  
  # Global before/after hooks
  config.before(:suite) do
    puts "=" * 60
    puts "Bangla Language Portal API Test Suite"
    puts "=" * 60
    puts "Make sure your Go server is running on http://localhost:8080"
    puts "=" * 60
  end
  
  config.after(:suite) do
    puts "=" * 60
    puts "API tests completed!"
    puts "=" * 60
  end
end

# Custom matchers for API testing
RSpec::Matchers.define :have_api_structure do |expected_keys|
  match do |actual|
    expected_keys.all? { |key| actual.key?(key.to_s) }
  end
  
  failure_message do |actual|
    "Expected response to have keys: #{expected_keys.join(', ')}"
  end
end

# Helper methods
module ApiHelpers
  def base_url
    'http://localhost:8080'
  end
  
  def json_headers
    { 'Content-Type' => 'application/json' }
  end
  
  def expect_success(response)
    expect(response.code).to be_between(200, 299)
  end
  
  def expect_error(response, status_code = 400)
    expect(response.code).to eq(status_code)
    expect(response.parsed_response).to have_key('error')
  end
  
  def expect_pagination(response)
    data = response.parsed_response
    expect(data).to have_key('pagination')
    pagination = data['pagination']
    expect(pagination).to have_api_structure(['page', 'total_pages', 'total_items', 'items_per_page'])
  end
end

# Test data setup helpers
module TestDataHelpers
  # Known test data IDs from test database (words.test.db)
  TEST_GROUP_ID = 1   # basic_greetings group in test database
  TEST_WORD_ID = 1    # First word in test database
  
  def get_test_group_id
    TEST_GROUP_ID
  end
  
  def get_test_word_id
    TEST_WORD_ID
  end
  
  def create_test_session(group_id = nil, activity_id = 1)
    group_id ||= get_test_group_id
    payload = {
      group_id: group_id,
      study_activity_id: activity_id
    }
    
    response = HTTParty.post(
      "#{base_url}/api/study_activities",
      headers: json_headers,
      body: payload.to_json
    )
    
    response.code == 200 ? response.parsed_response['id'] : 1
  end
  
  def setup_test_data
    # Create a few test study sessions for testing
    session_id_1 = create_test_session(TEST_GROUP_ID, 1)
    session_id_2 = create_test_session(TEST_GROUP_ID, 2)
    
    # Add some word reviews to the sessions
    if session_id_1 && session_id_1 != 1
      review_word(session_id_1, TEST_WORD_ID, true)
      review_word(session_id_1, TEST_WORD_ID + 1, false)
    end
    
    if session_id_2 && session_id_2 != 1
      review_word(session_id_2, TEST_WORD_ID + 2, true)
    end
    
    puts "Created test study sessions: #{session_id_1}, #{session_id_2}"
  end
end

# Common API endpoint helpers
module EndpointHelpers
  def get_dashboard_stats
    HTTParty.get("#{base_url}/api/dashboard/quick-stats")
  end
  
  def get_last_study_session
    HTTParty.get("#{base_url}/api/dashboard/last_study_session")
  end
  
  def get_study_progress
    HTTParty.get("#{base_url}/api/dashboard/study_progress")
  end
  
  def get_words(page = 1)
    HTTParty.get("#{base_url}/api/words?page=#{page}")
  end
  
  def get_word(word_id)
    HTTParty.get("#{base_url}/api/words/#{word_id}")
  end
  
  def get_groups(page = 1)
    HTTParty.get("#{base_url}/api/groups?page=#{page}")
  end
  
  def get_group(group_id)
    HTTParty.get("#{base_url}/api/groups/#{group_id}")
  end
  
  def get_group_words(group_id, page = 1)
    HTTParty.get("#{base_url}/api/groups/#{group_id}/words?page=#{page}")
  end
  
  def create_study_activity(group_id, activity_id)
    payload = {
      group_id: group_id,
      study_activity_id: activity_id
    }
    
    HTTParty.post(
      "#{base_url}/api/study_activities",
      headers: json_headers,
      body: payload.to_json
    )
  end
  
  def review_word(session_id, word_id, correct)
    payload = { "Correct": correct }
    
    HTTParty.post(
      "#{base_url}/api/study_sessions/#{session_id}/words/#{word_id}/review",
      headers: json_headers,
      body: payload.to_json
    )
  end
  
  def reset_history
    HTTParty.post("#{base_url}/api/reset_history")
  end
  
  def full_reset
    HTTParty.post("#{base_url}/api/full_reset")
  end
  
  def health_check
    HTTParty.get("#{base_url}/health")
  end
end

# Response validation helpers
module ValidationHelpers
  def expect_word_structure(word)
    # Based on Backend_Technical_Specs.md lines 166-174
    expect(word).to have_api_structure(['bengali', 'parts_of_speech', 'english', 'correct_count', 'wrong_count'])
  end
  
  def expect_word_detail_structure(word)
    # Based on Backend_Technical_Specs.md lines 187-202
    expect(word).to have_api_structure(['id', 'bengali', 'english', 'stats', 'items'])
    expect(word['stats']).to have_api_structure(['correct_count', 'wrong_count'])
    expect(word['items']).to be_an(Array)
  end
  
  def expect_group_word_structure(word)
    # Based on Backend_Technical_Specs.md lines 247-254 (group words don't have parts_of_speech)
    expect(word).to have_api_structure(['id', 'bengali', 'english', 'correct_count', 'wrong_count'])
  end
  
  def expect_group_structure(group)
    # Based on Backend_Technical_Specs.md lines 212-216
    expect(group).to have_api_structure(['id', 'name', 'words_count'])
  end
  
  def expect_group_detail_structure(group)
    # Based on Backend_Technical_Specs.md lines 230-238
    expect(group).to have_api_structure(['id', 'name', 'words', 'stats'])
    expect(group['stats']).to have_api_structure(['total_word_count'])
    expect(group['words']).to be_an(Array)
  end
  
  def expect_dashboard_stats_structure(stats)
    # Based on Backend_Technical_Specs.md lines 98-107
    expect(stats).to have_api_structure(['total_words', 'total_groups', 'words_learned', 'sessions_completed', 'current_streak', 'accuracy_rate'])
  end
  
  def expect_study_activity_structure(activity)
    # Based on Backend_Technical_Specs.md lines 112-119
    expect(activity).to have_api_structure(['id', 'name', 'thumbnail_url', 'description'])
  end
  
  def expect_study_session_structure(session)
    # Based on Backend_Technical_Specs.md lines 127-135
    expect(session).to have_api_structure(['id', 'activity_name', 'group_name', 'start_time', 'end_time', 'review_items_count'])
  end
  
  def expect_review_structure(review)
    # Based on Backend_Technical_Specs.md lines 300-307
    expect(review).to have_api_structure(['success', 'word_id', 'study_session_id', 'correct', 'created_at'])
  end
  
  def expect_reset_structure(response)
    # Based on Backend_Technical_Specs.md lines 268-272 and 278-282
    expect(response).to have_api_structure(['success', 'message'])
  end
  
  def expect_last_study_session_structure(session)
    # Based on Backend_Technical_Specs.md lines 72-80
    expect(session).to have_api_structure(['id', 'group_id', 'created_at', 'study_activities_id', 'group_name'])
  end
  
  def expect_study_progress_structure(progress)
    # Based on Backend_Technical_Specs.md lines 87-92
    expect(progress).to have_api_structure(['total_words_studied', 'total_available_words'])
  end
  
  def expect_pagination_structure(pagination)
    # Based on Backend_Technical_Specs.md lines 137-142 (study_sessions)
    # and lines 175-180 (words), lines 218-223 (groups)
    # Note: some use 'current_page', others use 'page' - need to handle both
    valid_keys = ['current_page', 'page', 'total_pages', 'total_items', 'items_per_page']
    expect(pagination.keys.all? { |key| valid_keys.include?(key) }).to be true
    expect(pagination).to have_key('total_pages')
    expect(pagination).to have_key('total_items')
    expect(pagination).to have_key('items_per_page')
  end
end

RSpec.configure do |config|
  config.include ApiHelpers
  config.include TestDataHelpers
  config.include EndpointHelpers
  config.include ValidationHelpers
end
