package ports

import "pokemon-api/internal/core/domain"

// PokemonRepository defines the interface for Pokemon data persistence
type PokemonRepository interface {
	Create(pokemon *domain.Pokemon) error
	GetByID(id uint) (*domain.Pokemon, error)
	GetByName(name string) (*domain.Pokemon, error)
	List() ([]*domain.Pokemon, error)
}

// PokemonAPIClient defines the interface for external PokeAPI integration
type PokemonAPIClient interface {
	GetPokemonData(identifier string) (*domain.ExternalPokemonResponse, error)
}

// PokemonService defines the interface for Pokemon business logic
type PokemonService interface {
	CreatePokemon(req *domain.CreatePokemonRequest) (*domain.Pokemon, error)
	CreatePokemonFlexible(req *domain.FlexiblePokemonRequest) (*domain.Pokemon, error)
	GetPokemon(id uint) (*domain.Pokemon, error)
	ListPokemon() ([]*domain.Pokemon, error)
}
