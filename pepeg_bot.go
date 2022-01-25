package main

import (
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/rahagi/pepeg-bot2/config"
	"github.com/rahagi/pepeg-bot2/filter"
	"github.com/rahagi/pepeg-bot2/handler"
	"github.com/rahagi/pepeg-bot2/internal/bot"
	"github.com/rahagi/pepeg-bot2/internal/irc"
	"github.com/rahagi/pepeg-bot2/markov/generator"
	"github.com/rahagi/pepeg-bot2/markov/trainer"
)

var Version = ""

func initBot(cfg *config.Config, r *redis.Client) {
	log.Printf("starting pepeg-bot2 version: %s", Version)

	// IRC client initialization
	log.Printf("connecting to (%s)\n", cfg.IRCAddr)
	f := filter.NewFromFile(cfg.BannedWordsListPath)
	ircClient := irc.NewClient(cfg.Username, cfg.OAuth, cfg.Channel, cfg.IRCAddr, f)
	log.Printf("connected to (%s)\n", cfg.IRCAddr)

	// Bot initialization
	g := generator.NewGenerator(r)
	t := trainer.NewTrainer(r)
	b := bot.NewBot(ircClient, cfg.EnableLogging, g, t, cfg.LearningOnlyMode)
	b.RegisterHandler("--version", handler.MakeVersionHandler(Version))
	b.RegisterHandler("--echo", handler.MakeEchoHandler())
	b.Init()
}

func initTrain(cfg *config.Config, r *redis.Client, cmd []string) {
	tData := ""
	if len(cmd) <= 2 {
		tData = "training"
	} else {
		tData = cmd[1]
	}
	log.Printf("connecting to redis (%s)\n", cfg.RedisHostname)
	t := trainer.NewTrainerWithData(r, tData)
	t.Train()
}

func main() {
	cmd := os.Args
	cfg := config.BuildConfig()
	r := redis.NewClient(&redis.Options{
		Addr: cfg.RedisHostname,
	})
	switch {
	case len(cmd) <= 1:
		initBot(cfg, r)
	case cmd[1] == "train":
		initTrain(cfg, r, cmd)
	}
}
