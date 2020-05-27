package heimdallr

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/pkg/errors"
	"log"
)

//LogIfError logs an error to the admin log channel and the console
func LogIfError(s *discordgo.Session, err error) {
	if err == nil {
		return
	}
	errorMsg := fmt.Sprintf("%+v", err)
	log.Println(errorMsg)
	_, err = s.ChannelMessageSend(Config.AdminLogChannel, fmt.Sprintf("```\n%s```", errorMsg))
	if err != nil {
		errorMsg := fmt.Sprintf("%+v", errors.Wrap(err, "logging error in admin log channel failed"))
		log.Println(errorMsg)
	}
}

//LogMessage logs a message to the admin log channel and the console
func LogMessage(s *discordgo.Session, message string) {
	log.Println(message)
	_, err := s.ChannelMessageSend(Config.AdminLogChannel, message)
	if err != nil {
		errorMsg := fmt.Sprintf("%+v", errors.Wrap(err, "logging error in admin log channel failed"))
		log.Println(errorMsg)
	}
}
