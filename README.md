# system-deploy

*A systemd inspired system configuration and deployment tool*

`system-deploy` is my personal server management and deployment tool. It's inspired by systemd's unit files and the deployment tool/script used @safing. It currently supports copying files/directories, installing packages and installing/enabling systemd unit files. `system-deploy` is meant to be executed periodically and supports running different actions and tasks when changes are detected.

The compiled binary itself includes help and documentation for almost all supported operations and even some examples.

## Installation

`system-deploy` is written in <img src="https://golang.org/lib/godoc/images/go-logo-blue.svg" alt="Go" height="13"> and thus requires a working go installation to be compiled from source.

```bash
go install github.com/ppacher/system-deploy/cmd/deploy
```
This will install the `deploy` command into your `$GOBIN` or `$GOPATH/bin` if the former is not set.

## Concepts

`TBD`

## Contributing

Any contributions to the `system-deploy` project are welcome! Just fork the repository and create a PR with your changes. It's recommended to discuss planned changes in an [issue](https://github.com/ppacher/system-deploy/issues) first.

## License

`system-deploy` is licensed under a BSD 3-clause license. See [LICENSE](LICENSE) for more information.
