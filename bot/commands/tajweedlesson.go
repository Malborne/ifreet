package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var tajweedLessonCommand = command{
	"tajweedlesson",
	tajweedlesson,
	"Displays information about the Tajweed Lesson.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//tajweedlesson gives information about the weekly Sunday calls.
func tajweedlesson(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Tajweed Lesson",
		Description: "This server has a weekly Tajweed Lesson on Saturdays, where we explain the rules of Tajweed. There are two sessions to accomodate different time zones. The first session is at 18:00 AST (UTC+3) and the second session is 20:00 PST (UTC-5).",
		Color:       0x00AA00,
	})
	return errors.Wrap(err, "sending embed failed")
}