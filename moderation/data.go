package moderation

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

// ModerationData stores the data about the moderation received from menu commands.
type moderationData struct {
	Reason string
	Notes  string
}

// ModLog is the DB structure for a moderation log.
type modLog struct {
	Action    action `bson:"action,int32"`
	GuildID   string `bson:"guild_id"`
	Notes     string `bson:"notes,string"`
	Reason    string `bson:"reason,string"`
	Timestamp string `bson:"timestamp,string"`
	UserID    string `bson:"user_id,string"`
}
