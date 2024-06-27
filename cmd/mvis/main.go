package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Print(err)
	}
	// TODO: FIGURE OUT A WAY TO DO THIS PROGRAMATICALLY (WRITE TO ENV)
	val, err := refreshAccessToken(os.Getenv("CLIENT_ID"), os.Getenv("CLIENT_SECRET"), os.Getenv("REFRESH_TOKEN"))
	if err != nil {
		log.Print("err getting token" + err.Error())
	}

	log.Print()
	log.Print(val)
}
