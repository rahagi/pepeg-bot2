package config

import (
	"os"

	"github.com/joho/godotenv"
)

type config struct {
	Username string
	OAuth    string
	Channel  string
	IRCAddr  string
}

func BuildConfig() *config {
	godotenv.Load()
	c := new(config)
	c.Username = os.Getenv("USERNAME")
	c.OAuth = os.Getenv("OAUTH")
	c.Channel = os.Getenv("CHANNEL")
	c.IRCAddr = os.Getenv("IRC_ADDR")
	return c
}
