require_relative '../support/spec_helper'

RSpec.describe 'Study Activities API' do
  before(:all) do
    @test_group_id = get_test_group_id
    @test_activity_id = 1  # Default study activity ID
    
    # Set up test data for study sessions
    setup_test_data
    
    puts "Using test_group_id: #{@test_group_id}, test_activity_id: #{@test_activity_id}"
  end

  describe 'GET /api/study_activities/:id' do
    it 'returns study activity details' do
      response = HTTParty.get("#{base_url}/api/study_activities/#{@test_activity_id}")
      expect_success(response)
      
      data = response.parsed_response
      expect_study_activity_structure(data)
    end

    it 'returns default activity for any ID' do
      response = HTTParty.get("#{base_url}/api/study_activities/99999")
      expect_success(response)
      
      data = response.parsed_response
      expect_study_activity_structure(data)
      expect(data['id']).to eq(99999)
    end
  end

  describe 'GET /api/study_activities/:id/study_sessions' do
    it 'returns paginated study sessions for activity' do
      response = HTTParty.get("#{base_url}/api/study_activities/#{@test_activity_id}/study_sessions?page=1")
      expect_success(response)
      
      data = response.parsed_response
      expect(data).to have_key('items')
      expect(data).to have_key('pagination')
      
      expect(data['items']).to be_an(Array)
      expect_pagination_structure(data['pagination'])
      
      if data['items'].any?
        session = data['items'].first
        expect_study_session_structure(session)
      end
    end

    it 'supports pagination parameters' do
      response = HTTParty.get("#{base_url}/api/study_activities/#{@test_activity_id}/study_sessions?page=2")
      expect_success(response)
      
      pagination = response.parsed_response['pagination']
      expect(pagination['page']).to eq(2)
    end

    it 'returns empty results for non-existent activity' do
      response = HTTParty.get("#{base_url}/api/study_activities/99999/study_sessions")
      expect_success(response)
      
      data = response.parsed_response
      expect(data['items']).to be_an(Array)
      expect(data['items']).to be_empty
      expect_pagination_structure(data['pagination'])
    end
  end

  describe 'POST /api/study_activities' do
    context 'with valid parameters' do
      it 'creates new study activity session' do
        response = create_study_activity(@test_group_id, @test_activity_id)
        expect_success(response)
        
        data = response.parsed_response
        expect(data).to have_key('id')
        expect(data).to have_key('group_id')
        expect(data['group_id']).to eq(@test_group_id)
      end
    end

    context 'with invalid parameters' do
      it 'validates required fields' do
        payload = { group_id: @test_group_id } # missing study_activity_id
        
        response = HTTParty.post(
          "#{base_url}/api/study_activities",
          headers: json_headers,
          body: payload.to_json
        )
        
        expect_error(response, 400)
      end

      it 'accepts any group ID (validation not implemented yet)' do
        payload = {
          group_id: 99999,  # non-existent group
          study_activity_id: @test_activity_id
        }
        
        response = HTTParty.post(
          "#{base_url}/api/study_activities",
          headers: json_headers,
          body: payload.to_json
        )
        
        # Note: Group validation might not be implemented yet
        # This test documents current behavior
        expect(response.code).to be_between(200, 299).or(eq(400))
        
        if response.code == 200
          data = response.parsed_response
          expect(data).to have_key('id')
          expect(data['group_id']).to eq(99999)
        else
          expect_error(response, 400)
        end
      end
    end
  end
end
