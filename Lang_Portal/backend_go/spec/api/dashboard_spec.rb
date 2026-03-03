require_relative '../support/spec_helper'

RSpec.describe 'Dashboard API' do
  describe 'GET /api/dashboard/last_study_session' do
    it 'returns last study session info' do
      response = get_last_study_session
      expect_success(response)
      
      data = response.parsed_response
      
      # Handle the case where no study sessions exist
      if data && data['message'] == 'No study sessions found'
        # This is valid when no sessions exist
        expect(data).to have_key('message')
      elsif data && data.keys.any?
        # When sessions exist, expect the full structure
        expect_last_study_session_structure(data)
      else
        # Handle empty/null response
        expect(data).to be_nil.or(eq({}))
      end
    end
  end

  describe 'GET /api/dashboard/study_progress' do
    it 'returns study progress statistics' do
      response = get_study_progress
      expect_success(response)
      
      data = response.parsed_response
      expect_study_progress_structure(data)
      
      expect(data['total_words_studied']).to be_a(Integer)
      expect(data['total_available_words']).to be_a(Integer)
    end
  end

  describe 'GET /api/dashboard/quick-stats' do
    it 'returns quick statistics' do
      response = get_dashboard_stats
      expect_success(response)
      
      data = response.parsed_response
      expect_dashboard_stats_structure(data)
      
      expect(data['total_words']).to be_a(Integer)
      expect(data['total_groups']).to be_a(Integer)
      expect(data['words_learned']).to be_a(Integer)
      expect(data['sessions_completed']).to be_a(Integer)
      expect(data['current_streak']).to be_a(Integer)
      # accuracy_rate might be Integer or Float depending on implementation
      expect(data['accuracy_rate']).to be_a(Numeric).or be_nil
    end
  end
end
