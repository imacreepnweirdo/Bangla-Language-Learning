require_relative '../support/spec_helper'

RSpec.describe 'Error Handling API' do
  base_url = 'http://localhost:8080'

  it 'returns 404 for non-existent endpoints' do
    response = HTTParty.get("#{base_url}/api/nonexistent")
    expect(response.code).to eq(404)
  end

  it 'handles invalid JSON in POST requests' do
    response = HTTParty.post(
      "#{base_url}/api/study_activities",
      headers: { 'Content-Type' => 'application/json' },
      body: 'invalid json'
    )
    
    expect(response.code).to eq(400)
  end

  it 'handles CORS preflight requests' do
    response = HTTParty.options(
      "#{base_url}/api/words",
      headers: {
        'Origin' => 'http://localhost:3000',
        'Access-Control-Request-Method' => 'GET'
      }
    )
    
    expect(response.code).to eq(204)
    expect(response.headers['Access-Control-Allow-Origin']).to eq('*')
  end
end
