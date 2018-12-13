/*
Package hermetic provides facilities for proposing changes to a Git repository
hosted on GitHub. The proposed changes, ultimatedly delivered as a GitHub Pull
Request, are "hermetic" in the sense that the upstream repository will be
cloned to a unique location and changes pushed to a user GitHub repository
will overwrite any changes through force pushing.

This package was designed to make it easy to connect the KEP process to the
enhancements tracking process

	toLocation := filepath.Join(os.TempDir(), "an-unused-path")
	withBranchName := "kep-title-targets-release-1.15"

	whereGitHubTokenExists := "good/path/for/a/secret"
	githubHandle, err := settings.FindPrincipal()
	if err != nil {
		// graceful handling
	}

	token, err := NewProvideTokenFromPath(whereGitHubTokenExists)
	if err != nil {
		// graceful handling
	}

	kubernetes := "kubernetes"
	enhancementsTracking := "enhancements-tracking"

	repo, err := hermetic.Fork(githubHandle, token, kubernetes, enhancementsTracking, toLocation, withBranchName)
	if err != nil {
		// graceful handling
	}

	kepSummaryLocation := "path-to-kep-summary.md"
	releaseTrackingDir := "release-1.15/proposed"

	kepSummaryBytes, err := ioutil.ReadFile(kepSummaryLocation)
	if err != nil {
		// graceful handling
	}

	err = repo.Add(kepSummaryLocation, releaseTrackingDir)
	if err != nil {
		// graceful handling
	}

	prUrl, err := repo.CreatePR("Track delivery of KEP title", string(kepSummaryBytes))
	if err != nil {
		// graceful handling
	}
	
	fmt.Printf("proposed tracking KEP in PR: %s\n", prUrl)	
*/
package hermetic
