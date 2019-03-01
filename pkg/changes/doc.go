/*

Ch-ch-ch-ch-changes
Turn and face the strange
Ch-ch-changes

- David Bowie

The changes package is all about making changes to various parts of the Kubernetes project.
Concrete examples include requesting an API review, or requesting approval to enter the release
process. Submitting changes to the KEP content itself is also a primary motivator for this
package. Today all changes to related Kubernetes processes must be `hermetic` in the sense that a
clean copy of the location to submit the changes is created (e.g. forking a repsoitory). This is
in contrast to `inplace` changes which would involve only git (and not GitHub) operations in the
general case; submitting changes to KEP content through Prow is an example of a possible producer of
`inplace` changes.
*/
package changes
