## [ Systemd ]
Install and manage systemd unit files.

## Options

   **Install** ([]string)  
      Path to a systemd unit file to install.
      Multiple files can be split using a space character. May be specified
      multiple times.

   **AutoEnable** (bool)  
      Whether or not to automatically enable all installed units. (Default:
      "no")

   **EnableNow** (bool)  
      If AutoEnable is true, or Enable option is set, EnableNow controls if
      those units should be started immediately. (Default: "no")

   **Enable** ([]string)  
      A list of systemd units to enable

   **InstallDirectory** (string)  
      Path to the systemd unit directoy used to install units. (Default:
      "/etc/systemd/system")


## Contact

*Patrick Pacher <patrick.pacher@gmail.com>*  
https://github.com/ppacher/system-deploy  
