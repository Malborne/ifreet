package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var getSheetCommand = command{
	"getsheet",
	commandGetSheet,
	"gives a link to the Google docs sheet for a particular student",
	[]string{
		"<user>",
	},
	[]string{
		"@username",
		"240114490929053696",
	},
}

//commandGetSheet gives a link to the Google docs sheet for a particular student
func commandGetSheet(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string), s)

	guildID := m.GuildID
	guild, err := heimdallr.GetGuild(s, guildID)
	if err != nil {
		return err
	}
	member, err := heimdallr.GetMember(s, guildID, userID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No member was found with ID %s.", userID))
		return errors.Wrap(err, "sending message failed")
	}
	user := member.User

	author, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", userID))
		return errors.Wrap(err, "getting author failed")
	}

	if !heimdallr.IsHelper(author, guild) && (m.Author.ID != user.ID) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Only helper and the sheet owner are allowed to get links for the sheet."))
		return errors.Wrap(err, "sending message failed")
	}

	student, err := heimdallr.GetStudent(userID)
	if err != nil {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The student is not registered in the database."))
		return errors.Wrap(err, "getting the sheetLink failed")
	}
	if student.ID != "" {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(student.SheetLink))
	} else {
		_, _ = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The student is not registered in the database. Make sure you add the student first."))
		return errors.Wrap(err, "getting the sheetLink failed")
	}
	return errors.Wrap(err, "sending message failed")

}
