require_relative '../support/spec_helper'

RSpec.describe 'Words API' do
  before(:all) do
    @test_word_id = get_test_word_id
    puts "Using test_word_id: #{@test_word_id}"
  end

  describe 'GET /api/words' do
    it 'returns paginated words list' do
      response = get_words(1)
      expect_success(response)
      
      data = response.parsed_response
      expect(data).to have_key('items')
      expect(data).to have_key('pagination')
      
      expect(data['items']).to be_an(Array)
      expect_pagination_structure(data['pagination'])
      
      # Check word structure if items exist
      if data['items'].any?
        word = data['items'].first
        expect_word_structure(word)
      end
    end

    it 'supports pagination parameters' do
      response = get_words(2)
      expect_success(response)
      
      pagination = response.parsed_response['pagination']
      expect(pagination['page']).to eq(2)
    end
  end

  describe 'GET /api/words/:id' do
    context 'when word exists' do
      it 'returns word details with stats and groups' do
        response = get_word(@test_word_id)
        expect_success(response)
        
        data = response.parsed_response
        expect_word_detail_structure(data)
      end
    end

    context 'when word does not exist' do
      let(:word_id) { 99999 }

      it 'returns 404 error' do
        response = get_word(word_id)
        expect_error(response, 404)
      end
    end
  end
end
