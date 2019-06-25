package main

import (
	"log"
	"os"
)

func main() {
	// Get environment variables.
	token := os.Getenv("DISCORD_HACK_2019_TOKEN")
	mongoURI := os.Getenv("DISCORD_HACK_2019_MONGODB")

	// If vars are empty/unset, quit.
	if token == "" || mongoURI == "" {
		log.Fatalln("The env vars need to be set")
	}
}
