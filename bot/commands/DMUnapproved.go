package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

//DMUnapprovedCommand DMs all the unapproved users
var DMUnapprovedCommand = command{
	"dmunapproved",
	commandDMUnapproved,
	"DM unapproved users to let them know that they will lose acess to the server and should contact on of the mods and direct them to the #welcome.",
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
	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("There are %d members in this guild", len(guild.Members)))
	for _, member := range guild.Members {

		if !isApproved(member) && !member.User.Bot {
			count = count + 1
			userChannel, err := s.UserChannelCreate(member.User.ID)
			if err != nil {
				return errors.Wrap(err, "creating private channel failed")
			}
			_, err = s.ChannelMessageSend(userChannel.ID, fmt.Sprintf(
				"You are an unapproved member of Quran Learning Center Server and you do not have access to most of the server. If you would like to have access to the server, please contact one of the moderators in the %s channel below to be approved.\n\n\nhttps://discord.gg/R6jKWT\n\nKeep in mind that if you stay for longer than a week without getting approved, you will risk being kicked out of the server.\n\nYou cannot reply to this message.", heimdallr.Config.WelcomeChannel))
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Does NOT ACCEPT DMs", member.User.Mention()))
				// return errors.Wrap(err, "sending message failed")
			}
		}
	}

	if count == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No unapproved users found."))
		return errors.Wrap(err, "sending message failed")
	} else {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sucessfully sent messages to %d user(s)", count))
		return errors.Wrap(err, "sending message failed")
	}

}
