package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	"gitlab.com/NorwegianLanguageLearning/heimdallr/bot"
	"strings"
)

var welcomeMessageCommand = command{
	"welcomemessage",
	commandWelcomeMessage,
	"Handles the welcome message.",
	[]string{
		"show",
		"set <message>",
	},
	[]string{
		"show",
		"set \"This is an example welcome message\"",
	},
}

func commandWelcomeMessage(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	if show, _ := args.Bool("show"); show {
		return showWelcomeMessage(s, m)
	}
	return setWelcomeMessage(s, m, args["<message>"].(string))
}

func setWelcomeMessage(s *discordgo.Session, m *discordgo.MessageCreate, welcomeMessage string) error {
	if strings.Count(welcomeMessage, "%s") > 1 {
		_, err := s.ChannelMessageSend(m.ChannelID, "The welcome message can only have a single placeholder for the user mention.")
		return errors.Wrap(err, "sending message failed")
	}

	heimdallr.Config.WelcomeMessage = welcomeMessage
	err := heimdallr.Config.SaveConfig("config.toml")
	if err != nil {
		return err
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}

func showWelcomeMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	welcomeMessage := heimdallr.Config.WelcomeMessage
	var err error
	if welcomeMessage == "" {
		_, err = s.ChannelMessageSend(m.ChannelID, "No welcome message set.")
	} else {
		_, err = s.ChannelMessageSend(m.ChannelID, heimdallr.Config.WelcomeMessage)
	}
	return errors.Wrap(err, "sending message failed")
}
