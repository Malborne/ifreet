package heimdallr

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//UserJoinHandler handles new users joining the server, and will welcome them.
func NewMemberJoinHandler(s *discordgo.Session, g *discordgo.GuildMemberAdd) {
	// type GuildChannelCreateData struct {
	// 	Name                 string                 `json:"name"`
	// 	Type                 ChannelType            `json:"type"`
	// 	Topic                string                 `json:"topic,omitempty"`
	// 	Bitrate              int                    `json:"bitrate,omitempty"`
	// 	UserLimit            int                    `json:"user_limit,omitempty"`
	// 	RateLimitPerUser     int                    `json:"rate_limit_per_user,omitempty"`
	// 	Position             int                    `json:"position,omitempty"`
	// 	PermissionOverwrites []*PermissionOverwrite `json:"permission_overwrites,omitempty"`
	// 	ParentID             string                 `json:"parent_id,omitempty"`
	// 	NSFW                 bool                   `json:"nsfw,omitempty"`
	// }

	// ID    string                  `json:"id"`
	// Type  PermissionOverwriteType `json:"type"`
	// Deny  int64                   `json:"deny,string"`
	// Allow int64                   `json:"allow,string"`
	// deniedPermissions := []int64{0x0000000001, 0x0000000400, 0x0000000800, 0x0000001000, 0x0000004000, 0x0000008000, 0x0000010000, 0x0000020000, 0x0000040000, 0x0000080000, 0x0080000000, 0x0800000000, 0x1000000000}
	deniedPermissions := []int64{0x0000000400, 0x0000000800}

	permissionObjects := DenyPermissions(s, Config.UserRole, deniedPermissions)
	s.ChannelMessageSend(Config.AdminChannel, fmt.Sprintf("There are %d permission objtects", len(permissionObjects)))
	data := discordgo.GuildChannelCreateData{Name: g.User.Username, Type: discordgo.ChannelTypeGuildText, Position: 4, PermissionOverwrites: permissionObjects, ParentID: "715788591766437898", NSFW: false}
	newChannel, err := s.GuildChannelCreateComplex(g.GuildID, data)
	// newChannel, err := s.GuildChannelCreate(g.GuildID, g.User.Username, discordgo.ChannelTypeGuildText)
	if err != nil {
		LogIfError(s, errors.Wrap(err, "Creating New Channel failed"))

	}

	// deniedPermissions := []int{0x0000000400, 0x0000000800}

	// allowedUserPermissions := []int{0x400, 0x800, 0x10000}

	// ModPermissions := []int{0x1, 0x400, 0x800}
	// DenyPermissions(s, newChannel.ID, Config.FemaleOnlyRole, deniedPermissions)

	// allowPermissions(s, newChannel.ID, g.User.ID, discordgo.PermissionOverwriteTypeMember, allowedUserPermissions)
	// allowPermissions(s, newChannel.ID, Config.ModRole, discordgo.PermissionOverwriteTypeRole, ModPermissions)
	// allowPermissions(s, newChannel.ID, Config.TrialModRole, ModPermissions)

	// err = s.ChannelPermissionSet(newChannel.ID, Config.UserRole, discordgo.PermissionOverwriteTypeRole, 0, 0x0000000400)
	// err = s.ChannelPermissionSet(newChannel.ID, Config.FemaleOnlyRole, discordgo.PermissionOverwriteTypeRole, 0, 1024)

	// if err != nil {
	// 	LogIfError(s, errors.Wrap(err, "Changing permissions failed"))

	// }
	welcomeMessage := Config.WelcomeMessage
	if strings.Count(welcomeMessage, "%s") > 0 {
		welcomeMessage = fmt.Sprintf(welcomeMessage, g.User.Mention(), Config.RulesChannel)
	}

	_, err = s.ChannelMessageSend(newChannel.ID, welcomeMessage)
	LogIfError(s, errors.Wrap(err, "sending message failed"))
}

func DenyPermissions(s *discordgo.Session, roleID string, permissions []int64) []*discordgo.PermissionOverwrite {
	var permissionObjects = make([]*discordgo.PermissionOverwrite, len(permissions))
	for i, perm := range permissions {
		denied := discordgo.PermissionOverwrite{ID: roleID, Type: discordgo.PermissionOverwriteTypeRole, Deny: perm}

		permissionObjects[i] = &denied
		s.ChannelMessageSend(Config.AdminChannel, fmt.Sprintf("Denied Permission: %d", permissionObjects[i].Deny))

		// err := s.ChannelPermissionSet(channelID, roleID, discordgo.PermissionOverwriteTypeRole, 0, perm)
		// if err != nil {
		// 	LogIfError(s, errors.Wrap(err, "Changing permissions failed"))

		// }
		// time.Sleep(350 * time.Millisecond)

	}
	return permissionObjects
}

// func allowPermissions(s *discordgo.Session,userID string, targetType discordgo.PermissionOverwriteType, permissions []int) {
// 	for _, perm := range permissions {

// 		err := s.ChannelPermissionSet(channelID, userID, targetType, perm, 0)
// 		if err != nil {
// 			LogIfError(s, errors.Wrap(err, "Changing permissions failed"))

// 		}
// 	}
// }
