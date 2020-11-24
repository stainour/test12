package reader

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type fileWalker struct {
	rootPath string
}

type walkResult struct {
	path string
	err  error
}

func newFileWalker(rootPath string) (fileWalker, error) {
	info, err := os.Stat(rootPath)
	if err != nil {
		return fileWalker{}, fmt.Errorf("error accessing folder %s: %w", rootPath, err)
	}

	if !info.IsDir() {
		return fileWalker{}, fmt.Errorf("path %s is not folder", rootPath)
	}

	return fileWalker{rootPath: rootPath}, nil
}

func (fw *fileWalker) walk(ctx context.Context) chan walkResult {
	filePaths := make(chan walkResult, runtime.GOMAXPROCS(-1))

	go func() {
		err := filepath.Walk(fw.rootPath, func(path string, info os.FileInfo, err error) error {
			select {
			case <-ctx.Done():
				return errors.New("file reading was canceled")

			default:
				if err != nil {
					return fmt.Errorf("error accessing file %s: %w", path, err)
				}

				if info.IsDir() {
					return nil
				}
				filePaths <- walkResult{
					path: path,
				}

				return nil
			}
		})
		if err != nil {
			filePaths <- walkResult{
				err: err,
			}
		}

		close(filePaths)
	}()

	return filePaths
}
