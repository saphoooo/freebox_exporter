# Changelog

## [unreleased] - 2020-05-27

- Add build script
- Log Freebox Server uptime and firmware version
- Log G.INP data
- Log connection status, protocol and modulation
- Log XDSL stats
- Don't log incorrect values (previously logged as zero)
- Remove dead code

## [1.1.9] - 2020-04-24

- Log freeplug speeds and connectivity
- Go 1.14
- Remove Godeps and vendored files

## [1.1.7] - 2019-09-04

- Go 1.13

## [1.1.6] - 2019-09-04

- There is no more uncomfortable error message when the application renews its token
- Adding a `-fiber` flag so that Freebox fiber users do not capture DSL metrics, which are empty on this type of Freebox

## [1.1.4] - 2019-09-03

- Adding a `-debug` flag to have more verbose error logs

## [1.1.2] - 2019-08-25

- Improve error messages

## [1.1.1] - 2019-07-31

- New Dockerfile for amd64 arch: reduce the image size to 3mb

## [1.1.0] - 2019-07-31

- Fix temp metrics
- Add Godeps

## [1.0.1] - 2019-07-30

- Change error catching

## [1.0.0] - 2019-07-29

- Rewriting the application by adding a ton of unit tests
