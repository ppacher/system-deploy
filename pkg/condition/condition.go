package condition

import (
	"fmt"
	"strings"
)

// Condition describes a generic condtion.
type Condition struct {
	// Name is the name of the condition.
	Name string

	// Description holds a human readable description of
	// the condition.
	Description string

	// check the condition and return the evaluation
	// result.
	check func(value string) (bool, error)
}

// Instance is a specific instance of a condition that
// evaluages against a set of values.
type Instance struct {
	// Condition is the condition for this instance
	*Condition

	// Assertion can be set to true if the condition
	// represents an assertion. A failed assertion will
	// cause the task to fail while a failed condition
	// will disable the task.
	Assertion bool

	// Values holds the conditions values.
	Values []string
}

// Run checks c against all values. The results of checking
// each value is ANDed. If a value is prefixed with an
// exclamation mark the check is negated. Run aborts on the
// first values that evalutes to false.
func (instance *Instance) Run() error {
	for idx, v := range instance.Values {
		negate := false
		if strings.HasPrefix(v, "!") {
			negate = true
			v = v[1:]
		} else if strings.HasPrefix(v, "\\!") {
			// the exclamation mark is part of the
			// value and has been escaped.
			v = v[1:]
		}

		result, err := instance.check(v)
		if err != nil {
			return fmt.Errorf("%q: %w", instance.Values[idx], err)
		}

		switch {
		case negate && !result,
			!negate && result:
			return nil
		default:
			return fmt.Errorf("%q failed", instance.Values[idx])
		}
	}

	return nil
}
