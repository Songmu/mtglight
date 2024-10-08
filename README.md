mtglight
=======

[![Test Status](https://github.com/Songmu/mtglight/workflows/test/badge.svg?branch=main)][actions]
[![Coverage Status](https://codecov.io/gh/Songmu/mtglight/branch/main/graph/badge.svg)][codecov]
[![MIT License](https://img.shields.io/github/license/Songmu/mtglight)][license]
[![PkgGoDev](https://pkg.go.dev/badge/github.com/Songmu/mtglight)][PkgGoDev]

[actions]: https://github.com/Songmu/mtglight/actions?workflow=test
[codecov]: https://codecov.io/gh/Songmu/mtglight
[license]: https://github.com/Songmu/mtglight/blob/main/LICENSE
[PkgGoDev]: https://pkg.go.dev/github.com/Songmu/mtglight

The mtglight turn on and off the [Yeelight](https://yeelight.com/) bulb when the online meeting starts and ends in cooperation with the [OverSight](https://objective-see.org/products/oversight.html).

## Installation

```console
# Install the latest version. (Install it into ./bin/ by default).
% curl -sfL https://raw.githubusercontent.com/Songmu/mtglight/main/install.sh | sh -s

# Specify installation directory ($(go env GOPATH)/bin/) and version.
% curl -sfL https://raw.githubusercontent.com/Songmu/mtglight/main/install.sh | sh -s -- -b $(go env GOPATH)/bin [vX.Y.Z]

# In alpine linux (as it does not come with curl by default)
% wget -O - -q https://raw.githubusercontent.com/Songmu/mtglight/main/install.sh | sh -s [vX.Y.Z]

# go install
% go install github.com/Songmu/mtglight/cmd/mtglight@latest
```

## Author

[Songmu](https://github.com/Songmu)
