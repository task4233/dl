package dl

import (
	"bufio"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

const dlPath = "\"github.com/task4233/dl\""

var _ Cmd = (*Clean)(nil)

type Clean struct {
	dlPkgName   string
	removedIdxs *IntHeap
}

func NewClean() *Clean {
	return &Clean{
		dlPkgName:   "dl", // default package name
		removedIdxs: &IntHeap{},
	}
}

var (
	excludedFiles = []string{dlDir, ".git"}
)

// Run deletes all methods related to dl in ".go" files under the given directory path
func (c *Clean) Run(ctx context.Context, baseDir string) error {
	dlDirPath := filepath.Join(baseDir, dlDir)
	if _, err := os.Stat(dlDirPath); os.IsNotExist(err) {
		return fmt.Errorf(".dl directory doesn't exist. Please execute $ dl init .: %s", dlDirPath)
	}

	return walkDirWithValidation(ctx, baseDir, func(path string, info fs.DirEntry) error {
		for _, file := range excludedFiles {
			if strings.Contains(path, file) {
				return nil
			}
		}
		if err := c.Evacuate(ctx, baseDir, path); err != nil {
			return fmt.Errorf("failed to evacuate %s, %s", path, err.Error())
		}

		// might be good running concurrently? TODO(#7)
		fmt.Fprintf(os.Stderr, "remove dl from %s\n", path)
		return c.Sweep(ctx, path)
	})
}

// Sweep deletes all methods related to dl in a ".go" file.
// This method requires ".dl" directory to exist.
func (c *Clean) Sweep(ctx context.Context, targetFilePath string) error {
	// validation
	if !strings.HasSuffix(targetFilePath, ".go") {
		return fmt.Errorf("targetPath is not .go file: %s", targetFilePath)
	}

	fset := token.NewFileSet()
	fileAst, err := parser.ParseFile(fset, targetFilePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, decl := range fileAst.Decls {
		switch w := decl.(type) {
		case *ast.GenDecl:
			// check import alias
			if err := c.removeImportSpec(ctx, &w.Specs); err != nil {
				return err
			}
		case *ast.FuncDecl:
			// remove all methods
			if err := c.removeDlStmts(ctx, &w.Body.List); err != nil {
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
	// might be change to GOTMPDIR
	tmpFile, err := os.CreateTemp("", "_dl.go")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name())

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

func (c *Clean) removeImportSpec(ctx context.Context, specs *[]ast.Spec) error {
	var removedIdx int = -1

	for importSpecIdx, spec := range *specs {
		switch exp := spec.(type) {
		case *ast.ImportSpec:
			if exp.Path != nil && exp.Path.Value == dlPath {
				removedIdx = importSpecIdx
				if exp.Name != nil {
					c.dlPkgName = exp.Name.Name
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

func (c *Clean) removeDlStmts(ctx context.Context, statements *[]ast.Stmt) error {
	for idx, stmt := range *statements {
		if err := c.removeDlStmt(ctx, &stmt, idx); err != nil {
			return err
		}
	}

	for c.removedIdxs.Len() > 0 {
		idx := c.removedIdxs.Pop()
		*statements = append((*statements)[:idx], (*statements)[idx+1:]...)
	}

	return nil
}

func (c *Clean) removeDlStmt(ctx context.Context, stmt *ast.Stmt, idx int) error {
	switch exp := (*stmt).(type) {
	case *ast.AssignStmt:
		for _, expr := range exp.Rhs {
			if err := c.scanDlIdentInExpr(ctx, &expr, idx); err != nil {
				return err
			}
		}
	case *ast.BlockStmt:
		return c.removeDlStmts(ctx, &exp.List)
	case *ast.CommClause:
		{
			return c.removeDlStmts(ctx, &exp.Body)
		}
	case *ast.CaseClause:
		return c.removeDlStmts(ctx, &exp.Body)
	case *ast.DeclStmt:
		// As expressions are not allowed here, ignore this stmt
	case *ast.ExprStmt:
		return c.scanDlIdentInExpr(ctx, &exp.X, idx)
	case *ast.ForStmt:
		return c.removeDlStmts(ctx, &exp.Body.List)
	case *ast.GoStmt:
		return c.scanDlIdentInCallExpr(ctx, exp.Call, idx)
	case *ast.IfStmt:
		return c.removeDlStmts(ctx, &exp.Body.List)
	case *ast.SendStmt:
		// TODO: As expressions are not allowed here, ignore this stmt for now
	case *ast.SelectStmt:
		return c.removeDlStmts(ctx, &exp.Body.List)
	case *ast.SwitchStmt:
		return c.removeDlStmts(ctx, &exp.Body.List)
	case *ast.TypeSwitchStmt:
		return c.removeDlStmts(ctx, &exp.Body.List)
	case *ast.RangeStmt:
		if err := c.removeDlStmts(ctx, &exp.Body.List); err != nil {
			return err
		}
		if err := c.scanDlIdentInExpr(ctx, &exp.Key, idx); err != nil {
			return err
		}
	case *ast.ReturnStmt:
		for _, expr := range exp.Results {
			if err := c.scanDlIdentInExpr(ctx, &expr, idx); err != nil {
				return err
			}
		}
	default:
		Printf("not implemented: %#v\nplease report this bug to https://github.com/task4233/dl/issues/new/choose 🙏\n", exp)
	}
	return nil
}

func (c *Clean) scanDlIdentInExpr(ctx context.Context, expr *ast.Expr, idx int) error {
	switch x := (*expr).(type) {
	case *ast.CallExpr:
		return c.scanDlIdentInCallExpr(ctx, x, idx)
	case *ast.FuncLit:
		if err := c.removeDlStmts(ctx, &x.Body.List); err != nil {
			return err
		}
	}

	return nil
}

func (c *Clean) scanDlIdentInCallExpr(ctx context.Context, expr *ast.CallExpr, idx int) error {
	switch fun := expr.Fun.(type) {
	case *ast.SelectorExpr:
		switch x2 := fun.X.(type) {
		case *ast.Ident:
			if c.dlPkgName == x2.Name {
				c.removedIdxs.Push(idx)
			}
		}
	}

	return nil
}

// Evacuate copies ".go" files to under ".dl" directory.
// This method requires ".dl" directory to exist.
// This method doesn't allow to invoke with a file included in `excludeFiles`.
func (c *Clean) Evacuate(ctx context.Context, baseDirPath string, srcFilePath string) error {
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
