package config

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func readConfig() (*AppSettings, error) {
	var err error
	Settings, err = NewSettings("../config.yml")
	return Settings, err
}

func TestReadConfig(t *testing.T) {
	if _, err := readConfig(); err != nil {
		t.Errorf("Fatal error config file: %s", err)
	}
}

func TestSettingsTags(t *testing.T) {
	readConfig()
	got, wanted := Settings.Tags, viper.GetStringSlice("tags")
	if !reflect.DeepEqual(got, wanted) {
		t.Errorf("Settings Tags slice was incorrect, got: %v, want: %v", got, wanted)
	}
}

func TestSettingsFeeds(t *testing.T) {
	readConfig()
	got, wanted := Settings.Feeds, viper.GetStringSlice("feeds")
	if !reflect.DeepEqual(got, wanted) {
		t.Errorf("Settings Feeds slice was incorrect, got: %v, want: %v", got, wanted)
	}
}

func TestSettingsTrimStrings(t *testing.T) {
	readConfig()
	got, wanted := Settings.TrimStrings, viper.GetStringSlice("trim")
	if !reflect.DeepEqual(got, wanted) {
		t.Errorf("Settings TrimStrings slice was incorrect, got: %v, want: %v", got, wanted)
	}
}

func TestSettingsStoreBolt(t *testing.T) {
	readConfig()
	got, wanted := Settings.Store.Bolt.FilePath, viper.GetString("store.bolt.file")
	if got != wanted {
		t.Errorf("Settings Store.Bolt.FilePath was incorrect, got: %v, want: %v", got, wanted)
	}
}

func TestSettingsSlack(t *testing.T) {
	readConfig()

	got, wanted := Settings.Slack.Token, viper.GetString("slack.token")
	if got != wanted {
		t.Errorf("Settings Slack.Token was incorrect, got: %v, want: %v", got, wanted)
	}

	got, wanted = Settings.Slack.Channel, viper.GetString("slack.channel")
	if got != wanted {
		t.Errorf("Settings Slack.Channel was incorrect, got: %v, want: %v", got, wanted)
	}
}

func TestSettingsTelegram(t *testing.T) {
	readConfig()

	tgot, twanted := Settings.Telegram.Token, viper.GetString("telegram.token")
	if tgot != twanted {
		t.Errorf("Settings Telegram.Token was incorrect, got: %v, want: %v", tgot, twanted)
	}

	cgot, cwanted := Settings.Telegram.ChatID, viper.GetInt64("telegram.chat_id")
	if cgot != cwanted {
		t.Errorf("Settings Telegram.ChatID was incorrect, got: %v, want: %v", cgot, cwanted)
	}
}
