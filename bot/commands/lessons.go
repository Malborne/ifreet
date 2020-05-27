package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var lessonsCommand = command{
	"lessons",
	commandLessons,
	"Gives a link to the list of lessons.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//commandLessons gives a link to the list of lessons
func commandLessons(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://docs.google.com/document/d/1DTjcpMeRKse91rC7Tut5WwQspqGI3_OOchGuYEveEjc/edit?usp=sharing")
	return errors.Wrap(err, "sending message failed")
}
