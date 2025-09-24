package external

import (
	"net/http"
	"net/http/httptest"
	"pokemon-api/internal/core/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestPokeAPIClient_GetPokemonData(t *testing.T) {
	tests := []struct {
		name           string
		identifier     string
		mockResponse   string
		mockStatusCode int
		expectedResult *domain.ExternalPokemonResponse
		expectedError  string
	}{
		{
			name:       "successful request",
			identifier: "pikachu",
			mockResponse: `{
				"id": 25,
				"name": "pikachu",
				"height": 4,
				"weight": 60,
				"base_experience": 112,
				"types": [
					{
						"type": {
							"name": "electric"
						}
					}
				]
			}`,
			mockStatusCode: http.StatusOK,
			expectedResult: &domain.ExternalPokemonResponse{
				ID:             25,
				Name:           "pikachu",
				Height:         4,
				Weight:         60,
				BaseExperience: 112,
				Types: []struct {
					Type struct {
						Name string `json:"name"`
					} `json:"type"`
				}{
					{
						Type: struct {
							Name string `json:"name"`
						}{
							Name: "electric",
						},
					},
				},
			},
		},
		{
			name:       "pokemon with multiple types",
			identifier: "charizard",
			mockResponse: `{
				"id": 6,
				"name": "charizard",
				"height": 17,
				"weight": 905,
				"base_experience": 267,
				"types": [
					{
						"type": {
							"name": "fire"
						}
					},
					{
						"type": {
							"name": "flying"
						}
					}
				]
			}`,
			mockStatusCode: http.StatusOK,
			expectedResult: &domain.ExternalPokemonResponse{
				ID:             6,
				Name:           "charizard",
				Height:         17,
				Weight:         905,
				BaseExperience: 267,
				Types: []struct {
					Type struct {
						Name string `json:"name"`
					} `json:"type"`
				}{
					{
						Type: struct {
							Name string `json:"name"`
						}{
							Name: "fire",
						},
					},
					{
						Type: struct {
							Name string `json:"name"`
						}{
							Name: "flying",
						},
					},
				},
			},
		},
		{
			name:           "pokemon not found",
			identifier:     "nonexistent",
			mockStatusCode: http.StatusNotFound,
			expectedError:  "pokemon 'nonexistent' not found",
		},
		{
			name:           "server error",
			identifier:     "pikachu",
			mockStatusCode: http.StatusInternalServerError,
			expectedError:  "PokeAPI returned status 500",
		},
		{
			name:           "invalid JSON response",
			identifier:     "pikachu",
			mockResponse:   `{"invalid": json}`,
			mockStatusCode: http.StatusOK,
			expectedError:  "failed to decode PokeAPI response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/pokemon/"+tt.identifier, r.URL.Path)
				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != "" {
					w.Write([]byte(tt.mockResponse))
				}
			}))
			defer server.Close()

			client := NewPokeAPIClient(server.URL)
			result, err := client.GetPokemonData(tt.identifier)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedResult.ID, result.ID)
				assert.Equal(t, tt.expectedResult.Name, result.Name)
				assert.Equal(t, tt.expectedResult.Height, result.Height)
				assert.Equal(t, tt.expectedResult.Weight, result.Weight)
				assert.Equal(t, tt.expectedResult.BaseExperience, result.BaseExperience)
				assert.Equal(t, len(tt.expectedResult.Types), len(result.Types))
			}
		})
	}
}

func TestPokeAPIClient_GetPokemonData_NameNormalization(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "uppercase name",
			input:    "PIKACHU",
			expected: "pikachu",
		},
		{
			name:     "mixed case name",
			input:    "ChArIzArD",
			expected: "charizard",
		},
		{
			name:     "name with spaces",
			input:    "  pikachu  ",
			expected: "pikachu",
		},
		{
			name:     "already lowercase",
			input:    "squirtle",
			expected: "squirtle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/pokemon/"+tt.expected, r.URL.Path)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"id": 1,
					"name": "` + tt.expected + `",
					"height": 1,
					"weight": 1,
					"base_experience": 1,
					"types": []
				}`))
			}))
			defer server.Close()

			client := NewPokeAPIClient(server.URL)
			result, err := client.GetPokemonData(tt.input)

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result.Name)
		})
	}
}

func TestPokeAPIClient_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(15 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewPokeAPIClient(server.URL)
	result, err := client.GetPokemonData("pikachu")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to make request to PokeAPI")
	assert.Nil(t, result)
}

func TestPokeAPIClient_Constructor(t *testing.T) {
	baseURL := "https://pokeapi.co/api/v2"
	client := NewPokeAPIClient(baseURL)

	pokeClient := client.(*pokeAPIClient)
	assert.Equal(t, baseURL, pokeClient.baseURL)
	assert.NotNil(t, pokeClient.httpClient)
	assert.Equal(t, 10*time.Second, pokeClient.httpClient.Timeout)
}

func TestPokeAPIClient_URLConstruction(t *testing.T) {
	identifier := "test-pokemon"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/pokemon/test-pokemon", r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewPokeAPIClient(server.URL + "/api")
	_, err := client.GetPokemonData(identifier)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pokemon 'test-pokemon' not found")
}
