package commands

import (
	"time"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
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
	members, err := s.GuildMembers(m.GuildID, "", 1000)
	if err != nil {
		return err
	}
	for _, member := range members {
		joinedAt, err := member.JoinedAt.Parse()
		if err != nil {
			return errors.Wrap(err, "parsing joinedAt failed")
		}
		if joinedAt.Before(time.Now().AddDate(0, 0, -7)) && !isApproved(member) && !member.User.Bot && !hasRole(member, heimdallr.Config.ServerBoosterRole) {
			err := s.GuildMemberDeleteWithReason(m.GuildID, member.User.ID, "Stayed in the server for at least 7 days without gaining the User role")
			if err != nil {
				return err
			}
		}
	}
	return nil
}
