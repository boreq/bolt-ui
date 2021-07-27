package commands

import (
	"crypto/rand"
	"encoding/hex"
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
		{
			Name:        "insecure-token",
			Type:        guinea.Bool,
			Default:     false,
			Description: "Disables token validation",
		},
	},
	ShortDescription: "a web user interface for the Bolt database",
	Description: `
Thanks to bolt-ui you are able to explore a Bolt database using a web
interface. To access the web interface access the address printed out by the
program. Make sure that the address includes the token query parameter.
`,
}

var log = logging.New("main")

func run(c guinea.Context) error {
	conf, err := newConfig(c)
	if err != nil {
		return errors.New("could not create the config")
	}

	if conf.InsecureCORS {
		log.Warn("insecure-cors option enabled")
	}

	if conf.InsecureToken {
		log.Warn("insecure-token option enabled")
	}

	service, err := wire.BuildService(conf)
	if err != nil {
		return errors.Wrap(err, "could not create a service")
	}

	printInfo(conf)

	return service.HTTPServer.Serve(conf.ServeAddress)
}

func newConfig(c guinea.Context) (*config.Config, error) {
	token, err := generateSecureToken()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a secure token")
	}

	conf := config.Default()
	conf.DatabaseFile = c.Arguments[0]
	conf.Token = token
	conf.InsecureCORS = c.Options["insecure-cors"].Bool()
	conf.InsecureToken = c.Options["insecure-token"].Bool()

	return conf, nil
}

func printInfo(conf *config.Config) {
	addr := conf.ServeAddress
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}

	addr = "http://" + addr
	if !conf.InsecureToken {
		addr = fmt.Sprintf("%s/?token=%s", addr, conf.Token)
	}

	fmt.Println("------------")
	fmt.Println()
	fmt.Println(addr)
	fmt.Println()
	fmt.Println("------------")
}

const tokenLength = 32

func generateSecureToken() (string, error) {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "failed to read random bytes")
	}
	return hex.EncodeToString(b), nil
}
