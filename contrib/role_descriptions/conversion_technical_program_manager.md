# Technical Program Manager, Content Conversion

## Introduction

The [KEP Process][] is emerging as the primary mechanism for change management within
the Kubernetes Project. The technical foundation of the KEP Process is [munging][]
files that are checked into Git. The process was originally based on a single [flat file][],
however, the limitations of that approach have encouraged the project to move to a directory
based format. We therefore need someone to lead the conversion of existing and proposed content
to the new format.

[KEP Process]: https://github.com/kubernetes/enhancements/blob/master/keps/0001-kubernetes-enhancement-proposal-process.md
[munging]: https://en.wikipedia.org/wiki/Mung_(computer_term)
[flat file]: https://en.wikipedia.org/wiki/Flat-file_database

## Description of Role

The technical program manager for KEP content conversion will:

- maintain and extend the [conversion helper][]
- work with KEP authors to update KEPs which fail automatic conversion
- provide conversion status updates to the Kubernetes community; and develop
  automation to aid in such reporting
- ultimate responsibility for ensuring that the [existing][], and [proposed][],
  content is smoothly converted to the [updated KEP format][]

[conversion helper]: https://github.com/calebamiles/keps/tree/master/helpers/convert
[existing]: https://github.com/kubernetes/enhancements/tree/master/keps
[proposed]: https://github.com/kubernetes/enhancements/pulls
[updated KEP format]: https://github.com/calebamiles/keps/blob/master/pkg/keps/metadata/kep_metadata.go#L19

## Desired Experience

A successful maintainer for this role should:

- have experience with Go development. One need not be a "professional developer"
  but should be comfortable writing and small Go programs and extending large
  programs
- have experience with Git. Particularly with workflows such as [extracting][] a
  sub directory of a repository without losing the history of changes
- have experience with planning within the Kubernetes community. Specifically the
  maintainer will have likely served on a Release Team

[extracting]: https://help.github.com/en/articles/splitting-a-subfolder-out-into-a-new-repository

## Mentoring Offered

A successful maintainer will receive mentoring and guidance in:

- Go development
- planning in large open source projects

## Expected Time Commitment

5 - 10 hours per week
