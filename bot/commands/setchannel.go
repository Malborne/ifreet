package commands

import (
	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var setChannelCommand = command{
	"setchannel",
	commandSetChannel,
	"Sets the current channel as the specified type.",
	[]string{
		"(welcome | log | admin | adminlog | bot)",
	},
	[]string{
		"welcome",
		"log",
		"admin",
		"adminlog",
		"bot",
	},
}

//commandSetChannel sets the welcome, log, admin, etc. channels
func commandSetChannel(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	channelTypes := []string{"welcome", "log", "admin", "adminlog", "bot"}
	var channelType string
	for _, potentialChannelType := range channelTypes {
		if present, _ := args.Bool(potentialChannelType); present {
			channelType = potentialChannelType
			break
		}
	}

	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		return errors.Wrap(err, "getting channel failed")
	}

	switch channelType {
	case "welcome":
		heimdallr.Config.WelcomeChannel = channel.ID
	case "log":
		heimdallr.Config.LogChannel = channel.ID
	case "adminlog":
		heimdallr.Config.AdminLogChannel = channel.ID
	case "bot":
		heimdallr.Config.BotChannel = channel.ID
	case "admin":
		heimdallr.Config.AdminChannel = channel.ID
	default:
		return nil
	}
	err = heimdallr.Config.SaveConfig("config.toml")
	if err != nil {
		return err
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}
