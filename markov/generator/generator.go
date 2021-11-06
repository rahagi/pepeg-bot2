package generator

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rahagi/pepeg-bot2/markov/config"
	"github.com/rahagi/pepeg-bot2/markov/helper/common"
)

type probabilitySlice []string

type Generator interface {
	Generate(seed string, maxWords int) (string, error)
}

type generator struct {
	r *redis.Client
}

func NewGenerator(r *redis.Client) Generator {
	return &generator{r}
}

func (g *generator) Generate(seed string, maxWords int) (string, error) {
	ctx, cancle := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancle()
	s := strings.Split(seed, config.WORD_SEPARATOR)
	if len(s) < config.CHAIN_LEN {
		return "", fmt.Errorf("generator: seed cannot be shorter than CHAIN_LEN")
	}
	chain := s[:config.CHAIN_LEN]
	res := strings.Join(chain, " ") + " "
	for i := 0; i < maxWords; i++ {
		key := common.MakeKey(chain)
		cmd := g.r.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
			Min: "-inf",
			Max: "+inf",
		})
		z, err := cmd.Result()
		if err != nil {
			return res, nil
		}
		arr := generateProbabilitySlice(z)
		next := pickRandomByProbability(arr)
		if next == config.STOP_TOKEN || len(arr) == 0 {
			break
		}
		res += next + " "
		mid := strings.Split(key, config.KEY_SEPARATOR)[1:]
		chain = append(mid, next)
	}
	return res, nil
}

func generateProbabilitySlice(z []redis.Z) probabilitySlice {
	var res probabilitySlice
	for _, v := range z {
		normalizedScore := (int(v.Score) / 10) + 1
		for i := 0; i < normalizedScore; i++ {
			res = append(res, v.Member.(string))
		}
	}
	return res
}

func pickRandomByProbability(p probabilitySlice) string {
	rand.Seed(time.Now().Unix())
	n := len(p)
	if n <= 0 {
		return ""
	}
	rand := rand.Intn(n)
	return p[rand]
}
