package heimdallr

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//UserJoinHandler handles new users joining the server, and will welcome them.
func NewMemberJoinHandler(s *discordgo.Session, g *discordgo.GuildMemberAdd) {

	var DeniedPermissions int64 = discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionManageMessages | discordgo.PermissionMentionEveryone | discordgo.PermissionCreateInstantInvite | discordgo.PermissionAttachFiles | discordgo.PermissionEmbedLinks | discordgo.PermissionUseExternalEmojis

	var ModPermissions int64 = discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionManageMessages | discordgo.PermissionMentionEveryone | discordgo.PermissionCreateInstantInvite
	var UserPermissions int = discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages

	var permissionObjects = make([]*discordgo.PermissionOverwrite, 5)

	permissionObjects[0] = &discordgo.PermissionOverwrite{ID: Config.UserRole, Type: discordgo.PermissionOverwriteTypeRole, Deny: DeniedPermissions}
	permissionObjects[1] = &discordgo.PermissionOverwrite{ID: Config.FemaleOnlyRole, Type: discordgo.PermissionOverwriteTypeRole, Deny: DeniedPermissions}
	permissionObjects[2] = &discordgo.PermissionOverwrite{ID: g.GuildID, Type: discordgo.PermissionOverwriteTypeRole, Deny: DeniedPermissions}

	permissionObjects[3] = &discordgo.PermissionOverwrite{ID: Config.ModRole, Type: discordgo.PermissionOverwriteTypeRole, Allow: ModPermissions}
	permissionObjects[4] = &discordgo.PermissionOverwrite{ID: Config.TrialModRole, Type: discordgo.PermissionOverwriteTypeRole, Allow: ModPermissions}

	data := discordgo.GuildChannelCreateData{Name: g.User.ID, Type: discordgo.ChannelTypeGuildText, Position: 4, PermissionOverwrites: permissionObjects, ParentID: "715788591766437898", NSFW: false}
	newChannel, err := s.GuildChannelCreateComplex(g.GuildID, data)
	if err != nil {
		LogIfError(s, errors.Wrap(err, "Creating New Channel failed"))

	}

	err = s.ChannelPermissionSet(newChannel.ID, g.User.ID, discordgo.PermissionOverwriteTypeMember, UserPermissions, 0)

	welcomeMessage := Config.WelcomeMessage
	if strings.Count(welcomeMessage, "%s") > 0 {
		welcomeMessage = fmt.Sprintf(welcomeMessage, g.User.Mention(), Config.RulesChannel)
	}

	_, err = s.ChannelMessageSend(newChannel.ID, welcomeMessage)
	LogIfError(s, errors.Wrap(err, "sending message failed"))
}
