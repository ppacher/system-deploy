package utils

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/flynn/go-shlex"
)

type ExecOptions struct {
	Attrs      *syscall.SysProcAttr
	PipeOutput bool
	PipeInput  bool
	Env        map[string]string
}

// ExecCommand executes cmd and returns any error encountered.
func ExecCommand(ctx context.Context, workDir string, cmd string, opts *ExecOptions) error {
	parts, err := shlex.Split(cmd)
	if err != nil {
		return err
	}

	if len(parts) < 1 {
		return fmt.Errorf("invalid command")
	}

	c := exec.CommandContext(ctx, parts[0], parts[1:]...)

	if workDir != "" {
		c.Dir = workDir
	}

	var output bytes.Buffer

	c.Stderr = &output
	c.Stdout = &output
	c.Env = os.Environ()

	if opts != nil {
		c.SysProcAttr = opts.Attrs

		if opts.PipeOutput {
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
		}

		if opts.PipeInput {
			c.Stdin = os.Stdin
		}

		if opts.Env != nil {
			for k, v := range opts.Env {
				c.Env = append(c.Env, fmt.Sprintf("%s=%q", k, v))
			}
		}
	}

	if err := c.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			return fmt.Errorf("command returned %d\n%s", ws.ExitStatus(), output.String())
		}
		return fmt.Errorf("command failed: %w", err)
	}

	return nil
}
