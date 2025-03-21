package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// WTFIsMyIPData is a data representation with the same structure returned by https://wtfismyip.com/json
type WTFIsMyIPData struct {
	YourFuckingIPAddress   string `json:"YourFuckingIPAddress"`
	YourFuckingLocation    string `json:"YourFuckingLocation"`
	YourFuckingHostname    string `json:"YourFuckingHostname"`
	YourFuckingISP         string `json:"YourFuckingISP"`
	YourFuckingTorExit     bool   `json:"YourFuckingTorExit"`
	YourFuckingCity        string `json:"YourFuckingCity"`
	YourFuckingCountry     string `json:"YourFuckingCountry"`
	YourFuckingCountryCode string `json:"YourFuckingCountryCode"`
}

func getIpAddressInformation(ipv6 bool) (WTFIsMyIPData, error) {
	wtfismyipURL := "https://ipv4.wtfismyip.com/json"
	if ipv6 {
		wtfismyipURL = "https://ipv6.wtfismyip.com/json"
	}
	response, err := http.Get(wtfismyipURL)
	if err != nil {
		return WTFIsMyIPData{}, err
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return WTFIsMyIPData{}, err
	}

	data := WTFIsMyIPData{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return WTFIsMyIPData{}, err
	}

	return data, nil
}
