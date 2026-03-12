package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

type GeocodeHandler struct {
}

type nominatimResponse struct {
	PlaceID     int64    `json:"place_id"`
	Description string   `json:"display_name"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	Address     struct {
		City        string `json:"city"`
		Town        string `json:"town"`
		Village     string `json:"village"`
		Country     string `json:"country"`
		CountryCode string `json:"country_code"`
	} `json:"address"`
}

type geocodeResult struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
}

var client = &http.Client{Timeout: 10 * time.Second}

func (h *GeocodeHandler) Search(c *fiber.Ctx) error {
	q := c.Query("q")
	if len(q) < 2 {
		return c.JSON(fiber.Map{"items": []geocodeResult{}})
	}

	searchURL := fmt.Sprintf(
		"https://nominatim.openstreetmap.org/search?format=json&q=%s&addressdetails=1&featuretype=city&limit=5",
		url.QueryEscape(q),
	)

	req, err := http.NewRequestWithContext(c.Context(), "GET", searchURL, nil)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create request")
	}
	
	// Nominatim requires a User-Agent
	req.Header.Set("User-Agent", "TripListik-App/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return fiber.NewError(fiber.StatusBadGateway, "Failed to connect to geocoding service")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fiber.NewError(fiber.StatusBadGateway, "Geocoding service returned an error")
	}

	var nomRes []nominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&nomRes); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to parse geocoding response")
	}

	var results []geocodeResult
	for _, item := range nomRes {
		city := item.Address.City
		if city == "" {
			city = item.Address.Town
		}
		if city == "" {
			city = item.Address.Village
		}
		
		results = append(results, geocodeResult{
			ID:          fmt.Sprintf("%d", item.PlaceID),
			Description: item.Description,
			City:        city,
			Country:     item.Address.Country,
			CountryCode: item.Address.CountryCode,
			Lat:         item.Lat,
			Lon:         item.Lon,
		})
	}

	return c.JSON(fiber.Map{"items": results})
}
