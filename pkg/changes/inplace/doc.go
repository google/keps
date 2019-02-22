/*

Package inplace coordinates changes to existing repositories. A primary use case
is submitting changes to KEP content throughout the lifecycle of a single KEP.
Package inplace is also a concrete demonstration of how feature branches can be
used to coordinate changes.

Proposing a KEP for SIG sponsorship is modeled as

	path, err := keps.Path(runtime.TargetDir())
	// graceful error handling

	kep, err := keps.Open(path)
	// graceful error handling

	err = kep.SetState(states.Provisional)
	// graceful error handling

	err = kep.Persist()
	// graceful error handling

	repo, err := inplace.Open(
		githubHandle,
		token,
		kubernetes,
		enhancements,
		runtime.ContentRoot(),
		kep.FeatureBranch()
	)

	// graceful error handling

	err = repo.Add(kep.ContentDir(), fmt.Sprintf("proposes %s for SIG sponsorship", kep.Title()))
	// graceful error handling

	summary, err := kep.Summary()
	// graceful error handling

	err = repo.CreatePR(fmt.Sprintf("proposes %s for SIG sponsorship", kep.Title()), summary.Content())
	// graceful error handling


while requesting SIG approval modeled, nearly identically, as

	path, err := keps.Path(runtime.TargetDir())
	// graceful error handling

	kep, err := keps.Open(path)
	// graceful error handling

	err = kep.SetState(states.Implementable)
	// graceful error handling

	err = kep.Persist()
	// graceful error handling

	repo, err := inplace.Open(
		githubHandle,
		token,
		kubernetes,
		enhancements,
		runtime.ContentRoot(),
		kep.FeatureBranch()
	)

	// graceful error handling

	err = repo.Add(kep.ContentDir(), fmt.Sprintf("requests SIG approval for %s", kep.Title()))
	// graceful error handling

	summary, err := kep.Summary()
	// graceful error handling

	prUrl, err = repo.CreatePR(fmt.Sprintf("requests SIG approval for %s", kep.Title()), summary.Content())
	// graceful error handling


where the only change for these two lifecycle events are the desired KEP state

*/
package inplace
