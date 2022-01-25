package handler

import (
	"fmt"

	"github.com/rahagi/pepeg-bot2/internal/bot"
	"github.com/rahagi/pepeg-bot2/internal/constant"
	"github.com/rahagi/pepeg-bot2/internal/irc"
	"github.com/rahagi/pepeg-bot2/internal/irc/message"
)

func MakeVersionHandler(version string) bot.HandlerFunc {
	return func(i irc.IRCClient, p *message.Payload) error {
		message := fmt.Sprintf("@%s pepeg-bot2 version: %s", p.User, version)
		i.Chat(message)
		return nil
	}
}

func MakeEchoHandler() bot.HandlerFunc {
	return func(i irc.IRCClient, p *message.Payload) error {
		m := p.Message[len(constant.COMMAND_PREFIX)+len("echo")+1:]
		i.Chat(m)
		return nil
	}
}
