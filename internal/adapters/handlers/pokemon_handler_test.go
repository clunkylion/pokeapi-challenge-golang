package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"pokemon-api/internal/core/domain"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPokemonService struct {
	mock.Mock
}

func (m *MockPokemonService) CreatePokemon(req *domain.CreatePokemonRequest) (*domain.Pokemon, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pokemon), args.Error(1)
}

func (m *MockPokemonService) CreatePokemonFlexible(req *domain.FlexiblePokemonRequest) (*domain.Pokemon, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pokemon), args.Error(1)
}

func (m *MockPokemonService) GetPokemon(id uint) (*domain.Pokemon, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Pokemon), args.Error(1)
}

func (m *MockPokemonService) ListPokemon() ([]*domain.Pokemon, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Pokemon), args.Error(1)
}

func setupRouter(service *MockPokemonService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	handler := NewPokemonHandler(service)

	router.GET("/health", handler.HealthCheck)
	api := router.Group("/api/v1")
	{
		pokemon := api.Group("/pokemon")
		{
			pokemon.POST("", handler.CreatePokemonFlexible)
			pokemon.GET("/:id", handler.GetPokemon)
			pokemon.GET("", handler.ListPokemon)
		}
	}

	return router
}

func TestPokemonHandler_CreatePokemonFlexible(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		setupMock      func(*MockPokemonService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name: "successful creation with direct name",
			requestBody: map[string]interface{}{
				"name":  "pikachu",
				"type1": "electric",
			},
			setupMock: func(service *MockPokemonService) {
				service.On("CreatePokemonFlexible", mock.AnythingOfType("*domain.FlexiblePokemonRequest")).Return(&domain.Pokemon{
					ID:      1,
					Name:    "pikachu",
					Type1:   "electric",
					Height:  4,
					Weight:  60,
					BaseExp: 112,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":              float64(1),
				"name":            "pikachu",
				"type1":           "electric",
				"height":          float64(4),
				"weight":          float64(60),
				"base_experience": float64(112),
			},
		},
		{
			name: "successful creation with nested name",
			requestBody: map[string]interface{}{
				"pokemon": map[string]interface{}{
					"name": "charizard",
				},
				"type1": "fire",
			},
			setupMock: func(service *MockPokemonService) {
				service.On("CreatePokemonFlexible", mock.AnythingOfType("*domain.FlexiblePokemonRequest")).Return(&domain.Pokemon{
					ID:      2,
					Name:    "charizard",
					Type1:   "fire",
					Height:  17,
					Weight:  905,
					BaseExp: 267,
				}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody: map[string]interface{}{
				"id":              float64(2),
				"name":            "charizard",
				"type1":           "fire",
				"height":          float64(17),
				"weight":          float64(905),
				"base_experience": float64(267),
			},
		},
		{
			name: "duplicate pokemon",
			requestBody: map[string]interface{}{
				"name":  "pikachu",
				"type1": "electric",
			},
			setupMock: func(service *MockPokemonService) {
				service.On("CreatePokemonFlexible", mock.AnythingOfType("*domain.FlexiblePokemonRequest")).Return(nil, errors.New("pokemon with this name already exists"))
			},
			expectedStatus: http.StatusConflict,
			expectedBody: map[string]interface{}{
				"error": "pokemon with this name already exists",
			},
		},
		{
			name: "external API error",
			requestBody: map[string]interface{}{
				"name":  "invalid-pokemon",
				"type1": "fire",
			},
			setupMock: func(service *MockPokemonService) {
				service.On("CreatePokemonFlexible", mock.AnythingOfType("*domain.FlexiblePokemonRequest")).Return(nil, errors.New("failed to fetch Pokemon data"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "failed to fetch Pokemon data",
			},
		},
		{
			name: "invalid request body",
			requestBody: map[string]interface{}{
				"name": "pikachu",
			},
			setupMock: func(service *MockPokemonService) {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockPokemonService)
			tt.setupMock(mockService)
			router := setupRouter(mockService)

			body, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/api/v1/pokemon", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				for key, expectedValue := range tt.expectedBody {
					assert.Equal(t, expectedValue, response[key])
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPokemonHandler_GetPokemon(t *testing.T) {
	tests := []struct {
		name           string
		pokemonID      string
		setupMock      func(*MockPokemonService)
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:      "successful get",
			pokemonID: "1",
			setupMock: func(service *MockPokemonService) {
				service.On("GetPokemon", uint(1)).Return(&domain.Pokemon{
					ID:      1,
					Name:    "pikachu",
					Type1:   "electric",
					Height:  4,
					Weight:  60,
					BaseExp: 112,
				}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: map[string]interface{}{
				"id":              float64(1),
				"name":            "pikachu",
				"type1":           "electric",
				"height":          float64(4),
				"weight":          float64(60),
				"base_experience": float64(112),
			},
		},
		{
			name:      "pokemon not found",
			pokemonID: "999",
			setupMock: func(service *MockPokemonService) {
				service.On("GetPokemon", uint(999)).Return(nil, errors.New("pokemon not found"))
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "pokemon not found",
			},
		},
		{
			name:      "invalid pokemon ID",
			pokemonID: "invalid",
			setupMock: func(service *MockPokemonService) {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid Pokemon ID",
			},
		},
		{
			name:      "database error",
			pokemonID: "1",
			setupMock: func(service *MockPokemonService) {
				service.On("GetPokemon", uint(1)).Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "database connection failed",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockPokemonService)
			tt.setupMock(mockService)
			router := setupRouter(mockService)

			req, _ := http.NewRequest("GET", "/api/v1/pokemon/"+tt.pokemonID, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedBody != nil {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)

				for key, expectedValue := range tt.expectedBody {
					assert.Equal(t, expectedValue, response[key])
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPokemonHandler_ListPokemon(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockPokemonService)
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "successful list",
			setupMock: func(service *MockPokemonService) {
				pokemon := []*domain.Pokemon{
					{ID: 1, Name: "pikachu", Type1: "electric"},
					{ID: 2, Name: "charizard", Type1: "fire"},
				}
				service.On("ListPokemon").Return(pokemon, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name: "empty list",
			setupMock: func(service *MockPokemonService) {
				service.On("ListPokemon").Return([]*domain.Pokemon{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name: "database error",
			setupMock: func(service *MockPokemonService) {
				service.On("ListPokemon").Return(nil, errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockPokemonService)
			tt.setupMock(mockService)
			router := setupRouter(mockService)

			req, _ := http.NewRequest("GET", "/api/v1/pokemon", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedCount >= 0 {
				var response []map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response, tt.expectedCount)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestPokemonHandler_HealthCheck(t *testing.T) {
	mockService := new(MockPokemonService)
	router := setupRouter(mockService)

	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "healthy", response["status"])
	assert.Equal(t, "pokemon-api", response["service"])
}

func TestParseUint(t *testing.T) {
	tests := []struct {
		input    string
		expected uint
		hasError bool
	}{
		{"1", 1, false},
		{"999", 999, false},
		{"0", 0, false},
		{"invalid", 0, true},
		{"-1", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run("parse_"+tt.input, func(t *testing.T) {
			result, err := strconv.ParseUint(tt.input, 10, 32)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, uint(result))
			}
		})
	}
}
