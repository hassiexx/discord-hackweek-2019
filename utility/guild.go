package utility

import "github.com/bwmarrin/discordgo"

// ChannelExists checks whether a channel exists in a guild using the channel ID.
func ChannelExists(connection *discordgo.Session, guildID string, channelID string) (bool, error) {
	// Get guild channels.
	channels, err := connection.GuildChannels(guildID)
	if err != nil {
		return false, err
	}

	// Run through channels and find the channel.
	for _, channel := range channels {
		if channel.ID == channelID {
			return true, nil
		}
	}

	return false, nil
}
