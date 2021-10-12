package commands

import (
	"fmt"
	"strconv"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var ruleCommand = command{
	"rule",
	commandRule,
	"Quotes a rule from the server rules channel",
	[]string{
		"<number>",
	},
	[]string{
		"10",
		"100",
	},
}

//commandRule Quotes a rule
func commandRule(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	number, err := args.Int("<number>")

	if err != nil || number == 0 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Incorrect use of command. Type the rule number you wish to quote"))
		return errors.Wrap(err, "quoting rule failed")
	}

	if number > 16 {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Incorrect use of command. Please enter a value between 1 and 16"))
		return errors.Wrap(err, "clearing message failed")
	}

	for _, rule := range heimdallr.Config.Rules {
		RuleNumber, _ := strconv.Atoi(rule.Number)
		if RuleNumber == number {
			//Quote the rule
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Title: fmt.Sprintf("Rule No. %s:", rule.Number),

				Footer: &discordgo.MessageEmbedFooter{
					Text: fmt.Sprintf("(quoted by: %s)", m.Author.Username),
				},
				Description: fmt.Sprintf("%s\n<#734707679116918836>", rule.Text),
				Color:       0xFFFF00,
			})
			if err != nil {
				return errors.Wrap(err, "sending embed failed")
			}
			break
		}
	}

	return nil
}
