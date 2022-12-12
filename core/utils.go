package core

import "strings"

var DefaultReplacer *strings.Replacer

func init() {
	DefaultReplacer = strings.NewReplacer("\t", "", "\r", "", "\n", "")

}

func Filter(msg string, r ...string) string {
	replacer := DefaultReplacer
	if len(r) > 0 {
		replacer = strings.NewReplacer("\t", r[0], "\r", r[0], "\n", r[0])
	}
	return replacer.Replace(msg)
}
