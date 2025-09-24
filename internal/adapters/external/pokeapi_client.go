package external

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pokemon-api/internal/core/domain"
	"pokemon-api/internal/core/ports"
	"strings"
	"time"
)

type pokeAPIClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewPokeAPIClient(baseURL string) ports.PokemonAPIClient {
	return &pokeAPIClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *pokeAPIClient) GetPokemonData(identifier string) (*domain.ExternalPokemonResponse, error) {
	identifier = strings.ToLower(strings.TrimSpace(identifier))
	url := fmt.Sprintf("%s/pokemon/%s", c.baseURL, identifier)

	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make request to PokeAPI: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("pokemon '%s' not found", identifier)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("PokeAPI returned status %d", resp.StatusCode)
	}

	var pokemonData domain.ExternalPokemonResponse
	if err := json.NewDecoder(resp.Body).Decode(&pokemonData); err != nil {
		return nil, fmt.Errorf("failed to decode PokeAPI response: %w", err)
	}

	return &pokemonData, nil
}
