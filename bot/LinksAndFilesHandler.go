package heimdallr

import (
	"time"

	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//LinksAndFilesHandler checks if someone sends a link and deletes the message if the user is new
func LinksAndFilesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}
	guildID := m.GuildID
	guild, err := GetGuild(s, guildID)
	if err != nil {
		LogIfError(s, err)
	}

	author, err := GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", m.Author.ID))
		LogIfError(s, err)
	}
	if IsModOrHigher(author, guild) {
		return
	}
	joinedAt, err := author.JoinedAt.Parse()
	if err != nil {
		LogIfError(s, err)
	}
	if joinedAt.Before(time.Now().AddDate(0, 0, -1)) { //If they joined the server more than 24 ago, just ignore it
		return
	}
	if len(m.Attachments) > 0 { //sent a file
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are NOT allowed to send files yet. Please wait until you are on the server for a longer time."))
		if err != nil {
			LogIfError(s, err)
			return
		}
		s.ChannelMessageDelete(m.ChannelID, m.ID)

	}
	if len(m.Embeds) > 0 || strings.Contains(strings.ToLower(m.Content), "https://") || strings.Contains(strings.ToLower(m.Content), "http://") { //sent a link

		if strings.Contains(strings.ToLower(m.Content), "youtube.com") || strings.Contains(strings.ToLower(m.Content), "https://youtu.be/") { //Ignore YouTube videos
			return
		}

		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You are NOT allowed to send files yet. Please wait until you are on the server for a longer time."))
		if err != nil {
			LogIfError(s, err)
			return
		}
		s.ChannelMessageDelete(m.ChannelID, m.ID)

	}

}
