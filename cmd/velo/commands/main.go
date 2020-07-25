package commands

import (
	"github.com/boreq/velo/cmd/velo/commands/users"
	"github.com/boreq/guinea"
)

var MainCmd = guinea.Command{
	Run: runMain,
	Subcommands: map[string]*guinea.Command{
		"run":   &runCmd,
		"users": &users.UsersCmd,
	},
	ShortDescription: "a music streaming service",
	Description: `
Eggplant serves your music using a web interface.
`,
}

func runMain(c guinea.Context) error {
	return guinea.ErrInvalidParms
}
