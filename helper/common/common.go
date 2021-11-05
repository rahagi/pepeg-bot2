package common

import (
	"fmt"
	"strings"

	"github.com/rahagi/pepeg-bot2/constant"
)

// PickCommand parse chat message and return a valid command
func PickCommand(message string) (string, error) {
	command := ""
	if strings.HasPrefix(message, constant.COMMAND_PREFIX) {
		m := strings.Split(message, " ")
		if len(m) < 2 {
			return "", fmt.Errorf("helper/common: invalid command")
		}
		command = m[1]
	}
	return command, nil
}
