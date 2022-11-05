package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Malborne/ifreet/tree/master/bot/commands"
	"github.com/Malborne/ifreet/tree/master/bot/version"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/pkg/errors"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

func init() {
	fmt.Printf("Ifreet version: %s, commit: %s\n", version.VERSION, version.COMMIT)

	err := heimdallr.Config.LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("%+v\n", errors.WithMessage(err, "failed to load config"))
	}

	flag.StringVar(&token, "token", heimdallr.Config.Token, "The bot token that Ifreet should use.")
	flag.Parse()

	if token == "" {
		log.Fatalln("Flag '-token' or token in config file not set. This program cannot be used without a valid token.")
	}
}

func main() {
	// err := heimdallr.OpenDb("heimdallr.db")
	err := heimdallr.OpenDb(os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "failed to create bot"))
	}

	dg.AddHandler(commands.CommandHandler)
	dg.AddHandler(commands.ReactionApprove)
	dg.AddHandler(commands.ReactionPrompt)
	// dg.AddHandler(heimdallr.UserJoinHandler)
	// dg.AddHandler(heimdallr.UserLeaveHandler)
	dg.AddHandler(heimdallr.MemberBanAddHandler)
	dg.AddHandler(heimdallr.MessageHandler)
	dg.AddHandler(heimdallr.OnDeleteHandler)
	dg.AddHandler(heimdallr.NewMemberJoinHandler)
	dg.AddHandler(heimdallr.NewMemberLeaveHandler)
	dg.Identify.Intents = discordgo.IntentsAll
	err = dg.Open()
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "failed to open session"))
	}

	// go heimdallr.CheckPermissions(dg)

	//Checks if any users are still isolated and restores them if/when their duration has expired
	isolatedUsers, err := heimdallr.GetAllIsolatedUsers()
	if err != nil {
		heimdallr.LogIfError(dg, errors.Wrap(err, "getting isolated Users failed"))
	}
	guildID := heimdallr.Config.GuildID
	for _, isoUser := range isolatedUsers {
		member, _ := heimdallr.GetMember(dg, guildID, isoUser.UserID)
		currentTime := time.Now()
		if member != nil {
			if currentTime.After(isoUser.EndTime) {
				commands.RestoreUser(dg, member, guildID)
			} else {
				time.AfterFunc(isoUser.EndTime.Sub(currentTime), func() { commands.RestoreUser(dg, member, guildID) })

			}
		}
	}

	//Check if more than 6 days have passed since the new welcome channel was created.
	newChannels, newUsers, err := heimdallr.GetAllnewChannelsWithUsers()
	if err != nil {
		heimdallr.LogIfError(dg, errors.Wrap(err, "getting new channels failed"))

	}
	for index, channelID := range newChannels {
		newChannel, _ := dg.Channel(channelID)
		mewMember, err := heimdallr.GetMember(dg, heimdallr.Config.GuildID, newUsers[index])
		if err != nil { //The user is no longer a member of the server
			continue
		}
		messages, err := dg.ChannelMessages(channelID, 100, "", "", "")
		if err != nil {
			heimdallr.LogIfError(dg, errors.Wrap(err, "getting channel messages failed"))

			if strings.Contains(err.Error(), "HTTP 404 Not Found") { //Channel has probably been manually deleted
				err = heimdallr.RemoveNewChannel(mewMember.User.ID) // Remove the channel from the database
				heimdallr.LogIfError(dg, errors.Wrap(err, "unable to remove the channel from the database"))
			}
		}

		alreadyWarned := false
		for _, message := range messages { //This for loop check if the member has already bee warned
			if message.Author.Bot && strings.Contains(message.Content, "You are an unapproved member of Quran Learning Center") { //User has been warned before already
				alreadyWarned = true
				sentTime, _ := heimdallr.GetDateTimeFromID(message.ID)
				currentTime := time.Now()

				hoursPassed := currentTime.Sub(sentTime).Hours()
				if hoursPassed >= 24 { //User has stayed for more than 1 day on the server after being warned without getting approved
					heimdallr.KickMember(dg, mewMember) // kick the user
					break
				} else if hoursPassed < 24 { // less than 24 hours has passed since the member was warned, wait until they pass and then kick

					time.AfterFunc(sentTime.Add(24*time.Hour).Sub(currentTime), func() { heimdallr.KickMember(dg, mewMember) })
					break
				}
			}
		}
		if !alreadyWarned {
			for _, message := range messages { //This for loop is run if the user has not received a warning before
				if message.Author.Bot && strings.Contains(message.Content, "a server dedicated to aiding its members") {

					sentTime, _ := heimdallr.GetDateTimeFromID(message.ID)
					currentTime := time.Now()
					hoursPassed := currentTime.Sub(sentTime).Hours()

					if hoursPassed >= 144 { //User has stayed for more than 6 days on the server
						heimdallr.SendUnapprovedMessage(dg, newChannel, newUsers[index])
						time.AfterFunc(24*time.Hour, func() { heimdallr.KickMember(dg, mewMember) }) //called after one more day

					} else if hoursPassed < 144 { // Member has stayed for less than 6 days on the server

						time.AfterFunc(sentTime.Add(144*time.Hour).Sub(currentTime), func() {
							heimdallr.SendUnapprovedMessage(dg, newChannel, newUsers[index])
							time.AfterFunc(24*time.Hour, func() { heimdallr.KickMember(dg, mewMember) }) //called after one more day
						})
						break
					}
				}
			}
		}
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	err = dg.Close()
	if err != nil {
		log.Printf("%+v\n", errors.Wrap(err, "failed to close session"))
	}

	err = heimdallr.CloseDb()
	if err != nil {
		log.Printf("%+v\n", err)
	}
}
