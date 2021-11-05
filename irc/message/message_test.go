package message

import (
	"testing"
)

func TestMessageRegex(t *testing.T) {
	m := ":lacaribot!lacaribot@lacaribot.tmi.twitch.tv PRIVMSG #lacari :ACTION CursedS0uL lost 123456 rage in roulette and now has 177332380 rage! FeelsBadMan"
	expectedMessage := "ACTION CursedS0uL lost 123456 rage in roulette and now has 177332380 rage! FeelsBadMan"
	expectedUser := "lacaribot"
	matches := messageRegex.FindStringSubmatch(m)
	if len(matches) < 3 {
		t.Error("matches len should be 2")
	}
	user := matches[1]
	if user != expectedUser {
		t.Error("failed to match user")
	}
	message := matches[2]
	if message != expectedMessage {
		t.Error("failed to match message")
	}
}
