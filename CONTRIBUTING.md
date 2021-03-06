# Contributing to Wash

## Code of Conduct

Review our [Code of Conduct](./CODE_OF_CONDUCT.md).

## How to ask a question

Browse for existing issues. Otherwise, open a "[new issue](https://github.com/puppetlabs/wash/issues/new)" in this repo.

## How to report a bug

Open a "[Bug report](https://github.com/puppetlabs/wash/issues/new?template=bug-report.md)" issue in this repo.

## How to suggest a new feature

Open a "[Feature request](https://github.com/puppetlabs/wash/issues/new?template=feature-request.md)" issue in this repo.

## Development Environment

### Updating Dependencies

Versioning for the kubernetes and docker projects don't work well with Go modules. The best way to update dependencies is to update specific packages to a specific version using `go get <mod>@<tag>`.

It should always be safe to run `go get -u=patch` to pickup patches.

### Modifying APIs

When making changes to Wash's APIs, remember to update the inline swagger documentation. Instructions for regenerating the API docs are [here](./website/README.md#regenerate-swagger-docs).

### Local filesystem testing

A developer-only plugin is included in Wash that references the local filesystem. It can be enabled by setting the `WASH_LOCALFS` environment variable to the local file path to mount before starting Wash. This can be useful for testing changes to core functionality like the `fuse` or `api` modules.

## Submitting Changes
Fork the repo, make changes, file a Pull Request.

Contributions to this project require sign-off consistent with the [Developers Certificate of Origin](https://developercertificate.org). This can be as simple as using `git commit -s` on each commit.

### Guidelines

We try to follow common Go conventions as enforced by the compiler and several static analysis tools run in GitHub Actions on pull requests and commits.

#### File naming

File naming should follow camelCase. When grouping several files around a single concept, multiple names can be joined by hyphens.

For example, the Docker plugin has multiple files implementing components of a container. We name them
- container-log.go
- container-metadata.go
- container.go

### Releases

Releases are done through the GitHub's [Draft a new release](https://github.com/puppetlabs/wash/releases/new). A [GitHub Action](https://github.com/puppetlabs/wash/blob/main/.github/workflows/release.yml) will build and upload assets and open a PR against https://github.com/puppetlabs/homebrew-puppet to update to the new release.

Release announcements are:

- the release is automatically posted to #wash by the GitHub App; ask a moderator to re-post into #announcements.
- Puppet internal [release announcement process](https://confluence.puppetlabs.com/display/PM/Sending+Product+Release+Announcements) sends to puppet-announce@ and puppet-users@
