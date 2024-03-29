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
		"<Member> <gender>",
	},
	[]string{
		"@username male",
		"295207597929480192 female",
	},
}

//commandApprove gives a member the User role.
func commandApprove(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<Member>"].(string), s) //Changed from user to Member
	gender := getIDFromMaybeMention(args["<gender>"].(string), s)

	guildID := m.GuildID
	member, err := heimdallr.GetMember(s, guildID, userID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No member was found with ID %s.", userID))
		return errors.Wrap(err, "sending message failed")
	}
	if isApproved(member) {
		return nil
	}
	// user := member.User

	if strings.Contains(strings.ToLower(gender), "female") {
		err = s.GuildMemberRoleAdd(m.GuildID, userID, heimdallr.Config.UserRole)
		if err != nil {
			return errors.Wrap(err, "adding user role failed")
		}
		err = s.GuildMemberRoleAdd(m.GuildID, userID, heimdallr.Config.FemaleRole)
		if err != nil {
			return errors.Wrap(err, "adding Female role failed")
		}

	} else if strings.Contains(strings.ToLower(gender), "male") {
		err = s.GuildMemberRoleAdd(m.GuildID, userID, heimdallr.Config.UserRole)
		if err != nil {
			return errors.Wrap(err, "adding Member role failed")
		}
		err = s.GuildMemberRoleAdd(m.GuildID, userID, heimdallr.Config.MaleRole)
		if err != nil {
			return errors.Wrap(err, "adding Male role failed")
		}

	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The gender must be either Male or Female."))
		return errors.Wrap(err, "adding gender role failed")
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		return errors.Wrap(err, "adding reaction failed")
	}

	userChannelID, err := heimdallr.GetnewChannel(userID)
	if userChannelID != "" {
		_, err = s.ChannelDelete(userChannelID)
		heimdallr.LogIfError(s, errors.Wrap(err, "unable to delete the channel"))
		err = heimdallr.RemoveNewChannel(userID)
		heimdallr.LogIfError(s, errors.Wrap(err, "unable to remove the channel from the database"))
	}

	//Send a message in the LogChannel once a user is approved
	_, err = s.ChannelMessageSend(heimdallr.Config.LogChannel, fmt.Sprintf("User `%s` (%s) has been successfully approved by `%s` ✅", member.User.Username, member.Mention(), m.Author.String()))
	heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))

	approvalMessage := heimdallr.Config.ApprovalMessage
	if approvalMessage != "" {
		if strings.Count(approvalMessage, "%s") > 0 {
			approvalMessage = fmt.Sprintf(approvalMessage, member.Mention(), heimdallr.Config.BotChannel)
		}
		// userChannel, err := s.UserChannelCreate(member.User.ID)
		// if err != nil {
		// 	s.ChannelMessageSend(heimdallr.Config.LogChannel, fmt.Sprintf("New user %s Does NOT ACCEPT DMs", member.Mention()))
		// }
		// _, err = s.ChannelMessageSend(userChannel.ID, approvalMessage)
		_, err = s.ChannelMessageSend(heimdallr.Config.BotChannel, approvalMessage)
		heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))

	}

	return nil
}

//ReactionApprove approves a person if a mod reacts to their message with a green checkmark in the welcome channel
func ReactionApprove(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
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
	if !heimdallr.IsTrialModOrHigher(reactingMember, guild) {
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

	if strings.Contains(strings.ToLower(message.Content), "female") && strings.Contains(strings.ToLower(strings.Replace(strings.ToLower(message.Content), "female", "", -1)), "male") {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("More than one gender was  found in the content of the message. Use the `;approve` command instead"))
		heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
		return
	}

	if strings.Contains(strings.ToLower(message.Content), "female") {
		err = s.GuildMemberRoleAdd(m.GuildID, message.Author.ID, heimdallr.Config.UserRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
		}
		err = s.GuildMemberRoleAdd(m.GuildID, message.Author.ID, heimdallr.Config.FemaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding female role failed"))
		}

	} else if strings.Contains(strings.ToLower(message.Content), "male") {
		err = s.GuildMemberRoleAdd(m.GuildID, message.Author.ID, heimdallr.Config.UserRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
		}
		err = s.GuildMemberRoleAdd(m.GuildID, message.Author.ID, heimdallr.Config.MaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding male role failed"))
		}

	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The gender was not found in the content of the message. Please make sure that you react to a message that contains the gender."))
		heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
		return
	}

	userChannelID, err := heimdallr.GetnewChannel(member.User.ID)
	if userChannelID != "" {
		_, err = s.ChannelDelete(userChannelID)
		heimdallr.LogIfError(s, errors.Wrap(err, "unable to delete the channel"))
		err = heimdallr.RemoveNewChannel(member.User.ID)
		heimdallr.LogIfError(s, errors.Wrap(err, "unable to remove the channel from the database"))
	}

	//Send a message in the LogChannel once a user is approved
	_, err = s.ChannelMessageSend(heimdallr.Config.LogChannel, fmt.Sprintf("User `%s` (%s) has been successfully approved by `%s` ✅", member.User.Username, member.Mention(), reactingMember.User.String()))
	heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))

	approvalMessage := heimdallr.Config.ApprovalMessage
	if approvalMessage != "" {
		if strings.Count(approvalMessage, "%s") > 0 {
			approvalMessage = fmt.Sprintf(approvalMessage, message.Author.Mention(), heimdallr.Config.BotChannel)
		}
		// userChannel, err := s.UserChannelCreate(member.User.ID)
		// if err != nil {
		// 	s.ChannelMessageSend(heimdallr.Config.LogChannel, fmt.Sprintf("New User %s Does NOT ACCEPT DMs", member.Mention()))
		// }
		// _, err = s.ChannelMessageSend(userChannel.ID, approvalMessage)

		_, err = s.ChannelMessageSend(heimdallr.Config.BotChannel, approvalMessage)
		heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))

	}
}

func isApproved(m *discordgo.Member) bool {
	for _, role := range m.Roles {
		if role == heimdallr.Config.UserRole || role == heimdallr.Config.FemaleOnlyRole {
			return true
		}
	}
	return false
}

func isMale(m *discordgo.Member) bool {
	for _, role := range m.Roles {
		if role == heimdallr.Config.MaleRole {
			return true
		}
	}
	return false
}

func isFemale(m *discordgo.Member) bool {
	for _, role := range m.Roles {
		if role == heimdallr.Config.FemaleRole {
			return true
		}
	}
	return false
}

// func sendPostApprovalMessage() {
// 	message := ""
// }
