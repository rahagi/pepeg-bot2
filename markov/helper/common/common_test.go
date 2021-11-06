package common

import (
	"strings"
	"testing"
)

func TestMakeKey(t *testing.T) {
	s := "Hello World"
	key := MakeKey(strings.Split(s, " "))
	expected := "Hello:World"
	if key != expected {
		t.Errorf("expected: %s, got: %s", expected, key)
	}
}
