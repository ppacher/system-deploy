package platform

import (
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/unit"
)

// Known package manager binaries.
const (
	APT    = "apt"
	Snap   = "snap"
	Pacman = "pacman"
	Dnf    = "dnf"
	Brew   = "brew"
)

func init() {
	actions.MustRegister(actions.Plugin{
		Name:        "Platform",
		Description: "Run deploy tasks only on certain platforms.",
		Author:      "Patrick Pacher <patrick.pacher@gmail.com>",
		Website:     "https://github.com/ppacher/system-deploy",
		Options: []deploy.OptionSpec{
			{
				Name:        "OperatingSystem",
				Description: "Match on the operating system. Supported values are 'darwin', 'linux', 'bsd', 'windows'",
				Type:        deploy.StringType,
			},
			{
				Name:        "Distribution",
				Description: "Match on the distribution string. See lsb_release -a",
				Type:        deploy.StringType,
			},
			{
				Name:        "PackageManager",
				Description: "Match on the package manager. Detected package managers include `apt`, `snap`, `pacman`, `dnf` and `brew`",
				Type:        deploy.StringType,
			},
		},
		Setup: setupPlatform,
	})
}

func setupPlatform(task deploy.Task, sec unit.Section) (actions.Action, error) {
	matchOS, err := sec.GetString("OperatingSystem")
	if err != nil && !unit.IsNotSet(err) {
		return nil, err
	}

	matchDist, err := sec.GetString("Distribution")
	if err != nil && !unit.IsNotSet(err) {
		return nil, err
	}

	matchPkg, err := sec.GetString("PackageManager")
	if err != nil && !unit.IsNotSet(err) {
		return nil, err
	}

	return &matchPlatformAction{
		task:      task.FileName,
		matchDist: matchDist,
		matchOS:   matchOS,
		matchPkg:  matchPkg,
	}, nil
}

type matchType string

const (
	allow   matchType = "allow"
	deny    matchType = "deny"
	noMatch matchType = "no-match"
)

type matchPlatformAction struct {
	actions.Base

	task      string
	matchOS   string
	matchDist string
	matchPkg  string
}

func (a *matchPlatformAction) Name() string {
	return "Platform"
}

func (a *matchPlatformAction) Prepare(graph actions.ExecGraph) error {
	var verdict = noMatch

	mask := func() error {
		a.Debugf("Disabling task %s due to platform constraints", color.New(color.Bold).Sprint(a.task))
		return graph.MaskTask(a.task)
	}

	if a.matchOS != "" {
		t := match(runtime.GOOS, a.matchOS)

		if t == deny {
			return mask()
		}

		if t == allow {
			verdict = allow
		}
	}

	if a.matchPkg != "" {
		t := matchList(getPackageManagers(), a.matchPkg)
		if t == deny {
			return mask()
		}

		if t == allow {
			verdict = allow
		}
	}

	// by default we deny if not one of the conditions
	// matched.
	if verdict != allow {
		return mask()
	}

	return nil
}

// matchList matches values against a condition and returns
// whether or not the condition is met. If condition defines
// and equality check, the condition is met as soon as the
// matching value is found. If condition defines a non-equality
// check the condition must not match on all values.
func matchList(values []string, condition string) matchType {
	checkType, condition := parseCondition(condition)

	switch checkType {
	case '=':
		for _, v := range values {
			if strings.ToLower(v) == condition {
				return allow
			}
		}
		return noMatch

	case '!':
		for _, v := range values {
			if strings.ToLower(v) == condition {
				return deny
			}
		}
		return noMatch

	default:
		panic("Unexpected switch case for condition type: " + string(checkType))
	}
}

// match returns true if value matches condition. By default, condition is
// an equality check (=). To check for non-equality, prefix the condition value with "!".
func match(value string, condition string) matchType {
	checkType, condition := parseCondition(condition)
	lowerValue := strings.ToLower(value)

	switch checkType {
	case '=':
		if condition == lowerValue {
			return allow
		}
		return noMatch
	case '!':
		if condition == lowerValue {
			return deny
		}
		return noMatch
	default:
		panic("Unexpected switch case for condition type: " + string(checkType))
	}
}

func parseCondition(condition string) (byte, string) {
	var checkType byte
	switch condition[0] {
	case '!', '=':
		checkType = condition[0]
		condition = condition[1:]
	default:
		checkType = '='
	}

	return checkType, strings.ToLower(condition)
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
