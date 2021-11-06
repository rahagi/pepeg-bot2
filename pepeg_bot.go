package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/rahagi/pepeg-bot2/bot"
	"github.com/rahagi/pepeg-bot2/config"
	"github.com/rahagi/pepeg-bot2/irc"
	"github.com/rahagi/pepeg-bot2/irc/message"
	"github.com/rahagi/pepeg-bot2/markov/generator"
	"github.com/rahagi/pepeg-bot2/markov/trainer"
)

var Version = ""

func main() {
	cm := os.Args
	cfg := config.BuildConfig()
	r := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHostname,
	})
	switch {
	case len(cm) <= 1:
		log.Printf("Starting pepeg-bot2 version: %s", Version)

		// IRC client initialization
		log.Printf("connecting to (%s)\n", cfg.IRCAddr)
		ircClient := irc.NewClient(cfg.Username, cfg.OAuth, cfg.Channel, cfg.IRCAddr)
		log.Printf("connected to (%s)\n", cfg.IRCAddr)

		// Bot initialization
		g := generator.NewGenerator(r)
		bot_ := bot.NewBot(ircClient, cfg.EnableLogging, g)
		bot_.Handle("--version", func(i irc.IRCClient, p *message.Payload) error {
			message := fmt.Sprintf("@%s pepeg-bot2 version: %s", p.User, Version)
			i.Chat(message)
			return nil
		})
		bot_.Init()
	case cm[1] == "train":
		tData := ""
		if len(cm) <= 2 {
			tData = "training"
		} else {
			tData = cm[1]
		}
		log.Printf("connecting to redis (%s)\n", cfg.RedisHostname)
		t := trainer.NewTrainer(r, tData)
		t.Train()
	}
}
