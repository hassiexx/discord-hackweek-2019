package moderation

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hassieswift621/discord-hackweek-2019/core"
	"github.com/hassieswift621/discord-hackweek-2019/utility"
)

type ban struct {
	client     *core.DiscordClient
	connection *discordgo.Session
	message    *discordgo.Message
	args       []string
	menuData   *core.MenuCommandData
	modData    *moderationData
}

func (c *ban) execute() {
	// Get permissions for the user and bot.
	userHasPerms, err1 := utility.HasPermission(c.connection, c.message.Author.ID, c.message.ChannelID, discordgo.PermissionBanMembers)
	botHasPerms, err2 := utility.HasPermission(c.connection, c.message.Author.ID, c.message.ChannelID, discordgo.PermissionBanMembers)
	if err1 != nil || err2 != nil {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":x: | An internal error occurred")
		return
	}

	// If the user does not have permission, send message and bail out.
	if !userHasPerms {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | You require BAN permissions to perform this action")
		return
	}

	// If the bot does not have permission, send message and bail out.
	if !botHasPerms {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | The bot requires BAN permissions to perform this action")
		return
	}

	// If there are no mentions for banning users, bail out.
	if len(c.message.Mentions) == 0 {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | You need to mention at least one user")
		return
	}

	// Store mentions.
	c.modData.Mentions = c.message.Mentions

	// Check if the users to ban have a higher role than the bot.
	// We'll also check for self mentions at this point too.
	for i := 0; i < len(c.modData.Mentions); i++ {
		// Store user.
		user := c.modData.Mentions[i]

		if user.ID == c.message.Author.ID {
			_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | You cannot moderate yourself")
			return
		} else {
			botHigher, err := utility.IsBotHigher(c.connection, c.message.GuildID, user.ID)
			if err != nil {
				_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":x: | An internal error occurred")
				return
			}
			if !botHigher {
				_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | Cannot moderate "+
					user.Username+"#"+user.Discriminator+" because they have a higher role than me")
				return
			}
		}
	}

	// Get log channel.
	logCh, err := logChannel(c.message.GuildID)
	if err != nil {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":x: | An internal error occurred")
		return
	}

	// Check if the channel is set.
	if logCh == "" {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | The log channel is not set")
		return
	}

	// Check if the channel exists.
	exists, err := utility.ChannelExists(c.connection, c.message.GuildID, logCh)
	if err != nil {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":x: | An internal error occurred")
		return
	}
	if !exists {
		_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":exclamation: | The log channel does not exist, please re-set")
		return
	}

	// Store log channel ID.
	c.modData.LogChannelID = logCh

	// All good at this point, prepare command menu for moderator to input reason and any notes.

	// The menu title will be the list of users being banned as well as a message to cancel the ban.
	var menuTitle string = "Moderation menu - Ban\n------------------------------\n" +
		"Users to ban: " + c.modData.Mentions[0].Username + "#" + c.modData.Mentions[0].Discriminator
	for i := 1; i < len(c.modData.Mentions); i++ {
		menuTitle = menuTitle + ", " + c.modData.Mentions[i].Username + "#" + c.modData.Mentions[i].Discriminator
	}
	menuTitle = menuTitle + "\nEnter < cancel > at any time to cancel\n"
	c.menuData.Log = append(c.menuData.Log, menuTitle)

	// Register menu command.
	c.client.RegisterMenuCommand(c.message.ChannelID+"-"+c.message.Author.ID, func(message *discordgo.Message, args []string) {
		// Store new message and args.
		c.message = message
		c.args = args

		// Handle reason.
		c.handleReason()
	})

	// Send menu message.
	c.updateMenu("* Enter reason (this is also shown in the server's audit log)")
}

