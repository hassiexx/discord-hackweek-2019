package core

import "github.com/bwmarrin/discordgo"

// MenuCommand is the handler for a menu command.
type MenuCommand func(message *discordgo.Message, args []string)

// MenuCommandData stores data to be used between different stages of a menu command.
type MenuCommandData struct {
	Log     []string
	Message *discordgo.Message
}
