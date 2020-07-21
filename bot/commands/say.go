package commands

import (
	"fmt"
	"strings"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"

	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var sayCommand = command{
	"say",
	commandSay,
	"Says a message in a specific channel.",
	[]string{
		"<message> <channelID>",
	},
	[]string{
		"Hello #general",
	},
}

//commandSay sends a message in a specific channel.
func commandSay(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	message, _ := args.String("<message>")
	channelID, _ := args.String("<channelID>")
	channelID = strings.Trim(strings.Trim(channelID, "<#"), ">")
	guildID := m.GuildID
	author, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", m.Author.ID))
		return errors.Wrap(err, "sending message failed")
	}
	if hasRole(author, heimdallr.Config.AdminRole) {
		_, err := s.ChannelMessageSend(channelID, fmt.Sprintf(message))

		return errors.Wrap(err, "sending message failed")

	}

	return errors.Wrap(err, "sending message failed")
}
