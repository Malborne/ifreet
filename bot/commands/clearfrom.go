package commands

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var clearFromCommand = command{
	"clearfrom",
	commandClearFromMessage,
	"Clears a number of messages from a starting message.",
	[]string{
		"<startingID>",
	},
	[]string{
		"724820450080849971",
		"https://discordapp.com/channels/678795606906634281/724040027436351508/724820450080849971",
	},
}

//commandClearFromMessage Clears a number of messages from a starting message.
func commandClearFromMessage(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	startingID, _ := args.String("<startingID>")

	if strings.Contains(startingID, "https://discordapp.com/channels") {
		slices := strings.Split(m.Content, "/")
		startingID = slices[len(slices)-1]
	}

	if _, err := s.ChannelMessage(m.ChannelID, startingID); err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Failed to get the starting message Please make sure that you enter either the ID of the message or a link to the message."))
		return errors.Wrap(err, "deleting message failed")
	}

	messages, err := s.ChannelMessages(m.ChannelID, 100, "", startingID, "")
	if err != nil {
		return errors.Wrap(err, "getting messages failed")
	}

	number := len(messages)

	if number > 99 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You cannot delete more than 99 messages at a time"))
		return errors.Wrap(err, "deleting message failed")
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Are you sure you want to clear %d messages starting from https://discordapp.com/channels/678795606906634281/%s/%s ?\nThis cannot be undone. ✅/❌", number, m.ChannelID, startingID))
	if err != nil {
		return errors.Wrap(err, "sending message failed")
	}

	return nil
}
