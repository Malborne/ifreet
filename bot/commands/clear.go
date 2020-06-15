package commands

import (
	"fmt"

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

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Are you sure you want to clear %d messages? This cannot be undone.", number))
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

//Reaction Prompt Performs the clear action based on the response to the prompt
func ReactionPrompt(s *discordgo.Session, m *discordgo.MessageReactionAdd, prompt *discordgo.Message) {
	// guildID := m.GuildID
	// guild, err := heimdallr.GetGuild(s, guildID)

	// if err != nil {
	// 	return
	// }

	// if prompt.ID != m.MessageID {
	// 	return
	// }

	message, err := heimdallr.GetMessage(s, m.ChannelID, m.MessageID)
	if !message.Author.Bot {
		return
	}

	if m.Emoji.Name != "✅" && m.Emoji.Name != "❌" {
		//Output incorrect reactions
		return
	}

	if m.Emoji.Name == "✅" {
		_, err := s.ChannelMessageSendEmbed(heimdallr.Config.AdminLogChannel, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("Messages  were cleared. The command was made by ..."),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "**Username**",
					Value: "Username" + "#" + "123456",
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
		messages, err := s.ChannelMessages(prompt.ChannelID, 2, prompt.ID, prompt.ID, prompt.ID)
		if err != nil {
			heimdallr.LogIfError(s, err)
			return
		}
		//Delete the message and the command
		for mess := range messages {
			s.ChannelMessageDelete(prompt.ChannelID, messages[mess].ID)

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
