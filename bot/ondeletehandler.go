package heimdallr

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//OnDeleteHandler keeps a copy of deleted messages
func OnDeleteHandler(s *discordgo.Session, m *discordgo.MessageDelete) {

	message, err := GetFromArchive(m.ID)
	if err != nil {
		LogIfError(s, errors.Wrap(err, "Getting the message failed from the database."))
		return
	}
	if message.userID == "" { //The message is not logged in the archive
		return
	}
	author, err := s.User(message.userID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No user was found with ID %s.", message.userID))
		LogIfError(s, errors.Wrap(err, "Failed to get the author of the message."))
		return
	}
	if author.Bot {
		return
	}

	_, err = s.ChannelMessageSendEmbed(Config.ArchiveChannel, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("A message was deleted"),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Message Author",
				Value: author.String(),
			},
			{
				Name:  "Channel",
				Value: fmt.Sprintf("<#%s>", message.channelID),
			},
			{
				Name:  "Time sent",
				Value: fmt.Sprintf("%s", message.Time.Format(time.RFC1123)),
			},
			{
				Name:  "Message ID",
				Value: message.messageID,
			},
			{
				Name:  "Message Content",
				Value: message.content,
			},
		},
		Color: 0xEE0000,
	})
	if err != nil {
		LogIfError(s, errors.Wrap(err, "sending embed failed"))
		return

	}
	err = RemovefromArchive(message.messageID)
	if err != nil {
		LogIfError(s, errors.Wrap(err, "sending embed failed"))
		return

	}
}
