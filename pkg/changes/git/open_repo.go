package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofrs/flock"
	libgit "gopkg.in/src-d/go-git.v4"

	"github.com/calebamiles/keps/pkg/settings/cache"
)

func Open(p string) (Repo, error) {
	nameish := strings.Replace(p, string(filepath.Separator), "-", -1)
	nameish = strings.Replace(nameish, "-", "", 1) // just the leading `-`

	lockLocation := filepath.Join(cache.Dir(), fmt.Sprintf("%s-%s", nameish, repositoryFileLock))
	err := os.MkdirAll(filepath.Dir(lockLocation), os.ModePerm)
	if err != nil {
		return nil, err
	}

	repoLock := flock.New(lockLocation)
	locked, err := repoLock.TryLock()
	if err != nil {
		return nil, err
	}

	switch locked {
	case true:
		underlyingRepo, err := libgit.PlainOpen(p)
		if err != nil {
			return nil, err
		}

		r := &repository{
			underlying: underlyingRepo,
			localPath:  p,
			locker:     new(sync.Mutex),
		}

		return r, nil

	default:
		return nil, errors.New("could not obtain exlusive file lock when opening repository")
	}
}

const (
	repositoryFileLock = "kepTool.lock"
)
