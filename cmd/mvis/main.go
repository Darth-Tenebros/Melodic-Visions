package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	model "github.com/Darth-Tenebros/Melodic-Visions/internal/model"
	render_charts "github.com/Darth-Tenebros/Melodic-Visions/internal/render_charts"
	"github.com/Darth-Tenebros/Melodic-Visions/internal/spotify"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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
	time_range := "short_term"
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

	top := render_charts.AggregateArtistTotalDurationListened(tracks)

	// var topTen []ArtistDuration
	// for k, v := range top {
	// 	topTen = append(topTen, ArtistDuration{ArtistName: k, ArtiDuration: v})
	// }

	// keys := make([]string, len(topTen))
	// values := make([]int, len(topTen))
	var keys []string
	var values []int

	end := 0
	for key, value := range top {
		end++
		keys = append(keys, key)
		values = append(values, value)
		if end == 20 {
			break
		}
	}

	bar := barBasic(keys, values)
	f, _ := os.Create("bar.html")
	bar.Render(f)

	// sort.Slice(topTen, func(i, j int) bool {
	// 	return topTen[i].ArtiDuration > topTen[j].ArtiDuration
	// })

	// for i := 0; i < 15; i++ {
	// 	duration := time.Duration(topTen[i].ArtiDuration) * time.Millisecond
	// 	fmt.Printf("%s ==> %v\n", topTen[i].ArtistName, duration.Minutes())
	// }
}

// TODO: CLEAN UP RENDERING
func generateBarItems(values []int) []opts.BarData {
	items := make([]opts.BarData, 0)
	for i := 0; i < len(values); i++ {
		items = append(items, opts.BarData{Value: values[i]})
	}
	return items
}

func barBasic(keys []string, values []int) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic bar example", Subtitle: "This is the subtitle."}),
	)

	bar.SetXAxis(keys).
		AddSeries("Category A", generateBarItems(values))
	return bar
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
