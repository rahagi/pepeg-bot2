package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/rahagi/pepeg-bot2/helper/common"
	"github.com/rahagi/pepeg-bot2/helper/logger"
	"github.com/rahagi/pepeg-bot2/irc"
	"github.com/rahagi/pepeg-bot2/irc/message"
)

type HandlerFunc func(irc.IRCClient, *message.Payload) error

type Bot interface {
	// Init starts the bot by receiving messages from the IRC server
	Init()

	// Handle register a command handler
	Handle(command string, handler HandlerFunc)
}

type bot struct {
	IRCClient     irc.IRCClient
	Handlers      map[string]HandlerFunc
	EnableLogging bool
}

func NewBot(ircClient irc.IRCClient, enableLogging bool) Bot {
	handlerMap := make(map[string]HandlerFunc)
	return &bot{
		IRCClient:     ircClient,
		Handlers:      handlerMap,
		EnableLogging: enableLogging,
	}
}

func (b *bot) Init() {
	messages := b.IRCClient.Receive()
	for m := range messages {
		b.defaultHandler(m)
		command, _ := common.PickCommand(m.Message)
		if hf, ok := b.Handlers[command]; ok {
			if err := hf(b.IRCClient, m); err != nil {
				log.Printf("bot: failed to handle command %s: %v\n", m.Message, err)
			}
		}
	}
}

func (b *bot) Handle(command string, handler HandlerFunc) {
	b.Handlers[command] = handler
}

func (b *bot) defaultHandler(m *message.Payload) {
	switch {
	case strings.HasPrefix(m.Message, "PING"):
		b.IRCClient.Pong()
	default:
		if b.EnableLogging {
			fm := fmt.Sprintf("%s: %s", m.User, m.Message)
			logPath := fmt.Sprintf("./log/%s.log", b.IRCClient.GetChannel())
			logger.Tee(fm, logPath)
		}
	}
}
