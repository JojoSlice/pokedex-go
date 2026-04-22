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

func GetLocationArea(url string, locationAreaName string) (LocationArea, error) {

	fullUrl := url + locationAreaName
	res, err := http.Get(fullUrl)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Could not get location. \n Error: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Could not read location body.\n Error: %v", err)
	}

	locationArea := LocationArea{}
	err = json.Unmarshal(body, &locationArea)
	if err != nil {
		return LocationArea{}, fmt.Errorf("Could not decode location.\n Error: %v", err)
	}

	return locationArea, nil
}

type GetLocationAreasResponse struct {
	Count    int            `json:"count"`
	Next     *string        `json:"next"`
	Previous *string        `json:"previous"`
	Results  []LocationArea `json:"results"`
}

// --- Locations section types ---

type Location struct {
	ID          int                   `json:"id"`
	Name        string                `json:"name"`
	Region      Region                `json:"region"`
	Names       []Name                `json:"names"`
	GameIndices []GenerationGameIndex `json:"game_indices"`
	Areas       []LocationArea        `json:"areas"`
}

type LocationArea struct {
	ID                   int                   `json:"id"`
	Name                 string                `json:"name"`
	GameIndex            int                   `json:"game_index"`
	EncounterMethodRates []EncounterMethodRate `json:"encounter_method_rates"`
	Location             Location              `json:"location"`
	Names                []Name                `json:"names"`
	PokemonEncounters    []PokemonEncounter    `json:"pokemon_encounters"`
}

type EncounterMethodRate struct {
	EncounterMethod EncounterMethod           `json:"encounter_method"`
	VersionDetails  []EncounterVersionDetails `json:"version_details"`
}

type EncounterVersionDetails struct {
	Rate    int     `json:"rate"`
	Version Version `json:"version"`
}

type PokemonEncounter struct {
	Pokemon        Pokemon                  `json:"pokemon"`
	VersionDetails []VersionEncounterDetail `json:"version_details"`
}

type VersionEncounterDetail struct {
	MaxChance        int               `json:"max_chance"`
	Version          Version           `json:"version"`
	EncounterDetails []EncounterDetail `json:"encounter_details"`
}

type EncounterDetail struct {
	MinLevel        int                       `json:"min_level"`
	MaxLevel        int                       `json:"max_level"`
	Chance          int                       `json:"chance"`
	Method          EncounterMethod           `json:"method"`
	ConditionValues []EncounterConditionValue `json:"condition_values"`
}

type PalParkArea struct {
	ID                int                       `json:"id"`
	Name              string                    `json:"name"`
	Names             []Name                    `json:"names"`
	PokemonEncounters []PalParkEncounterSpecies `json:"pokemon_encounters"`
}

type PalParkEncounterSpecies struct {
	BaseScore      int            `json:"base_score"`
	Rate           int            `json:"rate"`
	PokemonSpecies PokemonSpecies `json:"pokemon_species"`
}

type Region struct {
	ID             int            `json:"id"`
	Locations      []Location     `json:"locations"`
	Name           string         `json:"name"`
	Names          []Name         `json:"names"`
	MainGeneration Generation     `json:"main_generation"`
	Pokedexes      []Pokedex      `json:"pokedexes"`
	VersionGroups  []VersionGroup `json:"version_groups"`
}

type Name struct {
	Name     string   `json:"name"`
	Language Language `json:"language"`
}

type GenerationGameIndex struct {
	GameIndex  int        `json:"game_index"`
	Generation Generation `json:"generation"`
}

// --- Referenced types ---

type EncounterMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type EncounterConditionValue struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Version struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Pokemon struct {
	ID                   int                `json:"id"`
	Name                 string             `json:"name"`
	BaseExperience       int                `json:"base_experience"`
	Height               int                `json:"height"`
	IsDefault            bool               `json:"is_default"`
	Order                int                `json:"order"`
	Weight               int                `json:"weight"`
	Abilities            []PokemonAbility   `json:"abilities"`
	Forms                []PokemonForm      `json:"forms"`
	GameIndices          []VersionGameIndex `json:"game_indices"`
	HeldItems            []PokemonHeldItem  `json:"held_items"`
	LocationAreaEncounters string           `json:"location_area_encounters"`
	Moves                []PokemonMove      `json:"moves"`
	Species              PokemonSpecies     `json:"species"`
	Sprites              PokemonSprites     `json:"sprites"`
	Stats                []PokemonStat      `json:"stats"`
	Types                []PokemonType      `json:"types"`
}

type PokemonAbility struct {
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
	Ability  Ability `json:"ability"`
}

type VersionGameIndex struct {
	GameIndex int     `json:"game_index"`
	Version   Version `json:"version"`
}

type PokemonHeldItem struct {
	Item           Item                       `json:"item"`
	VersionDetails []PokemonHeldItemVersion   `json:"version_details"`
}

type PokemonHeldItemVersion struct {
	Rarity  int     `json:"rarity"`
	Version Version `json:"version"`
}

type PokemonMove struct {
	Move                Move                 `json:"move"`
	VersionGroupDetails []PokemonMoveVersion `json:"version_group_details"`
}

type PokemonMoveVersion struct {
	LevelLearnedAt  int             `json:"level_learned_at"`
	VersionGroup    VersionGroup    `json:"version_group"`
	MoveLearnMethod MoveLearnMethod `json:"move_learn_method"`
	Order           int             `json:"order"`
}

type PokemonSprites struct {
	BackDefault      string `json:"back_default"`
	BackFemale       string `json:"back_female"`
	BackShiny        string `json:"back_shiny"`
	BackShinyFemale  string `json:"back_shiny_female"`
	FrontDefault     string `json:"front_default"`
	FrontFemale      string `json:"front_female"`
	FrontShiny       string `json:"front_shiny"`
	FrontShinyFemale string `json:"front_shiny_female"`
}

type PokemonStat struct {
	Stat     Stat `json:"stat"`
	Effort   int  `json:"effort"`
	BaseStat int  `json:"base_stat"`
}

type PokemonType struct {
	Slot int         `json:"slot"`
	Type PokemonTypeResource `json:"type"`
}

type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonForm struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Item struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Move struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type MoveLearnMethod struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonTypeResource struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Generation struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Pokedex struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type VersionGroup struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Language struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonSpecies struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
