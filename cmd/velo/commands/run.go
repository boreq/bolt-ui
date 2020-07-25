package commands

import (
	"encoding/json"
	"os"

	"github.com/boreq/errors"
	"github.com/boreq/guinea"
	"github.com/boreq/velo/internal/config"
	"github.com/boreq/velo/internal/wire"
)

var runCmd = guinea.Command{
	Run: runRun,
	Arguments: []guinea.Argument{
		{
			Name:        "config",
			Optional:    false,
			Multiple:    false,
			Description: "Path to the configuration file",
		},
	},
	ShortDescription: "serves your music",
}

func runRun(c guinea.Context) error {
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
