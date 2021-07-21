package heimdallr

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//UserJoinHandler handles new users joining the server, and will welcome them.
func NewMemberJoinHandler(s *discordgo.Session, g *discordgo.GuildMemberAdd) {

	var DeniedPermissions int64 = discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionManageMessages | discordgo.PermissionMentionEveryone | discordgo.PermissionCreateInstantInvite | discordgo.PermissionAttachFiles | discordgo.PermissionEmbedLinks | discordgo.PermissionUseExternalEmojis | discordgo.PermissionSendTTSMessages | 0x0080000000

	var ModPermissions int64 = discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages | discordgo.PermissionAddReactions | discordgo.PermissionManageMessages | discordgo.PermissionMentionEveryone | discordgo.PermissionCreateInstantInvite
	var UserPermissions int = discordgo.PermissionViewChannel | discordgo.PermissionReadMessageHistory | discordgo.PermissionSendMessages
	var DeniedUserPermissions int = int(DeniedPermissions) | UserPermissions
	var permissionObjects = make([]*discordgo.PermissionOverwrite, 5)

	permissionObjects[0] = &discordgo.PermissionOverwrite{ID: Config.UserRole, Type: discordgo.PermissionOverwriteTypeRole, Deny: DeniedPermissions}
	permissionObjects[1] = &discordgo.PermissionOverwrite{ID: Config.FemaleOnlyRole, Type: discordgo.PermissionOverwriteTypeRole, Deny: DeniedPermissions}
	permissionObjects[2] = &discordgo.PermissionOverwrite{ID: g.GuildID, Type: discordgo.PermissionOverwriteTypeRole, Deny: DeniedPermissions}

	permissionObjects[3] = &discordgo.PermissionOverwrite{ID: Config.ModRole, Type: discordgo.PermissionOverwriteTypeRole, Allow: ModPermissions}
	permissionObjects[4] = &discordgo.PermissionOverwrite{ID: Config.TrialModRole, Type: discordgo.PermissionOverwriteTypeRole, Allow: ModPermissions}

	data := discordgo.GuildChannelCreateData{Name: "welcome-" + g.User.Username, Type: discordgo.ChannelTypeGuildText, Position: 4, PermissionOverwrites: permissionObjects, ParentID: "715788591766437898", NSFW: false}
	newChannel, err := s.GuildChannelCreateComplex(g.GuildID, data)
	if err != nil {
		LogIfError(s, errors.Wrap(err, "Creating New Channel failed"))

	}

	err = s.ChannelPermissionSet(newChannel.ID, g.User.ID, discordgo.PermissionOverwriteTypeMember, UserPermissions, DeniedUserPermissions)

	//Prevent new members from seeing the old welcome channel
	err = s.ChannelPermissionSet(Config.WelcomeChannel, g.User.ID, discordgo.PermissionOverwriteTypeMember, 0, int(DeniedPermissions))

	//Add new channel to database
	err = AddNewChannel(g.User.ID, newChannel.ID)
	LogIfError(s, errors.Wrap(err, "failed to add the new channel info to the database."))

	//Welcome message
	welcomeMessage := Config.WelcomeMessage
	if strings.Count(welcomeMessage, "%s") > 0 {
		welcomeMessage = fmt.Sprintf(welcomeMessage, g.User.Mention(), Config.RulesChannel)
	}

	_, err = s.ChannelMessageSend(newChannel.ID, welcomeMessage)
	LogIfError(s, errors.Wrap(err, "sending message failed"))

	//Send a message to remind the user that they are still unapproved after some time
	time.AfterFunc(144*time.Hour, func() { sendUnapprovedMessage(s, newChannel, g.User) }) //called after 6 days

	//kick the user they are still unapproved after some time
	member, _ := GetMember(s, g.GuildID, g.User.ID)
	time.AfterFunc(168*time.Hour, func() { kickMember(s , member ) }) //called after 7 days
}

//NewMemberLeaveHandler wishes ex members goodbye and deletes the channel that was created for them
func NewMemberLeaveHandler(s *discordgo.Session, g *discordgo.GuildMemberRemove) {

	userChannelID, err := GetnewChannel(g.User.ID)
	if userChannelID != "" {
		_, err = s.ChannelDelete(userChannelID)
		LogIfError(s, errors.Wrap(err, "unable to delete the channel"))
		err = RemoveNewChannel(g.User.ID)
		LogIfError(s, errors.Wrap(err, "unable to remove the channel from the database"))
	}
	var name string
	if g.Nick != "" {
		name = g.Nick
	} else {
		name = g.User.Username
	}
	_, err = s.ChannelMessageSend(Config.LogChannel, fmt.Sprintf("User `%s` (%s) has left the server.", name, g.User.Mention()))
	LogIfError(s, errors.Wrap(err, "sending message failed"))

}

func sendUnapprovedMessage(s *discordgo.Session, newChannel *discordgo.Channel, user *discordgo.User) {
	_, err := s.ChannelMessageSend(newChannel.ID, fmt.Sprintf(
		"%s You are an unapproved member of Quran Learning Center Server and you do not have access to most of the server. If you would like to have access to the server, please contact one of the moderators using `@Moderator` to be approved.\n\nKeep in mind that if you stay for one more day without getting approved, you will risk being kicked out of the server.", user.Mention()))
	if err != nil {
		LogIfError(s, errors.Wrap(err, "sending message failed"))
	}
}
func kickMember(s *discordgo.Session, member *discordgo.Member) {
	
	if !isUserApproved(member) && !member.User.Bot && !hasRole(member, Config.ServerBoosterRole) {
		err := s.GuildMemberDeleteWithReason(member.GuildID, member.User.ID, "Stayed in the server for at least 7 days without gaining the User role")
		if err != nil {
			LogIfError(s, errors.Wrap(err, "kicking failed"))
		}
	}
}

func isUserApproved(m *discordgo.Member) bool {
	for _, role := range m.Roles {
		if role == Config.UserRole || role == Config.FemaleOnlyRole {
			return true
		}
	}
	return false
}
