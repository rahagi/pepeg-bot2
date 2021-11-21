package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"
	"time"

	"github.com/jpillora/backoff"
	"github.com/rahagi/pepeg-bot2/internal/irc/message"
)

// IRCClient
type IRCClient interface {
	// Chat send message (`PRIVMSG`) to a channel
	Chat(m string)

	// Receive return a read-only channel that receive a message
	// in from an IRC channel
	Receive() <-chan *message.Payload

	// Pong handle `PING` message sent by the server
	Pong()

	// Username return the username of this bot
	Username() string

	// Channel return a channel name this client connect to
	Channel() string
}

type ircClient struct {
	username string
	oauth    string
	channel  string
	addr     string
	c        net.Conn
	b        *backoff.Backoff
}

// NewClient open an IRC connection using `net.Dial`
// with `tcp` connection. After connection has been established,
// it continue authenticate the connection and join a channel
func NewClient(username, oauth, channel, addr string) IRCClient {
	client := &ircClient{
		username: username,
		oauth:    oauth,
		channel:  channel,
		addr:     addr,
		b: &backoff.Backoff{
			Min: 1 * time.Second,
			Max: 30 * time.Minute,
		},
	}
	client.initConn()
	client.auth()
	client.join()
	return client
}

func (i *ircClient) Chat(m string) {
	message := fmt.Sprintf("PRIVMSG #%s :%s", i.channel, m)
	if err := i.send(message); err != nil {
		log.Printf("client: failed to send a chat message: %v\n", err)
	}
}

func (i *ircClient) Receive() <-chan *message.Payload {
	messages := make(chan *message.Payload)
	tp := textproto.NewReader(bufio.NewReader(i.c))
	go func() {
		defer i.c.Close()
		for {
			rawMessage, err := tp.ReadLine()
			if err != nil {
				i.c.Close()
				i.initConn()
				tp = textproto.NewReader(bufio.NewReader(i.c))
				continue
			}
			m := message.BuildPayload(rawMessage)
			messages <- m
		}
	}()
	return messages
}

func (i *ircClient) Pong() {
	message := "PONG:tmi.twitch.tv"
	i.send(message)
}

func (i *ircClient) Username() string { return i.username }

func (i *ircClient) Channel() string { return i.channel }

func (i *ircClient) join() {
	message := fmt.Sprintf("JOIN #%s", i.channel)
	if err := i.send(message); err != nil {
		log.Fatalf("client: failed to join a channel: %v", err)
	}
	log.Printf("joined channel %s\n", i.channel)
}

func (i *ircClient) send(m string) error {
	_, err := i.c.Write([]byte(m + "\r\n"))
	if err != nil {
		return err
	}
	return nil
}

func (i *ircClient) auth() {
	messages := []string{
		"PASS " + i.oauth,
		"NICK " + i.username,
	}
	for _, m := range messages {
		if err := i.send(m); err != nil {
			log.Fatalf("client: failed to authenticate:  %v", err)
		}
	}
}

func (i *ircClient) initConn() {
	conn, err := net.Dial("tcp", i.addr)
	if err != nil {
		log.Printf("irc: failed to connect to server. attempting to reconnect...")
		time.Sleep(i.b.Duration())
		i.initConn()
		return
	}
	i.c = conn
	i.b.Reset()
}
