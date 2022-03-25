package dl

import (
	"bufio"
	"container/heap"
	"context"
	"errors"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

const (
	dlPackageUrl = `"github.com/task4233/dl/v2"`
	issueRequest = `\nPlease report this bug to https://github.com/task4233/dl/issues/new/choose if possible🙏\n`
)

var _ cmd = (*cleanCmd)(nil)

type cleanCmd struct {
	dlPkgName   string
	removedIdxs heap.Interface
	astFile     *ast.File
}

func newCleanCmd() *cleanCmd {
	return &cleanCmd{
		dlPkgName:   "dl", // default package name
		removedIdxs: newintHeap(nil),
		astFile:     nil,
	}
}

var excludedFiles = []string{dlDir, ".git"}

// Run deletes all methods related to dl in ".go" files under the given directory path
func (c *cleanCmd) Run(ctx context.Context, baseDir string) error {
	dlDirPath := filepath.Join(baseDir, dlDir)
	if _, err := os.Stat(dlDirPath); os.IsNotExist(err) {
		return fmt.Errorf(".dl directory doesn't exist. Please execute $ dl init .: %s", dlDirPath)
	}

	return walkDirWithValidation(ctx, baseDir, func(ctx context.Context, path string, info fs.DirEntry) error {
		for _, file := range excludedFiles {
			if strings.Contains(path, file) {
				return nil
			}
		}
		if err := c.Evacuate(ctx, baseDir, path); err != nil {
			return fmt.Errorf("failed to evacuate %s, %s", path, err.Error())
		}

		fmt.Fprintf(os.Stderr, "remove dl from %s\n", path)
		return c.Sweep(ctx, path)
	})
}

// Sweep deletes all methods related to dl in a ".go" file.
// This method requires ".dl" directory to exist.
func (c cleanCmd) Sweep(ctx context.Context, targetFilePath string) error {
	// validation
	if !strings.HasSuffix(targetFilePath, ".go") {
		return fmt.Errorf("targetPath is not .go file: %s", targetFilePath)
	}

	fset := token.NewFileSet()
	var err error
	c.astFile, err = parser.ParseFile(fset, targetFilePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	if err := c.removeDlFromAst(ctx); err != nil {
		return err
	}

	// overwriting
	// tmp file is created in $GOTMPDIR
	tmpFile, err := os.CreateTemp("", "_dl.go")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

	writer := bufio.NewWriter(tmpFile)
	defer writer.Flush()

	if err := format.Node(writer, fset, c.astFile); err != nil {
		return err
	}
	if err := os.Rename(tmpFile.Name(), targetFilePath); err != nil {
		return err
	}

	return nil
}

func (c cleanCmd) removeDlFromAst(ctx context.Context) error {
	var ok bool
	c.astFile, ok = astutil.Apply(c.astFile, func(cur *astutil.Cursor) bool {
		// if c.Node belongs importspec, remove import statement for dl
		found, err := c.findDlImportInImportSpec(ctx, cur)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed findDlImportInImportSpec: %v"+issueRequest, err)
			return true
		}
		if found {
			cur.Delete()
			return true
		}

		// if c.Node belongs ExprStmt, remove callExpr for dl
		found, err = c.findDlInvocationInStmt(ctx, cur)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed findDlImportInImportSpec: %v"+issueRequest, err)
			return true
		}
		if found {
			cur.Delete()
			return true
		}

		// if return false, traversing is stopped immediately
		return true
	}, nil).(*ast.File)

	// remove import spec when it's empty
	for idx, decl := range c.astFile.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if len(d.Specs) == 0 {
				c.removedIdxs.Push(idx)
			}
		}
	}
	if !ok {
		return errors.New("failed type conversion from any to *ast.File")
	}

	for c.removedIdxs.Len() > 0 {
		idx := c.removedIdxs.Pop().(int)
		if !(idx+1 < len(c.astFile.Decls)) {
			c.astFile.Decls = c.astFile.Decls[:idx]
		} else {
			c.astFile.Decls = append(c.astFile.Decls[:idx], c.astFile.Decls[idx+1:]...)
		}
	}

	return nil
}

func (c *cleanCmd) findDlImportInImportSpec(ctx context.Context, cr *astutil.Cursor) (bool, error) {
	switch node := cr.Node().(type) {
	case *ast.ImportSpec:
		return cr.Index() >= 0 && node.Path.Value == dlPackageUrl, nil
	}

	return false, nil
}

func (c *cleanCmd) findDlInvocationInStmt(ctx context.Context, cr *astutil.Cursor) (bool, error) {
	switch node := cr.Node().(type) {
	case *ast.ExprStmt:
		switch x := node.X.(type) {
		case *ast.CallExpr:
			return c.findDlInvocationInCallExpr(ctx, x, cr.Index())
		}
	case *ast.AssignStmt:
		for _, r := range node.Rhs {
			switch x := r.(type) {
			case *ast.CallExpr:
				return c.findDlInvocationInCallExpr(ctx, x, cr.Index())
			}
		}
	case *ast.ReturnStmt:
		for _, r := range node.Results {
			switch x := r.(type) {
			case *ast.CallExpr:
				return c.findDlInvocationInCallExpr(ctx, x, cr.Index())
			}
		}
	}
	return false, nil
}

func (c *cleanCmd) findDlInvocationInCallExpr(ctx context.Context, callExpr *ast.CallExpr, idx int) (bool, error) {
	fun, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return false, fmt.Errorf("fun is not *ast.SelectorExpr: %v", callExpr.Fun)
	}
	x2, ok := fun.X.(*ast.Ident)
	if !ok {
		return false, fmt.Errorf("x2 is not *ast.Ident: %v", fun.X)
	}

	// check node is in a slice
	return idx >= 0 && c.dlPkgName == x2.Name, nil
}

// Evacuate copies ".go" files to under ".dl" directory.
// This method requires ".dl" directory to exist.
// This method doesn't allow to invoke with a file included in `excludeFiles`.
func (c cleanCmd) Evacuate(ctx context.Context, baseDirPath string, srcFilePath string) error {
	// resolve path
	rel, err := filepath.Rel(baseDirPath, srcFilePath)
	if err != nil {
		return err
	}

	targetFilePath := filepath.Join(baseDirPath, ".dl", rel)
	parentDir := filepath.Join(targetFilePath, "..")
	if _, err := os.Stat(parentDir); os.IsNotExist(err) {
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return err
		}
	}

	return copyFile(ctx, targetFilePath, srcFilePath)
}
