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

func TestBuildCaseInsensitiveMatch(t *testing.T) {
	s := "hi"
	m := buildCaseInsensitiveMatch(s)
	expected := "[hH][iI]*"
	if m != expected {
		t.Errorf("expected: %s, got: %s", expected, m)
	}
}
