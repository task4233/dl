package delog

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

type Cli struct{}

func New() *Cli {
	return &Cli{}
}

func (c *Cli) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return errors.New("no argument")
	}
	switch args[0] {
	case "clean":
		if len(args) == 1 {
			args = append(args, "sandbox/hello.go") // TODO: fix here
		}
		return c.Clean(ctx, args[1])
	default:
		return fmt.Errorf("command %s is not implemented", args[0])
	}
}

// Clean deletes all methods related to delog in ".go" files under the given directory path
func (c *Cli) Clean(ctx context.Context, baseDir string) error {
	// TODO: check recursively all files if baseDir is a directory
	// just in case, run with an assumption: baseDir == a single file
	targetFilePath := baseDir

	if err := Clean(ctx, targetFilePath); err != nil {
		return err
	}

	return nil
}

var delogName string = "delog"

const delogPath string = "\"github.com/task4233/delog\""

// var delogName string = "fmt"
// const delogPath string = "\"fmt\""

// Clean deletes all methods related to delog
func Clean(ctx context.Context, targetPath string) error {
	// validation
	if !strings.HasSuffix(targetPath, ".go") {
		return fmt.Errorf("targetPath is not .go file: %s", targetPath)
	}

	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, targetPath, nil, 0)
	if err != nil {
		return err
	}

	var removedIdx int = -1

	for idx, decl := range fileAst.Decls {
		switch w := decl.(type) {
		case *ast.GenDecl:
			// check import alias
			if w.Tok.String() == "import" {
				if len(w.Specs) > 0 {
					if importSpec, ok := w.Specs[0].(*ast.ImportSpec); ok && importSpec != nil {
						if importSpec.Path != nil && importSpec.Path.Value == delogPath {
							removedIdx = idx
							if importSpec.Name != nil {
								delogName = importSpec.Name.Name
							}
						}
					}
				}
			}
		case *ast.FuncDecl:
			// remove all methods
			err := removeDelogStmt(&w.Body.List)
			if err != nil {
				return err
			}
		}
	}

	if removedIdx == 0 {
		fileAst.Decls = fileAst.Decls[1:]
	} else if removedIdx > 0 {
		fileAst.Decls = append(fileAst.Decls[:removedIdx], fileAst.Decls[removedIdx+1:]...)
	}

	// overwriting
	tmpFile, cleanUp, err := createTmpFile()
	if err != nil {
		return err
	}
	defer cleanUp()

	writer := bufio.NewWriter(tmpFile)
	defer writer.Flush()

	fset = token.NewFileSet()

	if err := format.Node(writer, fset, fileAst); err != nil {
		return err
	}
	if err := os.Rename(tmpFile.Name(), targetPath); err != nil {
		return err
	}

	return nil
}

func removeDelogStmt(statements *[]ast.Stmt) error {
	removeIdxs := []int{}

	for idx, stmt := range *statements {
		switch exp := stmt.(type) {
		case *ast.ExprStmt:
			switch x := exp.X.(type) {
			case *ast.CallExpr:
				switch fun := x.Fun.(type) {
				case *ast.SelectorExpr:
					switch x2 := fun.X.(type) {
					case *ast.Ident:
						if delogName == x2.Name {
							removeIdxs = append(removeIdxs, idx)
						}
					}
				}
				// TODO: add other cases
			}
		default:
			fmt.Printf("other type: %#v\n", exp)
		}
	}

	for idx := len(removeIdxs) - 1; idx >= 0; idx-- {
		*statements = append((*statements)[:idx], (*statements)[idx+1:]...)
	}

	return nil
}

// createTmpFile creates a temporary file and return *os.File and a cleanUp function
func createTmpFile() (f *os.File, fn func(), err error) {
	// might be change to GOTMPDIR
	f, err = os.CreateTemp("", "_delog.go")
	if err == nil {
		fn = func() {
			os.Remove(f.Name())
		}
	}

	return
}
