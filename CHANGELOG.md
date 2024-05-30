# CHANGELOG

## Unreleased

*Release date*

### API BREAKING

### BUG FIXES

* Ensure correct version of `protoc-gen-gocosmos` is installed ([#41](https://github.com/atomone-hub/govgen/pull/41)).

### DEPENDENCIES

### FEATURES

### STATE BREAKING

## v1.0.2

*May 25st, 2024*

### FEATURES

* Add reproducible builds ([#34](https://github.com/atomone-hub/govgen/pull/34))

## v1.0.1

*February 26st, 2024*

### FEATURES

* `InitChainer` auto stakes uniformly to validators at genesis ([#26](https://github.com/atomone-hub/govgen/pull/26)).

## v1.0.0

*February 21st, 2024*

### FEATURES

* Forked `x/gov` module from Cosmos SDK v0.46.16 for custom modifications
  ([#4](https://github.com/atomone-hub/govgen/pull/4)).
* Removed ability for validators to vote on proposals with delegations, they can only use their own stake
  ([#20](https://github.com/atomone-hub/govgen/pull/20)).
* Remove community spend proposal
  ([#13](https://github.com/atomone-hub/govgen/pull/13)).
* Adapt voting period to proposal type.
  ([#16](https://github.com/atomone-hub/govgen/pull/16)).
