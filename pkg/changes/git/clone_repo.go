package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofrs/flock"

	log "github.com/sirupsen/logrus"
	libgit "gopkg.in/src-d/go-git.v4"
	libgithttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"

	"github.com/calebamiles/keps/pkg/changes/auth"
	"github.com/calebamiles/keps/pkg/settings/cache"
)

func Clone(token auth.TokenProvider, repoUrl string, toLocation string) (Repo, error) {
	authToken, err := token.Value()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(toLocation); !os.IsNotExist(err) {
		log.Errorf("location: %s may exist already, refusing to overwrite", toLocation)
		return nil, fmt.Errorf("location: %s may exist already, refusing to overwrite", toLocation)
	}

	underlyingRepo, err := libgit.PlainClone(toLocation, false, &libgit.CloneOptions{
		Auth: &libgithttp.BasicAuth{
			Username: auth.ArbitraryUsername,
			Password: authToken,
		},
		URL: repoUrl,
	})

	if err != nil {
		log.Errorf("cloning `upstream` remote: %s", err)
		return nil, err
	}

	nameish := strings.Replace(toLocation, string(filepath.Separator), "-", -1)
	nameish = strings.Replace(nameish, "-", "", 1) // just the leading `-`

	lockLocation := filepath.Join(cache.Dir(), fmt.Sprintf("%s-%s", nameish, repositoryFileLock))
	err = os.MkdirAll(filepath.Dir(lockLocation), os.ModePerm)
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
		r := &repository{
			underlying: underlyingRepo,
			localPath:  toLocation,
			locker:     new(sync.Mutex),
		}

		return r, nil

	default:
		return nil, fmt.Errorf("could not obtain exlusive file lock when cloning repository")
	}
}
