package services

import (
	"errors"
	"pokemon-api/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPokemonRepository struct {
	mock.Mock
}

func (m *MockPokemonRepository) Create(pokemon *domain.Pokemon) error {
	args := m.Called(pokemon)
	return args.Error(0)
}

func (m *MockPokemonRepository) GetByID(id uint) (*domain.Pokemon, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.Pokemon), args.Error(1)
}

func (m *MockPokemonRepository) GetByName(name string) (*domain.Pokemon, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pokemon), args.Error(1)
}

func (m *MockPokemonRepository) List() ([]*domain.Pokemon, error) {
	args := m.Called()
	return args.Get(0).([]*domain.Pokemon), args.Error(1)
}

type MockPokemonAPIClient struct {
	mock.Mock
}

func (m *MockPokemonAPIClient) GetPokemonData(identifier string) (*domain.ExternalPokemonResponse, error) {
	args := m.Called(identifier)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ExternalPokemonResponse), args.Error(1)
}

func TestPokemonService_CreatePokemon(t *testing.T) {
	tests := []struct {
		name           string
		request        *domain.CreatePokemonRequest
		setupMocks     func(*MockPokemonRepository, *MockPokemonAPIClient)
		expectedError  string
		expectedResult *domain.Pokemon
	}{
		{
			name: "successful creation",
			request: &domain.CreatePokemonRequest{
				Name:  "pikachu",
				Type1: "electric",
				Type2: "",
			},
			setupMocks: func(repo *MockPokemonRepository, client *MockPokemonAPIClient) {
				repo.On("GetByName", "pikachu").Return(nil, errors.New("not found"))
				client.On("GetPokemonData", "pikachu").Return(&domain.ExternalPokemonResponse{
					ID:             25,
					Name:           "pikachu",
					Height:         4,
					Weight:         60,
					BaseExperience: 112,
				}, nil)
				repo.On("Create", mock.AnythingOfType("*domain.Pokemon")).Return(nil)
			},
			expectedResult: &domain.Pokemon{
				Name:    "pikachu",
				Type1:   "electric",
				Type2:   "",
				Height:  4,
				Weight:  60,
				BaseExp: 112,
			},
		},
		{
			name: "duplicate pokemon name",
			request: &domain.CreatePokemonRequest{
				Name:  "pikachu",
				Type1: "electric",
			},
			setupMocks: func(repo *MockPokemonRepository, client *MockPokemonAPIClient) {
				repo.On("GetByName", "pikachu").Return(&domain.Pokemon{
					ID:   1,
					Name: "pikachu",
				}, nil)
			},
			expectedError: "pokemon with this name already exists",
		},
		{
			name: "external API error",
			request: &domain.CreatePokemonRequest{
				Name:  "invalid-pokemon",
				Type1: "fire",
			},
			setupMocks: func(repo *MockPokemonRepository, client *MockPokemonAPIClient) {
				repo.On("GetByName", "invalid-pokemon").Return(nil, errors.New("not found"))
				client.On("GetPokemonData", "invalid-pokemon").Return(nil, errors.New("pokemon not found"))
			},
			expectedError: "failed to fetch Pokemon data: pokemon not found",
		},
		{
			name: "repository save error",
			request: &domain.CreatePokemonRequest{
				Name:  "charizard",
				Type1: "fire",
			},
			setupMocks: func(repo *MockPokemonRepository, client *MockPokemonAPIClient) {
				repo.On("GetByName", "charizard").Return(nil, errors.New("not found"))
				client.On("GetPokemonData", "charizard").Return(&domain.ExternalPokemonResponse{
					ID:             6,
					Name:           "charizard",
					Height:         17,
					Weight:         905,
					BaseExperience: 267,
				}, nil)
				repo.On("Create", mock.AnythingOfType("*domain.Pokemon")).Return(errors.New("database error"))
			},
			expectedError: "failed to save Pokemon: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPokemonRepository)
			mockClient := new(MockPokemonAPIClient)
			tt.setupMocks(mockRepo, mockClient)

			service := NewPokemonService(mockRepo, mockClient)
			result, err := service.CreatePokemon(tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Type1, result.Type1)
				assert.Equal(t, tt.expectedResult.Type2, result.Type2)
				assert.Equal(t, tt.expectedResult.Height, result.Height)
				assert.Equal(t, tt.expectedResult.Weight, result.Weight)
				assert.Equal(t, tt.expectedResult.BaseExp, result.BaseExp)
			}

			mockRepo.AssertExpectations(t)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestPokemonService_CreatePokemonFlexible(t *testing.T) {
	tests := []struct {
		name           string
		request        *domain.FlexiblePokemonRequest
		setupMocks     func(*MockPokemonRepository, *MockPokemonAPIClient)
		expectedError  string
		expectedResult *domain.Pokemon
	}{
		{
			name: "direct name format",
			request: &domain.FlexiblePokemonRequest{
				Name:  "pikachu",
				Type1: "electric",
			},
			setupMocks: func(repo *MockPokemonRepository, client *MockPokemonAPIClient) {
				repo.On("GetByName", "pikachu").Return(nil, errors.New("not found"))
				client.On("GetPokemonData", "pikachu").Return(&domain.ExternalPokemonResponse{
					Name:           "pikachu",
					Height:         4,
					Weight:         60,
					BaseExperience: 112,
				}, nil)
				repo.On("Create", mock.AnythingOfType("*domain.Pokemon")).Return(nil)
			},
			expectedResult: &domain.Pokemon{
				Name:    "pikachu",
				Type1:   "electric",
				Height:  4,
				Weight:  60,
				BaseExp: 112,
			},
		},
		{
			name: "nested pokemon.name format",
			request: &domain.FlexiblePokemonRequest{
				Type1: "fire",
				Pokemon: map[string]interface{}{
					"name": "charizard",
				},
			},
			setupMocks: func(repo *MockPokemonRepository, client *MockPokemonAPIClient) {
				repo.On("GetByName", "charizard").Return(nil, errors.New("not found"))
				client.On("GetPokemonData", "charizard").Return(&domain.ExternalPokemonResponse{
					Name:           "charizard",
					Height:         17,
					Weight:         905,
					BaseExperience: 267,
				}, nil)
				repo.On("Create", mock.AnythingOfType("*domain.Pokemon")).Return(nil)
			},
			expectedResult: &domain.Pokemon{
				Name:    "charizard",
				Type1:   "fire",
				Height:  17,
				Weight:  905,
				BaseExp: 267,
			},
		},
		{
			name: "pokemon as string format",
			request: &domain.FlexiblePokemonRequest{
				Type1: "water",
				Pokemon: map[string]interface{}{
					"pokemon": "squirtle",
				},
			},
			setupMocks: func(repo *MockPokemonRepository, client *MockPokemonAPIClient) {
				repo.On("GetByName", "squirtle").Return(nil, errors.New("not found"))
				client.On("GetPokemonData", "squirtle").Return(&domain.ExternalPokemonResponse{
					Name:           "squirtle",
					Height:         5,
					Weight:         90,
					BaseExperience: 63,
				}, nil)
				repo.On("Create", mock.AnythingOfType("*domain.Pokemon")).Return(nil)
			},
			expectedResult: &domain.Pokemon{
				Name:    "squirtle",
				Type1:   "water",
				Height:  5,
				Weight:  90,
				BaseExp: 63,
			},
		},
		{
			name: "missing name",
			request: &domain.FlexiblePokemonRequest{
				Type1: "grass",
			},
			setupMocks: func(repo *MockPokemonRepository, client *MockPokemonAPIClient) {
			},
			expectedError: "pokemon name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPokemonRepository)
			mockClient := new(MockPokemonAPIClient)
			tt.setupMocks(mockRepo, mockClient)

			service := NewPokemonService(mockRepo, mockClient)
			result, err := service.CreatePokemonFlexible(tt.request)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Type1, result.Type1)
				assert.Equal(t, tt.expectedResult.Height, result.Height)
				assert.Equal(t, tt.expectedResult.Weight, result.Weight)
				assert.Equal(t, tt.expectedResult.BaseExp, result.BaseExp)
			}

			mockRepo.AssertExpectations(t)
			mockClient.AssertExpectations(t)
		})
	}
}

