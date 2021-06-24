package heimdallr

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//UserJoinHandler handles new users joining the server, and will welcome them.
func NewMemberJoinHandler(s *discordgo.Session, g *discordgo.GuildMemberAdd) {

	newChannel, err := s.GuildChannelCreate(g.GuildID, g.User.Username, discordgo.ChannelTypeGuildText)
	if err != nil {
		LogIfError(s, errors.Wrap(err, "Creating New Channel failed"))

	}
	err = s.ChannelPermissionSet(newChannel.ID, Config.UserRole, discordgo.PermissionOverwriteTypeRole, 0, 1)
	if err != nil {
		LogIfError(s, errors.Wrap(err, "Changing permissions failed"))

	}
	welcomeMessage := Config.WelcomeMessage
	if strings.Count(welcomeMessage, "%s") > 0 {
		welcomeMessage = fmt.Sprintf(welcomeMessage, g.User.Mention(), Config.RulesChannel)
	}

	_, err = s.ChannelMessageSend(newChannel.ID, welcomeMessage)
	LogIfError(s, errors.Wrap(err, "sending message failed"))
}
