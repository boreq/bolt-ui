package tracker

import (
	"math/rand"
	"sync"
	"time"

	"github.com/boreq/errors"

	"github.com/oklog/ulid/v2"
)

type UUIDGenerator struct {
	entropy *ulid.MonotonicEntropy
	mutex   sync.Mutex
}

func NewUUIDGenerator() *UUIDGenerator {
	return &UUIDGenerator{
		entropy: ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0),
	}
}

func (g *UUIDGenerator) Generate() (string, error) {
	g.mutex.Lock()
	defer g.mutex.Unlock()

	t := time.Now()

	u, err := ulid.New(ulid.Timestamp(t), g.entropy)
	if err != nil {
		return "", errors.Wrap(err, "ulid generation failed")
	}

	return u.String(), nil
}
