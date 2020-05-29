package heimdallr

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//MemberBanAddHandler sends a message to the admin log channel that a user was banned.
func MemberBanAddHandler(s *discordgo.Session, e *discordgo.GuildBanAdd) {

	bannedUser := e.User
	guildID := e.GuildID

	guildBans, err := s.GuildBans(guildID)
	if err != nil {
		LogIfError(s, errors.Wrap(err, "getting guild bans failed"))
		return
	}

	var guildBan *discordgo.GuildBan
	for _, ban := range guildBans {
		if ban.User.ID == bannedUser.ID {
			guildBan = ban
		}
	}

	if guildBan == nil {
		LogIfError(s, errors.Wrap(err, "finding ban failed"))
		return
	}

	var reason string
	if guildBan.Reason != "" {
		reason = guildBan.Reason
	} else {
		reason = "N/A"
	}

	_, err = s.ChannelMessageSendEmbed(Config.AdminLogChannel, &discordgo.MessageEmbed{
		Title: "User was banned.",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "**Username**",
				Value: bannedUser.Username + "#" + bannedUser.Discriminator,
			},
			{
				Name:  "**User ID**",
				Value: bannedUser.ID,
			},
			{
				Name:  "**Reason**",
				Value: reason,
			},
		},
		Color: 0xEE0000,
	})
	if err != nil {
		LogIfError(s, errors.Wrap(err, "sending embed failed"))
		return
	}

	_, err = s.ChannelMessageSend(Config.AdminLogChannel, "**Here is a report you can submit to other servers who should be wary of this user:**\n\n"+
		"```\nUser banned from Quran Learning Center:"+
		"\n**·User:** "+bannedUser.Username+"#"+bannedUser.Discriminator+
		"\n**·User ID:** "+bannedUser.ID+
		"\n**Reason: **\n"+guildBan.Reason+"\n```"+
		"\nYou may also want to include more information as to why this user was banned.")
	if err != nil {
		LogIfError(s, errors.Wrap(err, "sending message failed"))
	}
}
