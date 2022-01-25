package filter

import (
	"log"
	"regexp"
	"strings"

	"github.com/rahagi/pepeg-bot2/internal/helper/common"
)

type Filter interface {
	CensorBannedWord(string) string
}

type filter struct {
	bannedWords []string
}

func NewFromFile(path string) Filter {
	b, err := common.ReadLines(path)
	if err != nil {
		log.Fatalf("filter: failed to read banned word list: %v", err)
	}
	return &filter{b}
}

func (f *filter) CensorBannedWord(message string) (result string) {
	result = message
	for _, b := range f.bannedWords {
		re := regexp.MustCompile(`(?i)` + b)
		result = re.ReplaceAllStringFunc(result, func(s string) string {
			return strings.Map(func(r rune) rune {
				if r == ' ' {
					return r
				} else {
					return '*'
				}
			}, s)
		})
	}
	return
}
