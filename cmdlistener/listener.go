package cmdlistener

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/hassieswift621/discord-hackweek-2019/core"
)

// Initialise performs the required initialisation for the command listener.
func Initialise(client *core.DiscordClient) {
	client.Connection().AddHandler(func(connection *discordgo.Session, event *discordgo.MessageCreate) {
		// Get message.
		message := event.Message

		// If message is from our bot, ignore.
		if message.Author.ID == connection.State.User.ID {
			return
		}

		// If message is from bot, ignore.
		if message.Author.Bot {
			return
		}

		// If there is a menu command in progress, execute that.
		if client.HasMenuCommand(message.ChannelID + "-" + message.Author.ID) {
			client.MenuCommand(message.ChannelID+"-"+message.Author.ID)(message, strings.Split(message.Content, " "))
			return
		}

		// If the message does not start with the bot prefix, ignore.
		// Prefix is m.
		if !strings.HasPrefix(message.Content, "m.") {
			return
		}

		// Remove the bot prefix and split message into args.
		args := strings.Split(strings.Replace(message.Content, "m.", "", 1), " ")

		// If the command exists, execute.
		if client.HasCommand(args[0]) {
			client.Command(args[0])(connection, message, args[1:])
		}
	})
}
