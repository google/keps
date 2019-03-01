package auth

import (
	"io/ioutil"
)

func NewProvideTokenFromPath(p string) (TokenProvider, error) {
	var value = func() (string, error) {
		tokenBytes, err := ioutil.ReadFile(p)
		if err != nil {
			return "", err
		}

		return string(tokenBytes), nil
	}

	t := &tokenFromPath{
		value: value,
	}

	return t, nil
}

type TokenProvider interface {
	Value() (string, error)
}

type tokenValueFunc func() (string, error)

type tokenFromPath struct {
	value tokenValueFunc
}

func (t *tokenFromPath) Value() (string, error) {
	return t.value()
}
