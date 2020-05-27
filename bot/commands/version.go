package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	"gitlab.com/NorwegianLanguageLearning/heimdallr/bot/version"
)

var versionCommand = command{
	"version",
	commandVersion,
	"Show version and commit information.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//commandVersion prints information about the program's current version and commit.
func commandVersion(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("**Heimdallr**\nVersion: *%s*\nCommit: *%s*", version.VERSION, version.COMMIT))
	return errors.Wrap(err, "sending message failed")
}
