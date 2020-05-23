---
layout: default
parent: Actions
title: Platform
nav_order: 1
---
# Platform

Run deploy tasks only on certain platforms.

## Options

   **OperatingSystem** (string)  
      Match on the operating system. Supported values are 'darwin', 'linux',
      'bsd', 'windows'

   **Distribution** (string)  
      Match on the distribution string. See lsb_release -a

   **PackageManager** (string)  
      Match on the package manager. Detected package managers include `apt`,
      `snap`, `pacman`, `dnf` and `brew`


## Contact

*Patrick Pacher <patrick.pacher@gmail.com>*  
https://github.com/ppacher/system-deploy  
