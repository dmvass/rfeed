package config

import (
	"path/filepath"

	"github.com/spf13/viper"
)

// Settings global config
var Settings *AppSettings

// Slack Settings consists from from slack settings
type slackSettings struct {
	Token, Channel string
}

// Telegram Settings consists from from slack settings
type telegramSettings struct {
	Token  string
	ChatID int64
}

// Bolt Settings consists from boltdb settings
type boltSettings struct {
	FilePath string
}

// Store Settings include all databases settings
type storeSettings struct {
	Bolt *boltSettings
}

// AppSettings consists from all options in config file
type AppSettings struct {
	Tags, Feeds, TrimStrings []string
	Interval                 int64
	Slack                    *slackSettings
	Telegram                 *telegramSettings
	Store                    *storeSettings
}

// NewSettings create settings from config file
func NewSettings(configPath string) (*AppSettings, error) {

	dir, file := filepath.Split(configPath)
	ext := filepath.Ext(file)

	viper.AddConfigPath(dir)
	viper.SetConfigName(file[:len(file)-len(ext)])
	viper.SetConfigType(ext[1:])

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	settings := &AppSettings{
		Tags:        viper.GetStringSlice("tags"),
		Feeds:       viper.GetStringSlice("feeds"),
		TrimStrings: viper.GetStringSlice("trim"),
		Interval:    viper.GetInt64("check_interval"),
		Slack: &slackSettings{
			Token:   viper.GetString("slack.token"),
			Channel: viper.GetString("slack.channel"),
		},
		Telegram: &telegramSettings{
			Token:  viper.GetString("telegram.token"),
			ChatID: viper.GetInt64("telegram.chat_id"),
		},
		Store: &storeSettings{
			Bolt: &boltSettings{
				FilePath: viper.GetString("store.bolt.file"),
			},
		},
	}
	return settings, nil
}