// HandleCancel handles the cancellation of the menu.
// It returns true if the command was cancelled.
func (c *ban) handleCancel() bool {
	if strings.ToLower(c.message.Content) == "cancel" {
		// Unregister menu command.
		c.client.UnregisterMenuCommand(c.message.ChannelID + "-" + c.message.Author.ID)

		// Send message.
		c.updateMenu("\n> Ban cancelled")

		return true
	}

	return false
}

// HandleReason handles the input for the reason.
func (c *ban) handleReason() {
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
	c.updateMenu("* Enter any additional notes or enter < none >")
}

// HandleNotes handles the input for the notes.
func (c *ban) handleNotes() {
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
func (c *ban) handleConfirmation() {
	// Handle cancel and bail out if cancelled.
	if c.handleCancel() {
		return
	}

	// If input does not equal 1 or 2, notify and bail out.
	if c.args[0] != "1" && c.args[0] != "2" {
		// Update menu.
		c.updateMenu("\nInvalid input\nEnter < 1 > to confirm\nEnter < 2 > to cancel")
		return
	}

	// If input is 2, cancel and bail out.
	if c.args[0] == "2" {
		// Unregister menu command.
		c.client.UnregisterMenuCommand(c.message.ChannelID + "-" + c.message.Author.ID)

		// Update menu.
		c.updateMenu("\n> Ban cancelled")

		return
	}

	// Input is 1.
	// For each user, ban, create mod log, send mod log and store in database.

	// Store ban count.
	var banCount int

	for i := 0; i < len(c.modData.Mentions); i++ {
		// Store user.
		user := c.modData.Mentions[i]

		// Store time of ban.
		timestamp := time.Now()

		// Ban user.
		err := c.connection.GuildBanCreateWithReason(c.message.GuildID, user.ID, c.modData.Reason, 0)
		if err != nil {
			// Send error message and go to next loop iteration.
			_, _ = c.connection.ChannelMessageSend(c.message.ChannelID, ":x: | Failed to ban "+user.Username+"#"+user.Discriminator)
			continue
		}

		// Send embed with log.
		embed := &discordgo.MessageEmbed{
			Color: int(actionColourBan),
			Author: &discordgo.MessageEmbedAuthor{
				Name:    c.message.Author.Username + "#" + c.message.Author.Discriminator,
				IconURL: c.message.Author.AvatarURL(""),
			},
			Title:       "Ban | Case ID: #",
			Description: "User: " + user.Username + "#" + user.Discriminator + "\nID: " + user.ID,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Reason",
					Value:  c.modData.Reason,
					Inline: false,
				},
				{
					Name:   "Notes",
					Value:  c.modData.Notes,
					Inline: false,
				},
			},
			Timestamp: timestamp.Format(time.RFC3339Nano),
		}

		// Send log in log channel.
		_, _ = c.connection.ChannelMessageSendEmbed(c.modData.LogChannelID, embed)

		// Increment ban count.
		banCount++

		// Sleep for 0.2s.
		time.Sleep(200 * time.Millisecond)
	}

	// Update menu.
	if banCount == len(c.modData.Mentions) {
		c.updateMenu("\n> Ban successful")
	} else if banCount > 0 {
		c.updateMenu("\n> Ban partially successful")
	} else {
		c.updateMenu("\n> Ban failed")
	}
}

// UpdateMenu updates the menu message with the specified content.
func (c *ban) updateMenu(content string) {
	// Create message content, use md as markdown for coloured formatting.
	message := "```md\n" + strings.Join(c.menuData.Log, "\n")
	if content != "" {
		message = message + "\n"
	}
	message = message + content + "```"

	// If menu message is not nil, first delete the menu.
	if c.menuData.Message != nil {
		_ = c.connection.ChannelMessageDelete(c.menuData.Message.ChannelID, c.menuData.Message.ID)
		c.menuData.Message = nil
	}

	// Send new menu.
	msg, _ := c.connection.ChannelMessageSend(c.message.ChannelID, message)
	if msg != nil {
		// Store new message.
		c.menuData.Message = msg
	}
}
