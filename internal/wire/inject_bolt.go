package wire

import (
	bolt2 "github.com/boreq/bolt-ui/adapters/bolt"
	"github.com/boreq/bolt-ui/internal/config"
	"github.com/google/wire"
	bolt "go.etcd.io/bbolt"
)

//lint:ignore U1000 because
var boltSet = wire.NewSet(
	newBolt,
)

func newBolt(conf *config.Config) (*bolt.DB, error) {
	return bolt2.NewBolt(conf.DatabaseFile)
}
