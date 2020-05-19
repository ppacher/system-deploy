---
layout: default
parent: Actions
title: Copy
nav_order: 1
---
# Copy
Copy files and folder and script updates

## Change Detection

The `Copy` action uses a Murmur3 hash to check whether or not a destination file
needs to be updated. In any case, `Copy` ensures the destination files mode bit
either match the value of FileMode= or the mode bits of the source file. See
FileMode= for more information.

## Bugs

Note that FileMode does not work when copying a directory recursively. Also,
directories are always copied without checking if an update is required.This
will be fixed in a later release.

## Options

   **Source** (string)  
      The source file to copy to Destination. (required)

   **Destination** (string)  
      The destination path where Source should be copied to. (required)

   **CreateDirectories** (bool)  
      If set to true, missing directories in Destination will be created.
      (Default: "no")

   **FileMode** (int)  
      The mode bits (before umask) to use for the destination file. If unset the
      source filesmode bits will be used. The destination files mode will be
      changed to match FileMode= even if the content is already correct.Note
      that Mode is ingnored when copying directories.

   **DirectoryMode** (int)  
      When creating Destination path (CreateDirectories=yes) the mode bits
      (before umask) for that directories. (Default: "0755")


## Example

```ini
[Task]
Description= Copy file foo to /server/custom/bin

[Copy]
Source=./assets/foo
Destination=/server/custom/bin
CreateDirectories=yes
FileMode=0600
```

## Contact

*Patrick Pacher <patrick.pacher@gmail.com>*  
https://github.com/ppacher/system-deploy  
