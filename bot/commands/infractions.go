package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	"gitlab.com/NorwegianLanguageLearning/heimdallr/bot"
	"time"
)

var infractionsCommand = command{
	"infractions",
	commandViewInfractions,
	"Shows the amount of infractions a user has.",
	[]string{
		"<user>",
	},
	[]string{
		"@username",
		"245267597929400192",
	},
}

//commandViewInfractions lists a user's infractions
func commandViewInfractions(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string))

	guildID := m.GuildID
	var user *discordgo.User
	var member *discordgo.Member
	var err error
	if member, err = heimdallr.GetMember(s, guildID, userID); err != nil {
		if user, err = s.User(userID); err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No user was found with ID %s.", userID))
			return errors.Wrap(err, "sending message failed")
		}
	} else {
		user = member.User
	}

	infractions, err := heimdallr.GetInfractions(userID)
	if err != nil {
		return err
	}
	numInfractions := len(infractions)
	// There is a limit of 25 fields in an embed.
	if numInfractions > 25 {
		infractions = infractions[numInfractions-25:]
	}
	var fields []*discordgo.MessageEmbedField
	for _, infraction := range infractions {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  infraction.Time.Format(time.RFC1123),
			Value: infraction.Reason,
		})
	}

	var pluralize string
	if len(infractions) == 1 {
		pluralize = ""
	} else {
		pluralize = "s"
	}

	title := fmt.Sprintf("User has %d infraction%s.", numInfractions, pluralize)
	if numInfractions > 25 {
		title += " Showing the last 25."
	}

	_, err = s.ChannelMessageSendEmbed(heimdallr.Config.AdminChannel, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.Username,
			IconURL: user.AvatarURL(""),
		},
		Title:  title,
		Fields: fields,
	})
	return errors.Wrap(err, "sending embed failed")
}
