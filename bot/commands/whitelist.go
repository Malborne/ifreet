package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var whitelistCommand = command{
	"mute",
	commandWhitelist,
	"adds a user to the whitelist so they can post links and files.",
	[]string{
		"<user>",
	},
	[]string{
		"@username",
		"245207597929480192",
	},
}

//commandMuteUser mues another user
func commandWhitelist(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string), s)
	// number, _ := args.Int("<no>")
	var user *discordgo.User

	guildID := m.GuildID

	guild, err := heimdallr.GetGuild(s, guildID)
	if err != nil {
		return err
	}

	infractor, err := heimdallr.GetMember(s, guildID, userID)

	if err != nil {
		user, err = s.User(userID)
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No user was found with ID %s.", userID))
			return errors.Wrap(err, "sending message failed")
		}
	} else {
		user = infractor.User
	}

	if userID == s.State.User.ID {
		// _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I'm not going to mute myself, silly. 😉"))
		return errors.Wrap(err, "sending message failed")
	}

	if heimdallr.IsModOrHigher(infractor, guild) {
		// _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You cannot mute the admin. 👎"))
		return errors.Wrap(err, "sending message failed")
	}

	if m.Author.ID == user.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, "You cannot whitelist yourself, silly. 😉")
		return errors.Wrap(err, "sending message failed")
	}

	//Add the user to the whitelist in the database
	// err = heimdallr.AddWhitelistedUser(*infractor.User, time.Now())
	// if err != nil {
	// 	return err
	// }

	_, err = s.ChannelMessageSendEmbed(heimdallr.Config.AdminLogChannel, &discordgo.MessageEmbed{
		Title: "User was whitelisted.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Username**",
				Value: user.Username + "#" + user.Discriminator,
			},
			{
				Name:  "**User ID**",
				Value: user.ID,
			},
		},
		Color: 0xEE0000,
	})

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	return errors.Wrap(err, "adding reaction failed")

}
