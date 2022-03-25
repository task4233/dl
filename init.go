package dl

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.uber.org/multierr"
)

var _ cmd = (*initCmd)(nil)

// initCmd structs for $ dl init.
type initCmd struct{}

// newInitCmd for running $ dl init.
func newInitCmd() *initCmd {
	return &initCmd{}
}

// Run prepares an environment for dl commands.
func (i *initCmd) Run(ctx context.Context, baseDir string) error {
	// check if hooks directory exists or not
	if _, err := os.Stat(filepath.Join(baseDir, ".git", "hooks")); os.IsNotExist(err) {
		return err
	}

	// As it prevents to miss created files when an error occurs,
	// it handles errors all together at last.
	err := i.addGitPreHookScript(ctx, baseDir)
	err = multierr.Append(err, i.addGitPostHookScript(ctx, baseDir))
	err = multierr.Append(err, i.createDlDirIfNotExist(ctx, baseDir))
	err = multierr.Append(err, i.addDlIntoGitIgnore(ctx, baseDir))
	if err != nil {
		// rollback
		if removeErr := newRemoveCmd().Run(ctx, baseDir); err != nil {
			err = multierr.Append(err, removeErr)
		}
		return err
	}

	return nil
}

func (i *initCmd) addGitPreHookScript(ctx context.Context, baseDir string) error {
	return i.insertCodesIfNotExist(ctx, filepath.Join(baseDir, ".git", "hooks", "pre-commit"), cleanCmdStr, preCommitScript)
}

func (i *initCmd) addGitPostHookScript(ctx context.Context, baseDir string) error {
	return i.insertCodesIfNotExist(ctx, filepath.Join(baseDir, ".git", "hooks", "post-commit"), restoreCmdStr, postCommitScript)
}

func (i *initCmd) addDlIntoGitIgnore(ctx context.Context, baseDir string) error {
	return i.insertCodesIfNotExist(ctx, filepath.Join(baseDir, ".gitignore"), dlDir, fmt.Sprintf("\n%s\n", dlDir))
}

func (*initCmd) insertCodesIfNotExist(ctx context.Context, targetFilePath string, checkedCodesIfExists string, addedCodes string) error {
	f, err := os.OpenFile(targetFilePath, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	// It checks if `checkedCodesIfExists` has been installed or not.
	// If so, not inserting codes.
	buf, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	if bytes.Contains(buf, []byte(checkedCodesIfExists)) {
		return nil
	}

	if _, err := fmt.Fprint(f, addedCodes); err != nil {
		return err
	}

	return os.Chmod(targetFilePath, 0755)
}

func (*initCmd) createDlDirIfNotExist(ctx context.Context, baseDir string) error {
	path := filepath.Join(baseDir, dlDir)
	if stat, err := os.Stat(path); err == nil {
		if stat.IsDir() {
			return nil
		}
		return fmt.Errorf("%s has been already existed as file. Please rename or delete it.", path)
	}

	return os.MkdirAll(path, 0755)
}
