package deploy

import "github.com/ppacher/system-deploy/pkg/unit"

type taskMetaOption struct {
	OptionSpec
	set func(val unit.Options, t *Task) error
	get func(t *Task) []string
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
			if val == nil {
				t.Description = ""
				return nil
			}

			var err error
			t.Description, err = val.GetString("Description")
			return err
		},
		get: func(t *Task) []string {
			if t.Description == "" {
				return nil
			}
			return []string{t.Description}
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
			if val == nil {
				t.StartMasked = false
				return nil
			}

			var err error
			t.StartMasked, err = val.GetBool("StartMasked")
			return err
		},
		get: func(t *Task) []string {
			if !t.StartMasked {
				return nil
			}

			return []string{"yes"}
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
			if val == nil {
				t.Disabled = false
				return nil
			}
			var err error
			t.Disabled, err = val.GetBool("Disabled")
			return err
		},
		get: func(t *Task) []string {
			if !t.Disabled {
				return nil
			}

			return []string{"yes"}
		},
	},
	{
		OptionSpec: OptionSpec{
			Name: "Environment",
			Description: "Configure one or more environment files that are loaded into the task and may be used during substitution. " +
				"Environment files are loaded in the order they are specified and later ones overwrite already existing values.",
			Type: StringSliceType,
		},
		set: func(val unit.Options, t *Task) error {
			if val == nil {
				t.EnvironmentFiles = nil
				return nil
			}

			files, err := val.GetRequiredStringSlice("Environment")
			if err != nil {
				return err
			}

			t.EnvironmentFiles = files
			return nil
		},
		get: func(t *Task) []string {
			return t.EnvironmentFiles
		},
	},
}
