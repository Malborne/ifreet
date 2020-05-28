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
	_, err := s.ChannelMessageSend(m.ChannelID, "https://docs.google.com/document/d/188RD6TQEleoPqugOOM68JNQvDmZPsc58o7K4u1TSRA4/edit?usp=sharing")
	return errors.Wrap(err, "sending message failed")
}
