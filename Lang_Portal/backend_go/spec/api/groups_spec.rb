require_relative '../support/spec_helper'

RSpec.describe 'Groups API' do
  before(:all) do
    @test_group_id = get_test_group_id
    puts "Using test_group_id: #{@test_group_id}"
  end

  describe 'GET /api/groups' do
    it 'returns paginated groups list' do
      response = get_groups(1)
      expect_success(response)
      
      data = response.parsed_response
      expect(data).to have_key('items')
      expect(data).to have_key('pagination')
      
      expect(data['items']).to be_an(Array)
      expect_pagination_structure(data['pagination'])
      
      if data['items'].any?
        group = data['items'].first
        expect_group_structure(group)
      end
    end
  end

  describe 'GET /api/groups/:id' do
    context 'when group exists' do
      it 'returns group details with stats' do
        response = get_group(@test_group_id)
        expect_success(response)
        
        data = response.parsed_response
        expect_group_detail_structure(data)
      end
    end

    context 'when group does not exist' do
      it 'returns 404 error' do
        response = get_group(99999)
        expect_error(response, 404)
      end
    end
  end

  describe 'GET /api/groups/:id/words' do
    context 'when group exists' do
      it 'returns paginated words for group' do
        response = get_group_words(@test_group_id, 1)
        expect_success(response)
        
        data = response.parsed_response
        expect(data).to have_key('words')
        expect(data).to have_key('pagination')
        
        expect(data['words']).to be_an(Array)
        expect_pagination_structure(data['pagination'])
        
        if data['words'].any?
          word = data['words'].first
          expect_group_word_structure(word)
        end
      end
    end
  end
end
