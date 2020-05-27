package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var sundayCallCommand = command{
	"sundaycall",
	commandSundaycall,
	"Displays information about the Sunday calls.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//commandSundaycall gives information about the weekly Sunday calls.
func commandSundaycall(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Sunday calls",
		Description: "This server has a weekly call on Sundays, where we have lessons. These are at 18:00 Norwegian time.",
		Color:       0x00AA00,
	})
	return errors.Wrap(err, "sending embed failed")
}
