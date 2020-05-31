package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var DMUnapprovedCommand = command{
	"dmunapproved",
	commandDMUnapproved,
	"DM unapproved users to let them know that they will lose acess to the server and should contact on of the mods and direct them to the #approval and verification channel.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//DM unapproved users to let them know that they will lose acess to the server and should contact on of the mods and direct them to the #approval and verification channel.
func commandDMUnapproved(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		return err
	}
	var count int = 0
	for _, member := range guild.Members {

		if !isApproved(member) && member.User.ID != s.State.User.ID && member.User.ID != 680170745627934742 && member.User.ID != 714905000668692541 { //IDs for Admin and supermod
			userChannel, err := s.UserChannelCreate(member.User.ID)
			if err != nil {
				return errors.Wrap(err, "creating private channel failed")
			}
			_, err = s.ChannelMessageSend(userChannel.ID, fmt.Sprintf(
				"You are an unapproved member of Learn/Memorize Quran Server and you are about to lose access to most of the server. If you still wish to retain access to the server, please contact one of the moderators in the #approval-and-verification channel below to be approved.\n  https://discord.gg/R6jKWT \n\nYou cannot reply to this message."))
			if err != nil {
				return errors.Wrap(err, "sending message failed")
			} else {
				count = count + 1
			}
		}
	}

	if count == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No unapproved users found."))
		return errors.Wrap(err, "sending message failed")
	} else {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sucessfully sent messages to %d user(s)", count))
		return errors.Wrap(err, "sending message failed")
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}
