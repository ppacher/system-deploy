package systemd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ppacher/system-deploy/pkg/change"
	"github.com/ppacher/system-deploy/pkg/utils"
)

// systemctl wraps the systemd systemctl command.
type systemctl struct {
	installDirectory string
}

func newClient(installDirectory string) (*systemctl, error) {
	if f, err := os.Stat(installDirectory); err != nil || !f.IsDir() {
		if err == nil {
			err = fmt.Errorf("not a directory")
		}
		return nil, fmt.Errorf("invalid installation directory %s: %w", installDirectory, err)
	}

	if _, err := exec.LookPath("systemctl"); err != nil {
		return nil, fmt.Errorf("failed to find systemctl binary")
	}

	cli := &systemctl{installDirectory: installDirectory}
	return cli, nil
}

// enable enables all units and returns at the first error
// encountered. If now is true all units will be started
// immediately (systemctl enable --now)
func (cli *systemctl) enable(now bool, units ...string) ([]string, error) {
	enabled := []string{}
	args := []string{"enable"}
	if now {
		args = append(args, "--now")
	}

	for _, unit := range units {
		if err := cli.systemctl("is-enabled", unit); err == nil {
			continue
		}

		if err := cli.systemctl(append(args, unit)...); err != nil {
			return nil, err
		}
		enabled = append(enabled, unit)
	}

	return enabled, nil
}

// install installs units to the installation directory. Only
// files that are either missing or have the wrong content
// are installed.
func (cli *systemctl) install(unitFiles ...string) ([]string, error) {
	var filesInstalled []string
	for _, unit := range unitFiles {
		changed, err := cli.copyUnitFile(unit)
		if err != nil {
			return nil, err
		}

		if changed {
			filesInstalled = append(filesInstalled, unit)
		}
	}
	return filesInstalled, nil
}

// runSystemCtl executes systemctl with args and returns an
// errof if it fails.
func (cli *systemctl) systemctl(args ...string) error {
	cmd := exec.Command("systemctl", args...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl: %w\n%s", err, string(output))
	}

	return nil
}

// reloadDaemon reloads the systemd deamon via systemctl daemon-reload
func (cli *systemctl) reloadDaemon() error {
	return cli.systemctl("daemon-reload")
}

// copyUnitFile copies file to the unit directory pointed to
// by cli.installDirectory.
func (cli *systemctl) copyUnitFile(file string) (bool, error) {
	targetFileName := filepath.Join(cli.installDirectory, filepath.Base(file))

	if update, err := change.FileUpdateNeeded(file, targetFileName); err != nil || !update {
		return update, err
	}

	if err := utils.CopyAtomicKeepMode(file, targetFileName, 0600); err != nil {
		return false, err
	}
	return true, nil
}
