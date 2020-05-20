---
layout: default
parent: Documentation
title: Tasks
nav_order: 1
---


## Tasks

Defines the different steps to perform during system deployment and configuration
update.
{: .fs-5 .fw-300 }

### Overview

*system-deploy* is heavily inspired by [systemd](https://systemd.io/) and re-uses the
format and some concepts around systemd's unit files. That is, every action that *system-deploy*
should perform is described in a unit file with a `.task` extension. Those files
follow an [INI style syntax](https://en.wikipedia.org/wiki/INI_file) and are split into
multiple sections that contain key-value pairs in the format `key=value`.  Refer to
[systemd.syntax(7)](https://www.freedesktop.org/software/systemd/man/systemd.syntax.html) for
more information about the file format, escaping options and comment lines.

Before diving into the details of a `.task` unit, here's an example of simple task copying
a file:

10-install-crontab.task
{: .code-example .fw-300}

```ini
[Task]
Description=Install crontab to /etc/cron.daily
Disabled=no

[Copy]
Source=./assets/crontab
Destination=/etc/cron.daily/my-crontab
FileMode=0755
CreateDirectories=no

[OnChange]
Run=systemctl restart cronie
```

### The Task Section

As you can see, the unit file starts with a `[Task]` section that defines metadata for
the `.task` file. In the above example a `Description=` is set an the task is marked as
not-disabled (`Disabled=no` is actually the default value and only specified for
documentation purposes). The task's description is mainly used for logging purposes while
`Disabled=` configures the task state. Refer to [Execution Graph](./30-execution-graph)
for more information on task states and their meaning.

Each `.task` file must contain exactly one `[Task]` section. Although allowed to be empty,
it is recommended to always specify at least a task `Description=`.  See [Task Properties](#task-properties)
below for a list of supported and accepted task metadata properties.

### The Action Sections

Following the `[Task]` section, one ore more *action sections* can be defined. Those sections
actually define what happens when the task is executed by *system-deploy*. In the above example,
a `[Copy]` and a `[OnChange]` action is defined. The `[Copy]` action will ensure the file
`/etc/cron.daily/my-crontab` is kept in sync with `./assets/crontab`. The `[OnChange]` action
instead invokes `systemctl restart cronie` in case one of the action sections reports a change
(see [Change Detection](#change-detection)).

As a result, each time *system-deploy* updates the file the cron daemon is restarted so changes
take effect immediately.

In addition, `.task` files may contain multiple section with the same name. As an example, a task
may contain multiple *[Copy]* sections to ensure multiple files are kept up-to-date.
The order of the action sections is important as *system-deploy* will execute them in 
that order and abort as soon as an error or failure condition is reported. 

In the above example, the [Copy](../actions/Copy) and the [OnChange](../actions/OnChange) have
been used.

[Available Actions](../actions/){: .btn .btn-outline .text-green-200 .fs-5 .mb-4 .mb-md-0}
{: align="center" }

### Task Properties

The `[Task]` section of each task file may contain the following properties:

| Property        | Type          | Default | Description                                        |
|:----------------|:--------------|:--------|:---------------------------------------------------|
| Description     | String        |         | A human readable description of what the task does |
| StartMasked     | Boolean       | no      | Wether or not the task is masked by default        |
| Disabled        | Boolean       | no      | Wether or not the task is disabled                 |

---

[Next: Drop-in Files](/docs/concepts/20-dropins){: .btn .btn-outline }
{: .fs-3 align="right"}
