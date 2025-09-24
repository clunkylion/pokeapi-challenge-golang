package domain

import "time"

type Pokemon struct {
	ID uint `json:"id" gorm:"primaryKey"`

	Name  string `json:"name" gorm:"unique;not null"`
	Type1 string `json:"type1" binding:"required"`
	Type2 string `json:"type2,omitempty"`

	Height  int `json:"height"`
	Weight  int `json:"weight"`
	BaseExp int `json:"base_experience"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreatePokemonRequest struct {
	Name  string `json:"name" binding:"required"`
	Type1 string `json:"type1" binding:"required"`
	Type2 string `json:"type2,omitempty"`
}

type FlexiblePokemonRequest struct {
	Name    string                 `json:"name,omitempty"`
	Type1   string                 `json:"type1" binding:"required"`
	Type2   string                 `json:"type2,omitempty"`
	Pokemon map[string]interface{} `json:"pokemon,omitempty"`
}

type ExternalPokemonResponse struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	BaseExperience int    `json:"base_experience"`
	Types          []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
}
