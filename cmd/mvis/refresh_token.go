package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func refreshAccessToken(CLIENT_ID, CLIENT_SECRET, REFRESH_TOKEN string) (string, error) {
	// Encode credentials
	TOKEN_URL := "https://accounts.spotify.com/api/token"

	credentials := CLIENT_ID + ":" + CLIENT_SECRET
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))

	// Prepare form data
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", REFRESH_TOKEN)

	// Prepare request
	req, err := http.NewRequest("POST", TOKEN_URL, strings.NewReader(data.Encode()))
	log.Print(req.Body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Basic "+encodedCredentials)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse JSON response to extract access token
	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	accessToken, ok := responseMap["access_token"].(string)
	if !ok {
		return "", fmt.Errorf("failed to retrieve access token from response")
	}

	return accessToken, nil
}
