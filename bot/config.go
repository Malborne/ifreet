package heimdallr

import (
	"os"

	"github.com/pkg/errors"

	"github.com/BurntSushi/toml"
)

//Config is the app config
var Config = BotConfig{}

//BotConfig contains config info
type BotConfig struct {
	Token           string `toml:"token"`
	WelcomeChannel  string `toml:"welcome_channel"`
	LogChannel      string `toml:"log_channel"`
	AdminChannel    string `toml:"admin_channel"`
	AdminLogChannel string `toml:"admin_log_channel"`
	ArchiveChannel  string `toml:"archive_channel"`

	BotChannel   string `toml:"bot_channel"`
	RulesChannel string `toml:"rules_channel"`

	AdminRole          string `toml:"admin_role_id"`
	SuperModRole       string `toml:"supermod_role_id"`
	ModRole            string `toml:"mod_role_id"`
	UserRole           string `toml:"user_role_id"`
	MaleRole           string `toml:"male_role_id"`
	FemaleRole         string `toml:"female_role_id"`
	VerifiedMaleRole   string `toml:"verified_male_role_id"`
	VerifiedFemaleRole string `toml:"verified_female_role_id"`
	MutedRole          string `toml:"muted_role_id"`

	Roles []Role `toml:"role"`

	WelcomeMessage  string `toml:"welcome_message"`
	ApprovalMessage string `toml:"approval_message"`

	CommandPrefix string `toml:"command_prefix"`
}

//Role is a struct containing details about user-assignable roles.
type Role struct {
	ID   string `toml:"id"`
	Name string `toml:"name"`
	Desc string `toml:"description"`
}

//LoadConfig loads the configuration file
func (conf *BotConfig) LoadConfig(f string) error {
	_, err := toml.DecodeFile(f, &conf)
	return errors.Wrap(err, "decoding config file failed")
}

//SaveConfig saves the configuration file
func (conf *BotConfig) SaveConfig(f string) error {
	file, err := os.Create(f)
	if err != nil {
		return errors.Wrap(err, "creating config file failed")
	}

	encoder := toml.NewEncoder(file)

	err = encoder.Encode(conf)
	if err != nil {
		return errors.Wrap(err, "encoding config file failed")
	}

	err = file.Close()
	if err != nil {
		return errors.Wrap(err, "closing config file failed")
	}
	return nil
}
