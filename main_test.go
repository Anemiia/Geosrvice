package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddressSearchHandler(t *testing.T) {
	// Создание нового запроса для /api/address/search
	reqBody := `{"query": "Москва, Красная площадь"}`
	req, err := http.NewRequest("POST", "/api/address/search", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatal(err)
	}

	// Мокаем сервер и добавляем обработчик
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddressSearchHandler)

	// Выполнение запроса
	handler.ServeHTTP(rr, req)

	// Проверка статуса
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %v", rr.Code)
	}

	// Проверка ответа (например, что в ответе есть адрес с городом "Москва")
	var result ResponseAddress
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if len(result.Addresses) == 0 || result.Addresses[0].City != "Москва" {
		t.Errorf("expected city 'Москва', got %v", result.Addresses[0].City)
	}
}

func TestAddressGeocodeHandler(t *testing.T) {
	// Создание нового запроса для /api/address/geocode
	reqBody := `{"lat": "41.6801619", "lng": "48.1714836"}`
	req, err := http.NewRequest("POST", "/api/address/geocode", bytes.NewBuffer([]byte(reqBody)))
	if err != nil {
		t.Fatal(err)
	}

	// Мокаем сервер и добавляем обработчик
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(AddressGeocodeHandler)

	// Выполнение запроса
	handler.ServeHTTP(rr, req)

	// Проверка статуса
	if rr.Code != http.StatusOK {
		t.Errorf("expected status code 200, got %v", rr.Code)
	}

	// Проверка ответа (например, что в ответе есть улица "Тельмана")
	var result ResponseAddress
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}

	if len(result.Addresses) == 0 || result.Addresses[0].Street != "Тельмана" {
		t.Errorf("expected street 'Тельмана', got %v", result.Addresses[0].Street)
	}
}
