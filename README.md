# A Library for Interacting with KEPs

Contained is a library for programatically interacting with Kubernetes Enhancement Proposal (KEP)
content. At a high level the library is organized as follows

```
.
├── go.mod
├── go.sum
├── helpers
│   ├── convert             (attempt to convert a flat file KEP to the new directory structure)
│   ├── initSigDirs         (create a playground for experimenting with this library)
│   └── renderSigList       (regenerate Kubernetes SIG information for this library)
├── implementation_plan     (a poorly maintained list of high level TODOs)
├── LICENSE
├── okrs                    (goals for this project)
│   └── 2018
├── pkg
│   ├── filter              (finding KEPs which match given criteria)
│   ├── index               (a high level summary of all KEPs)
│   ├── keps                (the KEP object model)
│   ├── porcelain           (interacting with Git repositories on GitHub)
│   ├── settings            (configuration for this library)
│   ├── sigs                (basic Kubernetes SIG information)
│   └── workflow            (management of a single KEP)
├── teaching_notes.md       (longer explainations of concepts used in the library)
└── wish_list.md            (ideas for new contributors)
```
