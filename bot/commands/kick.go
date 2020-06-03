package commands

import (
	"fmt"
	"time"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var kickCommand = command{
	"kick",
	commandKickUser,
	"Kicks a user from the server.",
	[]string{
		"<user> <reason>",
	},
	[]string{
		"@username \"Did something wrong\"",
		"245267597929480102 \"Did something wrong\"",
	},
}

//commandKickUser kicks a user from the server.
func commandKickUser(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string))
	reason, _ := args.String("<reason>")

	guildID := m.GuildID
	member, err := heimdallr.GetMember(s, guildID, userID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No member was found with ID %s.", userID))
		return errors.Wrap(err, "sending message failed")
	}
	user := member.User
	if userID == user.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to let you kick yourself, silly. 😉")
		return errors.Wrap(err, "sending message failed")
		if user.ID == "550664345302859786" { // Wasan's ID
			_, err := s.ChannelMessageSend(m.ChannelID, "I'm looking at you, وسن. I had to make this because of you 😒")
			return errors.Wrap(err, "sending message failed")
		}

	}
	if userID == s.State.User.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to kick myself, silly. 😉")
		return errors.Wrap(err, "sending message failed")
	}

	err = s.GuildMemberDeleteWithReason(guildID, userID, reason)
	if err != nil {
		return errors.Wrap(err, "kick failed")
	}
	err = heimdallr.AddInfraction(*user, heimdallr.Infraction{
		Reason: fmt.Sprintf("Received a kick with reason: %s", reason),
		Time:   time.Now(),
	})
	if err != nil {
		return err
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}
