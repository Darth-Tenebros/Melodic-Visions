package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	charts "github.com/Darth-Tenebros/Melodic-Visions/internal/charts"
	model "github.com/Darth-Tenebros/Melodic-Visions/internal/model"
	"github.com/Darth-Tenebros/Melodic-Visions/internal/spotify"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print(err)
	}
	// TODO: FIGURE OUT A WAY TO DO THIS PROGRAMATICALLY (WRITE TO ENV)
	// val, err := refreshAccessToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), os.Getenv("REFRESH_TOKEN"))
	// if err != nil {
	// 	log.Print("err getting token" + err.Error())
	// }
	// log.Print()
	// log.Print(val)

	limit := 50
	time_range := "long_term"
	offset := 0
	TOP_ITEMS_URL := "https://api.spotify.com/v1/me/top/"
	reqUrl := fmt.Sprintf("%stracks?time_range=%s&limit=%d&offset=%d", TOP_ITEMS_URL, time_range, limit, offset)

	result, err := spotify.GetUserTopItems(os.Getenv("ACCESS_TOKEN"), reqUrl)
	if err != nil {
		fmt.Println(err)
	}

	// database, err := sql.Open("sqlite3", "/home/yolisa/Documents/Projects/Melodic-Visions/data/spotify_data")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer database.Close()

	// track_repo := repository.NewTrackRepository(database)

	// time_listened := 0
	var tracks []model.Track
	for {

		data := convertItemsToTracks(result.Items)
		tracks = append(tracks, data...)

		if err != nil {
			log.Print(err)
		}

		if strings.Contains(result.Next, "http") {
			reqUrl = result.Next
			result, err = spotify.GetUserTopItems(os.Getenv("ACCESS_TOKEN"), reqUrl)
		} else {
			break
		}
	}

	// for _, track := range tracks {
	// 	_, err := track_repo.InsertTrack("long_term", track)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }

	top := charts.AggregateArtistTotalDurationListened(tracks)

	var topTen []ArtistDuration
	for k, v := range top {
		topTen = append(topTen, ArtistDuration{ArtistName: k, ArtiDuration: v})
	}

	sort.Slice(topTen, func(i, j int) bool {
		return topTen[i].ArtiDuration > topTen[j].ArtiDuration
	})

	for i := 0; i < 15; i++ {
		duration := time.Duration(topTen[i].ArtiDuration) * time.Millisecond
		fmt.Printf("%s ==> %v\n", topTen[i].ArtistName, duration.Minutes())
	}
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
	ArtistName   string
	ArtiDuration int
}
