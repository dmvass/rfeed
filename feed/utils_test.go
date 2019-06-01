package feed

import (
	"fmt"
	"testing"

	conf "github.com/dmvass/rfeed/config"
)

func init() {
	var err error
	// Read settings from config file
	conf.Settings, err = conf.NewSettings("../config.yml")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

func TestTrim(t *testing.T) {
	saved := conf.Settings.TrimStrings
	defer func() {
		conf.Settings.TrimStrings = saved
	}()

	text := "The Go Programming Language!"
	tables := []struct {
		Text        string
		TrimStrings []string
		Wanted      string
	}{
		{text, []string{"Go"}, "The  Programming Language!"},
		{text, []string{"Go Programming"}, "The  Language!"},
		{text, []string{"Go Language"}, text},
		{text, []string{"go"}, text},
		{text, []string{"Language!"}, "The Go Programming "},
		{text + text, []string{"Go Programming"}, "The  Language!The  Language!"},
		{text, []string{"Go", "Language", "The", " "}, "Programming!"},
	}

	for _, table := range tables {
		conf.Settings.TrimStrings = table.TrimStrings
		trimedText := Trim(table.Text)
		if trimedText != table.Wanted {
			t.Errorf(
				"Trim of %v was incorrect, got: %s, want: %s",
				table.TrimStrings, trimedText, table.Wanted,
			)
		}
	}
}
