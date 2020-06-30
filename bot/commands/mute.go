package commands

import (
	"fmt"
	"time"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var muteCommand = command{
	"mute",
	commandMuteUser,
	"mutes a user.",
	[]string{
		"<user>",
	},
	[]string{
		"@username",
		// "@username 3 minutes",
		"245207597929480192",
	},
}

//commandMuteUser mues another user
func commandMuteUser(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string))
	// number, _ := args.Int("<no>")
	var user *discordgo.User

	guildID := m.GuildID

	guild, err := heimdallr.GetGuild(s, guildID)
	if err != nil {
		return err
	}

	infractor, err := heimdallr.GetMember(s, guildID, userID)

	if err != nil {
		user, err = s.User(userID)
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No user was found with ID %s.", userID))
			return errors.Wrap(err, "sending message failed")
		}
	} else {
		user = infractor.User
	}

	author, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", userID))
		return errors.Wrap(err, "sending message failed")
	}

	if userID == s.State.User.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I'm not going to mute myself, silly. 😉"))
		return errors.Wrap(err, "sending message failed")
	}

	if heimdallr.IsAdminOrHigher(infractor, guild) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You cannot mute the admin. 👎"))
		return errors.Wrap(err, "sending message failed")
	}

	if m.Author.ID == user.ID && userID == "550664345302859786" { // Wasan's ID
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to let you mute yourself, silly. 😉 I'm looking at you, وسن. I had to make this because of you 😒")
		return errors.Wrap(err, "sending message failed")
	} else if m.Author.ID == user.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to let you mute yourself, silly. 😉")
		return errors.Wrap(err, "sending message failed")
	}

	if isOneLowerThanTwo(author, infractor) {
		// _, _ = s.ChannelMessageSend(heimdallr.Config.AdminLogChannel, fmt.Sprintf("%s the infractor has rank of: %s and %s the author has rank of: %s", infractor.Mention(), getHighestRole(infractor), author.Mention(), getHighestRole(author)))
		_, err := s.ChannelMessageSend(m.ChannelID, "You cannot mute a user that has the same or a role higher than you")
		return errors.Wrap(err, "sending message failed")
	}

	//Add the muted user's roles to the database
	err = heimdallr.AddMutedUser(*infractor.User, time.Now(), getRoleIDs(infractor))
	if err != nil {
		return err
	}

	//Remove all the other user roles
	for _, role := range infractor.Roles {
		if role != heimdallr.Config.ServerBoosterRole {
			err = s.GuildMemberRoleRemove(m.GuildID, infractor.User.ID, role)

			if err != nil {
				return errors.Wrap(err, "removing role failed")
			}
		}
	}
	//Add the muted role
	err = s.GuildMemberRoleAdd(guildID, userID, heimdallr.Config.MutedRole)
	if err != nil {
		return errors.Wrap(err, "adding user role failed")
	}

	if err != nil {
		return errors.Wrap(err, "getting user failed")
	}
	_, err = s.ChannelMessageSendEmbed(heimdallr.Config.AdminLogChannel, &discordgo.MessageEmbed{
		Title: "User was muted.",
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

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	userChannel, err := s.UserChannelCreate(userID)
	if err != nil {
		// s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Does NOT ACCEPT DMs but was successfully muted", infractor.Mention()))
		return nil
		// return errors.Wrap(err, "creating private channel failed")
	}
	_, err = s.ChannelMessageSend(userChannel.ID, fmt.Sprintf(
		"You have been muted in %s \n\nYou cannot reply to this message.",
		guild.Name,
	))
	if err != nil {
		// s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Does NOT ACCEPT DMs but has been muted", infractor.Mention()))
		return nil
		// return errors.Wrap(err, "sending message failed")
	}

	return errors.Wrap(err, "adding reaction failed")

}

//getRoleIDs returns the IDs of the roles of a given member
func getRoleIDs(m *discordgo.Member) string {
	var roleIDs = ""
	for _, role := range m.Roles {
		if role != heimdallr.Config.ServerBoosterRole {
			roleIDs = roleIDs + role + ","
		}
	}
	return roleIDs
}
