package feed

import (
	"strings"

	conf "github.com/vasilishin/rfeed/config"
)

// Trim substrings from text
func Trim(s string) string {
	for _, cut := range conf.Settings.TrimStrings {
		s = strings.Replace(s, cut, "", -1)
	}
	return s
}
