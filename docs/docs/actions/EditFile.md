---
layout: default
parent: Actions
title: EditFile
nav_order: 1
---
# EditFile
Manipulate existing files using SED like syntax

## Options

   **Sed** ([]string)  
      Apply an SED instruction on the target file. May be specified multiple
      times. Refer to https://github.com/rwtodd/Go.Sed for more information
      about the regexp syntax.

   **File** (string)  
      Path to the file to modify (required)

   **IgnoreMissing** (bool)  
      Check if the file exists and if not, don't do anything.


## Example

```ini
[Task]
Description= Permit root login via SSH

[EditFile]
File=/etc/ssh/sshd_config
Sed=s/#PermitRootLogin[ ]+no/PermitRootLogin yes/g

```

## Contact

*Patrick Pacher <patrick.pacher@gmail.com>*  
https://github.com/ppacher/system-deploy  
