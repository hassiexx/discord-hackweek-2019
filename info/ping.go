package info

import (
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/hassieswift621/discord-hackweek-2019/core"
)

type ping struct {
	client     *core.DiscordClient
	connection *discordgo.Session
	message    *discordgo.Message
	args       []string
}

func (c *ping) execute() {
	// Send pong message.
	pongMsg, err := c.connection.ChannelMessageSend(c.message.ChannelID, ":ping_pong: | Pong!")
	if err != nil {
		return
	}

	// Store times of the messages.
	pingMsgTime, _ := c.message.Timestamp.Parse()
	pongMsgTime, _ := pongMsg.Timestamp.Parse()

	// Calculate latency in ms.
	latency := pongMsgTime.Sub(pingMsgTime).Seconds() * 1000

	// Edit message.
	_, _ = c.connection.ChannelMessageEdit(pongMsg.ChannelID, pongMsg.ID,
		":ping_pong: | Pong in **"+strconv.FormatFloat(latency, 'f', 0, 64)+"ms**")
}
