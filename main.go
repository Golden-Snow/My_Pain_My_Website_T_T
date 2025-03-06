package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
	"os"
)

// startTime holds the time when the service was started (for uptime calculation)
var startTime time.Time

// CountryInfo defines the response structure for the info endpoint.
type CountryInfo struct {
	Name       string            `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flag       string            `json:"flag"`
	Capital    string            `json:"capital"`
	Cities     []string          `json:"cities"`
}

// RestCountry represents the JSON structure from the REST Countries API.
type RestCountry struct {
	Name struct {
		Common string `json:"common"`
	} `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flags      struct {
		PNG string `json:"png"`
	} `json:"flags"`
	Capital []string `json:"capital"`
}

// populationAPIResponse represents the response from the CountriesNow population API.
type populationAPIResponse struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
	Data  struct {
		Country          string `json:"country"`
		Code             string `json:"code"`
		Iso3             string `json:"iso3"`
		PopulationCounts []struct {
			Year  int `json:"year"`
			Value int `json:"value"`
		} `json:"populationCounts"`
	} `json:"data"`
}

// fetchRestCountry retrieves country data from the REST Countries API using the country code.
func fetchRestCountry(countryCode string) (*RestCountry, error) {
	url := fmt.Sprintf("http://129.241.150.113:8080/v3.1/alpha/%s", countryCode)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error calling REST Countries API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading REST Countries API response: %w", err)
	}

	var countries []RestCountry
	if err := json.Unmarshal(body, &countries); err != nil || len(countries) == 0 {
		return nil, fmt.Errorf("error parsing REST Countries API response: %w", err)
	}
	return &countries[0], nil
}

