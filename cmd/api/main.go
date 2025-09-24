package main

import (
	"log"
	"os"
	"pokemon-api/internal/adapters/external"
	"pokemon-api/internal/adapters/handlers"
	"pokemon-api/internal/adapters/repositories"
	"pokemon-api/internal/core/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "pokemon-api/docs"
)

// @title Pokemon API
// @version 1.0
// @description A REST API for managing Pokemon data with PokeAPI integration
// @host localhost:8080
// @BasePath /
func main() {
	dbHost := getEnv("DB_HOST", "localhost")
	dbUser := getEnv("DB_USER", "pokemon_user")
	dbPassword := getEnv("DB_PASSWORD", "pokemon_pass")
	dbName := getEnv("DB_NAME", "pokemon_db")
	dbPort := getEnv("DB_PORT", "5432")
	pokeAPIBaseURL := getEnv("POKEAPI_BASE_URL", "https://pokeapi.co/api/v2")

	dsn := "host=" + dbHost + " user=" + dbUser + " password=" + dbPassword + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	repo := repositories.NewPokemonRepository(db)
	pokemonRepo := repo.(*repositories.PokemonRepository)
	if err := pokemonRepo.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	apiClient := external.NewPokeAPIClient(pokeAPIBaseURL)
	service := services.NewPokemonService(repo, apiClient)
	handler := handlers.NewPokemonHandler(service)

	router := gin.Default()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := getEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	log.Fatal(router.Run(":" + port))
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
