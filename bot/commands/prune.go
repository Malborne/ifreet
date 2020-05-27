package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	heimdallr "gitlab.com/NorwegianLanguageLearning/heimdallr/bot"
	"time"
)

var pruneCommand = command{
	"prune",
	commandPrune,
	"Kicks unapproved users who joined more than a week ago.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//commandPrune kicks users who have stayed in the server for at least a week without being approved
func commandPrune(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		return err
	}
	for _, member := range guild.Members {
		joinedAt, err := member.JoinedAt.Parse()
		if err != nil {
			return errors.Wrap(err, "parsing joinedAt failed")
		}
		if joinedAt.Before(time.Now().AddDate(0, 0, -7)) && !isApproved(member) {
			err := s.GuildMemberDeleteWithReason(m.GuildID, member.User.ID, "Stayed in the server for at least 7 days without gaining the User role")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
