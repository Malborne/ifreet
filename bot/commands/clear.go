package commands

import (
	"fmt"
	"strconv"
	"strings"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var clearCommand = command{
	"clear",
	commandClearMessages,
	"Clears a number of messages.",
	[]string{
		"<number>",
	},
	[]string{
		"10",
		"100",
	},
}

//commandWarnUser warns another user and gives an infraction.
func commandClearMessages(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	number, _ := args.Int("<number>")

	guildID := m.GuildID

	// guild, err := heimdallr.GetGuild(s, guildID)
	// if err != nil {
	// 	return err
	// }

	author, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", author.User.ID))
		return errors.Wrap(err, "sending message failed")
	}

	if number > 100 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You cannot delete more than 100 messages at a time"))
		return errors.Wrap(err, "deleting message failed")
	}

	messages, err := s.ChannelMessages(m.ChannelID, number, m.ID, m.ID, m.ID)
	if err != nil {
		return errors.Wrap(err, "getting messages failed")
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Are you sure you want to clear %d messages starting from %s? This cannot be undone. ✅/❌", number, messages[0].ID))
	if err != nil {
		return errors.Wrap(err, "sending message failed")
	}
	// err = s.MessageReactionAdd(m.ChannelID, prompt.ID, "✅")
	// if err != nil {
	// 	return errors.Wrap(err, "adding reaction failed")
	// }
	// err = s.MessageReactionAdd(m.ChannelID, prompt.ID, "❌")
	// if err != nil {
	// 	return errors.Wrap(err, "adding reaction failed")
	// }

	return nil
}

//ReactionPrompt Performs the clear action based on the response to the prompt
func ReactionPrompt(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	// guildID := m.GuildID
	// guild, err := heimdallr.GetGuild(s, guildID)

	// if err != nil {
	// 	return
	// }

	// if prompt.ID != m.MessageID {
	// 	return
	// }

	reactingMember, err := heimdallr.GetMember(s, m.GuildID, m.UserID)
	if err != nil {
		heimdallr.LogIfError(s, err)
		return
	}
	message, err := heimdallr.GetMessage(s, m.ChannelID, m.MessageID)
	if !message.Author.Bot || message.Author.ID != s.State.User.ID {
		return
	}

	var number int = 0
	for _, word := range strings.Split(message.Content, " ") {
		if n, err := strconv.Atoi(word); err == nil {
			number = n
			break
		}

	}
	if m.Emoji.Name == "✅" {
		messages, err := s.ChannelMessages(m.ChannelID, number+1, message.ID, "", "")
		if err != nil {
			heimdallr.LogIfError(s, err)
			return
		}
		for mess := range messages {

			s.ChannelMessageDelete(message.ChannelID, messages[mess].ID)

		}
		s.ChannelMessageDelete(message.ChannelID, message.ID)
		_, err = s.ChannelMessageSendEmbed(heimdallr.Config.AdminLogChannel, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%d Messages  were cleared.", len(messages)),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "**Username**",
					Value: reactingMember.User.String(),
				},
			},
			Color: 0xEE0000,
		})
		if err != nil {
			heimdallr.LogIfError(s, err)
			return
		}
	}

	if m.Emoji.Name == "❌" {
		messages, err := s.ChannelMessages(message.ChannelID, 2, message.ID, message.ID, message.ID)
		if err != nil {
			heimdallr.LogIfError(s, err)
			return
		}
		//Delete the message and the command
		for mess := range messages {

			s.ChannelMessageDelete(message.ChannelID, messages[mess].ID)

		}
		// s.ChannelMessageDelete(prompt.ChannelID, prompt.ID)
	}

	reactingMember, err := heimdallr.GetMember(s, m.GuildID, m.UserID)
	if err != nil {
		heimdallr.LogIfError(s, err)
		return
	}
	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		heimdallr.LogIfError(s, err)
		return
	}
	if !heimdallr.IsModOrHigher(reactingMember, guild) {
		//Output, You don't have permissions to delete messages
		return
	}

}
