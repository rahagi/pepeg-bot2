package common

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rahagi/pepeg-bot2/markov/constant"
)

func MakeKey(s []string) string {
	k := ""
	for _, w := range s {
		k += fmt.Sprintf("%s%s", w, constant.KEY_SEPARATOR)
	}
	k = strings.TrimSuffix(k, constant.KEY_SEPARATOR)
	return k
}

func Sanitize(s string) string {
	ss := strings.Split(s, ": ")
	if len(ss) < 2 {
		return ""
	}
	r2 := regexp.MustCompile(`\x01(ACTION )?`)
	s = ss[1]
	s = r2.ReplaceAllString(s, "")
	return s
}
