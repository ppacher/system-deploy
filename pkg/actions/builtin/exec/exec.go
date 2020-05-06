package exec

import (
	"context"
	"fmt"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/unit"
	"github.com/ppacher/system-deploy/pkg/utils"
)

func init() {
	actions.MustRegister(actions.Plugin{
		Name:        "Exec",
		Description: "Execute one or more commands",
		Setup:       setupAction,
		Author:      "Patrick Pacher <patrick.pacher@gmail.com>",
		Website:     "https://github.com/ppacher/system-deploy",
		Options: []deploy.OptionSpec{
			{
				Name:        "Command",
				Type:        deploy.StringType,
				Description: "The command to execute.",
				Required:    true,
			},
			{
				Name:        "WorkingDirectory",
				Type:        deploy.StringType,
				Description: "The working directory for the command",
			},
			{
				Name:        "Chroot",
				Type:        deploy.StringType,
				Description: "Chroot for the command",
			},
			{
				Name:        "User",
				Type:        deploy.StringType,
				Description: "Execute the command as User (either name or ID)",
			},
			{
				Name:        "Group",
				Type:        deploy.StringType,
				Description: "Execute the command under Group (either name or ID)",
			},
			{
				Name:        "DisplayOutput",
				Type:        deploy.BoolType,
				Description: "Display command output",
				Default:     "no",
			},
			{
				Name:        "ForwardStdin",
				Type:        deploy.BoolType,
				Description: "Forward current stdin to the command",
				Default:     "no",
			},
			{
				Name:        "Environment",
				Type:        deploy.StringSliceType,
				Description: "Add environment variables for the command. The value should follow the format KEY=VALUE",
			},
		},
	})
}

func setupAction(task deploy.Task, sec unit.Section) (actions.Action, error) {
	cmd, err := sec.GetString("Command")
	if err != nil {
		return nil, err
	}

	workDir, err := sec.GetString("WorkingDirectory")
	if err != nil {
		if !unit.IsNotSet(err) {
			return nil, err
		}

		workDir = task.Directory
	}

	chroot, err := sec.GetString("Chroot")
	if err != nil && !unit.IsNotSet(err) {
		return nil, err
	}

	uid := -1
	gid := -1

	userName, err := sec.GetString("User")
	if err != nil && !unit.IsNotSet(err) {
		return nil, err
	} else if err == nil {
		u, err := user.Lookup(userName)
		if err != nil {
			u, err = user.LookupId(userName)
		}

		if err != nil {
			return nil, fmt.Errorf("user %q does not exist", userName)
		}

		uid64, err := strconv.ParseInt(u.Uid, 0, 32)
		if err != nil {
			return nil, err
		}
		gid64, err := strconv.ParseInt(u.Gid, 0, 32)
		if err != nil {
			return nil, err
		}

		uid = int(uid64)
		gid = int(gid64)
	}

	groupName, err := sec.GetString("Group")
	if err != nil && !unit.IsNotSet(err) {
		return nil, err
	} else if err == nil {
		grp, err := user.LookupGroup(groupName)
		if err != nil {
			grp, err = user.LookupGroupId(groupName)
		}

		if err != nil {
			return nil, fmt.Errorf("group %q does not exist", groupName)
		}

		gid64, err := strconv.ParseInt(grp.Gid, 0, 32)
		if err != nil {
			return nil, err
		}

		gid = int(gid64)

		if uid == -1 {
			cur, err := user.Current()
			if err != nil {
				return nil, err
			}

			uid64, err := strconv.ParseInt(cur.Uid, 0, 32)
			if err != nil {
				return nil, err
			}
			uid = int(uid64)
		}
	}

	pipeOut, err := sec.GetBool("DisplayOutput")
	if err != nil && !unit.IsNotSet(err) {
		return nil, fmt.Errorf("invalid setting for option 'DisplayOutput': %w", err)
	}

	pipeIn, err := sec.GetBool("ForwardStdin")
	if err != nil && !unit.IsNotSet(err) {
		return nil, fmt.Errorf("invalid setting for option 'ForwardStdin': %w", err)
	}

	var environ map[string]string
	envList := sec.GetStringSlice("Environment")
	if len(envList) > 0 {
		environ = make(map[string]string)

		for _, val := range envList {
			parts := strings.Split(val, "=")
			if len(parts) < 2 {
				return nil, fmt.Errorf("invalid value for option 'Environment'")
			}

			key := parts[0]
			value := strings.Join(parts[1:], "=")

			environ[key] = value
		}
	}

	a := &action{
		taskDir: workDir,
		chroot:  chroot,
		cmd:     cmd,
		user:    uid,
		group:   gid,
		pipeIn:  pipeIn,
		pipeOut: pipeOut,
		environ: environ,
	}

	return a, nil
}

type action struct {
	actions.Base

	taskDir string
	chroot  string
	user    int
	group   int
	cmd     string
	environ map[string]string
	pipeOut bool
	pipeIn  bool
}

func (a *action) Name() string {
	return fmt.Sprintf("Running %q", strings.Split(a.cmd, "\n")[0])
}

// Prepare does nothing for exec.
func (a *action) Prepare(_ actions.ExecGraph) error {
	return nil
}

func (a *action) Execute(ctx context.Context) (bool, error) {
	opts := &utils.ExecOptions{
		Attrs:      &syscall.SysProcAttr{},
		PipeInput:  a.pipeIn,
		PipeOutput: a.pipeOut,
		Env:        a.environ,
	}

	hasAttrs := false

	if a.chroot != "" {
		opts.Attrs.Chroot = a.chroot
	}

	if a.user != -1 && a.group != -1 {
		opts.Attrs.Credential = &syscall.Credential{
			Uid:         uint32(a.user),
			Gid:         uint32(a.group),
			NoSetGroups: true,
		}
		hasAttrs = true
	}

	if !hasAttrs {
		opts.Attrs = nil
	}

	if err := utils.ExecCommand(ctx, a.taskDir, a.cmd, opts); err != nil {
		return false, err
	}

	return true, nil
}
