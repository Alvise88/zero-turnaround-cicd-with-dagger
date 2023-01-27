package util

import (
	"io/fs"
	"strings"

	"dagger.io/dagger"
)

func ToCommand(cmd string) []string {
	return strings.Split(cmd, " ")
}

// create a copy of an embed directory
func CopyEmbedDir(e fs.FS, dir *dagger.Directory) (*dagger.Directory, error) {
	err := fs.WalkDir(e, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(e, path)
		if err != nil {
			return err
		}

		dir = dir.WithNewFile(path, string(content))

		return nil
	})
	if err != nil {
		return nil, err
	}
	return dir, nil
}
