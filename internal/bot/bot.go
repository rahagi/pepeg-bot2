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
	MARKOV_MAX_WORDS       = 30
	MARKOV_DEFAULT_COUNTER = 25
)

type HandlerFunc func(irc.IRCClient, *message.Payload) error

type HandlerMap map[string]HandlerFunc

type Bot interface {
	// Init starts the bot by receiving messages from the IRC server
	Init()

	// Register register a command handler
	RegisterHandler(cmd string, h HandlerFunc)
}

type bot struct {
	i irc.IRCClient
	h HandlerMap
	g generator.Generator
	t trainer.Trainer

	countUntilGenerate int
	enableLogging      bool
	learningOnly       bool
}

// NewBot create new bot
func NewBot(ircClient irc.IRCClient, enableLogging bool, g generator.Generator, t trainer.Trainer, learningOnly bool) Bot {
	handlerMap := make(HandlerMap)
	return &bot{
		i: ircClient,
		h: handlerMap,
		g: g,
		t: t,

		countUntilGenerate: MARKOV_DEFAULT_COUNTER,
		learningOnly:       learningOnly,
		enableLogging:      enableLogging,
	}
}

func (b *bot) Init() {
	messages := b.i.Receive()
	for p := range messages {
		b.countUntilGenerate -= 1
		b.defaultHandler(p)
		cmd, _ := common.PickCommand(p.Message)
		if hf, ok := b.h[cmd]; ok {
			if err := hf(b.i, p); err != nil {
				log.Printf("bot: failed to handle command %s: %v\n", p.Message, err)
			}
		}
	}
}

func (b *bot) RegisterHandler(cmd string, h HandlerFunc) {
	b.h[cmd] = h
}

func (b *bot) defaultHandler(p *message.Payload) {
	if strings.HasPrefix(p.Message, "PING") {
		b.i.Pong()
		return
	}
	if b.countUntilGenerate <= 0 {
		if !b.learningOnly {
			m := b.generateMarkov(p)
			b.i.Chat(m)
		}
		b.resetCounter()
	}
	if p.User != "" {
		fm := p.String()
		go b.t.AddChain(fm)
		if b.enableLogging {
			b.log(p)
		}
	}
}

func (b *bot) resetCounter() {
	b.countUntilGenerate = MARKOV_DEFAULT_COUNTER
}

func (b *bot) generateMarkov(p *message.Payload) string {
	res, err := b.g.Generate(p.String(), MARKOV_MAX_WORDS)
	if err != nil {
		return ""
	}
	return res
}

func (b *bot) log(p *message.Payload) {
	logPath := fmt.Sprintf("./log/%s.log", b.i.Channel())
	logger.Tee(p.String(), logPath)
}
