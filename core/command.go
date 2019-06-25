package core

import "github.com/bwmarrin/discordgo"

// Command is the handler for a bot command.
type Command func(connection *discordgo.Session, message *discordgo.Message, args []string)
