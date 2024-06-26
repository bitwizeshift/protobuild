# ![`protobuild`](docs/src/logo.png)

![Continuous Integration](https://img.shields.io/github/actions/workflow/status/bitwizeshift/protobuild/.github%2Fworkflows%2Fcontinuous-integration.yaml?logo=github)
[![GitHub Release](https://img.shields.io/github/v/release/bitwizeshift/protobuild?include_prereleases)][github-releases]
[![Gitter Channel](https://img.shields.io/badge/matrix-%23protobuild-darkcyan?logo=gitter)][gitter-channel]
[![readthedocs](https://img.shields.io/badge/docs-readthedocs-blue?logo=readthedocs&logoColor=white)][docs]
[![Godocs](https://img.shields.io/badge/docs-godocs-blue?logo=go&logoColor=white)][go-docs]

[gitter-channel]: https://matrix.to/#/#protobuild:gitter.im
[docs]: https://bitwizeshift.github.io/protobuild/
[go-docs]: https://bitwizeshift.github.io/protobuild/pkg/github.com/bitwizeshift/protobuild
[github-releases]: https://github.com/bitwizeshift/protobuild/releases

> [!NOTE]
> `protobuild` is a **work-in-progress** tool, so some documentation may be
> out-of-date or may represent the desired future state of the tool until the
> initial `v1.0` release.

Protobuild is the missing coordinator/build-system for [protobuf] projects.

This offers an easy, data-driven build-system for generating protobuf
definitions.

`protobuild` is **free** for teams and companies to use, and always will be.
No subscription is required, like _some other tools_.

> [!IMPORTANT]
> This is not an official supported Google product.

## Teaser

> [!NOTE]
> This section contains output from content that is not yet on `master`.

```bash
$ protobuild generate my-project
target my-project depends on protocolbuffers/protobuf, google/fhir

protocolbuffers/protobuf
  go .............................. ✔ (skipped)
  go-grpc ......................... ✔ (not-needed)
  python .......................... ✔
  c++ ............................. ✔ (skipped)

google/fhir
  go .............................. ✔ (skipped)
  go-grpc ......................... ✔ (not-needed)
  python .......................... ✔
  c++ ............................. ✔ (skipped)

my-project
  go .............................. ✔
  go-grpc ......................... ✔
  python .......................... ✔
  c++ ............................. ✔

Generation successful
```

## Quick Links

* [❓ Why?](#why)
* [🥅 Project Goals](#project-goals)
* [📖 Documentation](https://bitwizeshift.github.io/protobuild)
  * [📦 Installation](https://bitwizeshift.github.io/protobuild)
  * [🚀 Getting Started](https://bitwizeshift.github.io/protobuild)
  * [🙋‍♂️ Contributing](https://bitwizeshift.github.io/protobuild)
* [⚙️ Go Docs](https://bitwizeshift.github.io/protobuild/pkg/github.com/bitwizeshift/protobuild)
* [⚖️ License](#license)

## Why?

Anyone who has worked in a large project with a lot of protocol buffers knows
what a pain it is to coordinate the generation of protobuf definitions. It
becomes more complicated as more target-languages and plugins get into the mix,
and overall has no _good_ solution to this.

Some products exist in the market, like [`buf`], which privatizes and
centralizes registries behind a **paid subscription**. This forces a separate
and orthogonal system for something that can easily be done locally; and
`protobuild` aims to make this as easily and painlessly as possible.

## Project Goals

The `protobuild` project aims to provide a safe, easy, and **free** mechanism
to control building and generating Protocol buffers.

Below are some goals of this project:

* [ ] Allow user-definitions of external protobuf projects
* [ ] Allow custom `git`-driven registries of public protobuf projects, so that
      teams may centralize their definitions.
* [ ] Dependency tracing of protobuf projects.
* [ ] Manage recipes for `protoc` plugin installations to better centralize this.
* [ ] GitHub action support
* [ ] Provide hosted JSON Schema definitions of the project files

## License

This project is **dual-licensed** under both MIT and Apache-2.0, at the
user's choice.

[protobuf]: https://protobuf.dev "Protocol Buffers"
[`buf`]: https://buf.build/ "buf.build"
