package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var addStudentCommand = command{
	"addstudent",
	commandAddStudent,
	"adds a student to the database.",
	[]string{
		"<user> <circle> <sheetLink>",
	},
	[]string{
		"@username \"Ali Bin Abi Talib\" https://docs.google.com/document/d/188RD6TQEleoPqugOOM68JNQvDmZPsc58o7K4u1TSRA4/edit?usp=sharing",
		"240114490929053696 \"Ali Bin Abi Talib\" https://docs.google.com/document/d/188RD6TQEleoPqugOOM68JNQvDmZPsc58o7K4u1TSRA4/edit?usp=sharing",
	},
}

//commandAddStudent adds a student to the database
func commandAddStudent(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string), s)
	circle, _ := args.String("<circle>")
	sheetLink, _ := args.String("<sheetLink>")

	guildID := m.GuildID
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

	if !hasRole(author, heimdallr.Config.CricleLeaderRole) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Only Circle leaders are allowed to add students."))
		return errors.Wrap(err, "sending message failed")
	}

	err = heimdallr.AddStudent(*user, circle, sheetLink)
	if err != nil {
		return errors.Wrap(err, "adding student failed")
	}

	_, err = s.ChannelMessageSendEmbed(heimdallr.Config.LogChannel, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("User was added to %s cirlce by %s.", circle, author.User.Username+"#"+author.User.Discriminator),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Username**",
				Value: user.Username + "#" + user.Discriminator,
			},
			{
				Name:  "**User ID**",
				Value: userID,
			},
			{
				Name:  "**Sheet Link**",
				Value: sheetLink,
			},
		},
		Color: 0xEE0000,
	})

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")

}
