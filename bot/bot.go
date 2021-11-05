package bot

import (
	"fmt"
	"strings"

	"github.com/rahagi/pepeg-bot2/helper/logger"
	"github.com/rahagi/pepeg-bot2/irc"
)

type Bot interface {
	// Init starts the bot by receiving messages from the IRC server
	Init()
}

type bot struct {
	IRCClient irc.IRCClient
}

func NewBot(ircClient irc.IRCClient) Bot {
	return &bot{ircClient}
}

func (b *bot) Init() {
	messages := b.IRCClient.Receive()
	for m := range messages {
		switch {
		case strings.HasPrefix(m.Message, "PING"):
			b.IRCClient.Pong()
		default:
			fm := fmt.Sprintf("%s: %s", m.User, m.Message)
			logPath := fmt.Sprintf("./log/%s.log", b.IRCClient.GetChannel())
			logger.Tee(fm, logPath)
		}
	}
}
