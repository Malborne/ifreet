package commands

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	heimdallr "github.com/Malborne/ifreet/tree/master/bot"
	"github.com/bwmarrin/discordgo"
	"github.com/docopt/docopt-go"
	"github.com/google/shlex"
	"github.com/pkg/errors"
)

var parser = &docopt.Parser{
	SkipHelpFlags: true,
	HelpHandler:   docopt.NoHelpHandler,
}

type handler func(*discordgo.Session, *discordgo.MessageCreate, docopt.Opts) error

type command struct {
	Name        string
	Handler     handler
	Description string
	Usages      []string
	Examples    []string
}

func (command *command) getCombinedUsage() string {
	var format string
	if len(command.Usages) > 1 {
		format = "(%s)"
	} else {
		format = "%s"
	}
	return fmt.Sprintf(format, strings.Join(command.Usages, " | "))
}

func (command *command) getFullCommand() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", command.Name, command.getCombinedUsage()))
}

func (command *command) getShortHelpMessage() string {
	return fmt.Sprintf("`%s` %s", command.getFullCommand(), command.Description)
}

func (command *command) getFullExamples() []string {
	var examples []string
	for _, example := range command.Examples {
		exampleString := strings.TrimSpace(fmt.Sprintf("%s%s %s", getCommandPrefix(), command.Name, example))
		examples = append(examples, fmt.Sprintf("`%s`", exampleString))
	}
	return examples
}

func (command *command) getFullHelpMessage() string {
	var usages []string
	for _, usage := range command.Usages {
		usages = append(usages, fmt.Sprintf("  %s  %s", command.Name, usage))
	}
	return fmt.Sprintf("%s\n\nUsage:\n%s\n", command.Description, strings.Join(usages, "\n"))
}

func (command *command) getFullHelpEmbedWithExamples(includeDescription bool, color int) *discordgo.MessageEmbed {
	var usages []string
	for _, usage := range command.Usages {
		usageString := strings.TrimSpace(fmt.Sprintf("%s %s", command.Name, usage))
		usages = append(usages, fmt.Sprintf("`%s`", usageString))
	}
	embed := discordgo.MessageEmbed{
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Usage",
				Value: strings.Join(usages, "\n"),
			},
			{
				Name:  "Examples",
				Value: strings.Join(command.getFullExamples(), "\n"),
			},
		},
	}
	if includeDescription {
		embed.Author = &discordgo.MessageEmbedAuthor{Name: command.Description}
	}
	if color != -1 {
		embed.Color = color
	}
	return &embed
}

func (command *command) handle(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
	return command.Handler(s, m, args)
}

func (command *command) parse(args []string) (docopt.Opts, error) {
	return parser.ParseArgs(command.getFullHelpMessage(), args, "")
}

var userCommands []command
var helperCommands []command
var trialModeratorCommands []command
var moderatorCommands []command
var superModeratorCommands []command
var adminCommands []command
var ownerCommands []command
var nameToCommand = map[string]command{}

func init() {
	userCommands = []command{
		helpCommand,
		quoteCommand,
		tajweedLessonCommand,
		arabicLessonCommand,
		roleCommand,
		// versionCommand,
		lessonsCommand,
		arabiclessonsCommand,
		isolateCommand,
		// searchResourcesCommand,
		getSheetCommand,
	}

	helperCommands = []command{
		addStudentCommand,
		removeStudentCommand,
		mystudentsCommand,
	}

	trialModeratorCommands = []command{
		warnCommand,
		infractionsCommand,
		approveCommand,
		verifyCommand,
		muteCommand,
		unmuteCommand,
	}
	moderatorCommands = []command{
		kickCommand,
		banCommand,
	}

	superModeratorCommands = []command{
		clearCommand,
		clearFromCommand,
	}

	adminCommands = []command{
		welcomeMessageCommand,
		approvalMessageCommand,
		pruneCommand,
		DMUnapprovedCommand,
		DMUnverifiedCommand,
		sayCommand,
		removeInfractionCommand,
		// castSihrCommand,
	}

	ownerCommands = []command{
		setRoleCommand,
		setChannelCommand,
	}

	requireRoleForCommands("helper", helperCommands)
	requireRoleForCommands("trial moderator", trialModeratorCommands)
	requireRoleForCommands("moderator", moderatorCommands)
	requireRoleForCommands("moderator", moderatorCommands)
	requireRoleForCommands("admin", adminCommands)
	requireRoleForCommands("admin", ownerCommands)
	requireRoleForCommands("supermoderator", superModeratorCommands)

	var commands []command
	commands = append(commands, userCommands...)
	commands = append(commands, helperCommands...)
	commands = append(commands, trialModeratorCommands...)
	commands = append(commands, moderatorCommands...)
	commands = append(commands, superModeratorCommands...)
	commands = append(commands, adminCommands...)
	commands = append(commands, ownerCommands...)
	for _, command := range commands {
		nameToCommand[command.Name] = command
	}
}

