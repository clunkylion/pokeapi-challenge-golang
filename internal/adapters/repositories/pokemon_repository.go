package repositories

import (
	"errors"
	"pokemon-api/internal/core/domain"
	"pokemon-api/internal/core/ports"

	"gorm.io/gorm"
)

type PokemonRepository struct {
	db *gorm.DB
}

func NewPokemonRepository(db *gorm.DB) ports.PokemonRepository {
	return &PokemonRepository{db: db}
}

func (r *PokemonRepository) Create(pokemon *domain.Pokemon) error {
	return r.db.Create(pokemon).Error
}

func (r *PokemonRepository) GetByID(id uint) (*domain.Pokemon, error) {
	var pokemon domain.Pokemon
	err := r.db.First(&pokemon, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pokemon not found")
		}
		return nil, err
	}
	return &pokemon, nil
}

func (r *PokemonRepository) GetByName(name string) (*domain.Pokemon, error) {
	var pokemon domain.Pokemon
	err := r.db.Where("name = ?", name).First(&pokemon).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pokemon not found")
		}
		return nil, err
	}
	return &pokemon, nil
}

func (r *PokemonRepository) List() ([]*domain.Pokemon, error) {
	var pokemon []*domain.Pokemon
	err := r.db.Find(&pokemon).Error
	if err != nil {
		return nil, err
	}
	return pokemon, nil
}

func (r *PokemonRepository) Migrate() error {
	if err := r.db.AutoMigrate(&domain.Pokemon{}); err != nil {
		return r.db.Exec(`
			CREATE TABLE IF NOT EXISTS pokemons (
				id SERIAL PRIMARY KEY,
				name VARCHAR(255) UNIQUE NOT NULL,
				type1 VARCHAR(255),
				type2 VARCHAR(255),
				height INTEGER DEFAULT 0,
				weight INTEGER DEFAULT 0,
				base_exp INTEGER DEFAULT 0,
				created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
				updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
			)
		`).Error
	}
	return nil
}
