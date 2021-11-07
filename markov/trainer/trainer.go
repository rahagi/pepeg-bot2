package trainer

import (
	"bufio"
	"context"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rahagi/pepeg-bot2/markov/constant"
	"github.com/rahagi/pepeg-bot2/markov/helper/common"
)

type Trainer interface {
	Train()
	AddChain(s string)
}

type trainer struct {
	r            *redis.Client
	trainingData string
}

func NewTrainer(r *redis.Client) Trainer {
	return &trainer{r, ""}
}

func NewTrainerWithData(r *redis.Client, tData string) Trainer {
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
	log.Println("starting training model...")
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
	log.Println("training done")
}

func (t *trainer) AddChain(s string) {
	t.addChain(s)
}

func (t *trainer) addChain(s string) {
	s = common.Sanitize(s)
	m := strings.Split(s, constant.WORD_SEPARATOR)
	if len(m) > constant.CHAIN_LEN {
		for i, j := 0, constant.CHAIN_LEN; j <= len(m); i, j = i+1, j+1 {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			chain := common.MakeKey([]string{m[i], m[i+1]})
			next := ""
			if j == len(m) {
				next = constant.STOP_TOKEN
			} else {
				next = m[j]
			}
			t.r.ZIncrBy(ctx, chain, 1, next)
		}
	}
}