func requireRoleForCommands(role string, commands []command) {
	for i := range commands {
		commands[i] = requireRoleForCommand(role, commands[i])
	}
}

func requireRoleForCommand(role string, originalCommand command) command {
	privilegeChecker := getPrivilegeChecker(role)
	handler := originalCommand.Handler

	originalCommand.Handler = func(s *discordgo.Session, m *discordgo.MessageCreate, args docopt.Opts) error {
		guildID := m.GuildID
		guild, err := heimdallr.GetGuild(s, guildID)
		if err != nil {
			return err
		}
		member, err := heimdallr.GetMember(s, guildID, m.Author.ID)
		if err != nil {
			return err
		}

		if !privilegeChecker(member, guild) {
			_, err := s.ChannelMessageSend(m.ChannelID, "You don't have the necessary role to do that.")
			return errors.Wrap(err, "sending message failed")
		}
		return handler(s, m, args)
	}
	return originalCommand
}

func getPrivilegeChecker(role string) func(*discordgo.Member, *discordgo.Guild) bool {
	switch role {
	case "helper":
		return heimdallr.IsHelper
	case "trial moderator":
		return heimdallr.IsTrialModOrHigher
	case "moderator":
		return heimdallr.IsModOrHigher
	case "supermoderator":
		return heimdallr.IsSuperModOrHigher
	case "admin":
		return heimdallr.IsAdminOrHigher
	case "owner":
		return heimdallr.IsOwner
	default:
		return nil
	}
}

func isOneLowerThanTwo(member1 *discordgo.Member, member2 *discordgo.Member) bool {
	if getHighestRole(member1) >= getHighestRole(member2) { // Higher number means lower rank
		return true
	}
	return false

}

func getHighestRole(m *discordgo.Member) int {
	var highestRole int = 6
	for _, role := range m.Roles {
		switch role {
		case heimdallr.Config.UserRole:
			if highestRole > 5 {
				highestRole = 5
			}
		case heimdallr.Config.TrialModRole:
			if highestRole > 4 {
				highestRole = 4
			}
		case heimdallr.Config.ModRole:
			if highestRole > 3 {
				highestRole = 3
			}
		case heimdallr.Config.SuperModRole:
			if highestRole > 2 {
				highestRole = 2
			}
		case heimdallr.Config.AdminRole:
			if highestRole > 1 {
				highestRole = 1
			}
		}
	}
	return highestRole
}

func getCommandPrefix() string {
	commandPrefix := heimdallr.Config.CommandPrefix
	if commandPrefix == "" {
		commandPrefix = ";"
	}
	return commandPrefix
}

//CommandHandler provides help.
func CommandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	commandPrefix := getCommandPrefix()
	content := strings.TrimSpace(m.Content)
	if strings.HasPrefix(content, commandPrefix) {
	OUTER:
		for _, content := range splitCommands(content) {
			content = content[len(commandPrefix):]
			if len(content) == 0 {
				continue
			}
			commandName := strings.Split(content, " ")[0]
			command, ok := nameToCommand[commandName]
			if ok {
				args, err := shlex.Split(content)
				if err != nil {
					embed := command.getFullHelpEmbedWithExamples(false, 15797003)
					embed.Title = "Incorrect use of command"
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
					heimdallr.LogIfError(s, errors.Wrap(err, "sending embed failed"))
					return
				}
				opts, err := command.parse(args[1:])
				if err != nil {
					embed := command.getFullHelpEmbedWithExamples(false, 15797003)
					embed.Title = "Incorrect use of command"
					_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
					heimdallr.LogIfError(s, errors.Wrap(err, "sending embed failed"))
					return
				}
				err = command.handle(s, m, opts)
				if err != nil {
					heimdallr.LogIfError(s, err)
					_, err := s.ChannelMessageSend(m.ChannelID, "Something went wrong, sorry!")
					heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))
				}
			} else {
				// If there is only one character or characters that aren't letters,
				// they probably didn't mean to type a command. Might be e.g. an emoji
				if len(commandName) == 1 {
					continue
				}
				for _, r := range commandName {
					if !unicode.IsLetter(r) {
						continue OUTER
					}
				}
				_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Unknown command `%s`. Type `%shelp` for a list of commands.", commandName, commandPrefix))
				heimdallr.LogIfError(s, errors.Wrap(err, "sending message failed"))
			}
		}
	}
}

func splitCommands(s string) []string {
	numCommands := 0
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, getCommandPrefix()) {
			numCommands++
		}
	}
	if numCommands == len(lines) {
		return lines
	}
	return []string{s}
}

// Allows both mentions and plain IDs
func getIDFromMaybeMention(maybeMention string, s *discordgo.Session) string {
	re, err := regexp.Compile(`<@[!&]?(\d+)>`)

	if err != nil {
		heimdallr.LogIfError(s, err)

	}
	if submatch := re.FindStringSubmatch(maybeMention); len(submatch) == 2 {
		return submatch[1]
	}
	return maybeMention
}
