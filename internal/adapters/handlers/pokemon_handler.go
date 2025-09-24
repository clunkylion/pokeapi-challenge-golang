package handlers

import (
	"net/http"
	"pokemon-api/internal/core/domain"
	"pokemon-api/internal/core/ports"
	"strconv"

	"github.com/gin-gonic/gin"
)

type pokemonHandler struct {
	service ports.PokemonService
}

func NewPokemonHandler(service ports.PokemonService) *pokemonHandler {
	return &pokemonHandler{
		service: service,
	}
}


// @Summary Create a new Pokemon
// @Description Create a new Pokemon with data from PokeAPI
// @Tags pokemon
// @Accept json
// @Produce json
// @Param pokemon body domain.CreatePokemonRequest true "Pokemon data"
// @Success 201 {object} domain.Pokemon
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/pokemon [post]
func (h *pokemonHandler) CreatePokemon(c *gin.Context) {
	var req domain.CreatePokemonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pokemon, err := h.service.CreatePokemon(&req)
	if err != nil {
		if err.Error() == "pokemon with this name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pokemon)
}

// @Summary Create a new Pokemon with flexible name format
// @Description Create a new Pokemon supporting both direct name and nested pokemon.name formats
// @Tags pokemon
// @Accept json
// @Produce json
// @Param pokemon body domain.FlexiblePokemonRequest true "Pokemon data with flexible name"
// @Success 201 {object} domain.Pokemon
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/pokemon [post]
func (h *pokemonHandler) CreatePokemonFlexible(c *gin.Context) {
	var req domain.FlexiblePokemonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pokemon, err := h.service.CreatePokemonFlexible(&req)
	if err != nil {
		if err.Error() == "pokemon with this name already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pokemon)
}

// @Summary Get Pokemon by ID
// @Description Retrieve a Pokemon by its ID
// @Tags pokemon
// @Accept json
// @Produce json
// @Param id path int true "Pokemon ID"
// @Success 200 {object} domain.Pokemon
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/pokemon/{id} [get]
func (h *pokemonHandler) GetPokemon(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid Pokemon ID"})
		return
	}

	pokemon, err := h.service.GetPokemon(uint(id))
	if err != nil {
		if err.Error() == "pokemon not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pokemon)
}

// @Summary List all Pokemon
// @Description Retrieve all Pokemon from the database
// @Tags pokemon
// @Accept json
// @Produce json
// @Success 200 {array} domain.Pokemon
// @Failure 500 {object} map[string]string
// @Router /api/v1/pokemon [get]
func (h *pokemonHandler) ListPokemon(c *gin.Context) {
	pokemon, err := h.service.ListPokemon()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pokemon)
}

// @Summary Health check endpoint
// @Description Check if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *pokemonHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "pokemon-api"})
}
