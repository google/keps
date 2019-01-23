package keps

import (
	"path/filepath"

	"github.com/calebamiles/keps/pkg/sigs"
)

// Path returns the path to a possible KEP directory or an error if
// SIG routing information cannot be determined
//
// example:
//
//	contentRoot := "/tmp/keps-sandbox/content"
//
//	(1)
//	kepDir, _ := Path(contentRoot, "large-value-delivered-incrementally")
//	println(kepDir) // /tmp/keps-sandbox/content/kubernetes-wide/large-value-delivered-incrementally
//
//	(2)
//	kepDir, _ := Path(contentRoot, "/tmp/keps-sandbox/content/large-value-delivered-incrementally")
//	println(kepDir) // /tmp/keps-sandbox/content/kubernetes-wide/large-value-delivered-incrementally
//
//	(3)
//	kepDir, _ := Path(contentRoot, "sig-node/functional-value-delivered-incrementally")
//	println(kepDir) // /tmp/keps-sandbox/content/sig-node/sig-wide/functional-value-delivered-incrementally
//
//	(4)
//	kepDir, _ := Path(contentRoot, "/tmp/keps-sandbox/content/sig-node/sig-wide/functional-value-delivered-incrementally")
//	println(kepDir) // /tmp/keps-sandbox/content/sig-node/sig-wide/functional-value-delivered-incrementally
//
//	(5)
//	kepDir, _ := Path(contentRoot, "sig-node/kubelet/kubelet-specific-value-delivered-incrementally")
//	println(kepDir) // /tmp/keps-sandbox/content/sig-node/kubelet/kubelet-specific-value-delivered-incrementally
//
//	(6)
//	kepDir, _ := Path(contentRoot, "/tmp/keps-sandbox/content/sig-node/kubelet/kubelet-specific-value-delivered-incrementally")
//	println(kepDir) // /tmp/keps-sandbox/content/sig-node/kubelet/kubelet-specific-value-delivered-incrementally
//
//	an absolute path is returned with the intention of using the result as input to standard
//	library calls to functions such as os.Open so that caller location is irrelevant
func Path(contentRoot string, p string) (string, error) {
	switch {
	case filepath.IsAbs(p):
		// addresses examples: {2, 4, 6}
		_, err := filepath.Rel(contentRoot, p) // create error if p cannot be made relative to contentRoot
		if err != nil {
			return "", err
		}

		return p, nil
	default:
		// address examples: {1, 3, 5}

		routingInfo, err := sigs.BuildRoutingFromPath(contentRoot, p)
		if err != nil {
			return "", err
		}

		return routingInfo.ContentDir(), nil
	}
}
