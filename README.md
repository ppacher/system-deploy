<p align="center">
    <img src="docs/assets/logo.png">
</p>

*A systemd inspired system configuration and deployment tool*

[![Go](https://github.com/ppacher/system-deploy/workflows/Go/badge.svg)](https://github.com/ppacher/system-deploy/actions)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ppacher/system-deploy?style=flat-square)
![Stable: Beta](https://img.shields.io/badge/Stable-BETA-yellowgreen?style=flat-square)
![GitHub](https://img.shields.io/github/license/ppacher/system-deploy?style=flat-square)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/ppacher/system-deploy?label=Release&style=flat-square)

`system-deploy` is my personal server management and deployment tool. It's inspired by systemd's unit files and the deployment tool/script used @safing. It currently supports copying files/directories, installing packages and installing/enabling systemd unit files. `system-deploy` is meant to be executed periodically and supports running different actions and tasks when changes are detected.

The compiled binary itself includes help and documentation for almost all supported operations and even some examples.

**Checkout [system-conf](https://github.com/ppacher/system-conf) for a systemd inspired configuration system for Go projects.**

## Getting Started 

Please checkout the documentation at [ppacher.github.io/system-deploy](https://ppacher.github.io/system-deploy).

## Installation from Source

`system-deploy` is written in <img src="https://golang.org/lib/godoc/images/go-logo-blue.svg" alt="Go" height="13"> and thus requires a working go installation to be compiled from source.

```bash
go install github.com/ppacher/system-deploy/cmd/system-deploy
```
This will install the `deploy` command into your `$GOBIN` or `$GOPATH/bin` if the former is not set.


## Contributing

Any contributions to the `system-deploy` project are welcome! Just fork the repository and create a PR with your changes. It's recommended to discuss planned changes in an [issue](https://github.com/ppacher/system-deploy/issues) first.

## License

`system-deploy` itself is licensed under a BSD 3-clause license. See [LICENSE](LICENSE) for more information.

Note that the bineries distributed via the release page are licensed under GPL-3 because the `EditFile` action is compiled against a GPL-3 licensed library.
