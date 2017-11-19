package config

import "github.com/spf13/viper"

// Settings global config
var Settings *AppSettings

// SlackSettings consists from from slack settings
type slackSettings struct {
	Token, Channel string
}

// boltSettings consists from boltdb settings
type boltSettings struct {
	FilePath string
}

// StoreSettings include all databases settings
type storeSettings struct {
	Bolt *boltSettings
}

// AppSettings consists from all options in config file
type AppSettings struct {
	Tags, Feeds, TrimStrings []string
	Slack                    *slackSettings
	Store                    *storeSettings
}

// NewSettings create settings from config file
func NewSettings(name, path string) (*AppSettings, error) {
	viper.SetConfigName(name)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	settings := &AppSettings{
		Tags:        viper.GetStringSlice("tags"),
		Feeds:       viper.GetStringSlice("feeds"),
		TrimStrings: viper.GetStringSlice("trim"),
		Slack: &slackSettings{
			Token:   viper.GetString("slack.token"),
			Channel: viper.GetString("slack.channel"),
		},
		Store: &storeSettings{
			Bolt: &boltSettings{
				FilePath: viper.GetString("store.bolt.file"),
			},
		},
	}
	return settings, nil
}
