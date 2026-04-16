package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetLocationAreas(url string, c *Cache) (GetLocationAreasResponse, error) {
	var response GetLocationAreasResponse

	if entry, ok := c.Get(url); ok {
		if err := json.Unmarshal(entry.val, &response); err != nil {
			return GetLocationAreasResponse{}, fmt.Errorf("Could not read cache entry.\n Error: %v", err)
		}
		return response, nil
	}

	res, err := http.Get(url)
	if err != nil {
		return GetLocationAreasResponse{}, fmt.Errorf("Could not get map locations.\n Error: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return GetLocationAreasResponse{}, fmt.Errorf("Could not read map body.\n Error: %v", err)
	}
	c.Set(url, cacheEntry{createdAt: time.Now(), val: body})

	if err := json.Unmarshal(body, &response); err != nil {
		return GetLocationAreasResponse{}, fmt.Errorf("Could not read map locations.\n Error: %v", err)
	}

	return response, nil
}

type GetLocationAreasResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

type LocationArea struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int                   `json:"game_index"`
	Location             NamedAPIResource      `json:"location"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

type NamedAPIResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EncounterMethodRate struct {
	EncounterMethod NamedAPIResource       `json:"encounter_method"`
	VersionDetails  []VersionEncounterRate `json:"version_details"`
}

type VersionEncounterRate struct {
	Rate    int              `json:"rate"`
	Version NamedAPIResource `json:"version"`
}

type Name struct {
	Name     string           `json:"name"`
	Language NamedAPIResource `json:"language"`
}

type PokemonEncounter struct {
	Pokemon        NamedAPIResource       `json:"pokemon"`
	VersionDetails []PokemonVersionDetail `json:"version_details"`
}

type PokemonVersionDetail struct {
	MaxChance        int               `json:"max_chance"`
	Version          NamedAPIResource  `json:"version"`
	EncounterDetails []EncounterDetail `json:"encounter_details"`
}

type EncounterDetail struct {
	MinLevel        int                `json:"min_level"`
	MaxLevel        int                `json:"max_level"`
	Chance          int                `json:"chance"`
	Method          NamedAPIResource   `json:"method"`
	ConditionValues []NamedAPIResource `json:"condition_values"`
}
