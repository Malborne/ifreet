package commands

import (
	"fmt"
	"time"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var banCommand = command{
	"ban",
	commandBanUser,
	"Bans a user from the server.",
	[]string{
		"<user> <reason>",
	},
	[]string{
		"@username \"Did something wrong\"",
		"245167597929480192 \"Did something wrong\"",
	},
}

//commandBanUser bans a user from the server.
func commandBanUser(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string))
	reason, _ := args.String("<reason>")

	guildID := m.GuildID
	member, err := heimdallr.GetMember(s, guildID, userID)
	var user *discordgo.User
	if err != nil {
		user, err = s.User(userID)
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No user was found with ID %s.", userID))
			return errors.Wrap(err, "sending message failed")
		}
	} else {
		user = member.User
	}
	if userID == s.State.User.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to ban myself, silly. 😉")
		return errors.Wrap(err, "sending message failed")
	}

	err = s.GuildBanCreateWithReason(guildID, userID, reason, 0)
	if err != nil {
		return errors.Wrap(err, "ban failed")
	}
	err = heimdallr.AddInfraction(*user, heimdallr.Infraction{
		Reason: fmt.Sprintf("Received a ban with reason: %s", reason),
		Time:   time.Now(),
	})
	if err != nil {
		return err
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}
