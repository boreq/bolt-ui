package commands

import (
	"github.com/boreq/guinea"
	"github.com/boreq/velo/cmd/velo/commands/users"
)

var MainCmd = guinea.Command{
	Run: runMain,
	Subcommands: map[string]*guinea.Command{
		"run":            &runCmd,
		"users":          &users.UsersCmd,
		"default_config": &defaultConfigCmd,
	},
	ShortDescription: "a music streaming service",
	Description: `
Eggplant serves your music using a web interface.
`,
}

func runMain(c guinea.Context) error {
	return guinea.ErrInvalidParms
}
