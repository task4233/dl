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
	cleanCmd   = "dl clean"
	restoreCmd = "dl restore"
	dlDir      = ".dl"
	msg        = "dl %s: The instant logger package for debug.\n"
)

type Cmd interface {
	Run(ctx context.Context, baseDir string) error
}

// DeLog structs.
type DeLog struct {
}

// New for running dl package with CLI.
func New() *DeLog {
	return &DeLog{}
}

// Run executes each method for dl package.
func (d *DeLog) Run(ctx context.Context, version string, args []string) error {
	if len(args) == 0 {
		d.usage(version, "")
		return errors.New("no command is given.")
	}
	if len(args) == 1 {
		args = append(args, ".")
	}

	switch args[0] {
	case "clean":
		return NewClean().Run(ctx, args[1])
	case "init":
		return NewInit().Run(ctx, args[1])
	case "remove":
		return NewRemove().Run(ctx, args[1])
	case "restore":
		return NewRestore().Run(ctx, args[1])
	default:
		return d.usage(version, args[0])
	}
}

func (d *DeLog) usage(version, invalidCmd string) error {
	fmt.Fprintf(os.Stderr, msg+`Usage: dl [command]
Commands:
clean <dir>                 deletes logs used this package.
init <dir>                  add dl command into pre-commit script.
remove <dir>                remove dl command from pre-commit script.
restore <dir>               restore removed logs.
`, version)
	return fmt.Errorf("%s is not implemented.", invalidCmd)
}
