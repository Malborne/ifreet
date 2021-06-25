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

	var ModPermissions int64 = discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionManageMessages | discordgo.PermissionMentionEveryone | discordgo.PermissionCreateInstantInvite
	var UserPermissions int64 = discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages

	// newChannel, err := s.GuildChannelCreate(g.GuildID, g.User.Username, discordgo.ChannelTypeGuildText)
	var permissionObjects = make([]*discordgo.PermissionOverwrite, 5)

	permissionObjects[0] = &discordgo.PermissionOverwrite{ID: Config.UserRole, Type: discordgo.PermissionOverwriteTypeRole, Deny: 0x1111111111}
	permissionObjects[1] = &discordgo.PermissionOverwrite{ID: Config.FemaleOnlyRole, Type: discordgo.PermissionOverwriteTypeRole, Deny: 0x1111111111}

	permissionObjects[2] = &discordgo.PermissionOverwrite{ID: g.GuildID, Type: discordgo.PermissionOverwriteTypeRole, Deny: 0x1111111111}
	permissionObjects[3] = &discordgo.PermissionOverwrite{ID: g.User.Username, Type: discordgo.PermissionOverwriteTypeMember, Allow: UserPermissions}
	permissionObjects[4] = &discordgo.PermissionOverwrite{ID: Config.ModRole, Type: discordgo.PermissionOverwriteTypeRole, Allow: ModPermissions}

	// s.ChannelMessageSend(Config.AdminChannel, fmt.Sprintf("There are %d permission objtects", len(permissionObjects)))
	data := discordgo.GuildChannelCreateData{Name: g.User.Username, Type: discordgo.ChannelTypeGuildText, Position: 4, PermissionOverwrites: permissionObjects, ParentID: "715788591766437898", NSFW: false}
	newChannel, err := s.GuildChannelCreateComplex(g.GuildID, data)
	if err != nil {
		LogIfError(s, errors.Wrap(err, "Creating New Channel failed"))

	}
	err = s.ChannelPermissionSet(newChannel.ID, Config.UserRole, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionViewChannel)
	err = s.ChannelPermissionSet(newChannel.ID, "678795606906634281", discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionViewChannel|discordgo.PermissionReadMessageHistory)
	// err = s.ChannelPermissionSet(newChannel.ID, "678795606906634281", discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionReadMessageHistory)

	// err = s.ChannelPermissionSet(newChannel.ID, g.User.ID, discordgo.PermissionOverwriteTypeMember, 0x0000000400, 0)

	// deniedPermissions := []int{0x0000000400, 0x0000000800}

	allowedUserPermissions := []int64{discordgo.PermissionViewChannel, discordgo.PermissionSendMessages, discordgo.PermissionReadMessageHistory}

	// ModPermissions := []int{0x1, 0x400, 0x800}
	// DenyPermissions(s, newChannel.ID, Config.FemaleOnlyRole, deniedPermissions)
	newChannel.PermissionOverwrites = append(newChannel.PermissionOverwrites, &discordgo.PermissionOverwrite{ID: g.User.ID, Type: discordgo.PermissionOverwriteTypeMember, Allow: discordgo.PermissionViewChannel})

	allowPermissions(s, newChannel, g.User.ID, discordgo.PermissionOverwriteTypeMember, allowedUserPermissions)
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
		// s.ChannelMessageSend(Config.AdminChannel, fmt.Sprintf("Denied Permission: %d", permissionObjects[i].Deny))

		// err := s.ChannelPermissionSet(channelID, roleID, discordgo.PermissionOverwriteTypeRole, 0, perm)
		// if err != nil {
		// 	LogIfError(s, errors.Wrap(err, "Changing permissions failed"))

		// }
		// time.Sleep(350 * time.Millisecond)

	}
	return permissionObjects
}

func allowPermissions(s *discordgo.Session, newChannel *discordgo.Channel, userID string, targetType discordgo.PermissionOverwriteType, permissions []int64) {
	// var permissionObjects = make([]*discordgo.PermissionOverwrite, len(permissions))

	for _, perm := range permissions {
		// allowed := discordgo.PermissionOverwrite{ID: userID, Type: discordgo.PermissionOverwriteTypeMember, Allow: perm, Deny: 0}
		newChannel.PermissionOverwrites = append(newChannel.PermissionOverwrites, &discordgo.PermissionOverwrite{ID: userID, Type: discordgo.PermissionOverwriteTypeMember, Allow: perm, Deny: 0})

		// permissionObjects[i] = &denied
		// err := s.ChannelPermissionSet(newChannel.ID, userID, targetType, perm, 0)
		// if err != nil {
		// 	LogIfError(s, errors.Wrap(err, "Changing permissions failed"))

		// }
	}

}