func TestPokemonService_GetPokemon(t *testing.T) {
	tests := []struct {
		name           string
		pokemonID      uint
		setupMocks     func(*MockPokemonRepository)
		expectedError  string
		expectedResult *domain.Pokemon
	}{
		{
			name:      "successful get",
			pokemonID: 1,
			setupMocks: func(repo *MockPokemonRepository) {
				repo.On("GetByID", uint(1)).Return(&domain.Pokemon{
					ID:      1,
					Name:    "pikachu",
					Type1:   "electric",
					Height:  4,
					Weight:  60,
					BaseExp: 112,
				}, nil)
			},
			expectedResult: &domain.Pokemon{
				ID:      1,
				Name:    "pikachu",
				Type1:   "electric",
				Height:  4,
				Weight:  60,
				BaseExp: 112,
			},
		},
		{
			name:      "pokemon not found",
			pokemonID: 999,
			setupMocks: func(repo *MockPokemonRepository) {
				repo.On("GetByID", uint(999)).Return((*domain.Pokemon)(nil), errors.New("pokemon not found"))
			},
			expectedError: "pokemon not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPokemonRepository)
			mockClient := new(MockPokemonAPIClient)
			tt.setupMocks(mockRepo)

			service := NewPokemonService(mockRepo, mockClient)
			result, err := service.GetPokemon(tt.pokemonID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Type1, result.Type1)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPokemonService_ListPokemon(t *testing.T) {
	tests := []struct {
		name          string
		setupMocks    func(*MockPokemonRepository)
		expectedError string
		expectedCount int
	}{
		{
			name: "successful list",
			setupMocks: func(repo *MockPokemonRepository) {
				pokemon := []*domain.Pokemon{
					{ID: 1, Name: "pikachu", Type1: "electric"},
					{ID: 2, Name: "charizard", Type1: "fire"},
				}
				repo.On("List").Return(pokemon, nil)
			},
			expectedCount: 2,
		},
		{
			name: "empty list",
			setupMocks: func(repo *MockPokemonRepository) {
				repo.On("List").Return([]*domain.Pokemon{}, nil)
			},
			expectedCount: 0,
		},
		{
			name: "repository error",
			setupMocks: func(repo *MockPokemonRepository) {
				repo.On("List").Return(([]*domain.Pokemon)(nil), errors.New("database error"))
			},
			expectedError: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPokemonRepository)
			mockClient := new(MockPokemonAPIClient)
			tt.setupMocks(mockRepo)

			service := NewPokemonService(mockRepo, mockClient)
			result, err := service.ListPokemon()

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Len(t, result, tt.expectedCount)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestPokemonService_ExtractPokemonName(t *testing.T) {
	service := &pokemonService{}

	tests := []struct {
		name     string
		input    *domain.FlexiblePokemonRequest
		expected string
	}{
		{
			name: "direct name",
			input: &domain.FlexiblePokemonRequest{
				Name: "  PIKACHU  ",
			},
			expected: "pikachu",
		},
		{
			name: "nested pokemon.name",
			input: &domain.FlexiblePokemonRequest{
				Pokemon: map[string]interface{}{
					"name": "  CHARIZARD  ",
				},
			},
			expected: "charizard",
		},
		{
			name: "pokemon as string",
			input: &domain.FlexiblePokemonRequest{
				Pokemon: map[string]interface{}{
					"pokemon": "  SQUIRTLE  ",
				},
			},
			expected: "squirtle",
		},
		{
			name:     "empty input",
			input:    &domain.FlexiblePokemonRequest{},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.extractPokemonName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
