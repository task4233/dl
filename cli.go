package dl

import (
	"context"
	"errors"
	"fmt"
	"os"
)

// Although these values will be casted to []byte, they are declared as constants
// because that is't happened frequently.
const (
	preCommitScript = `#!/bin/sh
dl clean .
git add .
`
	postCommitScript = `#!/bin/sh
dl restore .
`
	cleanCmdStr   = "dl clean"
	restoreCmdStr = "dl restore"
	dlDir         = ".dl"
	msg           = "dl %s: The instant logger package for debug.\n"
)

type cmd interface {
	Run(ctx context.Context, baseDir string) error
}

// CLI structs.
type CLI struct {
}

// New for running dl package with CLI.
func New() *CLI {
	return &CLI{}
}

// Run executes each method for dl package.
func (d *CLI) Run(ctx context.Context, version string, args []string) error {
	if len(args) == 0 {
		d.usage(version, "")
		return errors.New("no command is given.")
	}
	if len(args) == 1 {
		args = append(args, ".")
	}

	switch args[0] {
	case "clean":
		return newCleanCmd().Run(ctx, args[1])
	case "init":
		return newInitCmd().Run(ctx, args[1])
	case "remove":
		return newRemoveCmd().Run(ctx, args[1])
	case "restore":
		return newRestoreCmd().Run(ctx, args[1])
	default:
		return d.usage(version, args[0])
	}
}

func (d *CLI) usage(version, invalidCmd string) error {
	fmt.Fprintf(os.Stderr, msg+`Usage: dl [command]
Commands:
clean <dir>                 deletes logs used this package.
init <dir>                  add dl command into pre-commit script.
remove <dir>                remove dl command from pre-commit script.
restore <dir>               restore removed logs.
`, version)
	return fmt.Errorf("%s is not implemented.", invalidCmd)
}
