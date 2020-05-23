package deploy

import (
	"fmt"
	"os"

	"github.com/a8m/envsubst/parse"
	"github.com/ppacher/system-deploy/pkg/utils/envfile"
)

// Envsubst applies environment substition on value.
func (tsk *Task) Envsubst(file, value string) (string, error) {
	r := &parse.Restrictions{
		NoEmpty: false,
		NoUnset: true,
	}

	p := parse.New(file, tsk.Environment, r)
	return p.Parse(value)
}

// LoadEnv loads all environment files specified in t.EnvironmentFiles
// and populates t.Environment.
func LoadEnv(t *Task) error {
	env := make(map[string]string)

	for _, file := range t.EnvironmentFiles {
		var err error
		env, err = loadEnv(file, env)
		if err != nil {
			return err
		}
	}

	t.Environment = make([]string, 0, len(env))
	for key, value := range env {
		t.Environment = append(t.Environment, fmt.Sprintf("%s=%s", key, value))
	}

	return nil
}

// ApplyEnvironment applies environment variable substitution
// to all unit options. If tsk.Environment is nil, ApplyEnvironment
// tries to call LoadEnv(tsk) first. Note that ApplyEnvironment does
// not substitude variables in the tasks meta section.
func ApplyEnvironment(tsk *Task) (*Task, error) {
	tsk = tsk.Clone()

	if tsk.Environment == nil {
		if err := LoadEnv(tsk); err != nil {
			return nil, err
		}
	}

	for idx, sec := range tsk.Sections {
		for optIdx, opt := range sec.Options {
			var err error
			opt.Value, err = tsk.Envsubst(tsk.FileName, opt.Value)
			if err != nil {
				return nil, err
			}

			sec.Options[optIdx] = opt
		}

		tsk.Sections[idx] = sec
	}

	return tsk, nil
}

func loadEnv(file string, env map[string]string) (map[string]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	p := envfile.NewWithConfig(file, f, envfile.Config{
		Env:                env,
		EnableSubstitution: true,
	})

	if err := p.Parse(); err != nil {
		return nil, err
	}

	return p.Env(), nil
}
