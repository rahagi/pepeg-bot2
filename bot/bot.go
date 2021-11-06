package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/rahagi/pepeg-bot2/helper/common"
	"github.com/rahagi/pepeg-bot2/helper/logger"
	"github.com/rahagi/pepeg-bot2/irc"
	"github.com/rahagi/pepeg-bot2/irc/message"
	"github.com/rahagi/pepeg-bot2/markov/generator"
	"github.com/rahagi/pepeg-bot2/markov/trainer"
)

const (
	MARKOV_MAX_WORDS       = 20
	MARKOV_DEFAULT_COUNTER = 15
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
	G             generator.Generator
	T             trainer.Trainer

	countUntilGenerate int
	learningOnly       bool
}

func NewBot(ircClient irc.IRCClient, enableLogging bool, g generator.Generator, t trainer.Trainer, learningOnly bool) Bot {
	handlerMap := make(map[string]HandlerFunc)
	return &bot{
		IRCClient:     ircClient,
		Handlers:      handlerMap,
		EnableLogging: enableLogging,
		G:             g,
		T:             t,

		countUntilGenerate: MARKOV_DEFAULT_COUNTER,
		learningOnly:       learningOnly,
	}
}

func (b *bot) Init() {
	messages := b.IRCClient.Receive()
	for m := range messages {
		b.countUntilGenerate -= 1
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
	case b.countUntilGenerate <= 0:
		if !b.learningOnly {
			res, err := b.G.Generate(m.Message, MARKOV_MAX_WORDS)
			if err != nil {
				return
			}
			b.IRCClient.Chat(res)
			b.resetCounter()
		}
	default:
		if m.User != "" {
			fm := fmt.Sprintf("%s: %s", m.User, m.Message)
			go b.T.AddChain(fm)
			if b.EnableLogging {
				logPath := fmt.Sprintf("./log/%s.log", b.IRCClient.GetChannel())
				logger.Tee(fm, logPath)
			}
		}
	}
}

func (b *bot) resetCounter() {
	b.countUntilGenerate = MARKOV_DEFAULT_COUNTER
}
