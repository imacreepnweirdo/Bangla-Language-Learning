require_relative '../support/spec_helper'

RSpec.describe 'Reset API' do
  describe 'POST /api/reset_history' do
    it 'resets study history successfully' do
      response = reset_history
      
      expect_success(response)
      
      data = response.parsed_response
      expect(data).to have_key('success')
      expect(data).to have_key('message')
      expect(data['success']).to be true
      expect(data['message']).to eq('Study history reset successfully')
    end
  end

  describe 'POST /api/full_reset' do
    it 'performs full system reset successfully' do
      response = full_reset
      
      expect_success(response)
      
      data = response.parsed_response
      expect(data).to have_key('success')
      expect(data).to have_key('message')
      expect(data['success']).to be true
      expect(data['message']).to eq('System has been fully reset')
    end
  end
end
