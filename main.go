package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
	dg.AddHandler(heimdallr.UserJoinHandler)
	dg.AddHandler(heimdallr.UserLeaveHandler)
	dg.AddHandler(heimdallr.MemberBanAddHandler)
	dg.AddHandler(heimdallr.MessageHandler)
	dg.AddHandler(heimdallr.OnDeleteHandler)
	// dg.AddHandler(heimdallr.NewMemberJoinHandler)

	dg.Identify.Intents = discordgo.IntentsAll
	err = dg.Open()
	if err != nil {
		log.Fatalf("%+v\n", errors.Wrap(err, "failed to open session"))
	}

	// go heimdallr.CheckPermissions(dg)

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
