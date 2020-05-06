package builtin

import (
	// Import all built-in actions
	_ "github.com/ppacher/system-deploy/pkg/actions/builtin/copy"
	_ "github.com/ppacher/system-deploy/pkg/actions/builtin/editfile"
	_ "github.com/ppacher/system-deploy/pkg/actions/builtin/exec"
	_ "github.com/ppacher/system-deploy/pkg/actions/builtin/onchange"
	_ "github.com/ppacher/system-deploy/pkg/actions/builtin/platform"
	_ "github.com/ppacher/system-deploy/pkg/actions/builtin/systemd"
)
