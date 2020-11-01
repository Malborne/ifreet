package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

//DMUnverifiedCommand DM Unverified users to let them know that they will lose acess to the server and should contact on of the mods and direct them to the #approval and verification channel.
var DMUnverifiedCommand = command{
	"dmunverified",
	commandDMUnverified,
	"DM Unverified users to let them know that they will lose acess to the gender specific channels and should contact on of the mods and directs them to the #welcome.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//DM Unverified users to let them know that they will lose acess to the server and should contact on of the mods and direct them to the #approval and verification channel.
func commandDMUnverified(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		return err
	}
	var count int = 0
	for _, member := range guild.Members {

		if !isVerified(member) && member.User.ID != s.State.User.ID {
			userChannel, err := s.UserChannelCreate(member.User.ID)
			if err != nil {
				return errors.Wrap(err, "creating private channel failed")
			}
			_, err = s.ChannelMessageSend(userChannel.ID, fmt.Sprintf(
				"You are an Unverified member of Learn/Memorize Quran Server and you are about to lose access to the gender specific channels on the server. If you still wish to retain access to those channels, please contact one of the moderators in the %s below to be verified.\n  https://discord.gg/R6jKWT \n\nYou cannot reply to this message.", heimdallr.Config.WelcomeChannel))
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Does NOT ACCEPT DMs", member.User.Mention()))
				// return errors.Wrap(err, "sending message failed")
			} else {
				count = count + 1
			}
		}
	}

	if count == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No Unverified users found."))
		return errors.Wrap(err, "sending message failed")
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Sucessfully sent messages to %d user(s)", count))
	return errors.Wrap(err, "sending message failed")

}
