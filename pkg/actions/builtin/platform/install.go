package platform

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/unit"
)

func init() {
	actions.MustRegister(actions.Plugin{
		Name:        "InstallPackages",
		Author:      "Patrick Pacher <patrick.pacher@gmail.com>",
		Website:     "https://github.com/ppacher/system-deploy",
		Description: "Install software packages using various package managers. For more control on the installation behavior use the Exec section instead.",
		Options: []deploy.OptionSpec{
			{
				Name:        "AptPkgs",
				Description: "Packages to install if APT is available",
				Type:        deploy.StringSliceType,
			},
			{
				Name:        "PacmanPkgs",
				Description: "Packages to install if Pacman is available",
				Type:        deploy.StringSliceType,
			},
			// TODO(ppacher): add support for DNF
			// TODO(ppacher): add support for snap
			// TODO(ppacher): add support for arch-linux AUR (maybe using yay?)
			/*
				{
					Name:        "DnfPkgs",
					Description: "Packages to install if DNF is available",
					Type:        deploy.StringSliceType,
				},
				{
					Name:        "SnapPkgs",
					Description: "Packages to install if Snap is available",
					Type:        deploy.StringSliceType,
				},
			*/
		},
		Setup: setupInstallAction,
	})
}

func setupInstallAction(task deploy.Task, sec unit.Section) (actions.Action, error) {
	aptPkgs := getPackages("AptPkgs", sec)
	pacmanPkgs := getPackages("PacmanPkgs", sec)
	dnfPkgs := getPackages("DnfPkgs", sec)
	snapPkgs := getPackages("SnapPkgs", sec)

	if len(aptPkgs) == 0 && len(pacmanPkgs) == 0 && len(dnfPkgs) == 0 && len(snapPkgs) == 0 {
		return nil, fmt.Errorf("no packages to install")
	}

	return &installAction{
		aptPkgs:    aptPkgs,
		pacmanPkgs: pacmanPkgs,
		dnfPkgs:    dnfPkgs,
		snapPkgs:   snapPkgs,
	}, nil
}

func getPackages(configKey string, sec unit.Section) []string {
	var pkgs []string
	pkgOpts := sec.GetStringSlice(configKey)

	for _, p := range pkgOpts {
		pkgs = append(pkgs, strings.Fields(p)...)
	}

	return pkgs
}

type installAction struct {
	actions.Base

	aptPkgs    []string
	pacmanPkgs []string
	dnfPkgs    []string
	snapPkgs   []string
}

func (ia *installAction) Name() string {
	return "Installing packages"
}

func (ia *installAction) Prepare(graph actions.ExecGraph) error {
	return nil
}

func (ia *installAction) Run(ctx context.Context) (bool, error) {
	managers := getPackageManagers()

	for _, m := range managers {
		var err error

		switch m {
		case Pacman:
			if len(ia.pacmanPkgs) == 0 {
				continue
			}
			err = installPacman(ctx, ia.pacmanPkgs...)

		case APT:
			if len(ia.aptPkgs) == 0 {
				continue
			}
			err = installApt(ctx, ia.aptPkgs...)

		default:
			continue
		}

		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func installPacman(ctx context.Context, pkgs ...string) error {
	args := []string{
		"-S",
		"--noconfirm",
	}

	args = append(args, pkgs...)
	cmd := exec.CommandContext(ctx, "pacman", args...)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install packages: %w\n%s", err, string(output))
	}

	return nil
}

func installApt(ctx context.Context, pkg ...string) error {
	args := []string{
		"install",
		"-y",
	}
	args = append(args, pkg...)

	cmd := exec.CommandContext(ctx, "apt", args...)
	cmd.Env = os.Environ()

	cmd.Env = append(cmd.Env, "DEBCONF_FRONTEND='noninteractive'")

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install packages: %w\n%s", err, string(output))
	}

	return nil
}
