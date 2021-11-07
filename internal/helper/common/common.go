package common

import (
	"fmt"
	"strings"

	"github.com/rahagi/pepeg-bot2/internal/constant"
)

// PickCommand parse chat message and return a valid command
func PickCommand(message string) (string, error) {
	cmd := ""
	if strings.HasPrefix(message, constant.COMMAND_PREFIX) {
		m := strings.Split(message, " ")
		if len(m) < 2 {
			return "", fmt.Errorf("helper/common: invalid command")
		}
		cmd = m[1]
	}
	return cmd, nil
}
