/*
Package hermetic provides facilities for proposing changes to a Git repository
hosted on GitHub. The proposed changes, ultimatedly delivered as a GitHub Pull
Request, are "hermetic" in the sense that the upstream repository will be
cloned to a unique location and changes pushed to a user GitHub repository
will overwrite any changes through force pushing.
*/
package hermetic
