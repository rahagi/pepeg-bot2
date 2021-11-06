package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/rahagi/pepeg-bot2/internal/helper/common"
	"github.com/rahagi/pepeg-bot2/internal/helper/logger"
	"github.com/rahagi/pepeg-bot2/internal/irc"
	"github.com/rahagi/pepeg-bot2/internal/irc/message"
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
	IRCClient irc.IRCClient
	Handlers  map[string]HandlerFunc
	G         generator.Generator
	T         trainer.Trainer

	countUntilGenerate int
	enableLogging      bool
	learningOnly       bool
}

func NewBot(ircClient irc.IRCClient, enableLogging bool, g generator.Generator, t trainer.Trainer, learningOnly bool) Bot {
	handlerMap := make(map[string]HandlerFunc)
	return &bot{
		IRCClient: ircClient,
		Handlers:  handlerMap,
		G:         g,
		T:         t,

		countUntilGenerate: MARKOV_DEFAULT_COUNTER,
		learningOnly:       learningOnly,
		enableLogging:      enableLogging,
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
	if strings.HasPrefix(m.Message, "PING") {
		b.IRCClient.Pong()
	}
	if b.countUntilGenerate <= 0 && !b.learningOnly {
		res, err := b.G.Generate(m.Message, MARKOV_MAX_WORDS)
		if err != nil {
			return
		}
		b.IRCClient.Chat(res)
		b.resetCounter()
	}
	if m.User != "" {
		go b.T.AddChain(m.Format())
		if b.enableLogging {
			logPath := fmt.Sprintf("./log/%s.log", b.IRCClient.GetChannel())
			logger.Tee(m.Format(), logPath)
		}
	}
}

func (b *bot) resetCounter() {
	b.countUntilGenerate = MARKOV_DEFAULT_COUNTER
}
