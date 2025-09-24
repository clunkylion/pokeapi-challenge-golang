package services

import (
	"errors"
	"fmt"
	"pokemon-api/internal/core/domain"
	"pokemon-api/internal/core/ports"
	"strings"
)

type pokemonService struct {
	repository ports.PokemonRepository
	apiClient  ports.PokemonAPIClient
}

func NewPokemonService(repository ports.PokemonRepository, apiClient ports.PokemonAPIClient) ports.PokemonService {
	return &pokemonService{
		repository: repository,
		apiClient:  apiClient,
	}
}

func (s *pokemonService) CreatePokemon(req *domain.CreatePokemonRequest) (*domain.Pokemon, error) {
	existingPokemon, err := s.repository.GetByName(req.Name)
	if err == nil && existingPokemon != nil {
		return nil, errors.New("pokemon with this name already exists")
	}

	externalData, err := s.apiClient.GetPokemonData(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Pokemon data: %w", err)
	}

	pokemon := &domain.Pokemon{
		Name:    externalData.Name,
		Type1:   req.Type1,
		Type2:   req.Type2,
		Height:  externalData.Height,
		Weight:  externalData.Weight,
		BaseExp: externalData.BaseExperience,
	}

	if err := s.repository.Create(pokemon); err != nil {
		return nil, fmt.Errorf("failed to save Pokemon: %w", err)
	}

	return pokemon, nil
}

func (s *pokemonService) CreatePokemonFlexible(req *domain.FlexiblePokemonRequest) (*domain.Pokemon, error) {
	pokemonName := s.extractPokemonName(req)
	if pokemonName == "" {
		return nil, errors.New("pokemon name is required")
	}

	standardReq := &domain.CreatePokemonRequest{
		Name:  pokemonName,
		Type1: req.Type1,
		Type2: req.Type2,
	}

	return s.CreatePokemon(standardReq)
}

func (s *pokemonService) GetPokemon(id uint) (*domain.Pokemon, error) {
	return s.repository.GetByID(id)
}

func (s *pokemonService) ListPokemon() ([]*domain.Pokemon, error) {
	return s.repository.List()
}

func (s *pokemonService) extractPokemonName(input *domain.FlexiblePokemonRequest) string {
	if input.Name != "" {
		return strings.ToLower(strings.TrimSpace(input.Name))
	}

	if input.Pokemon != nil {
		if nameValue, exists := input.Pokemon["name"]; exists {
			if nameStr, ok := nameValue.(string); ok {
				return strings.ToLower(strings.TrimSpace(nameStr))
			}
		}

		if nameStr, ok := input.Pokemon["pokemon"].(string); ok {
			return strings.ToLower(strings.TrimSpace(nameStr))
		}
	}

	return ""
}
