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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	p := strings.Split(seed, constant.WORD_SEPARATOR)
	if len(p) < 1 {
		return ""
	}
	k := r.Keys(ctx, fmt.Sprintf("%s*", p[0]))
	res, _ := k.Result()
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
