package heimdallr

import (
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
)

//GetMember gets a member from the state or from the API if it's not available in the state
func GetMember(s *discordgo.Session, guildID, userID string) (*discordgo.Member, error) {
	member, err := s.State.Member(guildID, userID)
	if err != nil {
		member, err = s.GuildMember(guildID, userID)
		if err != nil {
			err = errors.Wrap(err, "getting member failed")
		}
	}
	return member, err
}

//GetGuild gets a guild from the state or from the API if it's not available in the state
func GetGuild(s *discordgo.Session, guildID string) (*discordgo.Guild, error) {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		guild, err = s.Guild(guildID)
		if err != nil {
			err = errors.Wrap(err, "getting guild failed")
		}
	}
	return guild, nil
}

//GetMessage gets a message from the state or from the API if it's not available in the state
func GetMessage(s *discordgo.Session, channelID string, messageID string) (*discordgo.Message, error) {
	message, err := s.State.Message(channelID, messageID)
	if err != nil {
		message, err = s.ChannelMessage(channelID, messageID)
		if err != nil {
			err = errors.Wrap(err, "getting message failed")
		}
	}
	return message, nil
}
