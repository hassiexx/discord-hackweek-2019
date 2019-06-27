package db

// Collection is the type for database collection name constants.
type Collection string

const (
	CollectionModLogs     Collection = "mod_logs"
	CollectionModSettings Collection = "mod_settings"
)
