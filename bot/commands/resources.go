package commands

import (
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/pkg/errors"
	heimdallr "gitlab.com/NorwegianLanguageLearning/heimdallr/bot"
	"strings"
)

var searchResourcesCommand = command{
	"resources",
	commandResources,
	"Get and search in resources.",
	[]string{
		"search <search-term>",
		"get <name-or-id>",
		"add <name> <content> [<tags>]...",
	},
	[]string{
		"search verbs",
		"get 1",
		"add Ordbok https://ordbok.uib.no/",
		"add Ordbok https://ordbok.uib.no/ dictionary bokmål nynorsk",
	},
}

//commandViewInfractions lists a user's infractions
func commandResources(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	if search, _ := args.Bool("search"); search {
		return resourcesSearch(s, m, args)
	} else if get, _ := args.Bool("get"); get {
		return resourcesGet(s, m, args)
	} else if add, _ := args.Bool("add"); add {
		return resourcesAdd(s, m, args)
	}
	return nil
}

func resourcesSearch(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	searchTerm, _ := args.String("<search-term>")
	resources, err := heimdallr.SearchResources(strings.Split(searchTerm, " "))
	if err != nil {
		return err
	}
	numResources := len(resources)
	if numResources == 0 {
		_, err = s.ChannelMessageSend(m.ChannelID, "No results found.")
		return errors.Wrap(err, "sending message failed")
	}
	// There is a limit of 25 fields in an embed.
	if numResources > 25 {
		resources = resources[numResources-25:]
	}
	var fields []*discordgo.MessageEmbedField
	for _, resource := range resources {
		content := resource.Content
		if len(resource.Content) > 200 && numResources > 1 {
			content = fmt.Sprintf("%s...", content[:200])
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("%s (id: %d)", resource.Name, resource.ID),
			Value: content,
		})
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Fields: fields,
	})
	return errors.Wrap(err, "sending embed failed")
}

func resourcesGet(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	id, err := args.Int("<name-or-id>")
	var resource *heimdallr.Resource
	if err == nil {
		resource, err = heimdallr.GetResourceByID(id)
	} else {
		name, _ := args.String("<name-or-id>")
		resource, err = heimdallr.GetResourceByName(name)
	}

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			_, err = s.ChannelMessageSend(m.ChannelID, "No resource with that name or id.")
			return errors.Wrap(err, "sending message failed")
		}
		return err
	}

	var fields []*discordgo.MessageEmbedField
	fields = append(fields, &discordgo.MessageEmbedField{
		Name:  fmt.Sprintf("%s (id: %d)", resource.Name, resource.ID),
		Value: resource.Content,
	})
	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
		Fields: fields,
	})
	return errors.Wrap(err, "sending embed failed")
}

func resourcesAdd(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	name, _ := args.String("<name>")
	content, _ := args.String("<content>")
	tags := args["<tags>"].([]string)

	guild, err := heimdallr.GetGuild(s, m.GuildID)
	if err != nil {
		return err
	}

	member, err := heimdallr.GetMember(s, m.GuildID, m.Author.ID)
	if err != nil {
		return err
	}

	if !heimdallr.IsModOrHigher(member, guild) {
		_, err := s.ChannelMessageSend(m.ChannelID, "Only moderators can add resources.")
		return errors.Wrap(err, "sending message failed")
	}

	resource := heimdallr.Resource{
		Name:    name,
		Content: content,
		Tags:    tags,
	}
	id, err := heimdallr.AddResource(resource)
	if err != nil {
		return err
	}

	_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Added resource with id %d", id))
	return errors.Wrap(err, "sending message failed")
}
