package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"time"

	model "github.com/Darth-Tenebros/Melodic-Visions/internal/model"
	"github.com/Darth-Tenebros/Melodic-Visions/internal/render_charts"
	"github.com/Darth-Tenebros/Melodic-Visions/internal/spotify"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const (
	refreshToken = "REFRESH_TOKEN"
	tokenEnvVar  = "ACCESS_TOKEN"
	expiryEnvVar = "TOKEN_EXPIRY"

	time_range = "short_term"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print(err)
	}

	// check if the the current token has expired
	expiryStr := os.Getenv(expiryEnvVar)
	fmt.Println(expiryStr)

	expiry, err := time.Parse(time.RFC3339, expiryStr)
	if err != nil {
		fmt.Println("Invalid expiry time:", err)
		return
	}

	duration := time.Now().Sub(expiry)
	if duration > time.Hour {
		fmt.Println("Token expired, refreshing...")
		newToken, newExpiry, err := refreshAccessToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), os.Getenv("REFRESH_TOKEN"))
		if err != nil {
			fmt.Printf("Error refreshing token: %v \n", err)
			return
		}

		err = updateEnvToken(newToken, newExpiry)
		if err != nil {
			fmt.Printf("Error updating environment variables: %v", err)
		}
	}

	// request Top User Items (tracks) from spotify
	tracks, trackIds, err := spotify.GetTopTracks(time_range, os.Getenv(tokenEnvVar))
	if err != nil {
		fmt.Println(err)
	}

	// get the number od tracks by each artist
	songsPerArtist := render_charts.AggregateArtistTotalTracks(tracks)

	topten := top10FromMap(songsPerArtist)

	var keys []string
	var values []int

	end := 0
	for key, value := range topten {
		end++
		keys = append(keys, key)
		values = append(values, value)
		// if end == 10 {
		// 	break
		// }
	}

	audioFeatures, err := paginateAudioFeatures(os.Getenv(tokenEnvVar), trackIds)
	if err != nil {
		fmt.Println(err)
	}

	acousticness, danceability, valence := retrieveAudioFeatures(audioFeatures)

	// render graphs
	bar := render_charts.BarBasic(keys, values)
	f, _ := os.Create("bar.html")
	bar.Render(f)

	pie := render_charts.PieBasic(topten)
	f, _ = os.Create("pie.html")
	pie.Render(f)

	scatter := render_charts.ScatterBasic(trackIds, acousticness, danceability, valence)
	f, _ = os.Create("scatter.html")
	scatter.Render(f)

}

func paginateAudioFeatures(accessToken string, trackids []string) ([]model.AudioFeature, error) {
	iter := 100
	var audio_features []model.AudioFeature

	for i := 0; i < len(trackids); i += iter {
		end := i + iter
		if end > len(trackids) {
			end = len(trackids)
		}

		batch := trackids[i:end]
		result, err := spotify.GetAudioFeatures(accessToken, batch)
		if err != nil {
			log.Print(err)
			return nil, err
		}

		audio_features = append(audio_features, result.AudioFeatures...)

	}

	return audio_features, nil
}

func appendToFile(filename string, text string) error {
	// Open the file with O_APPEND and O_WRONLY flags
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Write the text to the file
	_, err = file.WriteString(text)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}

func top10FromMap(input map[string]int) map[string]int {
	// Create a slice to hold the key-value pairs
	var kvSlice []ArtistDuration
	for k, v := range input {
		kvSlice = append(kvSlice, ArtistDuration{k, v})
	}

	// Sort the slice based on the values in descending order
	sort.Slice(kvSlice, func(i, j int) bool {
		return kvSlice[i].ArtistCount > kvSlice[j].ArtistCount
	})

	// Create a map to hold the top 10 key-value pairs
	top10Map := make(map[string]int)
	for i := 0; i < len(kvSlice) && i < 10; i++ {
		top10Map[kvSlice[i].ArtistName] = kvSlice[i].ArtistCount
	}

	return top10Map
}

func updateEnvToken(newToken string, newExpiry time.Time) error {
	cmd := fmt.Sprintf("export %s=%s; export %s=%s", tokenEnvVar, newToken, expiryEnvVar, newExpiry.Format(time.RFC3339))
	_, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		return fmt.Errorf("failed to update environment variable: %v", err)
	}

	os.Setenv(tokenEnvVar, newToken)
	os.Setenv(expiryEnvVar, newExpiry.Format(time.RFC3339))
	return nil
}

func retrieveAudioFeatures(items []model.AudioFeature) ([]float64, []float64, []float64) {
	var acousticness []float64
	var danceability []float64
	var valence []float64

	for _, feature := range items {
		acousticness = append(acousticness, feature.Acousticness)
		danceability = append(danceability, feature.Danceability)
		valence = append(valence, feature.Valence)
	}

	return acousticness, danceability, valence
}

type ArtistDuration struct {
	ArtistName  string
	ArtistCount int
}
