package repositories

import (
	"pokemon-api/internal/core/domain"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&domain.Pokemon{})
	assert.NoError(t, err)

	return db
}

func TestPokemonRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPokemonRepository(db)

	pokemon := &domain.Pokemon{
		Name:    "pikachu",
		Type1:   "electric",
		Type2:   "",
		Height:  4,
		Weight:  60,
		BaseExp: 112,
	}

	err := repo.Create(pokemon)
	assert.NoError(t, err)
	assert.NotZero(t, pokemon.ID)

	var count int64
	db.Model(&domain.Pokemon{}).Count(&count)
	assert.Equal(t, int64(1), count)
}

func TestPokemonRepository_Create_Duplicate(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPokemonRepository(db)

	pokemon1 := &domain.Pokemon{
		Name:  "pikachu",
		Type1: "electric",
	}
	pokemon2 := &domain.Pokemon{
		Name:  "pikachu",
		Type1: "electric",
	}

	err := repo.Create(pokemon1)
	assert.NoError(t, err)

	err = repo.Create(pokemon2)
	assert.Error(t, err)
}

func TestPokemonRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPokemonRepository(db)

	original := &domain.Pokemon{
		Name:    "charizard",
		Type1:   "fire",
		Type2:   "flying",
		Height:  17,
		Weight:  905,
		BaseExp: 267,
	}

	err := repo.Create(original)
	assert.NoError(t, err)

	found, err := repo.GetByID(original.ID)
	assert.NoError(t, err)
	assert.Equal(t, original.Name, found.Name)
	assert.Equal(t, original.Type1, found.Type1)
	assert.Equal(t, original.Type2, found.Type2)
	assert.Equal(t, original.Height, found.Height)
	assert.Equal(t, original.Weight, found.Weight)
	assert.Equal(t, original.BaseExp, found.BaseExp)
}

func TestPokemonRepository_GetByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPokemonRepository(db)

	found, err := repo.GetByID(999)
	assert.Error(t, err)
	assert.Equal(t, "pokemon not found", err.Error())
	assert.Nil(t, found)
}

func TestPokemonRepository_GetByName(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPokemonRepository(db)

	original := &domain.Pokemon{
		Name:    "squirtle",
		Type1:   "water",
		Height:  5,
		Weight:  90,
		BaseExp: 63,
	}

	err := repo.Create(original)
	assert.NoError(t, err)

	found, err := repo.GetByName("squirtle")
	assert.NoError(t, err)
	assert.Equal(t, original.Name, found.Name)
	assert.Equal(t, original.Type1, found.Type1)
	assert.Equal(t, original.Height, found.Height)
	assert.Equal(t, original.Weight, found.Weight)
	assert.Equal(t, original.BaseExp, found.BaseExp)
}

func TestPokemonRepository_GetByName_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPokemonRepository(db)

	found, err := repo.GetByName("nonexistent")
	assert.Error(t, err)
	assert.Equal(t, "pokemon not found", err.Error())
	assert.Nil(t, found)
}

func TestPokemonRepository_List(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPokemonRepository(db)

	pokemon1 := &domain.Pokemon{Name: "pikachu", Type1: "electric"}
	pokemon2 := &domain.Pokemon{Name: "charizard", Type1: "fire"}
	pokemon3 := &domain.Pokemon{Name: "squirtle", Type1: "water"}

	err := repo.Create(pokemon1)
	assert.NoError(t, err)
	err = repo.Create(pokemon2)
	assert.NoError(t, err)
	err = repo.Create(pokemon3)
	assert.NoError(t, err)

	list, err := repo.List()
	assert.NoError(t, err)
	assert.Len(t, list, 3)

	names := make([]string, len(list))
	for i, p := range list {
		names[i] = p.Name
	}
	assert.Contains(t, names, "pikachu")
	assert.Contains(t, names, "charizard")
	assert.Contains(t, names, "squirtle")
}

func TestPokemonRepository_List_Empty(t *testing.T) {
	db := setupTestDB(t)
	repo := NewPokemonRepository(db)

	list, err := repo.List()
	assert.NoError(t, err)
	assert.Empty(t, list)
}

func TestPokemonRepository_Migrate(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	repo := &PokemonRepository{db: db}
	err = repo.Migrate()
	assert.NoError(t, err)

	assert.True(t, db.Migrator().HasTable(&domain.Pokemon{}))
}
