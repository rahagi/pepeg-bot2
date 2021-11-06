package common

import (
	"fmt"
	"strings"

	"github.com/rahagi/pepeg-bot2/markov/config"
)

func MakeKey(s []string) string {
	k := ""
	for _, w := range s {
		k += fmt.Sprintf("%s%s", w, config.KEY_SEPARATOR)
	}
	k = strings.TrimSuffix(k, config.KEY_SEPARATOR)
	return k
}
