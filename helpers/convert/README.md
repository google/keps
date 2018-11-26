# KEP Conversion Helper

This small script is designed to help KEP process maintainers and SIG Leads
migrate existing KEP content to the new format described in [KEP1a][].

To run this tool:

```
go run to_current.go <workspace>/community/keps
```

where `<workspace>` is the location where you have previously cloned the [Kubernetes Community Repo][]. The tool
will create a temporary directory containing converted KEP content. A KEP process maintainer attempting to
convert the existing KEP content should

1. Referring to the documentation on [git filter-branch][], select only commits against the `keps` directory
   in clone of the Kubernetes Community Repo: 

   ```
   git filter-branch --prune-empty --subdirectory-filter keps
   ```

1. Create a directory named `content` to contain the existing KEPs

   ```
   mkdir content
   ```

1. Move the existing content to the new directory

   ```
   git mv sig-* content/
   git mv *.md content/
   ```

1. Run this tool against the newly created `content` directory

   ```
   go run to_current.go <workspace>/community/content
   ```

1. Fix errors in conversion, inspect converted result, and commit the result


Refer to the [GitHub documentation][] for more help with using `git filter-branch` to prepare a subdirectory for extraction
from a larger Git repository

[KEP1a]: https://github.com/kubernetes/community/blob/3b3f730761a7ab902672720ab1a254b1dd0aa387/keps/0001a-meta-kep-implementation.md
[Kubernetes Community Repo]: https://github.com/kubernetes/community/
[git filter-branch]: https://git-scm.com/docs/git-filter-branch
[GitHub documentation]: https://help.github.com/articles/splitting-a-subfolder-out-into-a-new-repository/
