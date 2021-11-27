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
		Description: fmt.Sprintf("This server has a weekly Arabic Lesson on Wednesdays & Saturdays at 2:30PM EDT (new york time)/7:30PM UK time.  in Lessons & Review VC, where we explain the grammar and learn some new vocabulary using **Al`arabiyyah Bayna Yadayk** textbook.\nYou can refer back to <#%s>.\nIf you want to see all the documents for all of the previous lessons, just type `;arabicdocs`.", "729392922444955720"),
		Color:       0x00AA00,
	})
	return errors.Wrap(err, "sending embed failed")
}
