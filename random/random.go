package random

import (
	"crypto/md5"
	"encoding/binary"
	"io"
	"math/rand"
	"time"
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
