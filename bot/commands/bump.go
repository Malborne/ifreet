package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var bumpCommand = command{
	"bump",
	commandBump,
	"Bumps the server on Disboard.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//commandLessons gives a link to the list of lessons
func commandBump(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "!d bump")
	return errors.Wrap(err, "sending message failed")
}
