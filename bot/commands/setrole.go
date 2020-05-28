package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet-bot/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var setRoleCommand = command{
	"setrole",
	commandSetRole,
	"Sets a server role to an administrative role.",
	[]string{
		"(mod | supermod | admin | user) <role>",
	},
	[]string{
		"mod Moderators",
		"supermod",
		"admin @admin",
	},
}

//commandSetRole sets the various roles in a server to mod, supermod, or admin
func commandSetRole(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	roleNameOrID, _ := args.String("<role>")
	roleTypes := []string{"mod", "supermod", "admin", "user"}
	var roleType string
	for _, potentialRoleType := range roleTypes {
		if present, _ := args.Bool(potentialRoleType); present {
			roleType = potentialRoleType
			break
		}
	}

	guildID := m.GuildID
	_, err := s.State.Role(guildID, roleNameOrID)
	var roleID string
	if err == nil {
		roleID = roleNameOrID
	} else {
		guild, err := heimdallr.GetGuild(s, guildID)
		if err != nil {
			return err
		}
		for _, role := range guild.Roles {
			if role.Name == roleNameOrID {
				roleID = role.ID
			}
		}
		if roleID == "" {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No role was found with ID or name %s.", roleNameOrID))
			return err
		}
	}
	switch roleType {
	case "mod":
		heimdallr.Config.ModRole = roleID
	case "supermod":
		heimdallr.Config.SuperModRole = roleID
	case "admin":
		heimdallr.Config.AdminRole = roleID
	case "user":
		heimdallr.Config.UserRole = roleID
	}
	err = heimdallr.Config.SaveConfig("config.toml")
	if err != nil {
		return err
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}
