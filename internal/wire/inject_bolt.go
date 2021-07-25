package wire

import (
	"github.com/boreq/velo/adapters"
	"github.com/boreq/velo/internal/config"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

//lint:ignore U1000 because
var boltSet = wire.NewSet(
	newBolt,
)

func newBolt(conf *config.Config) (*bolt.DB, error) {
	return adapters.NewBolt(conf.DatabaseFile)
}
