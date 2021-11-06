package trainer

import (
	"bufio"
	"context"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rahagi/pepeg-bot2/markov/config"
	"github.com/rahagi/pepeg-bot2/markov/helper/common"
)

type Trainer interface {
	Train()
}

type trainer struct {
	r            *redis.Client
	trainingData string
}

func NewTrainer(r *redis.Client, tData string) Trainer {
	return &trainer{r, tData}
}

func (t *trainer) Train() {
	d, err := os.Open(t.trainingData)
	if err != nil {
		log.Fatalf("trainer: cannot read directory %s: %v", t.trainingData, err)
	}
	defer d.Close()
	files, err := d.Readdir(0)
	if err != nil {
		log.Fatalf("trainer: cannot read directory %s: %v", t.trainingData, err)
	}
	for _, fi := range files {
		if !fi.IsDir() {
			p := path.Join(t.trainingData, fi.Name())
			f, err := os.Open(p)
			if err != nil {
				log.Printf("trainer: skipping file %s: %v", p, err)
			}
			defer f.Close()
			s := bufio.NewScanner(f)
			for s.Scan() {
				t.addChain(s.Text())
			}
		}
	}
}

func (t *trainer) addChain(s string) {
	s = sanitize(s)
	m := strings.Split(s, config.WORD_SEPARATOR)
	if len(m) > config.CHAIN_LEN {
		for i, j := 0, config.CHAIN_LEN; j <= len(m); i, j = i+1, j+1 {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			chain := common.MakeKey([]string{m[i], m[i+1]})
			next := ""
			if j == len(m) {
				next = config.STOP_TOKEN
			} else {
				next = m[j]
			}
			log.Printf("(%s): [%s]", chain, next)
			t.r.ZIncrBy(ctx, chain, 1, next)
		}
	}
}

func sanitize(s string) string {
	r := regexp.MustCompile(`(.*:\s)(.*)`)
	s = r.FindStringSubmatch(s)[2]
	return s
}
