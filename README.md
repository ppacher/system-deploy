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

As you can see, the task file starts with a <b><i>Task</i></b> section that contains a human readable description and may contain additional metadata for the task. The discription above is mainly for logging purposes. `StartMasked= no` tells system-deploy that this task is not masked from execution (masked = disabled) and will be executed. This is the default any only here for documentation purposes.  

## Contributing

Any contributions to the `system-deploy` project are welcome! Just fork the repository and create a PR with your changes. It's recommended to discuss planned changes in an [issue](https://github.com/ppacher/system-deploy/issues) first.

## License

`system-deploy` is licensed under a BSD 3-clause license. See [LICENSE](LICENSE) for more information.
