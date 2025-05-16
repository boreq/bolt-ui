package wire

import (
	"net/http"

	httpport "github.com/boreq/bolt-ui/ports/http"
	"github.com/google/wire"
)

//lint:ignore U1000 because
var httpSet = wire.NewSet(
	httpport.NewServer,
	httpport.NewHandler,
	httpport.NewTokenAuthProvider,
	wire.Bind(new(http.Handler), new(*httpport.Handler)),
	wire.Bind(new(httpport.AuthProvider), new(*httpport.TokenAuthProvider)),
)
