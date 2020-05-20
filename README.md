# system-deploy

*A systemd inspired system configuration and deployment tool*

[![Go](https://github.com/ppacher/system-deploy/workflows/Go/badge.svg)](https://github.com/ppacher/system-deploy/actions)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ppacher/system-deploy?style=flat-square)
![Stable: Beta](https://img.shields.io/badge/Stable-BETA-yellowgreen?style=flat-square)
![GitHub](https://img.shields.io/github/license/ppacher/system-deploy?style=flat-square)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/ppacher/system-deploy?label=Release&style=flat-square)

`system-deploy` is my personal server management and deployment tool. It's inspired by systemd's unit files and the deployment tool/script used @safing. It currently supports copying files/directories, installing packages and installing/enabling systemd unit files. `system-deploy` is meant to be executed periodically and supports running different actions and tasks when changes are detected.

The compiled binary itself includes help and documentation for almost all supported operations and even some examples.

## Installation

`system-deploy` is written in <img src="https://golang.org/lib/godoc/images/go-logo-blue.svg" alt="Go" height="13"> and thus requires a working go installation to be compiled from source.

```bash
go install github.com/ppacher/system-deploy/cmd/deploy
```
This will install the `deploy` command into your `$GOBIN` or `$GOPATH/bin` if the former is not set.

## Getting Started

`system-deploy` is modeled around `.task` files that follow systemd's unit file syntax. Those tasks can perform one or
even multiple actions on the system and may trigger subsequent actions or tasks. Here is an example of a task file that copies a crontab file and reloads the cron daemon if it has changed (i.e. a new version has been deployed):

```
[Task]
Description= Install cronfile and restart cronie
StartMasked = no

[Copy]
Source= ./crontab-backup
Destination= /etc/cron.d/
FileMode= 0600

[OnChange]
Run=systemctl restart cron
```

As you can see, the task file starts with a <b><i>Task</i></b> section that contains a human readable description and may contain additional metadata for the task. The discription above is mainly for logging purposes. `StartMasked= no` tells system-deploy that this task is not masked from execution (masked = disabled) and will be executed. This is the default and only here for documentation purposes. If set to yes (`StartMasked=yes`) then the task would be disabled and not executed if not enabled by another task (for example, by using `[OnChange] Unmask=<name>.task`).  

Following the <b><i>Task</i></b> section, one or more actions can be declared, each in it's own section. The example above defines the <b><i>Copy</i></b> action which ensures `./crontab-backup` is copied and kept update in `/etc/cron.d` with 0600 set as the file permissions.  

The next section is <b><i>OnChange</i></b> which can perform different actions depending on the outcome of all other sections of the task. This is another interesting concept in system-deploy. Each action section, after being executed, can set a `changed` flag on the task meaning that it actually performed some action. The <b><i>Copy</i></b> action for example performs a hash (Murmur3) comparsion of the target path to determine if it's already up-to-date or needs to be overwritten/created.

`system-deploy` currently contains a few built-in actions (see below) and is designed to be extended easily.

## Built-in Actions

`system-deploy` currently supports the following actions:

- [exec](./pkg/actions/builtin/exec/exec.md) allows you to execute arbitrary commands
- [copy](./pkg/actions/builtin/copy/copy.md) allows you to copy file and directories
- [onchange](./pkg/actions/builtin/onchange/onchange.md) supports running and modifying tasks if the current task changed something 
- [installpackages](./pkg/actions/builtin/platform/installpackages.md) allows you to install packages.
- [platform](./pkg/actions/builtin/platform/platform.md) allows you to enable/disable tasks based on the platform (package manager, os-release, ...)
- [systemd](./pkg/actions/builtin/systemd/systemd.md) allows you to install and enable systemd unit files.
- [editfile](./pkg/actions/builtin/editfile/editfile.md) allows you to manipulate files using SED like syntax.

## Documentation

Additional documentation on `system-deploy` can be found in the [docs](./docs/README.md) or at `https://system-deploy.ppacher.at`.

## Contributing

Any contributions to the `system-deploy` project are welcome! Just fork the repository and create a PR with your changes. It's recommended to discuss planned changes in an [issue](https://github.com/ppacher/system-deploy/issues) first.

## License

`system-deploy` itself is licensed under a BSD 3-clause license. See [LICENSE](LICENSE) for more information.

Note that the bineries distributed via the release page are licensed under GPL-3 because the `EditFile` action is compiled against a GPL-3 licensed library.
