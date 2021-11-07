package common

import (
	"fmt"
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
