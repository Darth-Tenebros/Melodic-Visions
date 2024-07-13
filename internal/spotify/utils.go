package spotify

import (
	"fmt"
	"log"
	"strings"

	"github.com/Darth-Tenebros/Melodic-Visions/internal/model"
)

func GetTopTracks(time_range string, token string) ([]model.Track, []string, error) {

	limit := 50
	offset := 0
	TOP_ITEMS_URL := "https://api.spotify.com/v1/me/top/"
	reqUrl := fmt.Sprintf("%stracks?time_range=%s&limit=%d&offset=%d", TOP_ITEMS_URL, time_range, limit, offset)

	// request Top User Items (tracks) from spotify
	result, err := GetUserTopItems(token, reqUrl)
	if err != nil {
		return []model.Track{}, nil, err
	}

	// paginate through results
	var tracks []model.Track
	var tracksIds []string
	for {

		data, ids := model.ConvertItemsToTracks(result.Items)
		tracks = append(tracks, data...)
		tracksIds = append(tracksIds, ids...)

		if err != nil {
			log.Println(err)
			return []model.Track{}, nil, err
		}

		if strings.Contains(result.Next, "http") {
			reqUrl = result.Next
			result, err = GetUserTopItems(token, reqUrl)
		} else {
			break
		}
	}

	return tracks, tracksIds, nil

}

func GetPaginatedAudioFeatures(accessToken string, trackids []string) ([]model.AudioFeature, error) {
	iter := 100
	var audio_features []model.AudioFeature

	for i := 0; i < len(trackids); i += iter {
		end := i + iter
		if end > len(trackids) {
			end = len(trackids)
		}

		batch := trackids[i:end]
		result, err := GetAudioFeatures(accessToken, batch)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		audio_features = append(audio_features, result.AudioFeatures...)

	}

	return audio_features, nil
}
