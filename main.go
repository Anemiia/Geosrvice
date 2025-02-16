package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log"
	"net/http"
)

const dadataURL = "https://suggestions.dadata.ru/suggestions/api/4_1/rs/suggest/address"

type RequestAddressSearch struct {
	Query string `json:"query"`
}

type RequestAddressGeocode struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

type Address struct {
	City   string `json:"city"`
	Street string `json:"street"`
}

type ResponseAddress struct {
	Addresses []*Address `json:"addresses"`
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Гео сервис запущен!"))
}

func AddressGeocodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод запроса не POST", http.StatusMethodNotAllowed)
		return
	}

	var req RequestAddressGeocode
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	result, err := GeocodeAddress(req.Lat, req.Lng)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func GeocodeAddress(lat, lng string) (*ResponseAddress, error) {
	dadataAPIKey := "d0e0b799107963c33b1f6e8cd96547b73b46bb16"

	reqBody, _ := json.Marshal(map[string]string{
		"lat": lat,
		"lon": lng,
	})

	req, err := http.NewRequest("POST", "https://suggestions.dadata.ru/suggestions/api/4_1/rs/geolocate/address", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Token "+dadataAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from DaData, status code: %d", resp.StatusCode)
	}

	var rawResponse struct {
		Suggestions []struct {
			Data Address `json:"data"`
		} `json:"suggestions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}
	var response ResponseAddress
	for _, suggestion := range rawResponse.Suggestions {
		response.Addresses = append(response.Addresses, &Address{
			City:   suggestion.Data.City,
			Street: suggestion.Data.Street,
		})
	}

	return &response, nil
}

func AddressSearchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не POST", http.StatusMethodNotAllowed)
		return
	}

	var req RequestAddressSearch
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
	}
	fmt.Println("Request body:", req)

	result, err := SearchAddress(req.Query)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func SearchAddress(query string) (*ResponseAddress, error) {
	dadataToken := "d0e0b799107963c33b1f6e8cd96547b73b46bb16"

	reqBody, _ := json.Marshal(map[string]string{"query": query})

	req, err := http.NewRequest("POST", dadataURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Token "+dadataToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch data from DaData, status code: %d", resp.StatusCode)
	}

	var rawResponse struct {
		Suggestions []struct {
			Data Address `json:"data"`
		} `json:"suggestions"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, err
	}

	var response ResponseAddress
	for _, suggestion := range rawResponse.Suggestions {
		response.Addresses = append(response.Addresses, &Address{
			City:   suggestion.Data.City,
			Street: suggestion.Data.Street,
		})
	}

	fmt.Println("Response body:", response)
	return &response, nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", HomeHandler)
	r.Post("/api/address/search", AddressSearchHandler)
	r.Post("/api/address/geocode", AddressGeocodeHandler)

	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
