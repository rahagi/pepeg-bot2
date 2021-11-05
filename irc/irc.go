package irc

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/textproto"

	"github.com/rahagi/pepeg-bot2/irc/message"
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

	// GetUsername return the username of this bot
	GetUsername() string

	// GetChannel return a channel name this client connect to
	GetChannel() string
}

type ircClient struct {
	Username string
	OAuth    string
	Channel  string
	Addr     string
	Conn     net.Conn
}

// NewClient open an IRC connection using `net.Dial`
// with `tcp` connection. After connection has been established,
// it continue authenticate the connection and join a channel
func NewClient(username, oauth, channel, addr string) IRCClient {
	client := &ircClient{
		Username: username,
		OAuth:    oauth,
		Channel:  channel,
		Addr:     addr,
	}
	client.initConn()
	client.auth()
	client.join()
	return client
}

func (i *ircClient) Chat(m string) {
	message := fmt.Sprintf("PRIVMSG #%s :%s", i.Channel, m)
	if err := i.send(message); err != nil {
		log.Printf("client: failed to send a chat message: %v\n", err)
	}
}

func (i *ircClient) Receive() <-chan *message.Payload {
	messages := make(chan *message.Payload)
	tp := textproto.NewReader(bufio.NewReader(i.Conn))
	go func() {
		defer i.Conn.Close()
		for {
			rawMessage, err := tp.ReadLine()
			if err != nil {
				i.reconnect()
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

func (i *ircClient) GetUsername() string { return i.Username }

func (i *ircClient) GetChannel() string { return i.Channel }

func (i *ircClient) join() {
	message := fmt.Sprintf("JOIN #%s", i.Channel)
	if err := i.send(message); err != nil {
		log.Fatalf("client: failed to join a channel: %v", err)
	}
	log.Printf("joined channel %s\n", i.Channel)
}

func (i *ircClient) send(m string) error {
	_, err := i.Conn.Write([]byte(m + "\r\n"))
	if err != nil {
		return err
	}
	return nil
}

func (i *ircClient) auth() {
	messages := []string{
		"PASS " + i.OAuth,
		"NICK " + i.Username,
	}
	for _, m := range messages {
		if err := i.send(m); err != nil {
			log.Fatalf("client: failed to authenticate:  %v", err)
		}
	}
}

func (i *ircClient) initConn() {
	conn, err := net.Dial("tcp", i.Addr)
	if err != nil {
		log.Fatalf("client: cannot connect to IRC server: %v\n", err)
	}
	i.Conn = conn
}

func (i *ircClient) reconnect() {
	log.Println("client: attempting to reconnect")
	i.Conn.Close()
	i.initConn()
}
