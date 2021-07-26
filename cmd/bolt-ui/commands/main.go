package commands

import (
	"github.com/boreq/guinea"
	"github.com/boreq/velo/internal/config"
	"github.com/boreq/velo/internal/wire"
	"github.com/pkg/errors"
)

var MainCmd = guinea.Command{
	Run: run,
	Arguments: []guinea.Argument{
		{
			Name:        "database",
			Optional:    false,
			Multiple:    false,
			Description: "Path to the database file",
		},
	},
	ShortDescription: "a web user interface for the Bolt database",
	Description: `
Thanks to bolt-ui you are able to explore a Bolt database using a web interface.
`,
}

func run(c guinea.Context) error {
	conf := config.Default()
	conf.DatabaseFile = c.Arguments[0]

	service, err := wire.BuildService(conf)
	if err != nil {
		return errors.Wrap(err, "could not create a service")
	}

	return service.HTTPServer.Serve(conf.ServeAddress)
}
