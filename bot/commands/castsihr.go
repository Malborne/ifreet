package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var castSihrCommand = command{
	"warn",
	commandcastSihr,
	"Casts Sihr on the user by changing their nickname to Mashoor",
	[]string{
		"<user>",
	},
	[]string{
		"@username",
		"245207597929480192",
	},
}

//Casts Sihr on the user by changing their nickname to Mashoor.
func commandcastSihr(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string))
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
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I'm not going to cast Sihr on myself, 👎\nI will cast Sihr on you instead Muhahahaha"))
		castSihr(author)
		return errors.Wrap(err, "sending message failed")
	}

	if heimdallr.IsAdminOrHigher(infractor, guild) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You cannot cast Sihr on the admin. 👎\nI will cast Sihr on you instead Muhahahaha"))
		castSihr(author)
		return errors.Wrap(err, "sending message failed")

	}

	if m.Author.ID == user.ID && userID == "550664345302859786" { // Wasan's ID
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to let you cast Sihr on yourself, silly. 😉 I'm looking at you, وسن. I had to make this because of you 😒")
		return errors.Wrap(err, "sending message failed")
	} else if m.Author.ID == user.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to let you cast Sihr on yourself, silly. 😉")
		return errors.Wrap(err, "sending message failed")
	}

	if isOneLowerThanTwo(author, infractor) {
		_, err := s.ChannelMessageSend(m.ChannelID, "You cannot cast Sihr on a user that has the same or a role higher than you.\nI will cast Sihr on you instead Muhahahaha")
		castSihr(author)
		return errors.Wrap(err, "sending message failed")
	}

	castSihr(infractor)

	// _, err = s.ChannelMessageSendEmbed(heimdallr.Config.AdminLogChannel, &discordgo.MessageEmbed{
	// 	Title: "User was warned.",
	// 	Fields: []*discordgo.MessageEmbedField{
	// 		{
	// 			Name:  "**Username**",
	// 			Value: user.Username + "#" + user.Discriminator,
	// 		},
	// 		{
	// 			Name:  "**User ID**",
	// 			Value: user.ID,
	// 		},
	// 	},
	// 	Color: 0xEE0000,
	// })

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	return errors.Wrap(err, "adding reaction failed")

}

func castSihr(member *discordgo.Member) {
	// _, err := member.User.setNickname("Mashoor💨")
	// return errors.Wrap(err, "casting Sihr failed")

}
