package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var websiteCommand = command{
	"website",
	commandWebsite,
	"Gives a link to our website.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//commandWebsite gives a link to our website
func commandWebsite(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	_, err := s.ChannelMessageSend(m.ChannelID, "https://norwegianlanguagelearning.no")
	return errors.Wrap(err, "sending message failed")
}
