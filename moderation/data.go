package moderation

import (
	"github.com/hassieswift621/discord-hackweek-2019/db"
	"go.mongodb.org/mongo-driver/mongo"
)

// Action is the type for moderation action constants.
type action int32

const (
	actionWarn action = iota + 1
	actionTempMute
	actionMute
	actionKick
	actionTempBan
	actionBan
)

// ActionColours is the type for moderator action embed colours.
type actionColour int

const (
	actionColourWarn actionColour = 0xf5e942
	actionColourKick actionColour = 0xf5a442
	actionColourBan  actionColour = 0xf54e42
)

// ModerationData stores the data about the moderation received from menu commands.
type moderationData struct {
	LogChannelID string
	Reason       string
	Notes        string
}

// ModSettings is the DB structure for moderation settings for a guild.
type modSettings struct {
	GuildID    string `bson:"guild_id,string,omitempty"`
	LogChannel string `bson:"log_channel,string,omitempty"`
}

// LogChannel gets the ID of the log channel if set, "" otherwise.
func logChannel(guildID string) (string, error) {
	var settings modSettings
	err := db.QueryOne(db.CollectionModSettings, modSettings{GuildID: guildID}, &settings)
	if err != nil && err == mongo.ErrNoDocuments {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return settings.LogChannel, nil
}
