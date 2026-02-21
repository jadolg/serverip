package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockResponse = WTFIsMyIPData{
	YourFuckingIPAddress:   "1.2.3.4",
	YourFuckingLocation:    "Somewhere",
	YourFuckingHostname:    "host.example.com",
	YourFuckingISP:         "ExampleISP",
	YourFuckingTorExit:     false,
	YourFuckingCity:        "ExampleCity",
	YourFuckingCountry:     "ExampleCountry",
	YourFuckingCountryCode: "EC",
}

func TestGetIpAddressInformation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewEncoder(w).Encode(mockResponse); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	oldURL := wtfismyipIPv4URL
	wtfismyipIPv4URL = server.URL
	defer func() { wtfismyipIPv4URL = oldURL }()

	got, err := getIpAddressInformation(context.Background(), false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if got.YourFuckingIPAddress != mockResponse.YourFuckingIPAddress {
		t.Errorf("expected IP %s, got %s", mockResponse.YourFuckingIPAddress, got.YourFuckingIPAddress)
	}
	if got.YourFuckingLocation != mockResponse.YourFuckingLocation {
		t.Errorf("expected Location %s, got %s", mockResponse.YourFuckingLocation, got.YourFuckingLocation)
	}
	if got.YourFuckingHostname != mockResponse.YourFuckingHostname {
		t.Errorf("expected Hostname %s, got %s", mockResponse.YourFuckingHostname, got.YourFuckingHostname)
	}
	if got.YourFuckingISP != mockResponse.YourFuckingISP {
		t.Errorf("expected ISP %s, got %s", mockResponse.YourFuckingISP, got.YourFuckingISP)
	}
	if got.YourFuckingTorExit != mockResponse.YourFuckingTorExit {
		t.Errorf("expected TorExit %v, got %v", mockResponse.YourFuckingTorExit, got.YourFuckingTorExit)
	}
	if got.YourFuckingCity != mockResponse.YourFuckingCity {
		t.Errorf("expected City %s, got %s", mockResponse.YourFuckingCity, got.YourFuckingCity)
	}
	if got.YourFuckingCountry != mockResponse.YourFuckingCountry {
		t.Errorf("expected Country %s, got %s", mockResponse.YourFuckingCountry, got.YourFuckingCountry)
	}
	if got.YourFuckingCountryCode != mockResponse.YourFuckingCountryCode {
		t.Errorf("expected CountryCode %s, got %s", mockResponse.YourFuckingCountryCode, got.YourFuckingCountryCode)
	}
}

func TestGetIpAddressInformationIPv6(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(mockResponse)
		if err != nil {
			slog.Error("failed to encode JSON", "error", err)
		}
	}))
	defer server.Close()

	oldURL := wtfismyipIPv6URL
	wtfismyipIPv6URL = server.URL
	defer func() { wtfismyipIPv6URL = oldURL }()

	_, err := getIpAddressInformation(context.Background(), true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
