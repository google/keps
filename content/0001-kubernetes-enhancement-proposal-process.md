# Kubernetes Enhancement Proposal Process

## Table of Contents
* [Kubernetes Enhancement Proposal Process](#kubernetes-enhancement-proposal-process)
   * [Table of Contents](#table-of-contents)
   * [Summary](#summary)
   * [Motivation](#motivation)
   * [Guide for Developers](#guide-for-developers)
      * [What type of work should be tracked by a KEP](#what-type-of-work-should-be-tracked-by-a-kep)
      * [KEP Directory Structure](#kep-directory-structure)
      * [KEP Metadata](#kep-metadata)
      * [KEP Workflow](#kep-workflow)
      * [Git and GitHub Implementation](#git-and-github-implementation)
      * [When to Merge](#when-to-merge)
      * [How to assign a number](#how-to-assign-a-number)
      * [KEP Editor Role](#kep-editor-role)
      * [Important Metrics](#important-metrics)
      * [Prior Art](#prior-art)
   * [Drawbacks](#drawbacks)
   * [Alternatives](#alternatives)
      * [Github issues vs. KEPs](#github-issues-vs-keps)
   * [Unresolved Questions](#unresolved-questions)
   * [Metadata](#metadata)

## Summary

A standardized development process for Kubernetes is proposed in order to

- provide a common structure for proposing changes to Kubernetes
- ensure that the motivation for a change is clear
- allow for the enumeration stability milestones and stability graduation
  criteria
- persist project information in a Version Control System (VCS) for future
  Kubernauts
- support the creation of _high value user facing_ information such as:
  - an overall project development roadmap
  - motivation for impactful user facing changes
- reserve GitHub issues for tracking work in flight rather than creating "umbrella"
  issues
- ensure community participants are successfully able to drive changes to
  completion across one or more releases while stakeholders are adequately
  represented throughout the process

This process is supported by a unit of work called a Kubernetes Enhancement Proposal or KEP.
A KEP attempts to combine aspects of a

- feature, and effort tracking document
- a product requirements document
- design document

into one file which is created incrementally in collaboration with one or more
Special Interest Groups (SIGs).

## Motivation

For cross project SIGs such as SIG PM and SIG Release an abstraction beyond a
single GitHub Issue or Pull request seems to be required in order to understand
and communicate upcoming changes to Kubernetes.  In a blog post describing the
[road to Go 2][], Russ Cox explains

> that it is difficult but essential to describe the significance of a problem
> in a way that someone working in a different environment can understand

as a project it is vital to be able to track the chain of custody for a proposed
enhancement from conception through implementation.

Without a standardized mechanism for describing important enhancements our
talented technical writers and product managers struggle to weave a coherent
narrative explaining why a particular release is important. Additionally for
critical infrastructure such as Kubernetes adopters need a forward looking road
map in order to plan their adoption strategy.

The purpose of the KEP process is to reduce the amount of "tribal knowledge" in
our community. By moving decisions from a smattering of mailing lists, video
calls and hallway conversations into a well tracked artifact this process aims
to enhance communication and discoverability.

A KEP is broken into sections which can be merged into source control
incrementally in order to support an iterative development process. An important
goal of the KEP process is ensuring that the process for submitting the content
contained in [design proposals][] is both clear and efficient. The KEP process
is intended to create high quality uniform design and implementation documents
for SIGs to deliberate.

[road to Go 2]: https://blog.golang.org/toward-go2
[design proposals]: /contributors/design-proposals


## Guide for Developers

### What type of work should be tracked by a KEP

The definition of what constitutes an "enhancement" is a foundational concern
for the Kubernetes project. Roughly any Kubernetes user or operator facing
enhancement should follow the KEP process: if an enhancement would be described
in either written or verbal communication to anyone besides the KEP author or
developer then consider creating a KEP.

Similarly, any technical effort (refactoring, major architectural change) that
will impact a large section of the development community should also be
communicated widely. The KEP process is suited for this even if it will have
zero impact on the typical user or operator.

As the local bodies of governance, SIGs should have broad latitude in describing
what constitutes an enhancement which should be tracked through the KEP process.
SIGs may find that helpful to enumerate what _does not_ require a KEP rather
than what does. SIGs also have the freedom to customize the KEP template
according to their SIG specific concerns. For example the KEP template used to
track API changes will likely have different subsections than the template for
proposing governance changes. However, as changes start impacting other SIGs or
the larger developer community outside of a SIG, the KEP process should be used
to coordinate and communicate.

Enhancements that have major impacts on multiple SIGs should use the KEP process.
A single SIG will own the KEP but it is expected that the set of approvers will span the impacted SIGs.
The KEP process is the way that SIGs can negotiate and communicate changes that cross boundaries.

KEPs will also be used to drive large changes that will cut across all parts of the project.
These KEPs will be owned by SIG-architecture and should be seen as a way to communicate the most fundamental aspects of what Kubernetes is.

### KEP Directory Structure

The canonical format for a KEP is changing from a
[single flat file](0000-kep-template.md) checked into source control to a
directory structure in order to

- support the extraction of metadata required for web rendering of KEPs from
  the metadata for the KEP process itself
- provide a location for diagrams or other associated content
- provide a location to store experience reports associated with a KEP
- make it easier to build tooling to validate KEP metadata
- further reinforce the idea that KEP process is iterative by breaking up the
  sections of a KEP into individual files
- allow KEP authors, reviewers, approvers, and editors to take advantage of the
  `OWNERS` process

A typical KEP under active development will have have the following directory
structure

```
├── 001-introduction.md (required)
├── 002-motivation.md (required)
├── 003-goals.md
├── 004-non-goals.md
├── 005-guide-for-developers.md (required)
├── 006-guide-for-teachers.md (required)
├── 007-guide-for-operators.md (required)
├── 008-guide-for-project-maintainers.md
├── 009-graduation-criteria.md (required)
├── 010-implementation-history.md
├── experience_reports (required)
│   ├── alpha-feedback.md
│   └── beta-feedback.md
├── features
│   ├── feature-01.md
│   ├── feature-02.md
│   └── feature-03.md
├── index.json (required)
├── metadata.yaml (required)
└── OWNERS
```

In most cases the directory structure should be created by dedicated tooling.


### KEP Metadata

There is a place in each KEP for a YAML document that has standard metadata.
This will be used to support tooling around filtering and display. It is also
critical to clearly communicate the status of a KEP. As KEPs are migrated from
a single flat file to the [Directory Structure](#kep-directory-structure) the
metadata will be extracted to a `metadata.yaml` file in the directory for the
KEP. KEP metadata has the following format

```
---
authors: # required
  - "calebamiles" # just a GitHub handle for now
  - "jbeda"
title: "Kubernetes Enhancement Proposal process"
number: 42 # required, use the number from the pull rquest marking the KEP as "accepted"
owning-sig: "sig-pm" # required
participating-sigs:
  - "sig-architecture"
  - "sig-contributor-experience"
approvers: # required
  - "bgrant0607" # just a GitHub handle for now
reviewers:
  - "justaugustus"  # just a GitHub handle for now
  - "jdumars"
editors:
  - null # generally omit empty/null fields
status: "active" # required
github:
  issues:
    - null # GitHub url
  pull_requests:
    - null # GitHub url
  projects:
    - project_id: null
      card_id: null
releases: # required
  - k8s_version: v1.9
    kep_status: "active"
    k8s_status: "alpha" # one of alpha|beta|GA
  - k8s_version: v1.10
    kep_status: "active"
    k8s_status: "alpha"
replaces:
  - kep_location: null
superseded-by:
  - kep_location: null
created: 2018-01-22 # in YYYY-MM-DD
updated: 2018-09-04
```

### KEP Workflow

A KEP has the following states

- `draft`: the author(s) would like to document some agreement by checking in a
  partial KEP but the content has not been approved by the SIG for implementation
- `accepted`: The owning SIG agrees with the `Motivation` and is committing to
  supporting further iteration on the KEP before implementation
- `implementable`: The approvers have approved the design and draft
  documentation of the KEP for implementation.
- `active`: The KEP is under active development
- `retired`: The functionality described by the KEP has been retired by the
  project
- `deferred`: The KEP is proposed but will not be considered for approval at
  this time.
- `rejected`: The approvers and authors have decided that this KEP is not moving
  forward.
  The KEP is kept around as a historical document.
- `withdrawn`: The KEP has been withdrawn by the authors.
- `superseded`: The KEP has been superseded by a new KEP.
  The `superseded-by` metadata value should point to the new KEP.

A KEP will typically take one of the following paths through state space


1. `draft -> accepted -> implementatble -> active`
1. `draft -> accepted -> implementatble -> active -> retired`
1. `draft -> rejected`
1. `draft -> deferred`

and the initial PR for a KEP should set the status to either `draft` or
`accepted` unless through previous communication with the SIG it has been agreed
to accept a fully formed KEP document in the initial PR. When merging a full
KEP in a single PR the status should likely be set to `implementable` or
`active` depending on whether work is already in flight.

### Git and GitHub Implementation

KEPs are checked into the community repo under the `/kep` directory.
In the future, as needed we can add SIG specific subdirectories.
KEPs in SIG specific subdirectories have limited impact outside of the SIG and can leverage SIG specific OWNERS files.

New KEPs can be checked in with a file name in the form of `draft-YYYYMMDD-my-title.md`.
As significant work is done on the KEP the authors can assign a KEP number.
This is done by taking the next number in the NEXT_KEP_NUMBER file, incrementing that number, and renaming the KEP.
No other changes should be put in that PR so that it can be approved quickly and minimize merge conflicts.
The KEP number can also be done as part of the initial submission if the PR is likely to be uncontested and merged quickly.

### When to Merge

In general Kubernetes uses a [lazy consensus][] process for decision making. The
correct time to merge a pull request is whenever consensus has been reached on a
section of a KEP with the goal being that large requests don't remain open
indefinitely. It is expected that a decision on the `Introduction` and `Motivation`
sections should happen relatively quickly while the detailed design documents
(`guide-for-{developers, teacher, operators, project-maintainers}.md`) may take
longer to merge as consensus forms.

[lazy consensus]: https://openoffice.apache.org/docs/governance/lazyConsensus.html

### How to assign a number

The `kep-number` should be taken from the pull request number which marks the
KEP as `accepted` which is a handy source of monotonically increasing numbers
and will hopefully mitigate "number squatting".


### KEP Editor Role

Taking a cue from the [Python PEP process][], we define the role of a KEP editor.
The job of an KEP editor is likely very similar to the [PEP editor responsibilities][] and will hopefully provide another opportunity for people who do not write code daily to contribute to Kubernetes.

In keeping with the PEP editors which

> Read the PEP to check if it is ready: sound and complete. The ideas must make
> technical sense, even if they don't seem likely to be accepted.
> The title should accurately describe the content.
> Edit the PEP for language (spelling, grammar, sentence structure, etc.), markup
> (for reST PEPs), code style (examples should match PEP 8 & 7).

KEP editors should generally not pass judgement on a KEP beyond editorial corrections.
KEP editors can also help inform authors about the process and otherwise help things move smoothly.

[Python PEP process]: https://www.python.org/dev/peps/pep-0001/
[PEP editor responsibilities]: https://www.python.org/dev/peps/pep-0001/#pep-editor-responsibilities-workflow

### Important Metrics

It is proposed that the primary metrics which would signal the success or
failure of the KEP process are

- how many "enhancements" are tracked with a KEP
- distribution of time a KEP spends in each state
- KEP rejection rate
- PRs referencing a KEP merged per week
- number of issued open which reference a KEP
- number of contributors who authored a KEP
- number of contributors who authored a KEP for the first time
- number of orphaned KEPs
- number of retired KEPs
- number of superseded KEPs

### Prior Art

The KEP process as proposed was essentially stolen from the [Rust RFC process][] which
itself seems to be very similar to the [Python PEP process][]

[Rust RFC process]: https://github.com/rust-lang/rfcs

## Drawbacks

Any additional process has the potential to engender resentment within the
community. There is also a risk that the KEP process as designed will not
sufficiently address the scaling challenges we face today. PR review bandwidth is
already at a premium and we may find that the KEP process introduces an unreasonable
bottleneck on our development velocity.

It certainly can be argued that the lack of a dedicated issue/defect tracker
beyond GitHub issues contributes to our challenges in managing a project as large
as Kubernetes, however, given that other large organizations, including GitHub
itself, make effective use of GitHub issues perhaps the argument is overblown.

The centrality of Git and GitHub within the KEP process also may place too high
a barrier to potential contributors, however, given that both Git and GitHub are
required to contribute code changes to Kubernetes today perhaps it would be reasonable
to invest in providing support to those unfamiliar with this tooling.

Expanding the proposal template beyond the single sentence description currently
required in the [features issue template][] may be a heavy burden for non native
English speakers and here the role of the KEP editor combined with kindness and
empathy will be crucial to making the process successful.

[features issue template]: https://git.k8s.io/features/ISSUE_TEMPLATE.md

## Alternatives

This KEP process is related to
- the generation of a [architectural roadmap][]
- the fact that the [what constitutes a feature][] is still undefined
- [issue management][]
- the difference between an [accepted design and a proposal][]
- [the organization of design proposals][]

this proposal attempts to place these concerns within a general framework.

[architectural roadmap]: https://github.com/kubernetes/community/issues/952
[what constitutes a feature]: https://github.com/kubernetes/community/issues/531
[issue management]: https://github.com/kubernetes/community/issues/580
[accepted design and a proposal]: https://github.com/kubernetes/community/issues/914
[the organization of design proposals]: https://github.com/kubernetes/community/issues/918

### Github issues vs. KEPs

The use of GitHub issues when proposing changes does not provide SIGs good
facilities for signaling approval or rejection of a proposed change to Kubernetes
since anyone can open a GitHub issue at any time. Additionally managing a proposed
change across multiple releases is somewhat cumbersome as labels and milestones
need to be updated for every release that a change spans. These long lived GitHub
issues lead to an ever increasing number of issues open against
`kubernetes/features` which itself has become a management problem.

In addition to the challenge of managing issues over time, searching for text
within an issue can be challenging. The flat hierarchy of issues can also make
navigation and categorization tricky. While not all community members might
not be comfortable using Git directly, it is imperative that as a community we
work to educate people on a standard set of tools so they can take their
experience to other projects they may decide to work on in the future. While
git is a fantastic version control system (VCS), it is not a project management
tool nor a cogent way of managing an architectural catalog or backlog; this
proposal is limited to motivating the creation of a standardized definition of
work in order to facilitate project management. This primitive for describing
a unit of work may also allow contributors to create their own personalized
view of the state of the project while relying on Git and GitHub for consistency
and durable storage.

## Unresolved Questions

- ~~How reviewers and approvers are assigned to a KEP~~ blunderbus like any
  other repo using the `OWNERS` process
- Example schedule, deadline, and time frame for each stage of a KEP
- Communication/notification mechanisms
- Review meetings and escalation procedure

## Metadata
```
---
kep-number: 1
title: Kubernetes Enhancement Proposal Process
authors:
  - "@calebamiles"
  - "@jbeda"
owning-sig: sig-architecture
participating-sigs:
  - sig-api-machinery
  - sig-apps
  - sig-auth
  - sig-aws
  - sig-azure
  - sig-cli
  - sig-cloud-provider
  - sig-cluster-lifecycle
  - sig-contributor-experience
  - sig-network
  - sig-node
  - sig-scheduling
reviewers:
  - name: "@timothysc"
approvers:
  - name: "@bgrant0607"
editor:
  name: "@jbeda"
creation-date: 2017-08-22
status: implementable
```
