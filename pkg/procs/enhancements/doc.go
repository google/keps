/*
Enhancements understands the details of connecting KEP lifecycle management with the
existing GitHub Pull Request review model used today in the Kubernetes community.
Specifically this package handles:
	- committing updated KEP content to the local clone of the enhancements repo
	- pushing local changes to a user fork of the enhancements repo
	- filing a pull request against the upstream enhancements repo
	- adding pull requests that advance a KEP in its lifecycle to its metadata
*/
package enhancements
