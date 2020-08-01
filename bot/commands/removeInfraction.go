package commands

import (
	"fmt"
	"time"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var removeInfractionCommand = command{
	"removeinfraction",
	commandRemoveInfraction,
	"Removes an infraction from a user.",
	[]string{
		"<user> <infractionID>",
	},
	[]string{
		"@username 22",
		"245207597929480192 20",
	},
}

//commandRemoveInfraction Removes an infraction from a user.
func commandRemoveInfraction(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string), s)
	infractionID, _ := args.String("<infractionID>")

	guildID := m.GuildID

	guild, err := heimdallr.GetGuild(s, guildID)
	if err != nil {
		return err
	}

	user, err := s.User(userID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No user was found with ID %s.", userID))
		return errors.Wrap(err, "sending message failed")
	}

	author, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", userID))
		return errors.Wrap(err, "sending message failed")
	}

	if !heimdallr.IsAdminOrHigher(author, guild) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You do not have the permission to remove infractions"))
		return errors.Wrap(err, "sending message failed")
	}

	infractions, err := heimdallr.GetInfractions(userID)
	if err != nil {
		return err
	}

	var infraction heimdallr.Infraction
	found := false
	for _, infrac := range infractions {
		if infrac.ID == infractionID {
			infraction = infrac
			found = true
		}
	}
	if !found {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Could not find the infraction for the user with the ID specified"))
		return errors.Wrap(err, "Finding infraction failed")
	}
	err = heimdallr.RemoveInfraction(infractionID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to remove the infraction from the Database"))
		heimdallr.LogIfError(s, err)
		return errors.Wrap(err, "Removing infraction failed")
	}
	_, err = s.ChannelMessageSendEmbed(heimdallr.Config.LogChannel, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("An infraction was removed by %s.", author.User.Username+"#"+author.User.Discriminator),

		Author: &discordgo.MessageEmbedAuthor{
			Name:    user.Username,
			IconURL: user.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  infraction.Time.Format(time.RFC1123),
				Value: infraction.Reason,
			},
		},
		Color: 0xEE0000,
	})

	if err != nil {
		return errors.Wrap(err, "Failed to send the message Embed to the log channel")
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	return errors.Wrap(err, "adding reaction failed")

}
