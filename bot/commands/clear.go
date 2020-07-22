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

//commandClearMessages Clears a number of messages.
func commandClearMessages(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	number, err := args.Int("<number>")

	if err != nil || number == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Incorrect use of command. Type the number of messages you wish to be deleted"))
		return errors.Wrap(err, "clearing message failed")
	}

	if number > 99 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You cannot delete more than 99 messages at a time"))
		return errors.Wrap(err, "clearing message failed")
	}

	author, err := heimdallr.GetMember(s, m.GuildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", m.Author.ID))
		return errors.Wrap(err, "getting the author failed")
	}
	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		return err
	}
	if m.ChannelID == heimdallr.Config.StaffChannel && !heimdallr.IsAdminOrHigher(author, guild) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You cannot clear messages in the <#%s>.", heimdallr.Config.StaffChannel))
		return errors.Wrap(err, "clearing messages failed")
	}

	messages, err := s.ChannelMessages(m.ChannelID, number, m.ID, "", "")
	if err != nil {
		return errors.Wrap(err, "getting messages failed")
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Are you sure you want to clear %d messages starting from https://discordapp.com/channels/678795606906634281/%s/%s ?\nThis cannot be undone. ✅/❌", number, m.ChannelID, messages[len(messages)-1].ID))
	if err != nil {
		return errors.Wrap(err, "sending message failed")
	}

	return nil
}

//ReactionPrompt Performs the clear action based on the response to the prompt
func ReactionPrompt(s *discordgo.Session, m *discordgo.MessageReactionAdd) {

	message, err := heimdallr.GetMessage(s, m.ChannelID, m.MessageID)
	if !message.Author.Bot || message.Author.ID != s.State.User.ID {
		return
	}
	if !strings.Contains(message.Content, "Are you sure you want to clear") {
		return
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
	if !heimdallr.IsSuperModOrHigher(reactingMember, guild) {
		// _, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You do NOT have permissions to delete messages"))
		heimdallr.LogIfError(s, err)
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
			if !messages[mess].Author.Bot && !strings.HasPrefix(messages[mess].Content, ";") {
				_, err = s.ChannelMessageSendEmbed(heimdallr.Config.ArchiveChannel, &discordgo.MessageEmbed{
					Title: fmt.Sprintf("This Message  was cleared by %s", reactingMember.User.String()),
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "**Message Author**",
							Value: messages[mess].Author.String(),
						},
						{
							Name:  "**Message Content**",
							Value: messages[mess].Content,
						},
						{
							Name:  "**Channel**",
							Value: fmt.Sprintf("<#%s>", messages[mess].ChannelID),
						},
					},
					Color: 0x00FF00,
				})
				if err != nil {
					heimdallr.LogIfError(s, err)
					return
				}
			}
			//TODO Add Delete message from database code here
			s.ChannelMessageDelete(message.ChannelID, messages[mess].ID)

		}
		s.ChannelMessageDelete(message.ChannelID, message.ID)
		_, err = s.ChannelMessageSendEmbed(heimdallr.Config.AdminLogChannel, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("%d Messages  were cleared by", number),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "**Username**",
					Value: reactingMember.User.String(),
				},
				{
					Name:  "**Channel**",
					Value: fmt.Sprintf("<#%s>", message.ChannelID),
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

}
