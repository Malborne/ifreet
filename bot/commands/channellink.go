package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	"gitlab.com/NorwegianLanguageLearning/heimdallr/bot"
	"regexp"
)

var channelLinkCommand = command{
	"channellink",
	commandChannelLink,
	"Get a link to a channel (useful for video chat).",
	[]string{
		"<channel>",
	},
	[]string{
		"441787683208036352",
		"#bot-chat",
		"tekst-for-søndagschatten",
	},
}

//commandVersion prints information about the program's current version and commit.
func commandChannelLink(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	channelIDOrName := args["<channel>"].(string)
	var channelID string
	if submatch := regexp.MustCompile(`<#(\d+)>`).FindStringSubmatch(channelIDOrName); len(submatch) == 2 {
		channelID = submatch[1]
	} else if submatch := regexp.MustCompile(`(\d+)`).FindStringSubmatch(channelIDOrName); len(submatch) == 2 {
		channelID = submatch[1]
	} else {
		guild, err := heimdallr.GetGuild(s, m.GuildID)
		if err != nil {
			return err
		}
		for _, channel := range guild.Channels {
			// If there are duplicate names, prioritize voice channels
			if channel.Name == channelIDOrName && (channelID == "" || channel.Type == discordgo.ChannelTypeGuildVoice) {
				channelID = channel.ID
			}
		}
		if channelID == "" {
			_, err = s.ChannelMessageSend(m.ChannelID, "Unknown channel.")
			return errors.Wrap(err, "sending message failed")
		}
	}
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<https://www.discordapp.com/channels/%s/%s>", m.GuildID, channelID))
	return errors.Wrap(err, "sending message failed")
}
