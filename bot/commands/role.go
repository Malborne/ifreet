package commands

import (
	"fmt"
	"strings"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var roleCommand = command{
	"role",
	commandRole,
	"Manage self-assignable roles.",
	[]string{
		"list",
		"get <role-name>...",
		"remove <role-name>...",
	},
	[]string{
		"list",
		"get sundaylesson",
		"get sundaylesson chitchat",
		"remove sundaylesson",
		"remove sundaylesson chitchat",
	},
}

//commandRole lets a user assign themselves a role.
func commandRole(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		return err
	}
	member, err := heimdallr.GetMember(s, m.GuildID, m.Author.ID)
	if err != nil {
		return err
	}
	if m.ChannelID != heimdallr.Config.BotChannel && !heimdallr.IsModOrHigher(member, guild) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Please use <#%s>", heimdallr.Config.BotChannel))
		return errors.Wrap(err, "sending message failed")
	}

	if list, _ := args.Bool("list"); list {
		return roleList(s, m)
	} else if get, _ := args.Bool("get"); get {
		return roleGet(s, m, args["<role-name>"].([]string))
	} else {
		return roleRemove(s, m, args["<role-name>"].([]string))
	}
}

func roleList(s *discordgo.Session, m *discordgo.MessageCreate) error {
	roleString := ""
	for _, role := range heimdallr.Config.Roles {
		roleString += fmt.Sprintf("**%s**: *%s*\n", role.Name, role.Desc)
	}
	if roleString == "" {
		roleString = "No roles available."
	}
	_, err := s.ChannelMessageSend(m.ChannelID, roleString)
	return errors.Wrap(err, "sending message failed")
}

func roleGet(s *discordgo.Session, m *discordgo.MessageCreate, roleNames []string) error {
	guildID := m.GuildID
	_, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		return err
	}
	success := true
	for _, roleName := range roleNames {
		var roleID string
		for _, role := range heimdallr.Config.Roles {
			if strings.ToLower(role.Name) == strings.ToLower(roleName) {
				roleID = role.ID
				break
			}
		}

		if roleID == "" {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No self-assignable role found with name `%s`.", roleName))
			heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))
			success = false
			continue
		}

		err = s.GuildMemberRoleAdd(guildID, m.Author.ID, roleID)

		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "adding role failed"))
			_, err := s.ChannelMessageSend(m.ChannelID, "Sorry, something went wrong.")
			heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))

			success = false
			continue
		}
	}
	if !success {
		return nil
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}

func roleRemove(s *discordgo.Session, m *discordgo.MessageCreate, roleNames []string) error {
	guildID := m.GuildID
	_, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		return err
	}
	success := true
	for _, roleName := range roleNames {
		var roleID string
		for _, role := range heimdallr.Config.Roles {
			if strings.ToLower(role.Name) == strings.ToLower(roleName) {
				roleID = role.ID
				break
			}
		}

		if roleID == "" {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No self-assignable role found with name `%s`.", roleName))
			heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))
			success = false
			continue
		}

		err = s.GuildMemberRoleRemove(guildID, m.Author.ID, roleID)

		if err != nil {
			heimdallr.LogIfError(s, errors.Wrap(err, "removing role failed"))
			_, err := s.ChannelMessageSend(m.ChannelID, "Sorry, something went wrong.")
			heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))
			success = false
			continue
		}
	}
	if !success {
		return nil
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}
