package commands

import (
	"fmt"
	"strings"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var approvalMessageCommand = command{
	"approvalmessage",
	commandApprovalMessage,
	"Handles the approval message.",
	[]string{
		"show",
		"set <message>",
	},
	[]string{
		"show",
		"set \"This is an example approval message\"",
	},
}

func commandApprovalMessage(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	if show, _ := args.Bool("show"); show {
		return showApprovalMessage(s, m)
	}
	return setApprovalMessage(s, m, args["<message>"].(string))
}

func setApprovalMessage(s *discordgo.Session, m *discordgo.MessageCreate, approvalMessage string) error {
	if strings.Count(approvalMessage, "%s") > 2 {
		_, err := s.ChannelMessageSend(m.ChannelID, "The approval message can only have two placeholders for the user mention and for #bot-commands.")
		return err
	}

	heimdallr.Config.ApprovalMessage = approvalMessage
	err := heimdallr.Config.SaveConfig("config.toml")
	if err != nil {
		return err
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}

func showApprovalMessage(s *discordgo.Session, m *discordgo.MessageCreate) error {
	approvalMessage := heimdallr.Config.ApprovalMessage
	var err error
	if approvalMessage == "" {
		_, err = s.ChannelMessageSend(m.ChannelID, "No approval message set.")

	} else {
		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(heimdallr.Config.ApprovalMessage, "<@Username>", heimdallr.Config.BotChannel))
	}
	return errors.Wrap(err, "sending message failed")
}
