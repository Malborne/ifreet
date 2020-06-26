package heimdallr

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//OnDeleteHandler keeps a copy of deleted messages
func OnDeleteHandler(s *discordgo.Session, m *discordgo.MessageDelete) {
	// if m.Author.Bot {
	// 	return
	// }
	s.ChannelMessageSend(Config.ArchiveChannel, fmt.Sprintf("Message \"%s\" deletion detected with author %s.", m.Message.Content, m.Author.ID))

	guildID := m.GuildID

	author, err := GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(Config.AdminLogChannel, fmt.Sprintf("Message Author with ID %s was not found.", m.Author.ID))
		LogIfError(s, err)
	}

	_, err = s.ChannelMessageSendEmbed(Config.AdminLogChannel, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("A message was deleted"),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Message Author",
				Value: author.User.Username + "#" + author.User.Discriminator,
			},
			{
				Name:  "Channel",
				Value: fmt.Sprintf("<#%s>", m.ChannelID),
			},
			{
				Name:  "Message ID",
				Value: m.ID,
			},
			{
				Name:  "Message Content",
				Value: m.Content,
			},
		},
		Color: 0xEE0000,
	})
	if err != nil {
		LogIfError(s, errors.Wrap(err, "sending embed failed"))
		return
	}
	_, err = s.ChannelMessageSendEmbed(Config.ArchiveChannel, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("A message was deleted"),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Message Author",
				Value: author.User.Username + "#" + author.User.Discriminator,
			},
			{
				Name:  "Channel",
				Value: fmt.Sprintf("<#%s>", m.ChannelID),
			},
			{
				Name:  "Message ID",
				Value: m.ID,
			},
			{
				Name:  "Message Content",
				Value: m.Content,
			},
		},
		Color: 0xEE0000,
	})
	if err != nil {
		LogIfError(s, errors.Wrap(err, "sending embed failed"))
		return
	}

}