// fetchCities retrieves the list of cities for the given country from the CountriesNow API.
func fetchCities(countryName string) ([]string, error) {
	url := "http://129.241.150.113:3500/api/v0.1/countries/cities"
	payload := map[string]string{"country": countryName}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling cities payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling CountriesNow cities API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading cities API response: %w", err)
	}

	var result struct {
		Error bool     `json:"error"`
		Msg   string   `json:"msg"`
		Data  []string `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error parsing cities API response: %w", err)
	}

	return result.Data, nil
}

// fetchPopulation retrieves the population data from the CountriesNow API using the official country name.
func fetchPopulation(countryName string) (*populationAPIResponse, error) {
	url := "http://129.241.150.113:3500/api/v0.1/countries/population"
	payload := map[string]string{"country": countryName}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling population payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating population request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error calling population API: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading population API response: %w", err)
	}

	// Removed debug logging as it is not required for production.
	var popResp populationAPIResponse
	if err := json.Unmarshal(body, &popResp); err != nil {
		return nil, fmt.Errorf("error parsing population API response: %w", err)
	}
	return &popResp, nil
}

// computeMean calculates the average of the given population values.
func computeMean(counts []struct {
	Year  int `json:"year"`
	Value int `json:"value"`
}) int {
	if len(counts) == 0 {
		return 0
	}
	sum := 0
	for _, c := range counts {
		sum += c.Value
	}
	return sum / len(counts)
}

// ==========================
// countryInfoHandler
// Handles requests to /countryinfo/v1/info/{code}?limit={cityLimit}
// Returns country details and an optional limited list of cities (alphabetically sorted).
// ==========================
func countryInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Extract country code from URL.
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 || pathParts[4] == "" {
		http.Error(w, "Country code not provided", http.StatusBadRequest)
		return
	}
	countryCode := strings.ToLower(pathParts[4])

	// Parse the "limit" query parameter for city list limit.
	limitStr := r.URL.Query().Get("limit")
	var cityLimit int
	if limitStr != "" {
		val, err := strconv.Atoi(limitStr)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
		cityLimit = val
	}

	// 1) Get official country data.
	restCountry, err := fetchRestCountry(countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2) Get list of cities for the country.
	countryName := restCountry.Name.Common
	cities, err := fetchCities(countryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Sort cities alphabetically.
	sort.Strings(cities)
	// Limit the city list if a valid limit is provided.
	if cityLimit > 0 && cityLimit < len(cities) {
		cities = cities[:cityLimit]
	}

	// 3) Construct response.
	info := CountryInfo{
		Name:       countryName,
		Continents: restCountry.Continents,
		Population: restCountry.Population,
		Languages:  restCountry.Languages,
		Borders:    restCountry.Borders,
		Flag:       restCountry.Flags.PNG,
		Cities:     cities,
	}
	if len(restCountry.Capital) > 0 {
		info.Capital = restCountry.Capital[0]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

// ==========================
// populationHandler
// Handles requests to /countryinfo/v1/population/{code}?limit=YYYY-YYYY
// Returns population data for the specified year range (if provided) and calculates the mean.
// ==========================
func populationHandler(w http.ResponseWriter, r *http.Request) {
	// Extract country code from URL.
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 5 || pathParts[4] == "" {
		http.Error(w, "Country code not provided", http.StatusBadRequest)
		return
	}
	countryCode := strings.ToLower(pathParts[4])

	// Parse the "limit" query parameter (expected format: YYYY-YYYY).
	limitStr := r.URL.Query().Get("limit")
	var startYear, endYear int
	var hasYearRange bool
	if limitStr != "" {
		parts := strings.Split(limitStr, "-")
		if len(parts) == 2 {
			sY, err1 := strconv.Atoi(parts[0])
			eY, err2 := strconv.Atoi(parts[1])
			if err1 != nil || err2 != nil {
				http.Error(w, "Invalid year range in limit parameter", http.StatusBadRequest)
				return
			}
			startYear, endYear = sY, eY
			hasYearRange = true
		} else {
			http.Error(w, "limit parameter must be in format YYYY-YYYY", http.StatusBadRequest)
			return
		}
	}

	// 1) Get official country data to determine the official country name.
	restCountry, err := fetchRestCountry(countryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 2) Get population data using the official country name.
	countryName := restCountry.Name.Common
	popResp, err := fetchPopulation(countryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 3) Filter population counts by the year range if specified.
	var filtered []struct {
		Year  int `json:"year"`
		Value int `json:"value"`
	}
	for _, pc := range popResp.Data.PopulationCounts {
		if hasYearRange {
			if pc.Year >= startYear && pc.Year <= endYear {
				filtered = append(filtered, pc)
			}
		} else {
			filtered = append(filtered, pc)
		}
	}

	// 4) Calculate the mean population.
	meanVal := computeMean(filtered)

	// 5) Construct and send response.
	response := map[string]interface{}{
		"mean":             meanVal,
		"populationCounts": filtered,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// getServiceStatus performs a simple GET request to an external URL and returns its HTTP status code.
// Returns 0 if an error occurs.
func getServiceStatus(url string) int {
	resp, err := http.Get(url)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return resp.StatusCode
}

// ==========================
// diagnosticsHandler
// Handles requests to /countryinfo/v1/status/
// Returns a JSON object with the status codes of external services, version, and uptime.
// ==========================
func diagnosticsHandler(w http.ResponseWriter, r *http.Request) {
	// Check status of CountriesNow API.
	countriesNowStatus := getServiceStatus("http://129.241.150.113:3500/api/v0.1/countries")
	// Check status of REST Countries API.
	restCountriesStatus := getServiceStatus("http://129.241.150.113:8080/v3.1/all")
	// Calculate uptime in seconds.
	uptime := int(time.Since(startTime).Seconds())

	response := map[string]interface{}{
		"countriesnowapi":  countriesNowStatus,
		"restcountriesapi": restCountriesStatus,
		"version":          "v1",
		"uptime":           uptime,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ==========================
// main
// Sets up file serving and endpoints, then starts the HTTP server.
// ==========================
func main() {
	// Record the service start time.
	startTime = time.Now()

	// Serve static files from the "./FrontEnd" directory.
	fs := http.FileServer(http.Dir("./FrontEnd"))
	http.Handle("/", fs)

	// Set up API endpoints.
	http.HandleFunc("/countryinfo/v1/info/", countryInfoHandler)
	http.HandleFunc("/countryinfo/v1/population/", populationHandler)
	http.HandleFunc("/countryinfo/v1/status/", diagnosticsHandler)

	// Listen on the port defined in the PORT environment variable if available.
	port := "8080"
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}
	fmt.Printf("Server running on port %s...\n", port)
	http.ListenAndServe(":" + port, nil)
}
