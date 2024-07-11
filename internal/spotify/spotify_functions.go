package spotify

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Darth-Tenebros/Melodic-Visions/internal/model"
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

func GetUserTopItems(accessToken, reqUrl string) (model.SpotifyResponse, error) {

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return model.SpotifyResponse{}, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return model.SpotifyResponse{}, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.SpotifyResponse{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return model.SpotifyResponse{}, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, body)
	}

	var result model.SpotifyResponse
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return model.SpotifyResponse{}, fmt.Errorf("err unmarshaling data")
	}

	return result, err
}

func GetAudioFeatures(accessToken string, tracksIds []string) (model.AudioFeatures, error) {

	ids := strings.Join(tracksIds, ",")

	GET_AUDIO_FEATURES_URL := "https://api.spotify.com/v1/audio-features?ids"
	fullUrl := fmt.Sprintf("%s=%s", GET_AUDIO_FEATURES_URL, ids)

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return model.AudioFeatures{}, fmt.Errorf("error setting up request for audio features: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return model.AudioFeatures{}, fmt.Errorf("error getting audio features: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.AudioFeatures{}, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return model.AudioFeatures{}, fmt.Errorf("request failed with status code %d: %s", resp.StatusCode, body)
	}

	var result model.AudioFeatures
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return model.AudioFeatures{}, fmt.Errorf("err unmarshaling data")
	}

	return result, nil
}

func Write_file(data string) error {

	file, err := os.OpenFile("data.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	defer file.Close()

	_, err = file.Write([]byte(data + "\n"))
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	return err
}
