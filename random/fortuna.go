package random

import (
	"encoding/binary"
	"math/rand"
	"sync"
	"time"

	"github.com/seehuhn/fortuna"
)

type FortunaWrap struct {
	sync.RWMutex
	rng  *fortuna.Accumulator
	sink chan<- time.Time
}

func NewFortunaWrap(init string) *FortunaWrap {
	rng, err := fortuna.NewRNG("./randomness")
	if err != nil {
		panic("cannot initialise the RNG: " + err.Error())
	}
	return &FortunaWrap{sync.RWMutex{}, rng, nil}
}

func (r *FortunaWrap) Close() {
	r.rng.Close()
	if r.sink != nil {
		close(r.sink)
	}
}

func (r *FortunaWrap) Update() {
	r.Lock()
	defer r.Unlock()
	if r.sink == nil {
		sink := r.rng.NewEntropyTimeStampSink()
		r.sink = sink
	}
	r.sink <- time.Now()
	r.SeedStdRand()
}

func (r *FortunaWrap) SeedStdRand() {
	randomness := make([]byte, 8)
	r.rng.Read(randomness) // always success
	seed := binary.BigEndian.Uint64(randomness)
	rand.Seed(int64(seed))
}
