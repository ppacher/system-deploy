package condition

import (
	"os/exec"
	"runtime"
)

// Known package manager binaries.
const (
	APT    = "apt"
	Snap   = "snap"
	Pacman = "pacman"
	Dnf    = "dnf"
	Brew   = "brew"
)

// HasPackageManager returns true if the package-manager
// pm is installed and reachable via $PATH.
func HasPackageManager(pm string) bool {
	for _, p := range getPackageManagers() {
		if p == pm {
			return true
		}
	}

	return false
}

// getPackageManagers searches for available package managers and
// returns a slice of package managers found.
func getPackageManagers() []string {
	switch runtime.GOOS {
	case "windows":
		return nil // TODO(ppacher): we could check for nuget

	case "darwin":
		// TODO(ppacher): are there more package managers for darwin?
		if _, err := exec.LookPath(Brew); err == nil {
			return []string{Brew}
		}
		return nil

	case "linux":
		pkg := []string{}

		for _, m := range []string{APT, Pacman, Dnf, Snap} {
			if _, err := exec.LookPath(m); err == nil {
				pkg = append(pkg, m)
			}
		}

		return pkg

	default:
		return nil
	}
}
