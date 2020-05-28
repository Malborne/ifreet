package commands

import (
	"fmt"
	"strings"

	heimdallr "github.com/Malborne/ifreet-bot/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var helpCommand = command{
	"help",
	commandHelp,
	"Displays help message for all or a specific command.",
	[]string{
		"[<command>]",
	},
	[]string{
		"",
		"role",
	},
}

//commandHelp sends information about the bot's command usage
func commandHelp(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	if args["<command>"] != nil {
		command, _ := args.String("<command>")
		return help(s, m, command)
	}
	return helpAll(s, m)
}

func help(s *discordgo.Session, m *discordgo.MessageCreate, commandName string) error {
	command, present := nameToCommand[commandName]
	if !present {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`%s` is not a valid command.", commandName))
		return errors.Wrap(err, "sending message failed")
	}
	_, err := s.ChannelMessageSendEmbed(m.ChannelID, command.getFullHelpEmbedWithExamples(true, -1))
	return errors.Wrap(err, "sending embed failed")
}

func helpAll(s *discordgo.Session, m *discordgo.MessageCreate) error {
	fields := []*discordgo.MessageEmbedField{{
		Name:  "User commands",
		Value: strings.Join(getHelpMessages(userCommands), "\n"),
	}}
	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		return err
	}
	member, err := heimdallr.GetMember(s, m.GuildID, m.Author.ID)
	if err != nil {
		return err
	}
	if heimdallr.IsModOrHigher(member, guild) {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Moderator commands",
			Value: strings.Join(getHelpMessages(moderatorCommands), "\n"),
		})
	}
	if heimdallr.IsAdminOrHigher(member, guild) {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Admin commands",
			Value: strings.Join(getHelpMessages(adminCommands), "\n"),
		})
	}
	if heimdallr.IsOwner(member, guild) {
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  "Owner commands",
			Value: strings.Join(getHelpMessages(ownerCommands), "\n"),
		})
	}
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Title:  "**Ifreet bot**",
		Fields: fields,
	})
	return errors.Wrap(err, "sending embed failed")
}

func getHelpMessages(commands []command) []string {
	var helpMessages []string
	for _, command := range commands {
		helpMessages = append(helpMessages, command.getShortHelpMessage())
	}
	return helpMessages
}
