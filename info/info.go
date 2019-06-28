package info

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hassieswift621/discord-hackweek-2019/core"
)

// Initialise registers info commands.
func Initialise(client *core.DiscordClient) {
	// Register info commands.

	// Ping.
	client.RegisterCommand("ping", func(connection *discordgo.Session, message *discordgo.Message, args []string) {
		(&ping{client: client, connection: connection, message: message, args: args}).execute()
	})
}
