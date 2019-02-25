# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

*NOTE*: Include in the changelog only what is interesting for the business and
service runtime.  
For instance, fixing the documentation and lintin SHOULD not be
included in the changelog document.

## [Unreleased] 

- Added in travisci and GolangCI linter
- Fix vet warnings

## [0.0.1] - 2019-02-17

Initial release of the orderbook library.
The functionality was used before as library.

- Standard price-time priority
- Supports both market and limit orders
- Supports order cancelling
- High performance (above 300k trades per second)
- Optimal memory usage
