package generator

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rahagi/pepeg-bot2/markov/constant"
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
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	seed = common.Sanitize(seed)
	s := strings.Split(seed, constant.WORD_SEPARATOR)
	if len(s) < constant.CHAIN_LEN {
		return "", fmt.Errorf("generator: seed cannot be shorter than CHAIN_LEN")
	}
	key := common.RandKeyBySeed(seed, g.r)
	res := common.NormalizeKey(key) + " "
	for i := 0; i < maxWords; i++ {
		cmd := g.r.ZRangeByScoreWithScores(ctx, key, &redis.ZRangeBy{
			Min: "-inf",
			Max: "+inf",
		})
		z, err := cmd.Result()
		if err != nil {
			return res, nil
		}
		arr := generateProbabilitySlice(z)
		next := common.PickRandomString(arr)
		if next == constant.STOP_TOKEN || len(arr) == 0 {
			break
		}
		res += next + " "
		mid := strings.Split(key, constant.KEY_SEPARATOR)[1:]
		key = common.MakeKey(append(mid, next))
	}
	res = strings.TrimSuffix(res, " ")
	if len(strings.Split(res, " ")) <= constant.CHAIN_LEN {
		res = ""
	}
	return res, nil
}

func generateProbabilitySlice(z []redis.Z) probabilitySlice {
	var res probabilitySlice
	for _, v := range z {
		normalizedScore := (int(v.Score) / constant.SCORE_MODIFIER) + 1
		for i := 0; i < normalizedScore; i++ {
			res = append(res, v.Member.(string))
		}
	}
	return res
}
