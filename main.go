package main

import (
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.com/NorwegianLanguageLearning/heimdallr/bot"
	"gitlab.com/NorwegianLanguageLearning/heimdallr/bot/commands"
	"gitlab.com/NorwegianLanguageLearning/heimdallr/bot/version"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	token string
)

func init() {
	fmt.Printf("Heimdallr version: %s, commit: %s\n", version.VERSION, version.COMMIT)

	err := heimdallr.Config.LoadConfig("config.toml")
	if err != nil {
		log.Fatalf("%+v\n", errors.WithMessage(err, "failed to load config"))
	}

	flag.StringVar(&token, "token", heimdallr.Config.Token, "The bot token that Heimdallr should use.")
	flag.Parse()

	if token == "" {
		log.Fatalln("Flag '-token' or token in config file not set. This program cannot be used without a valid token.")
	}
}

func main() {
	err := heimdallr.OpenDb("heimdallr.db")
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "failed to create bot"))
	}

	dg.AddHandler(commands.CommandHandler)
	dg.AddHandler(commands.ReactionApprove)
	dg.AddHandler(heimdallr.UserJoinHandler)
	dg.AddHandler(heimdallr.UserLeaveHandler)
	dg.AddHandler(heimdallr.MemberBanAddHandler)

	err = dg.Open()
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "failed to open session"))
	}

	go heimdallr.CheckPermissions(dg)

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
