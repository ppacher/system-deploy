package exec

import (
	"context"
	"fmt"
	"os/user"
	"strconv"
	"strings"
	"syscall"

	"github.com/ppacher/system-conf/conf"
	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/utils"
)

func init() {
	actions.MustRegister(actions.Plugin{
		Name:        "Exec",
		Description: "Execute one or more commands",
		Setup:       setupAction,
		Author:      "Patrick Pacher <patrick.pacher@gmail.com>",
		Website:     "https://github.com/ppacher/system-deploy",
		Options: []conf.OptionSpec{
			{
				Name:        "Command",
				Type:        conf.StringType,
				Description: "The command to execute.",
				Required:    true,
			},
			{
				Name:        "WorkingDirectory",
				Type:        conf.StringType,
				Description: "The working directory for the command",
			},
			{
				Name:        "Chroot",
				Type:        conf.StringType,
				Description: "Chroot for the command",
			},
			{
				Name:        "User",
				Type:        conf.StringType,
				Description: "Execute the command as User (either name or ID)",
			},
			{
				Name:        "Group",
				Type:        conf.StringType,
				Description: "Execute the command under Group (either name or ID)",
			},
			{
				Name:        "DisplayOutput",
				Type:        conf.BoolType,
				Description: "Display command output",
				Default:     "no",
			},
			{
				Name:        "ForwardStdin",
				Type:        conf.BoolType,
				Description: "Forward current stdin to the command",
				Default:     "no",
			},
			{
				Name:        "Environment",
				Type:        conf.StringSliceType,
				Description: "Add environment variables for the command. The value should follow the format KEY=VALUE",
			},
			{
				Name:        "ChangedOnExit",
				Type:        conf.IntType,
				Description: "If set, the task will be marked as changed/updated if Command= returns with the specified exit code.",
			},
			{
				Name:        "PristineOnExit",
				Type:        conf.IntType,
				Description: "If set, the task will be marked as unchanged/pristine if Command= returns with the specified exit code.",
			},
		},
	})
}

func resolveUserGroup(userName, groupName string) (uid uint32, gid uint32, err error) {
	var uidSet bool

	if userName != "" {
		u, err := user.Lookup(userName)
		if err != nil {
			u, err = user.LookupId(userName)
		}

		if err != nil {
			return 0, 0, fmt.Errorf("user %q does not exist", userName)
		}

		uid64, err := strconv.ParseInt(u.Uid, 0, 32)
		if err != nil {
			return 0, 0, err
		}
		gid64, err := strconv.ParseInt(u.Gid, 0, 32)
		if err != nil {
			return 0, 0, err
		}

		uid = uint32(uid64)
		gid = uint32(gid64)
		uidSet = true
	}

	if groupName != "" {
		grp, err := user.LookupGroup(groupName)
		if err != nil {
			grp, err = user.LookupGroupId(groupName)
		}

		if err != nil {
			return 0, 0, fmt.Errorf("group %q does not exist", groupName)
		}

		gid64, err := strconv.ParseInt(grp.Gid, 0, 32)
		if err != nil {
			return 0, 0, err
		}

		gid = uint32(gid64)

		if !uidSet {
			cur, err := user.Current()
			if err != nil {
				return 0, 0, err
			}

			uid64, err := strconv.ParseInt(cur.Uid, 0, 32)
			if err != nil {
				return 0, 0, err
			}
			uid = uint32(uid64)
		}
	}
	return uid, gid, nil
}

func setupAction(task deploy.Task, sec conf.Section) (actions.Action, error) {
	cmd, err := sec.GetString("Command")
	if err != nil {
		return nil, err
	}

	workDir, err := sec.GetString("WorkingDirectory")
	if err != nil {
		if !conf.IsNotSet(err) {
			return nil, err
		}

		workDir = task.Directory
	}

	chroot, err := sec.GetString("Chroot")
	if err != nil && !conf.IsNotSet(err) {
		return nil, err
	}

	userName, err := sec.GetString("User")
	if err != nil && !conf.IsNotSet(err) {
		return nil, err
	}

	groupName, err := sec.GetString("Group")
	if err != nil && !conf.IsNotSet(err) {
		return nil, err
	}

	pipeOut, err := sec.GetBool("DisplayOutput")
	if err != nil && !conf.IsNotSet(err) {
		return nil, fmt.Errorf("invalid setting for option 'DisplayOutput': %w", err)
	}

	pipeIn, err := sec.GetBool("ForwardStdin")
	if err != nil && !conf.IsNotSet(err) {
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

	var exitCode *int64
	ecChanged := false

	exitCodeChanged, err := sec.GetInt("ChangedOnExit")
	if err != nil && err != conf.ErrOptionNotSet {
		return nil, err
	}
	if err != conf.ErrOptionNotSet {
		exitCode = &exitCodeChanged
		ecChanged = true
	}

	exitCodeNotChanged, err := sec.GetInt("PristineOnExit")
	if err != nil && err != conf.ErrOptionNotSet {
		return nil, err
	}

	if err != conf.ErrOptionNotSet {
		if exitCode != nil {
			return nil, fmt.Errorf("cannot use ChangedOnExit and PristineOnExit at the same time")
		}

		exitCode = &exitCodeNotChanged
		ecChanged = false
	}

	a := &action{
		taskDir:         workDir,
		chroot:          chroot,
		cmd:             cmd,
		user:            userName,
		group:           groupName,
		pipeIn:          pipeIn,
		pipeOut:         pipeOut,
		environ:         environ,
		exitCode:        exitCode,
		exitCodeChanged: ecChanged,
	}

	return a, nil
}

type action struct {
	actions.Base

	taskDir         string
	chroot          string
	user            string
	group           string
	cmd             string
	environ         map[string]string
	pipeOut         bool
	pipeIn          bool
	exitCode        *int64
	exitCodeChanged bool
}

func (a *action) Name() string {
	return fmt.Sprintf("Running %q", strings.Split(a.cmd, "\n")[0])
}

// Prepare does nothing for exec.
func (a *action) Prepare(_ actions.ExecGraph) error {
	return nil
}

func (a *action) Execute(ctx context.Context) (bool, error) {
	var exitCode int64
	opts := &utils.ExecOptions{
		Attrs:      &syscall.SysProcAttr{},
		PipeInput:  a.pipeIn,
		PipeOutput: a.pipeOut,
		Env:        a.environ,
		ExitCode:   &exitCode,
	}

	hasAttrs := false

	if a.chroot != "" {
		opts.Attrs.Chroot = a.chroot
	}

	if a.user != "" || a.group != "" {
		uid, gid, err := resolveUserGroup(a.user, a.group)
		if err != nil {
			return false, err
		}

		opts.Attrs.Credential = &syscall.Credential{
			Uid:         uid,
			Gid:         gid,
			NoSetGroups: true,
		}
		hasAttrs = true
	}

	if !hasAttrs {
		opts.Attrs = nil
	}

	hasChanged := func() bool {
		if a.exitCode != nil {
			if *a.exitCode == exitCode {
				return a.exitCodeChanged
			}
			return !a.exitCodeChanged
		}

		return true
	}

	if err := utils.ExecCommand(ctx, a.taskDir, a.cmd, opts); err != nil {
		if _, ok := err.(*utils.ExitCodeError); ok && a.exitCode != nil {
			return hasChanged(), nil
		}

		return hasChanged(), err
	}

	return hasChanged(), nil
}
