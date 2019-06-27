package moderation

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/hassieswift621/discord-hackweek-2019/core"
	"github.com/hassieswift621/discord-hackweek-2019/utility"
)

type kick struct {
	client     *core.DiscordClient
	connection *discordgo.Session
	message    *discordgo.Message
	args       []string
	menuData   *core.MenuCommandData
	modData    *moderationData
}

func (c *kick) execute() {
	// If there are no mentions for kicking users, bail out.
	if len(c.message.Mentions) == 0 {
		return
	}

	// Get permissions for the user and bot.
	userHasPerms, err1 := utility.HasPermission(c.connection, c.message.Author.ID, c.message.ChannelID, discordgo.PermissionKickMembers)
	botHasPerms, err2 := utility.HasPermission(c.connection, c.message.Author.ID, c.message.ChannelID, discordgo.PermissionKickMembers)
	if err1 != nil || err2 != nil {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":x: | An internal error occurred")
		return
	}

	// If the user does not have permission, send message and bail out.
	if !userHasPerms {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | You require KICK permissions to perform this action")
		return
	}

	// If the bot does not have permission, send message and bail out.
	if !botHasPerms {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | The bot requires KICK permissions to perform this action")
		return
	}

	// TODO: Check if the log channel is set.

	// All good at this point, prepare command menu for moderator to input reason and any notes.

	// The menu title will be the list of users being kicked as well as a message to cancel the kick.
	var menuTitle string = "Moderation menu - Kick\n------------------------------\n" +
		"Users being kicked: " + c.message.Mentions[0].Username + "#" + c.message.Mentions[0].Discriminator
	for i := 1; i < len(c.message.Mentions); i++ {
		menuTitle = menuTitle + ", " + c.message.Mentions[i].Username + "#" + c.message.Mentions[i].Discriminator
	}
	menuTitle = "\nEnter < cancel > at any time to cancel the kick\n"
	c.menuData.Log = append(c.menuData.Log, menuTitle)

	// Register menu command.
	c.client.RegisterMenuCommand(c.message.GuildID+"-"+c.message.Author.ID, func(message *discordgo.Message, args []string) {
		// Store new message and args.
		c.message = message
		c.args = args

		// Handle reason.
		c.handleReason()
	})

	// Send menu message.
	c.updateMenu("* Enter reason (this is also shown in the server's audit log)...")
}

// HandleCancel handles the cancellation of the menu.
// It returns true if the command was cancelled.
func (c *kick) handleCancel() bool {
	if strings.ToLower(c.message.Content) == "cancel" {
		// Unregister menu command.
		c.client.UnregisterMenuCommand(c.message.ChannelID + "-" + c.message.Author.ID)

		// Send message.
		c.updateMenu("\n> Kick cancelled")

		return true
	}

	return false
}

// HandleReason handles the input for the reason.
func (c *kick) handleReason() {
	// Handle cancel and bail out if cancelled.
	if c.handleCancel() {
		return
	}

	// Store reason.
	c.modData.Reason = strings.Join(c.args, " ")

	// Update menu log.
	c.menuData.Log = append(c.menuData.Log, "* Reason: "+c.modData.Reason)

	// Register menu command for notes.
	c.client.RegisterMenuCommand(c.message.ChannelID+"-"+c.message.Author.ID, func(message *discordgo.Message, args []string) {
		// Store new message and args.
		c.message = message
		c.args = args

		// Handle notes.
		c.handleNotes()
	})

	// Update menu.
	c.updateMenu("* Enter any additional notes or enter 'none'")
}

// HandleNotes handles the input for the notes.
func (c *kick) handleNotes() {
	// Handle cancel and bail out if cancelled.
	if c.handleCancel() {
		return
	}

	// Store notes.
	// If notes are "none", then there are no notes.
	notes := strings.Join(c.args, " ")
	if strings.ToLower(notes) == "none" {
		c.modData.Notes = "N/A"
	} else {
		c.modData.Notes = notes
	}

	// Update menu log.
	c.menuData.Log = append(c.menuData.Log, "* Notes: "+c.modData.Notes)

	// Register menu command for confirmation.
	c.client.RegisterMenuCommand(c.message.ChannelID+"-"+c.message.Author.ID, func(message *discordgo.Message, args []string) {
		// Store new message and args.
		c.message = message
		c.args = args

		// Handle confirmation.
		c.handleConfirmation()
	})

	// Update menu.
	c.updateMenu("\nEnter < 1 > to confirm\nEnter < 2 > to cancel")
}

// HandleConfirmation handles the confirmation.
func (c *kick) handleConfirmation() {
	// Handle cancel and bail out if cancelled.
	if c.handleCancel() {
		return
	}

	// If input does not equal 1 or 2, notify and bail out.
	if c.args[0] != "1" && c.args[1] != "2" {
		// Update menu.
		c.updateMenu("\nInvalid input\nEnter < 1 > to confirm\nEnter < 2 > to cancel")
		return
	}

	// If input is 2, cancel and bail out.
	if c.args[0] == "2" {
		// Unregister menu command.
		c.client.UnregisterMenuCommand(c.message.ChannelID + "-" + c.message.Author.ID)

		// Update menu.
		c.updateMenu("\n> Kick cancelled")

		return
	}

	// Input is 1.
	// For each user kick.

}

// UpdateMenu updates the menu message with the specified content.
func (c *kick) updateMenu(content string) {
	// Create message content, use md as markdown for coloured formatting.
	message := "```md" + strings.Join(c.menuData.Log, "\n") + "\n" + content + "```"

	// If menu message is not nil, first delete the menu.
	if c.menuData.Message != nil {
		_ = c.connection.ChannelMessageDelete(c.menuData.Message.ChannelID, c.menuData.Message.ChannelID)
		c.menuData.Message = nil
	}

	// Send new menu.
	msg, _ := c.connection.ChannelMessageSend(c.message.ChannelID, message)
	if msg != nil {
		// Store new message.
		c.menuData.Message = msg
	}
}
