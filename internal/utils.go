package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func GetRecentlyPlayed(accessToken string) error {
	after := 1701381600000
	limit := 50
	RECENTLY_PLAYED_URL := "https://api.spotify.com/v1/me/player/recently-played"
	reqUrl := fmt.Sprintf("%s?limit=%d&after=%d", RECENTLY_PLAYED_URL, limit, after)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, body)
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(body, &responseMap); err != nil {
		return fmt.Errorf("failed to unmarshal response: %v", err)
	}

	fmt.Println("next:	", responseMap["next"])
	fmt.Println("after:	", responseMap["after"])
	return nil
}

func GetUserTopItems(accessToken string) error {
	limit := 50
	time_range := "long_term"
	offset := 0
	TOP_ITEMS_URL := "https://api.spotify.com/v1/me/top/"
	reqUrl := fmt.Sprintf("%stracks?time_range=%s&limit=%d&offset=%d", TOP_ITEMS_URL, time_range, limit, offset)

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, body)
	}

	var result SpotifyResponse
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return fmt.Errorf("err unmarshaling data")
	}

	for _, item := range result.Items {
		fmt.Println(item.Album.Name)
		fmt.Println(item.Name)
		fmt.Println()
	}

	return err
}
