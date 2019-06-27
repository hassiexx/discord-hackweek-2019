package moderation

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hassieswift621/discord-hackweek-2019/core"
)

// Initialise performs initialisation for the moderation commands by registering them.
func Initialise(client *core.DiscordClient) {
	// Register moderation commands.

	// Kick.
	client.RegisterCommand("kick", func(connection *discordgo.Session, message *discordgo.Message, args []string) {
		(&kick{client: client, connection: connection, message: message, args: args,
			menuData: &core.MenuCommandData{}, modData: &moderationData{}}).execute()
	})
}
