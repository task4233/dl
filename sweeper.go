package dl

import (
	"bufio"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const dlPath = "\"github.com/task4233/dl\""

type Sweeper struct {
	dlPkgName string
}

func NewSweeper() *Sweeper {
	return &Sweeper{
		dlPkgName: "dl", // default package name
	}
}

// Sweep deletes all methods related to dl in a ".go" file.
// This method requires ".dl" directory to exist.
func (d *Sweeper) Sweep(ctx context.Context, targetFilePath string) error {
	// validation
	if !strings.HasSuffix(targetFilePath, ".go") {
		return fmt.Errorf("targetPath is not .go file: %s", targetFilePath)
	}

	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, targetFilePath, nil, 0)
	if err != nil {
		return err
	}

	for _, decl := range fileAst.Decls {
		switch w := decl.(type) {
		case *ast.GenDecl:
			// check import alias
			if w.Tok.String() == "import" {
				if err := d.removeImportSpec(&w.Specs); err != nil {
					return err
				}
			}
		case *ast.FuncDecl:
			// remove all methods
			if err := d.removedlStmt(&w.Body.List); err != nil {
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
	if err := os.Rename(tmpFile.Name(), targetFilePath); err != nil {
		return err
	}

	return nil
}

func (d *Sweeper) removeImportSpec(specs *[]ast.Spec) error {
	var removedIdx int = -1

	for importSpecIdx, spec := range *specs {
		switch importSpec := spec.(type) {
		case *ast.ImportSpec:
			if importSpec.Path != nil && importSpec.Path.Value == dlPath {
				removedIdx = importSpecIdx
				if importSpec.Name != nil {
					d.dlPkgName = importSpec.Name.Name
				}
			}
		}
	}

	if removedIdx < 0 {
		return nil
	}
	*specs = append((*specs)[:removedIdx], (*specs)[removedIdx+1:]...)
	return nil
}

func (d *Sweeper) removedlStmt(statements *[]ast.Stmt) error {
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
						if d.dlPkgName == x2.Name {
							removedIdxs = append(removedIdxs, idx)
						}
					}
				}
				// TODO: add other cases
			}
		default:
			Printf("not implemented: %#v\nplease report this bug to https://github.com/task4233/dl/issues/new/choose 🙏", exp)
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
	f, err = os.CreateTemp("", "_dl.go")
	if err == nil {
		fn = func() {
			os.Remove(f.Name())
		}
	}

	return
}

// Evacuate copies ".go" files to under ".dl" directory.
// This method requires ".dl" directory to exist.
// This method doesn't allow to invoke with a file included in `excludeFiles`.
func (d *Sweeper) Evacuate(ctx context.Context, baseDirPath string, srcFilePath string) error {
	// resolve path
	rel, err := filepath.Rel(baseDirPath, srcFilePath)
	if err != nil {
		return err
	}
	targetPath := filepath.Join(baseDirPath, ".dl", rel)

	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	parentDir := filepath.Join(targetPath, "..")
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		if err := os.MkdirAll(parentDir, 0700); err != nil {
			return err
		}
	}
	dstFile, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

// Restore copies ".go" files to raw places.
// This method requires ".dl" directory to exist.
func (s *Sweeper) Restore(ctx context.Context, srcFilePath string, dstFilePath string) error {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}
