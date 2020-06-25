package heimdallr

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

//LinksAndFilesHandler checks if someone sends a link and deletes the message if the user is new
func LinksAndFilesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	if len(m.Attachments) > 0 { //sent a file
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("A file has been detected"))
		if err != nil {
			LogIfError(s, err)
			return
		}
	}
	if len(m.Embeds) > 0 { //sent a link
		// if strings.Contains(strings.ToLower(m.Content), "https://") || strings.Contains(strings.ToLower(m.Content), "http://") {

		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("A linke has been detected"))
		if err != nil {
			LogIfError(s, err)
			return
		}
		// s.ChannelMessageDelete(m.ChannelID, m.ID)

	}

}
