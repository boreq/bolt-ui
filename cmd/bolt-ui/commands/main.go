package commands

import (
	"encoding/json"
	"os"

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
	conf, err := loadConfig(c.Arguments[0])
	if err != nil {
		return errors.Wrap(err, "could not load the configuration")
	}

	service, err := wire.BuildService(conf)
	if err != nil {
		return errors.Wrap(err, "could not create a service")
	}

	return service.HTTPServer.Serve(conf.ServeAddress)
}

func loadConfig(path string) (*config.Config, error) {
	conf := config.Default()

	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "could not open the config file")
	}

	if err := json.NewDecoder(f).Decode(&conf); err != nil {
		return nil, errors.Wrap(err, "could not unmarshal the config")
	}

	return conf, nil
}
