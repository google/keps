
### Coding Philosophies

- The choice of whether to create an explicit internal contract on the
  narrow usage of an external dependency can be explored in tests at
  `pkg/keps/sections/sections_test.go` where the production usage of
  section addition is simulated by using a `metadatafakes.FakeKEP`, 
  presumably maintained by the `OWNERS` of `metadata.KEP`. The power
  of Go interfaces allows the sections package to have no knowledge
  of any package above it in the directory tree 
