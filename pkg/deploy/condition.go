package deploy

import (
	"github.com/ppacher/system-deploy/pkg/condition"
	"github.com/ppacher/system-deploy/pkg/unit"
	"github.com/sirupsen/logrus"
)

// EvaluateConditions evalutes all conditions of a task. If any
// condition fails the condition is returned along with an error.
// If all conditions are met, EvalutateConditions returns nil, nil.
func EvaluateConditions(t *Task) (*condition.Instance, error) {
	for _, cond := range t.Conditions {
		logrus.Debugf("%s: evaluating conditon %s", t.FileName, cond.Name)
		if err := cond.Run(); err != nil {
			return &cond, err
		}
	}

	return nil, nil
}

// RegisterCondition registers a new condition type for
// the task meta-section.
func RegisterCondition(cond condition.Condition) {
	logrus.Debugf("registering condition %s", cond.Name)

	condName := "Condition" + cond.Name
	assertName := "Assert" + cond.Name
	condSpec := OptionSpec{
		Name:        condName,
		Aliases:     []string{assertName},
		Description: cond.Description,
		Type:        StringSliceType,
	}
	// assertSpec is an internal-only option to handle
	// the assertName alias of condSpec.
	assertSpec := OptionSpec{
		Name:     assertName,
		Type:     StringSliceType,
		Internal: true,
	}

	var values []string
	getSetter := func(assert bool) func(val unit.Options, t *Task) error {
		return func(val unit.Options, t *Task) error {
			if val == nil {
				values = nil

				// remove the condition instance from the tasks instance
				// list.
				for idx, instance := range t.Conditions {
					if instance.Condition.Name == cond.Name &&
						instance.Assertion == assert {
						t.Conditions = append(t.Conditions[:idx], t.Conditions[idx+1:]...)
					}
				}

				return nil
			}

			name := condName
			if assert {
				name = assertName
			}

			values = val.GetStringSlice(name)
			instance := condition.Instance{
				Condition: &cond,
				Assertion: assert,
				Values:    values,
			}

			t.Conditions = append(t.Conditions, instance)
			what := "condition"
			if assert {
				what = "assertion"
			}
			logrus.Debugf("%s: added %s %s for values %v", t.FileName, what, instance.Name, instance.Values)

			return nil
		}
	}
	get := func(t *Task) []string {
		return values
	}

	taskOptions = append(taskOptions, taskMetaOption{
		OptionSpec: condSpec,
		get:        get,
		set:        getSetter(false),
	})

	taskOptions = append(taskOptions, taskMetaOption{
		OptionSpec: assertSpec,
		get:        get,
		set:        getSetter(true),
	})
}

// RegisterAllConditions registers all built-in conditions
// from the condition package.
func RegisterAllConditions() {
	for _, cond := range condition.BuiltinConditions {
		RegisterCondition(cond)
	}
}
