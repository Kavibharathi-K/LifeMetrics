package usda

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func SearchFoods(query string) ([]FoodSuggestion, error) {

	apiKey := "cjgdjdoOOnoqf9oDNB6mfAag1SbIGAedyFWf3bGa"

	endpoint := fmt.Sprintf(
		"https://api.nal.usda.gov/fdc/v1/foods/search?query=%s&dataType=Foundation&dataType=SR%%20Legacy&pageSize=10&api_key=%s",
		url.QueryEscape(query),
		apiKey,
	)

	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var searchResponse SearchResponse

	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	var suggestions []FoodSuggestion

	for _, food := range searchResponse.Foods {

		if shouldSkipFood(food.Description) {
			continue
		}

		suggestions = append(suggestions, MapFood(food))
	}

	return suggestions, nil
}

func shouldSkipFood(name string) bool {

	name = strings.ToLower(name)

	blockedWords := []string{
		"lunchmeat",
		"deli",
		"breaded",
		"fat-free",
		"sliced",
		"oven-roasted",
		"microwaved",
		"mesquite",
		"honey glazed",
	}

	for _, word := range blockedWords {
		if strings.Contains(name, word) {
			return true
		}
	}

	return false
}
