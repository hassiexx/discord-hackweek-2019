package main

import (
	"log"
	"os"

	"github.com/hassieswift621/discord-hackweek-2019/cmdlistener"
	"github.com/hassieswift621/discord-hackweek-2019/core"
	"github.com/hassieswift621/discord-hackweek-2019/db"
	"github.com/hassieswift621/discord-hackweek-2019/info"
	"github.com/hassieswift621/discord-hackweek-2019/moderation"
)

func main() {
	// Get environment variables.
	token := os.Getenv("DISCORD_HACK_WEEK_2019_TOKEN")
	mongoURI := os.Getenv("DISCORD_HACK_WEEK_2019_MONGODB")

	// If vars are empty/unset, quit.
	if token == "" || mongoURI == "" {
		log.Fatalln("The env vars need to be set")
	}

	// Connect to DB.
	err := db.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// Create client.
	client, err := core.NewClient(token)
	if err != nil {
		log.Fatalln(err)
	}

	// Initialise cmd listener.
	cmdlistener.Initialise(client)

	// Initialise moderation commands.
	moderation.Initialise(client)

	// Initialise info commands.
	info.Initialise(client)

	// Connect.
	err = client.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	// Block from exiting.
	select {}
}
