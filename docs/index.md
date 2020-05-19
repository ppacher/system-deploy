---
layout: default
title: Home
nav_order: 1
description: "system-deploy is a systemd inspired server deployment and configuration tool with no external dependencies"
permalink: /
---

# Simple server deployment and configuration
{: .fs-9 }

*system-deploy* is a systemd inspired server configuration tool with focus on simplicity and ease of use.
{: .fs-6 .fw-300 }

---

<div align="center" markdown="1">

[Get started now](#getting-started){: .btn .btn-primary .fs-5 .mb-4 .mb-md-0 .mr-2 } [View it on GitHub](https://github.com/ppacher/system-deploy){: .btn .fs-5 .mb-4 .mb-md-0 }


![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ppacher/system-deploy?style=flat-square)
![Stable: Beta](https://img.shields.io/badge/Stable-BETA-yellowgreen?style=flat-square)
![GitHub](https://img.shields.io/github/license/ppacher/system-deploy?style=flat-square)
![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/ppacher/system-deploy?label=Release&style=flat-square)

</div>

---

## Overview

`system-deploy` started as my personal server configuration and management tool. It's heavily inspired by systemd's unit files and is focused on
simplicity, ease of use. Thanks to the <img src="https://golang.org/lib/godoc/images/go-logo-blue.svg" alt="Go" width="30"> programming
language, the released binary is statically linked and does not require any external dependencies like heavy execution environments or
preinstalled packages.

#### Features

* systemd unit file syntax
* static binary without external dependencies
* detect changes and trigger actions
* multi-platform supported
* extensible and hackable

## Getting started

This section will give you a quick overview on how to work with system-deploy so you can jumpstart your server configuration.
For more information on system-deploy, examples and more documentation please checkout the sections on the left.

### Installation

To install *system-deploy*, just download the latest release binary for your platform from the release page and copy it to
a location inside your `$PATH`. 

[Release Page](https://github.com/ppacher/system-deploy/releases){: .btn .btn-green .fs-5 .mb-4 .mb-md-0 .mr-2}
{: align="center"}

For example, to install `system-deploy` to `/usr/local/bin`, use the following:

```bash
VERSION=0.1.0
DOWNLOAD_LINK=https://github.com/ppacher/system-deploy/releases/download/v${VERSION}/system-deploy_${VERSION}_linux_i386.tar.gz

# Download and unpack the archive
curl -o system-deploy.tar.gz ${DOWNLOAD_LINK}
tar xfz ./system-deploy.tar.gz

# Add execution permissions and copy to /usr/local/bin
chmod a+x ./system-deploy
sudo mv ./system-deploy /usr/local/bin 
```

Finally, try running system-deploy:

```bash
system-deploy help
```

### Create your first task

Now, let's create your first system-deploy task by creating the directory `~/.deploy/` and a task file:

~/.deploy/10-install-tools.task
{: .code-example .fw-300 }

```ini
[Task]
Description=Install various tools

[InstallPackages]
# AptPkgs will only be installed if `apt` is
# available (i.e. on Debia, Ubuntu and friends)
AptPkgs=vim git golang zsh restic ssmtp htop acl
AptPkgs=tree unattended-upgrades

# PacmanPkgs will onyl be installed if `pacman`
# is available (i.e. on Arch Linux, Manjaro, ...)
PacmanPkgs=vim go git acl zsh restic htop tree msmtp
```

Finally, execute all tasks in `~/.deploy`:

```bash
$ system-deploy ~/.deploy
INFO[0004] 10-install-tools.task         : updated 
```

Congrats, you just created your first deployment task! Time to read on:

[📖 Documentation](docs/concepts/){: .btn .btn-outline .text-green-200 .fs-5 .mb-4 .mb-md-0 .mr-2 }
[🧰 Actions](docs/actions/){: .btn .btn-outline .text-green-200 .fs-5 .mb-4 .mb-md-0 .mr-2 }
[📓 Examples](docs/examples/){: .btn .btn-outline .text-green-200 .fs-5 .mb-4 .mb-md-0 }
{: align="center"}

## About the project

*system-deploy* is &copy; 2020 by [Patrick Pacher](https://github.com/ppacher) and [Contributors](https://github.com/ppacher/system-deploy/graphs/contributors)

## License

*system-deploy* itself is licensed under a BSD 3-clause license. See [LICENSE](https://github.com/ppacher/system-deploy/tree/master/LICENSE) for more information.

Note that the binaries distributed via the [release page](https://github.com/ppacher/system-deploy/releases) are licensed under GPL-3 because the `EditFile` action is compiled against a GPL-3 licensed library.

#### Thank you to the contributors of system-deploy!

<ul class="list-style-none">
{% for contributor in site.github.contributors %}
  <li class="d-inline-block mr-1">
     <a href="{{ contributor.html_url }}"><img src="{{ contributor.avatar_url }}" width="32" height="32" alt="{{ contributor.login }}"/></a>
  </li>
{% endfor %}
</ul>