package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type config struct {
	Username      string
	OAuth         string
	Channel       string
	IRCAddr       string
	EnableLogging bool
	RedisHostname string
}

func BuildConfig() *config {
	godotenv.Load()
	c := new(config)
	c.Username = os.Getenv("USERNAME")
	c.OAuth = os.Getenv("OAUTH")
	c.Channel = os.Getenv("CHANNEL")
	c.IRCAddr = os.Getenv("IRC_ADDR")
	c.EnableLogging, _ = strconv.ParseBool(os.Getenv("ENABLE_LOGGING"))
	c.RedisHostname = os.Getenv("REDIS_HOSTNAME")
	return c
}
