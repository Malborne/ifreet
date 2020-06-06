package commands

import (
	"fmt"
	"time"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var warnCommand = command{
	"warn",
	commandWarnUser,
	"Warns a user / increases the number of infractions by one.",
	[]string{
		"<user> <reason>",
	},
	[]string{
		"@username \"Did something wrong\"",
		"245207597929480192 \"Did something wrong\"",
	},
}

//commandWarnUser warns another user and gives an infraction.
func commandWarnUser(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string))
	reason, _ := args.String("<reason>")
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

	author, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author was not found.", userID))
		return errors.Wrap(err, "sending message failed")
	}

	if userID == s.State.User.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I'm not going to warn myself, silly. 😉"))
		return errors.Wrap(err, "sending message failed")
	}

	if heimdallr.IsAdminOrHigher(infractor, guild) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You cannot warn the admin. 👎"))
		return errors.Wrap(err, "sending message failed")
	}

	if m.Author.ID == user.ID && userID == "550664345302859786" { // Wasan's ID
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to let you warn yourself, silly. 😉 I'm looking at you, وسن. I had to make this because of you 😒")
		return errors.Wrap(err, "sending message failed")
	} else if m.Author.ID == user.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to let you warn yourself, silly. 😉")
		return errors.Wrap(err, "sending message failed")
	}

	if isOneLowerThanTwo(author, infractor) {
		_, err := s.ChannelMessageSend(m.ChannelID, "You cannot warn a user that has the same or a role higher than you")
		return errors.Wrap(err, "sending message failed")
	}

	err = heimdallr.AddInfraction(*infractor.User, heimdallr.Infraction{Reason: reason, Time: time.Now()})
	if err != nil {
		return err
	}

	if err != nil {
		return errors.Wrap(err, "getting user failed")
	}
	_, err = s.ChannelMessageSendEmbed(heimdallr.Config.AdminLogChannel, &discordgo.MessageEmbed{
		Title: "User was warned.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Username**",
				Value: user.Username + "#" + user.Discriminator,
			},
			{
				Name:  "**User ID**",
				Value: user.ID,
			},
			{
				Name:  "**Reason**",
				Value: reason,
			},
		},
		Color: 0xEE0000,
	})

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	userChannel, err := s.UserChannelCreate(userID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Does NOT ACCEPT DMs but the warning was successfully registered", infractor.Mention()))
		return nil
		// return errors.Wrap(err, "creating private channel failed")
	}
	_, err = s.ChannelMessageSend(userChannel.ID, fmt.Sprintf(
		"You have received a warning in %s for the following reason: %s\n\nYou cannot reply to this message.",
		guild.Name, reason,
	))
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Does NOT ACCEPT DMs but the warning was successfully registered", infractor.Mention()))
		return nil
		// return errors.Wrap(err, "sending message failed")
	}

	return errors.Wrap(err, "adding reaction failed")

}
