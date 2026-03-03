require_relative '../support/spec_helper'

RSpec.describe 'Health Check API' do
  base_url = 'http://localhost:8080'

  describe 'GET /health' do
    it 'returns healthy status' do
      response = HTTParty.get("#{base_url}/health")
      
      expect(response.code).to eq(200)
      expect(response.parsed_response['status']).to eq('healthy')
      expect(response.parsed_response['database']).to eq('connected')
    end
  end
end
