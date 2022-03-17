package delog

import (
	"bufio"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

const delogPath = "\"github.com/task4233/delog\""

type sweeper struct {
	delogPkgName string
}

func NewDelog() *sweeper {
	return &sweeper{
		delogPkgName: "delog", // default package name
	}
}

// Clean deletes all methods related to delog
func (d *sweeper) Clean(ctx context.Context, targetPath string) error {
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

	for _, decl := range fileAst.Decls {
		switch w := decl.(type) {
		case *ast.GenDecl:
			// check import alias
			if w.Tok.String() == "import" {
				for importSpecIdx, spec := range w.Specs {
					if importSpec, ok := spec.(*ast.ImportSpec); ok && importSpec != nil {
						if importSpec.Path != nil && importSpec.Path.Value == delogPath {
							removedIdx = importSpecIdx
							if importSpec.Name != nil {
								d.delogPkgName = importSpec.Name.Name
							}
						}
					}
				}

				// in importing only delog
				if removedIdx == 0 {
					w.Specs = w.Specs[1:]
				} else {
					w.Specs = append(w.Specs[:removedIdx], w.Specs[removedIdx+1:]...)
				}
			}
		case *ast.FuncDecl:
			// remove all methods
			err := d.removeDelogStmt(&w.Body.List)
			if err != nil {
				return err
			}
		}
	}

	// if import spec is empty, remove import gen decl
	if len(fileAst.Decls) > 0 {
		if importDecl, ok := fileAst.Decls[0].(*ast.GenDecl); ok {
			if len(importDecl.Specs) == 0 {
				fileAst.Decls = fileAst.Decls[1:]
			}
		}
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

func (d *sweeper) removeDelogStmt(statements *[]ast.Stmt) error {
	removedIdxs := []int{}

	for idx, stmt := range *statements {
		switch exp := stmt.(type) {
		case *ast.ExprStmt:
			switch x := exp.X.(type) {
			case *ast.CallExpr:
				switch fun := x.Fun.(type) {
				case *ast.SelectorExpr:
					switch x2 := fun.X.(type) {
					case *ast.Ident:
						if d.delogPkgName == x2.Name {
							removedIdxs = append(removedIdxs, idx)
						}
					}
				}
				// TODO: add other cases
			}
		default:
			fmt.Printf("other type: %#v\n", exp)
		}
	}

	for idx := len(removedIdxs) - 1; idx >= 0; idx-- {
		*statements = append((*statements)[:removedIdxs[idx]], (*statements)[removedIdxs[idx]+1:]...)
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
