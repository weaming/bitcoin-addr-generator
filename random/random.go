package random

import (
	"crypto/md5"
	"encoding/binary"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/seehuhn/fortuna"
)

func InitRandSeed() {
	h := md5.New()
	// load the root randomness
	io.WriteString(h, "0a1b2f65373c93ea74ff8c0e84b3068feaad149b44efbdcb98061ffefa5dfb83")
	seed := binary.BigEndian.Uint64(h.Sum(nil))

	// different for every restart of program
	seed += uint64(time.Now().Unix())
	// fmt.Println("seed:", seed) // should not use in production

	// load the seed into global random
	rand.Seed(int64(seed))
}

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
