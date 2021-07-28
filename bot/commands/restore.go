package commands

import (
	"fmt"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var restoreCommand = command{
	"restore",
	commandRestoreUser,
	"restores an isolated user.",
	[]string{
		"<user>",
	},
	[]string{
		"@username",
		"245207597929480192",
	},
}

//commandRestoreUser restores an isolated user.
func commandRestoreUser(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string), s)
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

	if !isIsolated(infractor) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s is NOT isolated in the first place", user.Mention()))
		return errors.Wrap(err, "sending message failed")
	}

	author, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", userID))
		return errors.Wrap(err, "sending message failed")
	}

	if !heimdallr.IsModOrHigher(author, guild) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Only Moderators can restore an isolated user")
		return errors.Wrap(err, "sending message failed")
	}

	//Add all the other user roles
	roles, err := heimdallr.GetIsolatedUserRoles(infractor.User.ID)
	if err != nil {
		return errors.Wrap(err, "getting roles from the database failed")
	}

	for _, role := range roles {
		if role != heimdallr.Config.ServerBoosterRole {
			if role != "" {
				err = s.GuildMemberRoleAdd(m.GuildID, infractor.User.ID, role)
			}

			if err != nil {

				return errors.Wrap(err, fmt.Sprintf("adding role with ID %s failed", role))
			}
		}
	}
	//remove the muted user from the database
	err = heimdallr.RemoveIsolatedUser(infractor.User.ID)
	if err != nil {
		return errors.Wrap(err, "Removing the Muted user from the database failed")
	}

	//Remove the muted role
	err = s.GuildMemberRoleRemove(guildID, userID, heimdallr.Config.IsolatedRole)
	if err != nil {
		return errors.Wrap(err, "removing isolated role failed")
	}

	if err != nil {
		return errors.Wrap(err, "getting user failed")
	}
	_, err = s.ChannelMessageSendEmbed(heimdallr.Config.AdminLogChannel, &discordgo.MessageEmbed{
		Title: fmt.Sprintf("User was restored by %s.", author.User.Username+"#"+author.User.Discriminator),
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
		Color: 0xFFFF00,
	})

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")

	userChannel, err := s.UserChannelCreate(userID)
	if err != nil {
		s.ChannelMessageSend(heimdallr.Config.AdminLogChannel, fmt.Sprintf("%s Does NOT ACCEPT DMs but was sucessfully restored", infractor.Mention()))
		return nil
		// return errors.Wrap(err, "creating private channel failed")
	}
	_, err = s.ChannelMessageSend(userChannel.ID, fmt.Sprintf(
		"You have been restored in %s \n\nYou cannot reply to this message.",
		guild.Name,
	))
	if err != nil {
		return nil
	}

	return errors.Wrap(err, "adding reaction failed")

}
