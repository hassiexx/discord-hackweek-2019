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

	// Registered menu commands.
	// The ID for this will be ChannelUD-UserID to identify in progress commands accurately.
	// We only want one instance of a command running per user per channel.
	menuCommands map[string]MenuCommand
}

// NewClient creates a new instance of the Discord client.
func NewClient(token string) (*DiscordClient, error) {
	// Create connection.
	connection, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	// Create client.
	return &DiscordClient{
		commands:     make(map[string]Command),
		connection:   connection,
		menuCommands: make(map[string]MenuCommand),
	}, nil
}

// Connect opens the connection to the Discord gateway.
func (c *DiscordClient) Connect() error {
	return c.connection.Open()
}

// Command gets the specified command.
func (c *DiscordClient) Command(name string) Command {
	return c.commands[name]
}

// MenuCommand gets a menu command with the specified ID.
func (c *DiscordClient) MenuCommand(id string) MenuCommand {
	return c.menuCommands[id]
}

// Connection gets the Discord connection.
func (c *DiscordClient) Connection() *discordgo.Session {
	return c.connection
}

// HasCommand returns whether the specified command is registered.
func (c *DiscordClient) HasCommand(name string) bool {
	_, exists := c.commands[name]
	return exists
}

// HasMenuCommand returns whether a menu command is registered.
func (c *DiscordClient) HasMenuCommand(id string) bool {
	_, exists := c.menuCommands[id]
	return exists
}

// RegisterCommand registers a bot command.
func (c *DiscordClient) RegisterCommand(name string, command Command) {
	c.commands[name] = command
}

// RegisterMenuCommand registeres a menu command.
func (c *DiscordClient) RegisterMenuCommand(id string, command MenuCommand) {
	c.menuCommands[id] = command
}

// UnregisterMenuCommand unregisters a menu command.
func (c *DiscordClient) UnregisterMenuCommand(id string) {
	delete(c.menuCommands, id)
}
