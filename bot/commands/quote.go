package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var quoteCommand = command{
	"quote",
	commandQuote,
	"Quotes an earlier message.",
	[]string{
		"<message-id>",
	},
	[]string{
		"544277290610708140",
	},
}

//commandQuote quotes another used by message ID.
func commandQuote(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	messageID, _ := args.String("<message-id>")
	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		return err
	}

	var message *discordgo.Message
	var quoteChannel *discordgo.Channel
	// If we're very unlucky we might get the wrong message here,
	// because IDs aren't unique across channels. The alternative is to
	// require the user to submit the channel as well.
	for _, channel := range guild.Channels {
		message, err = s.ChannelMessage(channel.ID, messageID)
		if err == nil {
			quoteChannel = channel
			break
		}
	}
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No message was found with ID %s.", messageID))
		return errors.Wrap(err, "sending message failed")
	}

	permissions, err := s.UserChannelPermissions(m.Author.ID, quoteChannel.ID)
	if err != nil {
		return errors.Wrap(err, "getting permissions failed")
	}
	if !heimdallr.CanRead(permissions) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You don't have permission to quote this message."))
		return errors.Wrap(err, "sending message failed")
	}

	authorMember, err := heimdallr.GetMember(s, m.GuildID, message.Author.ID)
	var name string
	if err == nil && authorMember.Nick != "" {
		name = authorMember.Nick
	} else {
		name = message.Author.Username
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    name,
			IconURL: message.Author.AvatarURL(""),
		},
		Timestamp: string(message.Timestamp),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("(quoted by: %s) | #%s", m.Author.Username, quoteChannel.Name),
		},
		Description: fmt.Sprintf("%s\n[⁽ˢᵒᵘʳᶜᵉ⁾](https://discordapp.com/channels/%s/%s/%s)", message.Content, m.GuildID, message.ChannelID, message.ID),
	})
	if err != nil {
		return errors.Wrap(err, "sending embed failed")
	}
	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	return errors.Wrap(err, "deleting message failed")
}
