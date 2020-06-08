package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var ChangeVerifiedRolesCommand = command{
	"changeverified",
	commandChangeVerifiedRoles,
	"Switches everyone from Verified roles to Verified Male or Verified Female, depending on their gender.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//Switches everyone from Verified roles to Verified Male or Verified Female, depending on their gender.
func commandChangeVerifiedRoles(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		return err
	}
	var count int = 0
	for _, member := range guild.Members {

		if isVerified(member) && member.User.ID != s.State.User.ID && isMale(member) {
			err = s.GuildMemberRoleAdd(m.GuildID, member.User.ID, heimdallr.Config.VerifiedMaleRole)
			if err != nil {
				heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
				return
			} else {
				count = count + 1
			}

			err = s.GuildMemberRoleRemove(m.GuildID, member.User.ID, heimdallr.Config.MaleRole)
			if err != nil {
				heimdallr.LogIfError(s, errors.Wrap(err, "removing user role failed"))
				return
			}

			err = s.GuildMemberRoleRemove(m.GuildID, member.User.ID, heimdallr.Config.VerifiedRole)
			if err != nil {
				heimdallr.LogIfError(s, errors.Wrap(err, "removing user role failed"))
				return
			}

		} else if isVerified(member) && member.User.ID != s.State.User.ID && isFemale(member) {
			err = s.GuildMemberRoleAdd(m.GuildID, member.User.ID, heimdallr.Config.VerifiedFemaleRole)
			if err != nil {
				heimdallr.LogIfError(s, errors.Wrap(err, "adding user role failed"))
				return
			} else {
				count = count + 1
			}
			err = s.GuildMemberRoleRemove(m.GuildID, member.User.ID, heimdallr.Config.FemaleRole)
			if err != nil {
				heimdallr.LogIfError(s, errors.Wrap(err, "removing user role failed"))
				return
			}

			err = s.GuildMemberRoleRemove(m.GuildID, member.User.ID, heimdallr.Config.VerifiedRole)
			if err != nil {
				heimdallr.LogIfError(s, errors.Wrap(err, "removing user role failed"))
				return
			}
		}
	}

	if count == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No Users with the Old Verified role found."))
		return errors.Wrap(err, "sending message failed")
	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sucessfully changed the roles of %d user(s)", count))
		return errors.Wrap(err, "sending message failed")
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}
