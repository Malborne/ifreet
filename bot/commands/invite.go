package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	"/heimdallr/bot"
)

var inviteCommand = command{
	"invite",
	commandInvite,
	"Sends you an invite link.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//commandInvite creates an invite link to the server.
func commandInvite(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	invite := discordgo.Invite{
		MaxAge:    86400,
		MaxUses:   1,
		Temporary: false,
		Unique:    true,
	}
	createdInvite, err := s.ChannelInviteCreate(heimdallr.Config.WelcomeChannel, invite)
	if err != nil {
		return errors.Wrap(err, "creating an invite failed")
	}
	user := m.Author
	err = heimdallr.AddInvite(*user, *createdInvite)
	// Don't send the invite if we fail to log it in the database
	if err != nil {
		return err
	}
	userChannel, err := s.UserChannelCreate(user.ID)
	if err != nil {
		return errors.Wrap(err, "sending invite to user failed")
	}
	_, err = s.ChannelMessageSend(userChannel.ID, fmt.Sprintf(
		"https://discord.gg/%s", createdInvite.Code,
	))
	return errors.Wrap(err, "sending invite to user failed")
}
