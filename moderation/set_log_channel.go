package moderation

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/hassieswift621/discord-hackweek-2019/db"

	"github.com/bwmarrin/discordgo"
	"github.com/hassieswift621/discord-hackweek-2019/core"
	"github.com/hassieswift621/discord-hackweek-2019/utility"
)

type setLogChannel struct {
	client     *core.DiscordClient
	connection *discordgo.Session
	message    *discordgo.Message
	args       []string
}

func (c *setLogChannel) execute() {
	// If the user does not have manage server permissions, bail out.
	hasPerm, err := utility.HasPermission(c.connection, c.message.Author.ID, c.message.ChannelID, discordgo.PermissionManageServer)
	if err != nil {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":x: | An internal error occurred")
		return
	}
	if !hasPerm {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | You require the MANAGE SERVER permission to set the log channel")
		return
	}

	// If args are none, bail out.
	if len(c.args) == 0 {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | You need to specify the channel either by providing the ID or mentioning it")
		return
	}

	// Escape the mention and check if the channel exists.
	channelID := strings.Replace(strings.Replace(c.args[0], "<#", "", 1), ">", "", 1)

	exists, err := utility.ChannelExists(c.connection, c.message.GuildID, channelID)

	if err != nil {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":x: | An internal error occurred")
		return
	}
	if !exists {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | The channel does not exist")
		return
	}

	// Set log channel.
	err = db.UpsertOne(db.CollectionModSettings, modSettings{GuildID: c.message.GuildID},
		bson.D{{"$set", modSettings{LogChannel: channelID}}})
	if err != nil {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":x: | An internal error occurred")
		return
	}

	_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":white_check_mark: | Log channel set")
}
