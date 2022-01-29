# prolog
distributed transaction logger - based off work by Travis Jeffery from book "Distributed Services with Go"

The code is copy+pasted from https://github.com/travisjeffery/proglog with small edits here and there

I read a chapter and then integrate the code from the aforementioned repo (sometimes it doesn't line up perfectly)

I plan on using the experience earned here to build my own distributed systems for both fun and profit.

# next steps

Extract `distributed` from `log package`, create new package with just the `raft` coordination, then
make this into a GitHub template (or whatever you call it - make it a project template).

# notes

- there was some strangeness with the `log` package and wsl

- dependencies required
  - protobuf compiler and go plugin
  - kubernetes stuff (for deployment)
