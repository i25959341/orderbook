# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

*NOTE*: Include in the changelog only what is interesting for the business and
service runtime.  
For instance, fixing the documentation and lintin SHOULD not be
included in the changelog document.

## [0.2.5] - 2019-03-13

- Fix order done price for limit order

## [0.2.1] - 2019-03-13

- Added Depth method to get price levels for bids and asks
- Added Order method to get Order object by ID

## [0.2.0] - 2019-03-13

- Added parial quantity processed return value to Market and Limit orders

## [0.1.0] - 2019-03-01

- Added in travisci and GolangCI linter
- Fix vet warnings
- Added json.Marshaler and json.Unmarshaler interfaces
- Added total market price calculation for definite quantity (CalculateMarketPrice)

## [0.0.1] - 2019-02-17

Initial release of the orderbook library.
The functionality was used before as library.

- Standard price-time priority
- Supports both market and limit orders
- Supports order cancelling
- High performance (above 300k trades per second)
- Optimal memory usage
