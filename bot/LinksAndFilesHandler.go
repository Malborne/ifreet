package heimdallr

import (
	"time"

	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//LinksAndFilesHandler checks if someone sends a link and deletes the message if the user is new
func LinksAndFilesHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	guildID := m.GuildID
	guild, err := GetGuild(s, guildID)
	if err != nil {
		LogIfError(s, err)
	}

	author, err := GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", m.Author.ID))
		LogIfError(s, err)
	}

	if IsAdminOrHigher(author, guild) {
		return
	}
	if hasBannedWord(m.Content) {
		_, err = s.ChannelMessageSendEmbed(Config.LogChannel, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("A user attempted to post an inappropriate word"),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Message Author",
					Value: m.Author.Username + "#" + m.Author.Discriminator,
				},
				{
					Name:  "Channel",
					Value: fmt.Sprintf("<#%s>", m.ChannelID),
				},
				{
					Name:  "Message Content",
					Value: m.Content,
				},
			},
			Color: 0xEE0000,
		})
		if err != nil {
			LogIfError(s, errors.Wrap(err, "sending embed failed"))
			return
		}
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		if hasRole(author, Config.UserRole) {
			err = muteUser(s, author, guildID)
			if err != nil {
				LogIfError(s, errors.Wrap(err, "Muting user failed"))

			}
		}
		err = AddInfraction(*author.User, Infraction{Reason: "Attempting to post an inappropriate word", Time: time.Now()})
		if err != nil {
			LogIfError(s, errors.Wrap(err, "Adding infraction failed"))
			return
		}

	}

	if IsModOrHigher(author, guild) {
		return
	}
	joinedAt, err := author.JoinedAt.Parse()
	if err != nil {
		LogIfError(s, err)
	}
	if IsVerified(author) && joinedAt.Before(time.Now().Add(time.Minute*-60)) { //if verified and joined more than an hour ago, just ignore it
		return
	}
	if joinedAt.Before(time.Now().AddDate(0, 0, -1)) { //If they joined the server more than 24 ago, just ignore it
		return
	}
	if len(m.Attachments) > 0 { //sent a file

		_, err = s.ChannelMessageSendEmbed(Config.LogChannel, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("A user attempted to post a file."),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Message Author",
					Value: m.Author.Username + "#" + m.Author.Discriminator,
				},
				{
					Name:  "Channel",
					Value: fmt.Sprintf("<#%s>", m.ChannelID),
				},
			},
			Color: 0xEE0000,
		})
		if err != nil {
			LogIfError(s, errors.Wrap(err, "sending embed failed"))
			return
		}
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s You are NOT allowed to send files yet. Please wait until you are on the server for a longer time.", author.Mention()))
		if err != nil {
			LogIfError(s, errors.Wrap(err, "sending message failed"))
			return
		}

	}
	if len(m.Embeds) > 0 || strings.Contains(strings.ToLower(m.Content), "https://") || strings.Contains(strings.ToLower(m.Content), "http://") { //sent a link

		if strings.Contains(strings.ToLower(m.Content), "youtube.com") || strings.Contains(strings.ToLower(m.Content), "https://youtu.be/") { //Ignore YouTube videos
			return
		}

		_, err = s.ChannelMessageSendEmbed(Config.LogChannel, &discordgo.MessageEmbed{
			Title: fmt.Sprintf("A user attempted to post a link. Beaware of suspicious links don't click them"),
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Message Author",
					Value: m.Author.Username + "#" + m.Author.Discriminator,
				},
				{
					Name:  "Channel",
					Value: fmt.Sprintf("<#%s>", m.ChannelID),
				},
				{
					Name:  "Message Content",
					Value: m.Content,
				},
			},
			Color: 0xEE0000,
		})
		if err != nil {
			LogIfError(s, errors.Wrap(err, "sending embed failed"))
			return
		}
		s.ChannelMessageDelete(m.ChannelID, m.ID)

		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s You are NOT allowed to send links yet. Please wait until you are on the server for a longer time.", author.Mention()))
		if err != nil {
			LogIfError(s, err)
			return
		}

	}

}

func hasBannedWord(content string) bool {
	for _, word := range Config.BannedWords {
		if strings.Contains(strings.ToLower(content), word) {
			return true
		}
	}
	return false
}

func muteUser(s *discordgo.Session, infractor *discordgo.Member, GuildID string) error {

	guild, err := GetGuild(s, GuildID)

	//Add the muted user's roles to the database

	err = AddMutedUser(*infractor.User, time.Now(), getRoleIDs(infractor))
	if err != nil {

		return errors.Wrap(err, "Addin the user to the databaes failed")
	}

	//Remove all the other user roles
	for _, role := range infractor.Roles {
		err = s.GuildMemberRoleRemove(GuildID, infractor.User.ID, role)
		if err != nil {
			return errors.Wrap(err, "removing role failed")
		}

	}
	//Add the muted role
	err = s.GuildMemberRoleAdd(GuildID, infractor.User.ID, Config.MutedRole)
	if err != nil {
		return errors.Wrap(err, "adding user role failed")
	}

	if err != nil {
		return errors.Wrap(err, "getting user failed")
	}
	_, err = s.ChannelMessageSendEmbed(Config.LogChannel, &discordgo.MessageEmbed{
		Title: "User was automatically muted by Ifreet.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Username**",
				Value: infractor.User.Username + "#" + infractor.User.Discriminator,
			},
			{
				Name:  "**User ID**",
				Value: infractor.User.ID,
			},
		},
		Color: 0xEE0000,
	})

	userChannel, err := s.UserChannelCreate(infractor.User.ID)
	if err != nil {
		// s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Does NOT ACCEPT DMs but was successfully muted", infractor.Mention()))
		return nil
		// return errors.Wrap(err, "creating private channel failed")
	}
	_, err = s.ChannelMessageSend(userChannel.ID, fmt.Sprintf(
		"You have been muted for attempting to post an inappropriate word in %s \n\nIf you think there was a mistake, please contact one of the Moderators\n\nYou cannot reply to this message.",
		guild.Name,
	))
	if err != nil {
		// s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s Does NOT ACCEPT DMs but has been muted", infractor.Mention()))
		return nil
		// return errors.Wrap(err, "sending message failed")
	}

	return errors.Wrap(err, "Sending Message failed")
}

func getRoleIDs(m *discordgo.Member) string {
	var roleIDs = ""
	for _, role := range m.Roles {
		roleIDs = roleIDs + role + ","
	}
	return roleIDs
}
