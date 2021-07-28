package commands

import (
	"fmt"
	"strings"
	"time"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var isolateCommand = command{
	"isolate",
	commandIsolateme,
	"Isolates you from the server for a specified duration of time.",
	[]string{
		"<duration> <unit>",
	},
	[]string{
		"10 seconds",
		"3 minutes",
		"5 hours",
		"7 days",
	},
}

//commandIsolateme isolates a user for a specific period of time
func commandIsolateme(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	duration, _ := args.Int("<duration>")
	unit, _ := args.String("<unit>")
	// member := m.Member
	// var timer *time.Timer;
	acceptedUnits := []string{"s", "second", "seconds", "sec", "secs", "minute", "minutes", "min", "mins", "m", "h", "hour", "hours", "hr", "hrs", "day", "days", "d"}
	index, isUnitAccepted := stringInSlice(strings.ToLower(unit), acceptedUnits)

	if !isUnitAccepted {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Invalid Unit, please enter a correct unit either `seconds`, `minutes`, `hours`, or `days`"))
		return errors.Wrap(err, "Isolating user failed")
	} else {
		if index < 5 {
			unit = "seconds"
		} else if index < 10 {
			unit = "minutes"
		} else if index < 15 {
			unit = "hours"
		} else {
			unit = "days"
		}
	}

	guildID := m.GuildID

	guild, err := heimdallr.GetGuild(s, guildID)
	if err != nil {
		return err
	}

	member, err := heimdallr.GetMember(s, guildID, m.Member.User.ID)

	if heimdallr.IsAdminOrHigher(member, guild) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The Admin cannot be isolated."))
		return errors.Wrap(err, "Isolating the admin cannot be done")
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The duration of isolation is %d %s.", duration, unit))
	startTime := time.Now()
	var endTime time.Time
	if unit == "seconds" {
		endTime = startTime.Add(time.Duration(duration) * time.Second)

	} else if unit == "minutes" {
		endTime = startTime.Add(time.Duration(duration) * time.Minute)
	} else if unit == "hours" {
		endTime = startTime.Add(time.Duration(duration) * time.Hour)
	} else if unit == "days" {
		endTime = startTime.Add(time.Duration(24*duration) * time.Hour)
	} else {
		s.ChannelMessageSend(heimdallr.Config.AdminChannel, fmt.Sprintf("The unit used for the duration is incorrect"))
		return err
	}

	//Add the isolated user's roles to the database
	err = heimdallr.AddIsolatedUser(*member.User, startTime, endTime, getRoleIDs(member))
	if err != nil {
		return err
	}

	//Remove all the other user roles
	for _, role := range member.Roles {
		if role != heimdallr.Config.ServerBoosterRole {
			err = s.GuildMemberRoleRemove(m.GuildID, member.User.ID, role)

			if err != nil {
				heimdallr.LogIfError(s, err)
			}
		}
	}
	//Add the Isolated role
	err = s.GuildMemberRoleAdd(guildID, member.User.ID, heimdallr.Config.IsolatedRole)
	if err != nil {
		return errors.Wrap(err, "adding user role failed")
	}

	if err != nil {
		return errors.Wrap(err, "getting user failed")
	}
	_, err = s.ChannelMessageSendEmbed(heimdallr.Config.LogChannel, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("User isolated himself for %d %s", duration, unit),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Username**",
				Value: member.User.Username + "#" + member.User.Discriminator,
			},
			{
				Name:  "**User ID**",
				Value: member.User.ID,
			},
		},
		Color: 0xEE0000,
	})

	userChannel, err := s.UserChannelCreate(member.User.ID)
	if err != nil {
		return nil
	}
	_, err = s.ChannelMessageSend(userChannel.ID, fmt.Sprintf(
		"You have been isolated in %s for for %d %s\nYou will automatically be returned to the server after the duration expires. If you would like to return before that, please DM one of the moderators.\n\nYou cannot reply to this message.", guild.Name, duration, unit))
	if err != nil {
		return nil
	}

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	if err != nil {
		return errors.Wrap(err, "adding reaction failed")
	}

	//Restore the user after the timer expires

	time.AfterFunc(endTime.Sub(startTime), func() { restoreUser(s, member.User, guildID) })

	// if unit == "seconds" {
	// 	time.AfterFunc(time.Duration(duration)*time.Second, func() { restoreUser(s, member.User, guildID) })
	// } else if unit == "minutes" {
	// 	time.AfterFunc(time.Duration(duration)*time.Minute, func() { restoreUser(s, member.User, guildID) })

	// } else if unit == "hours" {
	// 	time.AfterFunc(time.Duration(duration)*time.Hour, func() { restoreUser(s, member.User, guildID) })

	// } else if unit == "days" {
	// 	time.AfterFunc(time.Duration(24*duration)*time.Hour, func() { restoreUser(s, member.User, guildID) })

	// }
	return nil
}

func stringInSlice(a string, list []string) (int, bool) {
	for index, b := range list {
		if b == a {

			return index, true
		}
	}
	return -1, false
}

func restoreUser(s *discordgo.Session, user *discordgo.User, guildID string) {
	//Add all the other user roles
	roles, err := heimdallr.GetIsolatedUserRoles(user.ID)
	if err != nil {
		heimdallr.LogIfError(s, errors.Wrap(err, "getting roles from the database failed"))
	}

	for _, role := range roles {
		if role != heimdallr.Config.ServerBoosterRole {
			if role != "" {
				err = s.GuildMemberRoleAdd(guildID, user.ID, role)
			}

			if err != nil {
				// s.ChannelMessageSend(heimdallr.Config.AdminLogChannel, fmt.Sprintf("No role with ID %s found", role))

				heimdallr.LogIfError(s, errors.Wrap(err, fmt.Sprintf("adding role with ID %s failed", role)))
			}
		}
	}
	//remove the isolated user from the database
	err = heimdallr.RemoveIsolatedUser(user.ID)
	if err != nil {
		heimdallr.LogIfError(s, errors.Wrap(err, "Removing the Isolated user from the database failed"))

	}

	//Remove the isolated role
	err = s.GuildMemberRoleRemove(guildID, user.ID, heimdallr.Config.IsolatedRole)
	if err != nil {
		heimdallr.LogIfError(s, errors.Wrap(err, "removing isolated role failed"))

	}

	_, err = s.ChannelMessageSendEmbed(heimdallr.Config.AdminLogChannel, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("%s was automatically restored to the server.", user.Username+"#"+user.Discriminator),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Username**",
				Value: user.Username + "#" + user.Discriminator,
			},
			{
				Name:  "**User ID**",
				Value: user.ID,
			},
		},
		Color: 0xEE0000,
	})

}
