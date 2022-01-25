package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Username            string
	OAuth               string
	Channel             string
	IRCAddr             string
	RedisHostname       string
	EnableLogging       bool
	LearningOnlyMode    bool
	BannedWordsListPath string
}

func BuildConfig() *Config {
	godotenv.Load()
	c := new(Config)
	c.Username = os.Getenv("USERNAME")
	c.OAuth = os.Getenv("OAUTH")
	c.Channel = os.Getenv("CHANNEL")
	c.IRCAddr = os.Getenv("IRC_ADDR")
	c.RedisHostname = os.Getenv("REDIS_HOSTNAME")
	c.EnableLogging, _ = strconv.ParseBool(os.Getenv("ENABLE_LOGGING"))
	c.LearningOnlyMode, _ = strconv.ParseBool(os.Getenv("LEARNING_ONLY_MODE"))
	c.BannedWordsListPath = os.Getenv("BANNED_WORDS_PATH")
	return c
}
