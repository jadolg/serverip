// Package main provides a web utility to identify the public exit IP address of the server.
package main

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"
)

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

var (
	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	wtfismyipIPv4URL = "https://ipv4.wtfismyip.com/json"
	wtfismyipIPv6URL = "https://ipv6.wtfismyip.com/json"
)

func getIpAddressInformation(ctx context.Context, ipv6 bool) (WTFIsMyIPData, error) {
	url := wtfismyipIPv4URL
	if ipv6 {
		url = wtfismyipIPv6URL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return WTFIsMyIPData{}, err
	}

	response, err := httpClient.Do(req)
	if err != nil {
		return WTFIsMyIPData{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("failed to close response body", "error", err)
		}
	}(response.Body)

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
