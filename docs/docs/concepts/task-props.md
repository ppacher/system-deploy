---
layout: default
parent: Documentation
title: Task Properties
nav_order: 3
---
# Task

The `[Task]` section must be available in each deploy unit file and contains
metadata for the task like a human-readable description or the tasks state
(Disabled or Masked). Users may also defined condition and assertions in the
tasks meta section that may disable or fail the task based on environmental
conditions. All properties starting with `Condition` will disable the task if
not matched, all properties starting with `Assert` will cause *system-deploy* to
throw an error and exit.

## Options

   **Description**= (string)  
      Defines a human readable description of the task's purpose

   **StartMasked**= (bool)  
      Set to true if the ask should be masked from execution (Default: "no")

   **Disabled**= (bool)  
      Set to true if the task should be disabled. A disabled task cannot be
      executed in any way (Default: "no")

   **Environment**= ([]string)  
      Configure one or more environment files that are loaded into the task and
      may be used during substitution. Environment files are loaded in the order
      they are specified and later ones overwrite already existing values.

   **ConditionOperatingSystem**= ([]string)  
   **AssertOperatingSystem**=  
      Match against the operating system. All values from GOOS are supported.

   **ConditionArchitecture**= ([]string)  
   **AssertArchitecture**=  
      Match against the architecture system-deploy was compiled for.

   **ConditionPackageManager**= ([]string)  
   **AssertPackageManager**=  
      Match against the installed package-managers.

   **ConditionFileExists**= ([]string)  
   **AssertFileExists**=  
      Test against the existence of a file.

   **ConditionDirectoryExists**= ([]string)  
   **AssertDirectoryExists**=  
      Test against the existence of a directory.


## Contact

**  
https://ppacher.github.io/system-deploy  
