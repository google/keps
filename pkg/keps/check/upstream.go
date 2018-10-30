package check

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/calebamiles/keps/pkg/keps/metadata"
	"github.com/calebamiles/keps/pkg/keps/states"
)

// TODO explain how this supports an important design philosophy: extension of verifyable trust
func ThatKEPExistsUpstream(meta metadata.KEP) error {
	var errs *multierror.Error

	expectedKEPURL := strings.Join([]string{upstreamKEPRepoURL, meta.ContentDir(), metadataFilename}, "/")
	resp, err := http.Get(expectedKEPURL)
	if err != nil {
		errs = multierror.Append(errs, err)
		return errs
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		errs = multierror.Append(errs, fmt.Errorf("KEP: %s with unique ID: %s not found upstream", meta.Title(), meta.UniqueID()))
		return errs
	}

	upstreamKEPBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errs = multierror.Append(errs, err)
		return errs
	}

	upstreamKEP, err := metadata.FromBytes(upstreamKEPBytes)
	if err != nil {
		errs = multierror.Append(errs, err)
		return errs
	}

	if upstreamKEP.UniqueID() != meta.UniqueID() {
		errs = multierror.Append(errs, fmt.Errorf("upstream KEP unique ID: %s, does not match given unique ID: %s", upstreamKEP.UniqueID(), meta.UniqueID()))
	}

	return errs.ErrorOrNil()
}

func ThatKEPHasBeenAcceptedUpstream(meta metadata.KEP) error {
	var errs *multierror.Error

	expectedKEPURL := strings.Join([]string{upstreamKEPRepoURL, meta.ContentDir(), metadataFilename}, "/")
	resp, err := http.Get(expectedKEPURL)
	if err != nil {
		errs = multierror.Append(errs, err)
		return errs
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		errs = multierror.Append(errs, fmt.Errorf("KEP: %s with unique ID: %s not found upstream", meta.Title(), meta.UniqueID()))
		return errs
	}

	upstreamKEPBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		errs = multierror.Append(errs, err)
		return errs
	}

	upstreamKEP, err := metadata.FromBytes(upstreamKEPBytes)
	if err != nil {
		errs = multierror.Append(errs, err)
		return errs
	}

	if upstreamKEP.UniqueID() != meta.UniqueID() {
		errs = multierror.Append(errs, fmt.Errorf("upstream KEP unique ID: %s, does not match given unique ID: %s", upstreamKEP.UniqueID(), meta.UniqueID()))
	}

	// TODO settle on state names + have KEP init create metadata with draft or proposal status
	if upstreamKEP.State() != states.Provisional {
		errs = multierror.Append(errs, fmt.Errorf("upstream KEP has state: %s, not 'Accepted'", upstreamKEP.State()))
	}

	return errs.ErrorOrNil()
}

const (
	// TODO update URL with final home
	upstreamKEPRepoURL = "https://raw.githubusercontent.com/calebamiles/keps/master/content"

	metadataFilename = "metadata.yaml"
)
