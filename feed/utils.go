package feed

import (
	"strings"

	"github.com/spf13/viper"
)

// Trim substrings from text
func Trim(s string) string {
	for _, cut := range viper.GetStringSlice("trim") {
		s = strings.Replace(s, cut, "", -1)
	}
	return s
}
