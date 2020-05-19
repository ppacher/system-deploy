---
layout: default
parent: Documentation
title: Drop-in Files
nav_order: 2
---


## Drop-in Files

Manipulate, extend and customize the behavior of pre-defined tasks.
{: .fs-5 .fw-300 }

Another concept and feature borrowed from systemd are drop-in files. Drop-in files allow
customization of `.task` units by adding, overwriting or removing propterties from defined
task sections.

For example, imagine the following `.task` file:

~/.deploy/10-install-crontab.task
{: .code-example .fw-300}

```ini
[Task]
Description=Install crontab file

[Copy]
Source=./assets/crontab
Destination=/etc/cron.daily
FileMode=0755

[OnChange]
Run=systemctl restart cronie
```

If that task would be executed on a system that uses `cron` instead of `cronie`, the `Run=`
command in the *OnChange* section would fail. To fix that, without touching the original
`.task` file we can create a drop-in in one of *system-deploy*'s search paths (see [Search Path](#search-path)):

~/.deploy/conf/10-install-crontab.task.d/10-fix-on-change.conf
{: .code-example .fw-300 }

```ini
[OnChange]
# Run is a []string type so we first need to clear out the
# previous value by setting in empty.
Run=
# Now add a new command
Run=systemctl restart cron
```

Once *system-deploy* loads the `10-install-crontab.task` file, it will automatically start
searching for drop-in files and will likely (depending on the [Search Path](#search-path))
find `10-install-crontab.task.d/10-fix-on-change.conf` and apply any changes to the task
file.

As another use-case, imagine that the above task is not going to make any sense at all
because the target system does not use or need the cron service. Since there's no point
in trying to install a crontab -- and failing because `/etc/cron.daily` does not exist --
let's disable the task completely.

~/.deploy/conf/10-install-crontab.task.d/10-disable.conf
{: .code-example .fw-300 }

```ini
[Task]
Disabled=yes
```

### Modifying Properties

As seen in the above example, properties can be altered by defining them in the appropriate
section in the drop-in file. For single-value properties, specifying the new value will
simply overwrite the existing one.

Like in the above example, `Disabled=yes` in the drop-in file simply overwrites any `Disabled`
property value the task might have. When working with list types (like `[]string`, `[]int` or
`[]float`) property values are always add the already existing set of values. To remove a single
value from a list property, you need to clear the whole list first by assigning an empty value
and the adding all other values back. The next example shows how to replace the value `nano`
with `vim` in a list of packages to install:

~/.deploy/80-install-tools.task
{: .code-example .fw-300}

```ini
[Task]
Description=Install various tools

[InstallPackages]
AptPkgs=htop nano tree acl
# The list can also be split accross multiple property
# assignments
AptPkgs=clamav ranger
```

Now, replace `nano` with `vim` by clearing and re-adding all other packages:

~/.deploy/conf/80-install-tools.task.d/replace-nano-with-vim.conf
{: .code-example .fw-300 }

```ini
[InstallPackages]
# Clear out _all_ values
AptPkgs=
# Re-add them with nano replaced by vim
AptPkgs=htop tree acl vim
AptPkgs=clamav ranger
```

Currently drop-in files require that modified *sections of a task file are unique*. That
means that it's not possible to apply a drop-in file to a section that occures more
than once because *system-deploy* cannot yet figure out which section should be updated.
{: .code-example .fs-3 .fw-300 .text-red-300 }

### Search Path

Like systemd, *system-deploy* searches various folders and locations for drop-in files.
It basically re-implements the whole search path that is supported by systemd releases.
Note that drop-in files must always have the `.conf` extension and must reside in a
directory with a `.d` suffix. If the unit name contains dashes the sub-strings (until
the dashes) are added to the search path as well. For example, when *system-deploy* 
searches, with a search path set to `/usr/lib/system-deploy:/etc/system-deploy`,
for drop-in files of a task called `foo-bar-baz.task`, the following paths
are checked in the below order (with later directories having priority):

1. `/usr/lib/system-deploy/task.d/`
2. `/usr/lib/system-deploy/foo-.task.d/`
3. `/usr/lib/system-deploy/foo-bar-.task.d/`
4. `/usr/lib/system-deploy/foo-bar-baz.task.d/`
5. `/etc/system-deploy/task.d/`
6. `/etc/system-deploy/foo-.task.d/`
7. `/etc/system-deploy/foo-bar-.task.d/`
8. `/etc/system-deploy/foo-bar-baz.task.d/`

From the above list, the file `/usr/lib/system-deploy/foo-bar-baz.task.d/10-overwrite.conf` would be overwritten
by the file `/etc/system-deploy/foo-.task.d/10-overwrite.conf`. Files with the same name are
completely replaced by the higher-priority version and are *not* merged.

By default, the search path of *system-deploy* is set to `./.config:/etc/system-deploy` so the `.config`
folder of the current working directory is checked as well. To specify you own search path (and overwriting
the default), use `--path`/`-p` when calling *system-deploy*.

```bash
system-deploy --path /server/cluster/overwrites  \
              --path /server/node/overwrites     \
              /server/deploy-tasks
```

---

[Next: Execution Graph](/docs/concepts/30-execution-graph){: .btn .btn-outline }
{: .fs-3 align="right"}
