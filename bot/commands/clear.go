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
		"500",
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

	prompt, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Are you sure you want to clear %s messages? This cannot be undone.", number))
	return errors.Wrap(err, "sending message failed")
	err = s.MessageReactionAdd(m.ChannelID, prompt.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
	err = s.MessageReactionAdd(m.ChannelID, prompt.ID, "❌")
	return errors.Wrap(err, "adding reaction failed")

	if 

	_, err = s.ChannelMessageSendEmbed(heimdallr.Config.AdminLogChannel, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s Messages  were cleared. The command was made by %s", number, author.Mention()),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Username**",
				Value: author.User.Username + "#" + author.User.Discriminator,
			},
		},
		Color: 0xEE0000,
	})

}
