package systemd

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/unit"
)

func init() {
	actions.MustRegister(actions.Plugin{
		Name:        "Systemd",
		Description: "Install and manage systemd unit files.",
		Author:      "Patrick Pacher <patrick.pacher@gmail.com>",
		Website:     "https://github.com/ppacher/system-deploy",
		Setup:       setupAction,
		Options: []deploy.OptionSpec{
			{
				Name:        "Install",
				Description: "Path to a systemd unit file to install.\nMultiple files can be split using a space character. May be specified multiple times.",
				Type:        deploy.StringSliceType,
			},
			{
				Name:        "AutoEnable",
				Description: "Whether or not to automatically enable all installed units.",
				Type:        deploy.BoolType,
				Default:     "no",
			},
			{
				Name:        "EnableNow",
				Description: "If AutoEnable is true, or Enable option is set, EnableNow controls if those units should be started immediately.",
				Type:        deploy.BoolType,
				Default:     "no",
			},
			{
				Name:        "Enable",
				Description: "A list of systemd units to enable",
				Type:        deploy.StringSliceType,
			},
			{
				Name:        "InstallDirectory",
				Description: "Path to the systemd unit directoy used to install units.",
				Type:        deploy.StringType,
				Default:     "/etc/systemd/system",
			},
		},
	})
}

func setupAction(task deploy.Task, sec unit.Section) (actions.Action, error) {
	installUnits := sec.GetStringSlice("Install")

	for idx := range installUnits {
		if !filepath.IsAbs(installUnits[idx]) {
			installUnits[idx] = filepath.Clean(filepath.Join(task.Directory, installUnits[idx]))
		}
	}

	autoEnable, err := sec.GetBool("AutoEnable")
	if err != nil && !unit.IsNotSet(err) {
		return nil, err
	}

	enableNow, err := sec.GetBool("EnableNow")
	if err != nil && !unit.IsNotSet(err) {
		return nil, err
	}

	enableUnits := sec.GetStringSlice("Enable")

	installDirectory, err := sec.GetString("InstallDirectory")
	if err != nil {
		if !unit.IsNotSet(err) {
			return nil, err
		}

		installDirectory = "/etc/systemd/system"
	}

	a := &systemdAction{
		installDirectory:    installDirectory,
		unitsToEnable:       enableUnits,
		enableNow:           enableNow,
		autoEnableInstalled: autoEnable,
		unitsToInstall:      installUnits,
	}
	return a, nil
}

type systemdAction struct {
	actions.Base

	unitsToInstall      []string
	unitsToEnable       []string
	autoEnableInstalled bool
	enableNow           bool
	installDirectory    string

	cli *systemctl
}

func (*systemdAction) Name() string { return "Systemd" }

func (a *systemdAction) Prepare(graph actions.ExecGraph) error {
	cli, err := newClient(a.installDirectory)
	if err != nil {
		return err
	}

	a.cli = cli
	return nil
}

func (a *systemdAction) Run(ctx context.Context) (bool, error) {
	var changed bool

	if len(a.unitsToInstall) > 0 {
		installed, err := a.cli.install(a.unitsToInstall...)
		if err != nil {
			return false, fmt.Errorf("failed to install units: %w", err)
		}

		if len(installed) > 0 {
			changed = true
			if err := a.cli.reloadDaemon(); err != nil {
				return false, fmt.Errorf("failed to reload systemd: %w", err)
			}
		}

		if a.autoEnableInstalled {
			enabled, err := a.cli.enable(a.enableNow, a.unitsToInstall...)
			if err != nil {
				return false, fmt.Errorf("failed to enable installed units: %w", err)
			}

			if len(enabled) > 0 {
				changed = true
			}
		}
	}

	if len(a.unitsToEnable) > 0 {
		enabled, err := a.cli.enable(a.enableNow, a.unitsToEnable...)
		if err != nil {
			return false, fmt.Errorf("failed to enable units: %w", err)
		}
		if len(enabled) > 0 {
			changed = true
		}
	}

	return changed, nil
}
