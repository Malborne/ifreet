package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var arabiclessonsCommand = command{
	"arabicdocs",
	commandArabicLessons,
	"Gives a link to the list of Arabic lessons.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//commandArabicLessons gives a link to the list of Arabic lessons
func commandArabicLessons(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	_, err := s.ChannelMessageSend(m.ChannelID, " https://docs.google.com/document/d/1vyekZNQVl0iD13QY414P2oy2f9c2XQpeXf5veGJLyco/edit")
	return errors.Wrap(err, "sending message failed")
}
