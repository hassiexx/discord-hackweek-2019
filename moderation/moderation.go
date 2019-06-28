package moderation

import (
	"github.com/bwmarrin/discordgo"
	"github.com/hassieswift621/discord-hackweek-2019/core"
)

// Initialise performs initialisation for the moderation commands by registering them.
func Initialise(client *core.DiscordClient) {
	// Register moderation commands.

	// Ban.
	client.RegisterCommand("ban", func(connection *discordgo.Session, message *discordgo.Message, args []string) {
		if isDM(connection, message) {
			return
		}
		(&ban{client: client, connection: connection, message: message, args: args,
			menuData: &core.MenuCommandData{}, modData: &moderationData{}}).execute()
	})

	// Kick.
	client.RegisterCommand("kick", func(connection *discordgo.Session, message *discordgo.Message, args []string) {
		if isDM(connection, message) {
			return
		}
		(&kick{client: client, connection: connection, message: message, args: args,
			menuData: &core.MenuCommandData{}, modData: &moderationData{}}).execute()
	})

	// Log Channel.
	client.RegisterCommand("logchannel", func(connection *discordgo.Session, message *discordgo.Message, args []string) {
		if isDM(connection, message) {
			return
		}
		(&setLogChannel{client: client, connection: connection, message: message, args: args}).execute()
	})

	// Warn.
	client.RegisterCommand("warn", func(connection *discordgo.Session, message *discordgo.Message, args []string) {
		if isDM(connection, message) {
			return
		}
		(&warn{client: client, connection: connection, message: message, args: args,
			menuData: &core.MenuCommandData{}, modData: &moderationData{}}).execute()
	})
}

func isDM(connection *discordgo.Session, message *discordgo.Message) bool {
	channel, err := connection.Channel(message.ChannelID)
	if err != nil {
		return true
	}
	if channel.Type != discordgo.ChannelTypeGuildText {
		return true
	}
	return false
}
