package message

import (
	"fmt"
	"regexp"
	"time"
)

var messageRegex = *regexp.MustCompile(`^:(.*)!.* PRIVMSG.* :(.*)$`)

// Payload contain user, message, and UNIX timestamp parsed from
// raw message received from the IRC server
type Payload struct {
	User      string
	Message   string
	Timestamp int64
}

// BuildPayload build `Payload`
func BuildPayload(rawMessage string) *Payload {
	matches := messageRegex.FindStringSubmatch(rawMessage)
	if len(matches) < 3 {
		return &Payload{
			User:    "",
			Message: rawMessage,
		}
	}
	user := matches[1]
	message := matches[2]
	return &Payload{
		User:      user,
		Message:   message,
		Timestamp: time.Now().UTC().Unix(),
	}
}

func (p *Payload) Format() string {
	return fmt.Sprintf("%s: %s", p.User, p.Message)
}
