package deploy

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/ppacher/system-deploy/pkg/unit"
)

// Section is a help section to be printed with a plugin
// definition.
type Section struct {
	Title       string
	Description string
}

// OptionSpec describes an option
type OptionSpec struct {
	// Name is the name of the option.
	Name string

	// Description is a human readable description of
	// the option.
	Description string

	// Type defines the type of the option.
	Type OptionType

	// Required may be set to true if deploy tasks must
	// specify this option.
	Required bool

	// Default may holds the default value for this option.
	// This value is only for help purposes and is NOT SET
	// as the default for that option.
	Default string
}

// AllowAny is a special option that can be used to disable
// option validation. Only use during development.
var AllowAny = []OptionSpec{}

// IsAllowAny returns true if spec is the constant AllowAny
// identifier.
func IsAllowAny(spec []OptionSpec) bool {
	return reflect.ValueOf(spec).Pointer() == reflect.ValueOf(AllowAny).Pointer()
}

// Validate validates if all unit options specified in sec conform
// to the specification options.
func Validate(sec unit.Section, options []OptionSpec) error {
	if IsAllowAny(options) {
		return nil
	}

	// build a lookup map for all options specs.
	lm := make(map[string]OptionSpec)
	for _, spec := range options {
		lm[strings.ToLower(spec.Name)] = spec
	}

	// group option values by option name.
	gv := make(map[string][]string)
	for _, opt := range sec.Options {
		n := strings.ToLower(opt.Name)
		gv[n] = append(gv[n], opt.Value)
	}

	// validate
	for name, values := range gv {
		spec, ok := lm[name]
		if !ok {
			// TODO(ppacher): we always use the lowercase version for the
			// error message here, use the original one instead.
			return fmt.Errorf("%s: %w", name, ErrOptionNotExists)
		}

		if err := ValidateOption(values, spec); err != nil {
			return fmt.Errorf("%s: %w", spec.Name, err)
		}

		// delete the spec from the lookup map
		// so any spec left-over may cause a Required
		// error.
		delete(lm, name)
	}

	// check if any option that is required is
	// missing completely
	for _, spec := range lm {
		if spec.Required {
			return fmt.Errorf("%s: %w", spec.Name, ErrOptionRequired)
		}
	}

	return nil
}

// ValidateOption validates if values matches spec.
func ValidateOption(values []string, spec OptionSpec) error {
	if len(values) > 1 && !spec.Type.IsSliceType() {
		return ErrOptionAllowedOnce
	}

	if spec.Required && len(values) == 0 {
		return ErrOptionRequired
	}

	for _, v := range values {
		// all occurences must have a value set
		// if the option is required.
		if spec.Required && v == "" {
			return ErrOptionRequired
		}

		// ensure the value matches the types expecations.
		if err := checkValue(v, spec.Type); err != nil {
			return err
		}
	}

	return nil
}

func checkValue(val string, optType OptionType) error {
	switch optType {
	case BoolType:
		if _, err := unit.ConvertBool(val); err != nil {
			return ErrInvalidBoolean
		}
	case StringSliceType, StringType:
		// we cannot validate anything here
		return nil
	case FloatSliceType, FloatType:
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return ErrInvalidFloat
		}
	case IntSliceType, IntType:
		// we support all number formats supported by ParseInt.
		// That is, hex (0xYY), binary (0bYY) and octal (0YY)
		if _, err := strconv.ParseInt(val, 0, 64); err != nil {
			return ErrInvalidNumber
		}
	}

	return nil
}
