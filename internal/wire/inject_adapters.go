package wire

import (
	"github.com/boreq/velo/adapters"
	"github.com/boreq/velo/application/auth"
	"github.com/boreq/velo/application/tracker"
	"github.com/google/wire"
)

//lint:ignore U1000 because
var adaptersSet = wire.NewSet(
	adapters.NewUUIDGenerator,
	wire.Bind(new(tracker.UUIDGenerator), new(*adapters.UUIDGenerator)),
	wire.Bind(new(auth.UUIDGenerator), new(*adapters.UUIDGenerator)),
)
