package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
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

	limit := 50
	time_range := "short_term"
	offset := 0
	TOP_ITEMS_URL := "https://api.spotify.com/v1/me/top/"
	reqUrl := fmt.Sprintf("%stracks?time_range=%s&limit=%d&offset=%d", TOP_ITEMS_URL, time_range, limit, offset)

	// request Top User Items (tracks) from spotify
	result, err := spotify.GetUserTopItems(os.Getenv("ACCESS_TOKEN"), reqUrl)
	if err != nil {
		fmt.Println(err)
	}

	// paginate through results
	var tracks []model.Track
	for {

		data := convertItemsToTracks(result.Items)
		tracks = append(tracks, data...)

		if err != nil {
			log.Print(err)
		}

		if strings.Contains(result.Next, "http") {
			reqUrl = result.Next
			result, err = spotify.GetUserTopItems(os.Getenv(tokenEnvVar), reqUrl)
		} else {
			break
		}
	}

	// for _, track := range tracks {
	// 	appendToFile("data.txt", track.ArtistName+"\n")
	// }

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

	bar := render_charts.BarBasic(keys, values)
	f, _ := os.Create("bar.html")
	bar.Render(f)

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

func convertItemToTrack(item model.Item) model.Track {
	artistName := ""
	if len(item.Artists) > 0 {
		artistName = item.Artists[0].Name
	}

	return model.Track{
		TrackName:  item.Name,
		AlbumName:  item.Album.Name,
		Duration:   item.DurationMs,
		ArtistName: artistName,
	}
}

func convertItemsToTracks(items []model.Item) []model.Track {
	tracks := make([]model.Track, len(items))
	for i, item := range items {
		tracks[i] = convertItemToTrack(item)
	}
	return tracks
}

type ArtistDuration struct {
	ArtistName  string
	ArtistCount int
}
