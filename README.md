[![Build Status](https://travis-ci.org/alexanderGugel/nerva.svg?branch=master)](https://travis-ci.org/alexanderGugel/nerva)
[![GitHub tag](https://img.shields.io/github/tag/alexanderGugel/nerva.svg?maxAge=2592000)](https://github.com/alexanderGugel/nerva)

![cool](public/cool.svg)

# nerva (PoC)

`nerva` is a RESTful, `package.json`-aware registry server. It aims to offer a
fast and easy to maintain alternative to `npm`. The registry server can either
be run as a public-facing service or as a private hosting facility. Clients,
such as `npm` or `ied`, can locate packages hosted on a system running `nerva`.

## General

`nerva` aims to be a configurable, but easy to maintain registry server. `nerva`
is primarily designed for medium-sized teams trying to implement a unified
package management system, while keeping control over package distribution and
access management (e.g. in larger corporations that require an on-premise
solution).

## Design

`nerva` implements the [CommonJS] [1] registry specification that defines a
standardized way of **identifying** and **locating** packages.

### Unified hosting

npm introduces lots of redundant concepts that already exist in ordinary version
control systems, such as Git tags.

nerva doesn't try to re-implement git, instead it consumes git repositories and
exposes them in a way that allows them to be consumed by Common JS registry
compliant package managers (such as the npm CLI itself).

Here is a breakdown of commands / concepts / features introduced by npm and
their respective equivalents in nerva:

* `npm tag`

  nerva uses git tags to locate packages. Every git tag pushed to nerva's
  storage directory can be used in an npm installation. E.g.
  `git tag experimental && git push nerva --tags` instructs nerva to generate
  a corresponding package tag. Therefore `npm install package@experimental`
  "just works". No more `npm tag experimental`!

* `npm publish`

  nerva's single source of truth is its `storage` directory. Nevertheless,
  `nerva` doesn't make **any** assumptions about how the storage directory
  itself is being managed. Usually it would be exposed via a SSH server.
  In that case, publishing a package via `git push nerva` is one possible way to
  make the package consumable.

  Alternatively setting up a CRON job or GitHub webhook that does a `git fetch`
  - in more or less - regular intervals is a viable alternative. In that case
  nerva's storage directory acts as a Git mirror.

* `npm install`

  `npm install` relies on the Common JS registry specification, which nerva
  implements to a large extent. nerva also dynamically generates checksums of
  package tarballs, thus enabling clients to verify the integrity of an
  installed dependency. In other words, `npm install` works just like with any
  registry server.

### Git as a database

nerva doesn't have any external dependencies. As such, nerva uses the `storage`
directory defined via `nerva -storage=/storage` or `backend.storage`.

The storage directory should be a "flat" directory containing Git repositories.
Although not required, it is recommended to use bare repositories.

The name of each package within the `storage` can be mapped to a hosted package.
Therefore `storage` defaults to the `./packages` directory in the current
working directory.

For instance, a `storage` directory might have the following structure:

    packages
    └── tape
        └── .git
    └── ied
        └── .git

    1 directory

In this case, nerva would expose the `tape` and `ied` packages. Git's object
database would be used in order to dynamically generate tarballs for the
requested versions.

### Upstream registries

If users running `npm install` try to install a package which hasn't been
"pushed" to nerva, nerva optionally redirects, but doesn't proxy, incoming
requests to an alternative "upstream" registry.

The default upstream registry is the publicly-facing npm registry
(`http://registry.npmjs.com`).

To configure an upstream registry, specify the registry's root URL via the
`-upstreamURL` flag, e.g. via `nerva -upstreamURL=http://registry.npmjs.com`.

Alternatively, you can specify the upstream as a "backend" in your
`.nerva.[yaml|json|...]` configuration file:

    backend:
      upstreamURL: http://registry.npmjs.com

## Motivation

Dependency management in Node.js is broken.

The idea of having a separate registry server as a redundant hosting facility
for projects that are already hosted on services such as GitHub is unnecessary.

`npm publish` is an unneeded complication. You already tag releases of using
Git, there shouldn't be a need for manually publishing individual versions of a
hosted package. An ideal registry server shouldn't try to re-implement version
control, instead it should "infer" arbitrary versions of a package.

![npm registry](registry-wall.svg)

As a consequence of npm's current architecture, package authors are able to
publish arbitrary, potentially dangerous tarballs. While the current npm
registry is unlikely to go away anytime soon, running a private npm registry
leads to similar problems. nerva tries to offer a compelling alternative for
running a private Common JS registry (such as npm enterprise).

## Project Status

The project is in **pre-alpha** stage. While the registry server is in a usable
state, the documentation is currently insufficient and the setup process
tedious.

At this stage, `nerva` is mostly a proof of concept.

## Development Setup

As a compiled binary without any significant external dependencies, `nerva` is
pretty easy to setup and very configurable.

During development, clone down the repository and build the `nerva` command.
Dependencies are being managed via `godep`. `nerva` has a dependency on
[`git2go`](https://github.com/libgit2/git2go), which provides libgit2 bindings.
Make sure to install [`libgit2`](https://github.com/libgit2/git2go) (e.g. via
`brew install libgit2`) before building the binary.

1. Install `libgit2`
2. Clone down:

  `go get github.com/alexanderGugel/nerva`

3. Run

  `nerva`

`nerva` comes with sensible defaults, but is very configurable. During
development, having a `.nerva` config file is usually not recommended.

## Production Setup

**Warning** As mentioned in the project status, `nerva` is currently mainly a
proof of concept. If you encounter any issues, please file an
[issue](https://github.com/alexanderGugel/nerva/issues).

For production usage, use one of the provided releases.

### Credits

* [Cool] [2] by [JMA] [3] from [the Noun Project] [4]

[1]: http://wiki.commonjs.org/wiki/Packages/Registry "CommonJS Registry Specification"
[2]: https://thenounproject.com/Mattebrooks/collection/objecticons/?i=63757 "Cool"
[3]: https://thenounproject.com/jmanwyl "JMA"
[4]: https://thenounproject.com/ "The Noun Project"
