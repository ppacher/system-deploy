package editfile

import (
	"context"
	"os"
	"strings"

	"github.com/ppacher/system-deploy/pkg/actions"
	"github.com/ppacher/system-deploy/pkg/change"
	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/unit"
	"github.com/ppacher/system-deploy/pkg/utils"
	"github.com/rwtodd/Go.Sed/sed"
)

func init() {
	actions.MustRegister(actions.Plugin{
		Name:        "EditFile",
		Author:      "Patrick Pacher <patrick.pacher@gmail.com>",
		Website:     "https://github.com/ppacher/system-deploy",
		Description: "Manipulate existing files using SED like syntax",
		Example:     example,
		Options: []deploy.OptionSpec{
			{
				Name:        "Sed",
				Type:        deploy.StringSliceType,
				Description: "Apply an SED instruction on the target file. May be specified multiple times. Refer to https://github.com/rwtodd/Go.Sed for more information about the regexp syntax.",
			},
			{
				Name:        "File",
				Type:        deploy.StringType,
				Description: "Path to the file to modify",
				Required:    true,
			},
			{
				Name:        "IgnoreMissing",
				Type:        deploy.BoolType,
				Description: "Check if the file exists and if not, don't do anything.",
			},
		},
		Setup: setup,
	})
}

type editAction struct {
	actions.Base

	source     string
	ignore     bool
	skip       bool
	engine     *sed.Engine
	hashBefore string
	mode       os.FileMode
}

func setup(task deploy.Task, section unit.Section) (actions.Action, error) {
	ignore := section.Options.GetBoolDefault("IgnoreMissing", false)
	seds := section.Options.GetStringSlice("Sed")
	source, err := section.Options.GetString("File")
	if err != nil {
		return nil, err
	}

	actions := strings.NewReader(strings.Join(seds, " "))
	engine, err := sed.New(actions)
	if err != nil {
		return nil, err
	}

	return &editAction{
		source: source,
		ignore: ignore,
		engine: engine,
	}, nil
}

func (action *editAction) Name() string {
	return "EditFile"
}

func (action *editAction) Prepare(graph actions.ExecGraph) error {
	// get the file hase before we try to update it
	// so we can detect any changes we did.
	checksum, err := change.FileChecksum(action.source)
	if err != nil {
		if !(os.IsNotExist(err) && action.ignore) {
			return err
		} else {
			action.skip = true
		}
	}
	action.hashBefore = checksum

	if !action.skip {
		action.mode, err = utils.FileMode(action.source)
		if err != nil {
			return err
		}
	}

	return nil
}

func (action *editAction) Execute(_ context.Context) (bool, error) {
	// return now if the source file does not exist and
	// IgnoreMissing= was set
	if action.skip {
		return false, nil
	}

	file, err := os.Open(action.source)
	if err != nil {
		return false, err
	}
	defer file.Close()

	pipe := action.engine.Wrap(file)
	if err := utils.CreateAtomic(action.source, action.mode, pipe); err != nil {
		return false, err
	}

	checksum, err := change.FileChecksum(action.source)
	if err != nil {
		return false, err
	}

	return checksum != action.hashBefore, nil
}

const example = `[Task]
Description= Permit root login via SSH

[EditFile]
File=/etc/ssh/sshd_config
Sed=s/#PermitRootLogin[ ]+no/PermitRootLogin yes/g
`
