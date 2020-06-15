---
layout: default
parent: Actions
title: Exec
nav_order: 1
---
# Exec

Execute one or more commands

## Options

   **Command**= (string)  
      The command to execute. (required)

   **WorkingDirectory**= (string)  
      The working directory for the command

   **Chroot**= (string)  
      Chroot for the command

   **User**= (string)  
      Execute the command as User (either name or ID)

   **Group**= (string)  
      Execute the command under Group (either name or ID)

   **DisplayOutput**= (bool)  
      Display command output (Default: "no")

   **ForwardStdin**= (bool)  
      Forward current stdin to the command (Default: "no")

   **Environment**= ([]string)  
      Add environment variables for the command. The value should follow the
      format KEY=VALUE

   **ChangedOnExit**= (int)  
      If set, the task will be marked as changed/updated if Command= returns
      with the specified exit code.

   **PristineOnExit**= (int)  
      If set, the task will be marked as unchanged/pristine if Command= returns
      with the specified exit code.


## Contact

*Patrick Pacher <patrick.pacher@gmail.com>*  
https://github.com/ppacher/system-deploy  
