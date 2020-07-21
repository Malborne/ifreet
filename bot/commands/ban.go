package commands

import (
	"fmt"
	"time"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
)

var banCommand = command{
	"ban",
	commandBanUser,
	"Bans a user from the server.",
	[]string{
		"<user> <reason>",
	},
	[]string{
		"@username \"Did something wrong\"",
		"245167597929480192 \"Did something wrong\"",
	},
}

//commandBanUser bans a user from the server.
func commandBanUser(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	userID := getIDFromMaybeMention(args["<user>"].(string), s)
	reason, _ := args.String("<reason>")

	guildID := m.GuildID
	guild, err := heimdallr.GetGuild(s, guildID)
	if err != nil {
		return err
	}

	member, err := heimdallr.GetMember(s, guildID, userID)
	var user *discordgo.User
	if err != nil {
		user, err = s.User(userID)
		if err != nil {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("No user was found with ID %s.", userID))
			return errors.Wrap(err, "sending message failed")
		}
	} else {
		user = member.User
	}
	author, err := heimdallr.GetMember(s, guildID, m.Author.ID)
	if err != nil {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Message Author with ID %s was not found.", userID))
		return errors.Wrap(err, "sending message failed")
	}
	if userID == s.State.User.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to ban myself, silly. 😉")
		return errors.Wrap(err, "sending message failed")
	}

	if heimdallr.IsAdminOrHigher(member, guild) {
		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("You cannot ban the admin. 👎"))
		return errors.Wrap(err, "sending message failed")
	}

	if m.Author.ID == user.ID && userID == "550664345302859786" { // Wasan's ID
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to let you ban yourself, silly. 😉\nI'm looking at you, وسن. I had to make this because of you 😒")
		return errors.Wrap(err, "sending message failed")

	} else if m.Author.ID == user.ID {
		_, err := s.ChannelMessageSend(m.ChannelID, "I'm not going to let you ban yourself, silly. 😉")
		return errors.Wrap(err, "sending message failed")
	}
	if isOneLowerThanTwo(author, member) {
		// _, _ = s.ChannelMessageSend(heimdallr.Config.AdminLogChannel, fmt.Sprintf("%s the infractor has rank of: %s and %s the author has rank of: %s", infractor.Mention(), getHighestRole(infractor), author.Mention(), getHighestRole(author)))
		_, err := s.ChannelMessageSend(m.ChannelID, "You cannot ban a user that has the same or a role higher than you")
		return errors.Wrap(err, "sending message failed")
	}
	err = s.GuildBanCreateWithReason(guildID, userID, reason, 0)
	if err != nil {
		return errors.Wrap(err, "ban failed")
	}
	err = heimdallr.AddInfraction(*user, heimdallr.Infraction{
		Reason: fmt.Sprintf("Received a ban with reason: %s", reason),
		Time:   time.Now(),
	})
	if err != nil {
		return err
	}
	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	return errors.Wrap(err, "adding reaction failed")
}
