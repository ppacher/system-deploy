package condition

import (
	"os"
	"os/user"
	"runtime"
	"strings"
)

// BuiltinConditions is a slice of all built-in conditions.
var BuiltinConditions = []Condition{
	{
		Name:        "OperatingSystem",
		Description: "Match against the operating system. All values from GOOS are supported.",
		check: func(value string) (bool, error) {
			return strings.EqualFold(runtime.GOOS, value), nil
		},
	},
	{
		Name:        "Architecture",
		Description: "Match against the architecture system-deploy was compiled for.",
		check: func(value string) (bool, error) {
			return strings.EqualFold(runtime.GOARCH, value), nil
		},
	},
	{
		Name:        "PackageManager",
		Description: "Match against the installed package-managers.",
		check: func(value string) (bool, error) {
			return HasPackageManager(value), nil
		},
	},
	{
		Name:        "FileExists",
		Description: "Test against the existence of a file.",
		check: func(path string) (bool, error) {
			stat, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					return false, nil
				}
				return false, err
			}

			return !stat.IsDir(), nil
		},
	},
	{
		Name:        "DirectoryExists",
		Description: "Test against the existence of a directory.",
		check: func(path string) (bool, error) {
			stat, err := os.Stat(path)
			if err != nil {
				if os.IsNotExist(err) {
					return false, nil
				}
				return false, err
			}

			return stat.IsDir(), nil
		},
	},
	{
		Name:        "UserExists",
		Description: "Test against the existence of a user or userid",
		check: func(value string) (bool, error) {
			_, err := user.Lookup(value)
			if err == nil {
				return true, nil
			}

			_, err = user.LookupId(value)
			if err == nil {
				return true, nil
			}
			return false, nil
		},
	},
	{
		Name:        "GroupExists",
		Description: "Test against the existence of a group or groupid",
		check: func(value string) (bool, error) {
			_, err := user.LookupGroup(value)
			if err == nil {
				return true, nil
			}

			_, err = user.LookupGroupId(value)
			if err == nil {
				return true, nil
			}
			return false, nil
		},
	},
}
