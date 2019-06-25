package core

import (
	"github.com/bwmarrin/discordgo"
)

// DiscordClient is the main interface of the bot.
type DiscordClient struct {
	// Registered commands.
	commands map[string]Command

	// Discord connection.
	connection *discordgo.Session
}

// NewClient creates a new instance of the Discord client.
func NewClient(token string) (*DiscordClient, error) {
	// Create connection.
	connection, err := discordgo.New(token)
	if err != nil {
		return nil, err
	}

	// Create client.
	return &DiscordClient{
		commands:   make(map[string]Command),
		connection: connection,
	}, nil
}

// RegisterCommand registers a bot command.
func (c *DiscordClient) RegisterCommand(name string, command Command) {
	c.commands[name] = command
}
