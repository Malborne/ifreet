package commands

import (
	"fmt"

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

//tajweedlesson gives information about the weekly Tajweed lesson.
func tajweedlesson(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:       "Tajweed Lesson",
		Description: fmt.Sprintf("This server has a weekly Tajweed Lesson on Saturdays, where we explain the rules of Tajweed. There are two sessions to accomodate different time zones. The first session is at 4:00 PM UTC/GMT and the second session is 2:00 AM UTC/GMT.\nYou can watch all the previous lessons in <#%s>.\nIf you want to see all the documents for all of the previous lessons, just type `;lessons`.", "679065435429142529"),
		Color:       0x00AA00,
	})
	return errors.Wrap(err, "sending embed failed")
}
