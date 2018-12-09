### Low Hanging Fruit

- create sketch handling the granting of exemptions. One (potentially) easy idea to play
  around with the UX would be to have each metadata check in the `keps/check` package
  return a typed error and let some exemption handler filter exempted errors out

  ```
  switch err {
  case exemptable.OwnerIsApprover:
	 return nil
  default:
	 return err
  }
  ```

  where:
    - exemptable errors are stored in their own `exemptable` package
    - an exemption handler is attached to the KEP and intercepts the result of
      the `#Check()`call

  ...just a thought to get going quickly

- create sketch of KEP events. Each KEP command should leave a record in the KEP metadata of being run
  eventually one could imagine computing interesting metrics about how long a KEP spends in each stage,
  or extending the event sources to include things like associated PRs merging. A reasonable place to
  start could be defining the event schema

  ```
  type kepEvent struct {
  	Principal string    `yaml:"principal"`
  	Time      time.Time `yaml:"time"`
  	Type      string    `yaml:"event_type"`
  }
  ```

  and a small set of compiled in events (e.g `{init, propose, accept, plan, approve} run`).

- vendor in hugo as part of a more general `keps-cli render` command that does the static site generation
  and builds a top level index for git/GitHub browsing and possibly other downstream tools in the future.
  It would also be great to get CD working that does the publishing automatically after new commits.
- logging: basically it's needed everywhere and needs to move to a structured format where it exists today
- testing: more are always needed
- integration with GitHub Projects for the API Review process: it should be possible to request/assign a specific
  reviewer who has a project board column dedicated to their state (e.g thockin/backlog, thockin/in-review). It is
  unlikely that there will be more than O(20) people doing API review anytime soon so it's probably fine to just
  hard code the reviewer <-> column mapping and save some tokens that would have otherwise been spent looking the
  information up (one could also imagine writing a simple generator for this in a better way than the SIG info
  generator works)
- basically anywhere there is a `TODO`
- extend idea of starting an enhancement to the general case where a KEP has multiple enhancements
  under `enhancements/`

### Higher Hanging Fruit

- prow integration for association of GitHub Issues and Pull Requests
