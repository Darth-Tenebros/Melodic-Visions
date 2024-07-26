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

	time_range = "long_term"
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

	duration := time.Since(expiry)
	if duration > time.Hour {
		fmt.Println("Token expired, refreshing...")
		newToken, newExpiry, err := refreshAccessToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), os.Getenv("REFRESH_TOKEN"))
		if err != nil {
			log.Printf("Error refreshing token: %v \n", err)
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
		log.Println(err)
	}

	// get the number od tracks by each artist
	songsPerArtist := render_charts.AggregateArtistTotalTracks(tracks)

	var topten map[string]interface{} = make(map[string]interface{})
	topten = top10FromMap(songsPerArtist)

	var keys []string
	var values []int

	for key, value := range topten {
		keys = append(keys, key)
		values = append(values, value.(int))
	}

	audioFeatures, err := spotify.GetPaginatedAudioFeatures(os.Getenv(tokenEnvVar), trackIds)
	if err != nil {
		log.Println(err)
	}

	acousticness, danceability, valence := model.RetrieveAudioFeatures(audioFeatures)

	// render graphs
	wordcloud := render_charts.WordCloudBasic(topten)
	f, _ := os.Create("wordcloud.html")
	wordcloud.Render(f)

	pie := render_charts.PieBasic(topten)
	f, _ = os.Create("pie.html")
	pie.Render(f)

	scatter := render_charts.ScatterBasic(trackIds, acousticness, danceability, valence)
	f, _ = os.Create("scatter.html")
	scatter.Render(f)

}

func top10FromMap(input map[string]int) map[string]interface{} {
	// Create a slice to hold the key-value pairs
	var kvSlice []model.ArtistDuration
	for k, v := range input {
		kvSlice = append(kvSlice, model.ArtistDuration{ArtistName: k, ArtistCount: v})
	}

	// Sort the slice based on the values in descending order
	sort.Slice(kvSlice, func(i, j int) bool {
		return kvSlice[i].ArtistCount > kvSlice[j].ArtistCount
	})

	// Create a map to hold the top 10 key-value pairs
	top10Map := make(map[string]interface{})
	for i := 0; i < len(kvSlice) && i < 15; i++ {
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
