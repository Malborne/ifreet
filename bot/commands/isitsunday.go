package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	"time"
)

var isItSundayCommand = command{
	"isitsunday",
	commandIsItSunday,
	"Veg's personal favourite.",
	[]string{
		"",
	},
	[]string{
		"",
	},
}

//commandIsItSunday tells a user if it's Sunday or not.
func commandIsItSunday(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	location, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		return errors.Wrap(err, "loading location Europe/Oslo failed")
	}
	t := time.Now().In(location)

	if t.Weekday() == time.Sunday {
		_, err := s.ChannelMessageSend(m.ChannelID, "***YES!***")
		return errors.Wrap(err, "sending message failed")
	}
	daysLeft := int(7 - t.Weekday())
	var pluralSuffix string
	if daysLeft > 1 {
		pluralSuffix = "s"
	} else {
		pluralSuffix = ""
	}
	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("***no, still %d day%s left :(***", daysLeft, pluralSuffix))
	return errors.Wrap(err, "sending message failed")
}
