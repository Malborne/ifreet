package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var verifyCommand = command{
	"verify",
	commandVerify,
	"Gives the user full access to the gender specific channels.",
	[]string{
		"<Member>",
	},
	[]string{
		"@username",
		"295207597929480192",
	},
}

//commandVerify gives a member the Verified Male/Female role.
func commandVerify(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<Member>"].(string)) //Changed from user to Member

	guildID := m.GuildID
	member, err := heimdallr.GetMember(s, guildID, userID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No member was found with ID %s.", userID))
		return errors.Wrap(err, "sending message failed")
	}
	if heimdallr.IsVerified(member) {
		return nil
	}

	if hasRole(member, heimdallr.Config.FemaleRole) {
		err = s.GuildMemberRoleAdd(m.GuildID, member.User.ID, heimdallr.Config.VerifiedFemaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
			return errors.Wrap(err, "adding role failed")
		}
		err = s.GuildMemberRoleRemove(m.GuildID, member.User.ID, heimdallr.Config.FemaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "removing user role failed"))
			return errors.Wrap(err, "Changing roles failed")
		}
	} else if hasRole(member, heimdallr.Config.MaleRole) {
		err = s.GuildMemberRoleAdd(m.GuildID, member.User.ID, heimdallr.Config.VerifiedMaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
			return errors.Wrap(err, "Changing roles failed")
		}
		err = s.GuildMemberRoleRemove(m.GuildID, member.User.ID, heimdallr.Config.MaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "removing user role failed"))
			return errors.Wrap(err, "Changing roles failed")
		}
	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s does NOT have any gender roles.", member.Mention()))
		return errors.Wrap(err, "sending message failed")
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}

//ReactionVerify verifies a person if a mod reacts to their message with a green checkmark in the welcome channel
func ReactionVerify(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if m.ChannelID != heimdallr.Config.WelcomeChannel {
		return
	}

	if m.Emoji.Name != "☑️" {
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
	if heimdallr.IsVerified(member) {
		return
	}

	if hasRole(member, heimdallr.Config.FemaleRole) {
		err = s.GuildMemberRoleAdd(m.GuildID, message.Author.ID, heimdallr.Config.VerifiedFemaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
			return
		}
		err = s.GuildMemberRoleRemove(m.GuildID, message.Author.ID, heimdallr.Config.FemaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "removing user role failed"))
			return
		}
	} else if hasRole(member, heimdallr.Config.MaleRole) {
		err = s.GuildMemberRoleAdd(m.GuildID, message.Author.ID, heimdallr.Config.VerifiedMaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
			return
		}
		err = s.GuildMemberRoleRemove(m.GuildID, message.Author.ID, heimdallr.Config.MaleRole)
		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "removing user role failed"))
			return
		}
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s does NOT have any gender roles.", member.Mention()))
		return
	}

}
func hasRole(m *discordgo.Member, r string) bool {
	for _, role := range m.Roles {
		if role == r {
			return true
		}
	}
	return false
}
