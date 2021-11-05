package main

import (
	"log"

	"github.com/rahagi/pepeg-bot2/bot"
	"github.com/rahagi/pepeg-bot2/config"
	"github.com/rahagi/pepeg-bot2/irc"
)

var Version = ""

func main() {
	log.Printf("Starting pepeg-bot2 version: %s", Version)
	// IRC client initialization
	cfg := config.BuildConfig()
	log.Printf("connecting to (%s)\n", cfg.IRCAddr)
	ircClient := irc.NewClient(cfg.Username, cfg.OAuth, cfg.Channel, cfg.IRCAddr)
	log.Printf("connected to (%s)\n", cfg.IRCAddr)

	// Bot initialization
	bot_ := bot.NewBot(ircClient)
	bot_.Init()
}
