package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Darth-Tenebros/Melodic-Visions/internal/spotify"
	"github.com/joho/godotenv"
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

	time_listened := 0
	for {

		for _, item := range result.Items {
			time_listened = item.DurationMs + time_listened
			spotify.Write_file(fmt.Sprintf("%s <:> %s", item.Name, item.Artists))
		}
		fmt.Println("SUCCESS!!")

		if err != nil {
			log.Print(err)
		}

		fmt.Println(result.Next)
		if strings.Contains(result.Next, "http") {
			reqUrl = result.Next
			result, err = spotify.GetUserTopItems(os.Getenv("ACCESS_TOKEN"), reqUrl)
		} else {
			break
		}
	}

	fmt.Println(time_listened)
	duration := time.Duration(time_listened) * time.Millisecond
	fmt.Printf("your listened for %f minutes\n", duration.Minutes())
	fmt.Printf("whch is %f hours", duration.Hours())

}
