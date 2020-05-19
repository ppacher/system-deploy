## Changelog

7ae5d8e Add .goreleaser.yml
7ac31fe Add DisableTask to TaskManager. Fixes #1
3a849fd Add EditFile support
885658a Add deploy task decoding tests
2b591fd Add documentation website
87f0c56 Add file utilitity
00fcb07 Add first concept documentation parts
6aa47c9 Add methods to load drop-in files and apply them to tasks
25ec4dd Add option validation
29f9665 Allow drop-ins to specify task options as well
88d9da3 Change docs URl to github pages
eb9fd1a Export ConvertBool in pkg/unit
7528a6d Extend cli to allow execution of actions directly
e948325 Fix example extension
4875204 Implement search for drop-in files
b387e41 Imported code
79ee7f5 Initial commit
0df1761 Merge branch 'feature/dropins'
e4d086e Refactor AddTask. Add support Disabled= option to Task section
921a67b Rename Run to Execute, make Prepare optional and add change detection to installpackages
8dd1fef Rename cmd package
ff3a7e4 Rename deploy/config.go to deploy/task.go
5dc3146 Update README
4227c0c Update README.md
98626db Update concept section in README
0898fe3 Update documentation
2f82281 Update goreleaser.yml
37a9fba Use strings.SplitAfter instead of Split
3b91f84 Use switch instead of if in platform actions
0545c77 Validate section before creating action
088b36a chore: go mod tidy
f5c43bd cli: Log error when no tasks are found
90733f7 cli: load drop-in files
