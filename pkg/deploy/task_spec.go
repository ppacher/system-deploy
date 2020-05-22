package deploy

import "github.com/ppacher/system-deploy/pkg/unit"

type taskMetaOption struct {
	OptionSpec
	set func(val unit.Options, t *Task) error
}

// TaskOptions returns a slice of OptionSpec that are
// allowed int he task meta section.
func TaskOptions() []OptionSpec {
	specs := make([]OptionSpec, len(taskOptions))

	for idx, spec := range taskOptions {
		specs[idx] = spec.OptionSpec
	}

	return specs
}

// TaskOptions defines all supported options for the task
// meta section.
var taskOptions = []taskMetaOption{
	{
		OptionSpec: OptionSpec{
			Name:        "Description",
			Description: "Defines a human readable description of the task's purpose",
			Type:        StringType,
		},
		set: func(val unit.Options, t *Task) error {
			str, err := val.GetString("Description")
			if err != nil {
				return err
			}

			t.Description = &str
			return nil
		},
	},
	{
		OptionSpec: OptionSpec{
			Name:        "StartMasked",
			Description: "Set to true if the ask should be masked from execution",
			Default:     "no",
			Type:        BoolType,
		},
		set: func(val unit.Options, t *Task) error {
			masked, err := val.GetBool("StartMasked")
			if err != nil {
				return err
			}

			t.StartMasked = &masked
			return nil
		},
	},
	{
		OptionSpec: OptionSpec{
			Name:        "Disabled",
			Description: "Set to true if the task should be disabled. A disabled task cannot be executed in any way",
			Default:     "no",
			Type:        BoolType,
		},
		set: func(val unit.Options, t *Task) error {
			disabled, err := val.GetBool("Disabled")
			if err != nil {
				return err
			}

			t.Disabled = &disabled
			return nil
		},
	},
}
