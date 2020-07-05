package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var arabicLessonCommand = command{
	"arabiclesson",
	arabiclesson,
	"Displays information about the Arabic Lesson.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//arabiclesson gives information about the weekly Arabic lesson.
func arabiclesson(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Arabic Lesson",
		Description: fmt.Sprintf("This server has a weekly Arabic Lesson on Sundays at 3:00 PM (UTC/GMT), where we explain the grammar and learn some new vocabulary using Al Madinah Arabic book.\nYou can refer back to <#%s>.\nIf you want to see all the documents for all of the previous lessons, just type `;arabicdocs`.", "729392922444955720"),
		Color:       0x00AA00,
	})
	return errors.Wrap(err, "sending embed failed")
}
