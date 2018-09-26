package sigs

import (
	"net/http"
	"io/ioutil"
	"fmt"

	"github.com/go-yaml/yaml"
)

func FetchUpstreamList() (*upstreamSIGList, error) {
	resp, err := http.Get(UpstreamSIGListURL)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("downloading SIG info: %s", err)
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading downloaded SIG info: %s", err)
	}

	sl := &upstreamSIGList{}
	err = yaml.Unmarshal(respBytes, sl)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling downloaded SIG info: %s", err)
	}

	return sl, nil
}
