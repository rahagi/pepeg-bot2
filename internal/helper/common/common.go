package common

import (
	"bufio"
	"fmt"
	"os"
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

func ReadLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var l []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		l = append(l, s.Text())
	}
	return l, s.Err()
}
