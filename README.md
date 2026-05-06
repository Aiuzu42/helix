# helix

> **Note:** This is a fork of [nicklaw5/helix](https://github.com/nicklaw5/helix)

A Twitch Helix API client written in Go (Golang).

[![Go Reference](https://pkg.go.dev/badge/github.com/Aiuzu42/helix.svg)](https://pkg.go.dev/github.com/Aiuzu42/helix/v2)

Twitch is always expanding and improving the available endpoints and features for the Helix API.
The maintainers of this package will make a best effort approach to implementing new changes
as they are released by the Twitch team.

See [here](SUPPORTED_ENDPOINTS.md) for a list of endpoints and features this package supports.

## Documentation & Examples

All documentation and usage examples for this package can be found [here](docs/README.md).
If you are looking for the Twitch API docs, see the [Twitch Developer website](https://dev.twitch.tv/docs/api).

## Support

Have a question? Need some assistance? Check out our dedicated channel in the
[Twitch API Discord](https://discord.gg/8HKVrmzczH).

## Supported Go Versions

Our support of Go versions is aligned with [Go's version release policy](https://golang.org/doc/devel/release#policy).
So we will support a major version of Go until there are two newer major releases.
We no longer support building this package with unsupported Go versions, as these contain security
vulnerabilities which will not be fixed.

## Contributions

PRs are very much welcome.
Where possible, please include unit tests for any code that is introduced by your PRs.
It's also helpful if you can include usage examples in the [docs](docs) directory.

## License

This package is distributed under the terms of the [MIT](LICENSE) license.
