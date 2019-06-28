package utility

import "github.com/bwmarrin/discordgo"

// HasPermission is a helper function to check whether a user has the required permissions.
func HasPermission(connection *discordgo.Session, userID string, channelID string, permissions ...int) (bool, error) {
	// Get user permissions.
	userPerms, err := connection.State.UserChannelPermissions(userID, channelID)
	if err != nil {
		return false, err
	}

	// Run through the permissions set and check if the user has the permission.
	for _, permission := range permissions {
		// Perform bitwise ops as per https://discordapp.com/developers/docs/topics/permissions
		// If the user does not have the permission, return false instantly.
		if userPerms&permission != permission {
			return false, nil
		}
	}

	return true, nil
}

// HasPermissionOr is a helper function to check whether a user has a required permission.
// This performs the same checks as HasPermission but is an OR function rather than AND
func HasPermissionOr(connection *discordgo.Session, userID string, channelID string, permissions ...int) (bool, error) {
	// Get user permissions.
	userPerms, err := connection.State.UserChannelPermissions(userID, channelID)
	if err != nil {
		return false, err
	}

	// Run through the permissions set and check if the user has a permission.
	for _, permission := range permissions {
		// Perform bitwise ops as per https://discordapp.com/developers/docs/topics/permissions
		// If the user has a permission, return true instantly.
		if userPerms&permission == permission {
			return true, nil
		}
	}

	return false, nil
}
