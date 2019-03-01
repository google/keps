package hermetic

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func RecursiveCopy(src string, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	hasTrailingSlash := strings.HasSuffix(dst, string(filepath.Separator))

	switch {
	// src == file; dst == file
	case !srcInfo.IsDir() && !hasTrailingSlash:
		from, err := os.Open(src)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Dir(dst), os.ModePerm)
		if err != nil {
			return err
		}

		to, err := os.Create(dst)
		if err != nil {
			return err
		}

		_, err = io.Copy(from, to)
		if err != nil {
			return err
		}

		err = to.Close()
		if err != nil {
			return err
		}

		err = from.Close()
		if err != nil {
			return err
		}

		return nil

	// src == file; dst == dir
	case !srcInfo.IsDir() && hasTrailingSlash:
		from, err := os.Open(src)
		if err != nil {
			return err
		}

		err = os.MkdirAll(dst, os.ModePerm)
		if err != nil {
			return err
		}

		to, err := os.Create(filepath.Join(dst, filepath.Base(src)))
		if err != nil {
			return err
		}

		_, err = io.Copy(from, to)
		if err != nil {
			return err
		}

		err = to.Close()
		if err != nil {
			return err
		}

		err = from.Close()
		if err != nil {
			return err
		}

		return nil

	case srcInfo.IsDir() && !hasTrailingSlash:
		return fmt.Errorf("cannot copy directory: %s into file: %s", src, dst)

	case srcInfo.IsDir() && hasTrailingSlash:
		shouldSkipDir := map[string]bool{
			".":  true,
			"..": true,
		}

		return filepath.Walk(src, func(path string, info os.FileInfo, prevErr error) error {
			if prevErr != nil {
				return prevErr
			}

			if shouldSkipDir[info.Name()] {
				return filepath.SkipDir
			}

			relPath, err := filepath.Rel(src, path)
			if err != nil {
				return err
			}

			if info.IsDir() && filepath.Base(src) == info.Name() {
				// relPath would be `.` here
				return os.MkdirAll(filepath.Join(dst, info.Name()), os.ModePerm)
			} else if info.IsDir() {
				// there's nothing more to do for directories
				return os.MkdirAll(filepath.Join(dst, filepath.Base(src), relPath), os.ModePerm)
			}

			from, err := os.Open(path)
			if err != nil {
				return err
			}

			to, err := os.Create(filepath.Join(dst, filepath.Base(src), relPath))
			if err != nil {
				return err
			}

			_, err = io.Copy(to, from)
			if err != nil {
				return err
			}

			err = to.Close()
			if err != nil {
				return err
			}

			err = from.Close()
			if err != nil {
				return err
			}

			return nil
		})

	default:
		return fmt.Errorf("do not know how to copy: %s into: %s", src, dst)
	}
}
