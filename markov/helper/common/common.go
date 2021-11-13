package common

import (
	"context"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rahagi/pepeg-bot2/markov/constant"
)

func MakeKey(s []string) string {
	k := ""
	for _, w := range s {
		k += fmt.Sprintf("%s%s", w, constant.KEY_SEPARATOR)
	}
	k = strings.TrimSuffix(k, constant.KEY_SEPARATOR)
	return k
}

func Sanitize(s string) string {
	ss := strings.Split(s, ": ")
	r2 := regexp.MustCompile(`\x01(ACTION )?`)
	if len(ss) >= 2 {
		s = strings.Join(ss[1:], " ")
	}
	s = r2.ReplaceAllString(s, "")
	return s
}

func RandKeyBySeed(seed string, r *redis.Client) string {
	p := strings.Split(seed, constant.WORD_SEPARATOR)
	if len(p) < 1 {
		return ""
	}
	m := buildCaseInsensitiveMatch(p[0])
	res := []string{}
	cursor := uint64(0)
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		s := r.ScanType(ctx, cursor, m, 10000, "ZSET")
		k, cursor, _ := s.Result()
		res = append(res, k...)
		if cursor == 0 {
			break
		}
	}
	return PickRandomString(res)
}

func PickRandomString(p []string) string {
	rand.Seed(time.Now().UnixNano())
	n := len(p)
	if n <= 0 {
		return ""
	}
	i := rand.Intn(n)
	return p[i]
}

func NormalizeKey(k string) string {
	return strings.ReplaceAll(k, constant.KEY_SEPARATOR, constant.WORD_SEPARATOR)
}

func buildCaseInsensitiveMatch(s string) string {
	res := ""
	for _, r := range s {
		c := string(r)
		res += fmt.Sprintf("[%s%s]", strings.ToLower(c), strings.ToUpper(c))
	}
	return res + "*"
}
