package utility

import (
	"github.com/bwmarrin/discordgo"
)

// IsBotHigher checks whether the bot has a higher role than the user to ban.
func IsBotHigher(connection *discordgo.Session, guildID string, userID string) (bool, error) {
	// Get guild roles.
	roles, err := connection.GuildRoles(guildID)
	if err != nil {
		return false, err
	}

	// Get bot as guild member.
	bot, err := connection.GuildMember(guildID, connection.State.User.ID)
	if err != nil {
		return false, err
	}

	// Get target user as guild member.
	user, err := connection.GuildMember(guildID, userID)
	if err != nil {
		return false, err
	}

	// Vars to store the highest role number for the bot and user.
	var botHighest int
	var userHighest int

	// Run through the guild roles and iterate over the bot and user roles,
	// to find the highest role for each.
	// It has to be done this way because DiscordGo only stores the role IDs for guild members,
	// so for efficiency we only iterate over the guild roles once.
	for _, role := range roles {
		// Bot.
		for _, botRole := range bot.Roles {
			if role.ID == botRole {
				// If this role is higher than the current one, store this.
				if role.Position > botHighest {
					botHighest = role.Position
				}

				// Break this inner loop because we've found the role,
				// prevent unnecessary iterations.
				break
			}
		}

		// User.
		for _, userRole := range user.Roles {
			if role.ID == userRole {
				// If this role is higher than the current one, store this.
				if role.Position > userHighest {
					userHighest = role.Position
				}

				// Break this inner loop because we've found the role,
				// prevent unnecessary iterations.
				break
			}
		}
	}

	return botHighest > userHighest, nil
}
