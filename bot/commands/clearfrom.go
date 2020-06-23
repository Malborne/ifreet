package commands

import (
	"fmt"
	"strings"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
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

//commandWarnUser warns another user and gives an infraction.
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
	_, erro := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message ID received successfully"))
	if erro != nil {
		return errors.Wrap(erro, "sending message failed")
	}

	messages, err := s.ChannelMessages(m.ChannelID, 100, "", m.ID, "")
	if err != nil {
		return errors.Wrap(err, "getting messages failed")
	}
	guildID := m.GuildID

	number := len(messages)

	author, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", author.User.ID))
		return errors.Wrap(err, "sending message failed")
	}

	if number > 99 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You cannot delete more than 99 messages at a time"))
		return errors.Wrap(err, "deleting message failed")
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Are you sure you want to clear %d messages starting from https://discordapp.com/channels/678795606906634281/%s/%s ? This cannot be undone. ✅/❌", number, m.ChannelID, messages[len(messages)-1].ID))
	if err != nil {
		return errors.Wrap(err, "sending message failed")
	}

	return nil
}
