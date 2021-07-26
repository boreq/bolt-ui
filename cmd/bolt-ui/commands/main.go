package commands

import (
	"fmt"
	"strings"

	"github.com/boreq/guinea"
	"github.com/boreq/velo/internal/config"
	"github.com/boreq/velo/internal/wire"
	"github.com/boreq/velo/logging"
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
	Options: []guinea.Option{
		{
			Name:        "insecure-cors",
			Type:        guinea.Bool,
			Default:     false,
			Description: "Disables CORS",
		},
	},
	ShortDescription: "a web user interface for the Bolt database",
	Description: `
Thanks to bolt-ui you are able to explore a Bolt database using a web interface.
`,
}

var log = logging.New("main")

func run(c guinea.Context) error {
	conf := config.Default()
	conf.DatabaseFile = c.Arguments[0]
	conf.InsecureCORS = c.Options["insecure-cors"].Bool()

	if conf.InsecureCORS {
		log.Warn("insecure-cors option enabled")
	}

	service, err := wire.BuildService(conf)
	if err != nil {
		return errors.Wrap(err, "could not create a service")
	}

	printInfo(conf)

	return service.HTTPServer.Serve(conf.ServeAddress)
}

func printInfo(conf *config.Config) {
	addr := conf.ServeAddress
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}

	addr = "http://" + addr

	fmt.Println("------------")
	fmt.Println()
	fmt.Printf("Starting listening on %s\n", addr)
	fmt.Println()
	fmt.Println("------------")
}
