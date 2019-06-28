package utility

import "github.com/bwmarrin/discordgo"

// ChannelExists checks whether a channel exists in a guild using the channel ID.
func ChannelExists(connection *discordgo.Session, guildID string, channelID string) (bool, error) {
	// Get guild.
	guild, err := connection.State.Guild(guildID)

	// If there was an error, get the guild via REST.
	if err != nil {
		guild, err = connection.Guild(guildID)
		if err != nil {
			return false, err
		}
	}

	// Run through channels and find the channel.
	for _, channel := range guild.Channels {
		if channel.ID == channelID {
			return true, nil
		}
	}

	return false, nil
}
