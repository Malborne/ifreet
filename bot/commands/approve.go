package commands

import (
	"fmt"
	"strings"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var approveCommand = command{
	"approve",
	commandApprove,
	"Gives the user full access to the server.",
	[]string{
		"<Member>",
	},
	[]string{
		"@username",
		"295207597929480192",
	},
}

//commandApprove gives a member the User role.
func commandApprove(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<Member>"].(string)) //Changed from user to Member

	guildID := m.GuildID
	member, err := heimdallr.GetMember(s, guildID, userID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No member was found with ID %s.", userID))
		return errors.Wrap(err, "sending message failed")
	}
	if isApproved(member) {
		return nil
	}
	user := member.User
	err = s.GuildMemberRoleAdd(guildID, userID, heimdallr.Config.UserRole)
	if err != nil {
		return errors.Wrap(err, "adding user role failed")
	}
	approvalMessage := heimdallr.Config.ApprovalMessage
	if approvalMessage != "" {
		if strings.Count(approvalMessage, "%s") > 0 {
			approvalMessage = fmt.Sprintf(approvalMessage, user.Mention(), heimdallr.Config.BotChannel)

		}
		_, err := s.ChannelMessageSend(m.ChannelID, approvalMessage)
		return errors.Wrap(err, "sending message failed")
	}
	return nil
}

//ReactionApprove approves a person if a mod reacts to their message with a green checkmark in the welcome channel
func ReactionApprove(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.ChannelID != heimdallr.Config.WelcomeChannel {
		return
	}

	if m.Emoji.Name != "✅" {
		return
	}

	reactingMember, err := heimdallr.GetMember(s, m.GuildID, m.UserID)
	if err != nil {
		heimdallr.LogIfError(s, err)
		return
	}
	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		heimdallr.LogIfError(s, err)
		return
	}
	if !heimdallr.IsModOrHigher(reactingMember, guild) {
		return
	}

	message, err := heimdallr.GetMessage(s, m.ChannelID, m.MessageID)
	if err != nil {
		heimdallr.LogIfError(s, err)
		return
	}
	member, err := heimdallr.GetMember(s, m.GuildID, message.Author.ID)
	if err != nil {
		heimdallr.LogIfError(s, err)
		return
	}
	if isApproved(member) {
		return
	}
	err = s.GuildMemberRoleAdd(m.GuildID, message.Author.ID, heimdallr.Config.UserRole)
	if err != nil {
		heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
		return
	}

	if isApproved(member) && strings.Contains(strings.ToLower(message.Content), "male") { // true
		err = s.GuildMemberRoleAdd(m.GuildID, message.Author.ID, heimdallr.Config.MaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
			return
		}
	} else if isApproved(member) && strings.Contains(strings.ToLower(message.Content), "female") {
		err = s.GuildMemberRoleAdd(m.GuildID, message.Author.ID, heimdallr.Config.FemaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
			return
		}
	}
	approvalMessage := heimdallr.Config.ApprovalMessage
	if approvalMessage != "" {
		if strings.Count(approvalMessage, "%s") > 0 {
			approvalMessage = fmt.Sprintf(approvalMessage, message.Author.Mention())
		}
		_, err := s.ChannelMessageSend(m.ChannelID, approvalMessage)
		heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))
	}
}

func isApproved(m *discordgo.Member) bool {
	for _, role := range m.Roles {
		if role == heimdallr.Config.UserRole {
			return true
		}
	}
	return false
}

func isVerified(m *discordgo.Member) bool {
	for _, role := range m.Roles {
		if role == heimdallr.Config.VerifiedRole {
			return true
		}
	}
	return false
}
