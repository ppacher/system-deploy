package deploy

import (
	"reflect"
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

func IsAllowAny(spec []OptionSpec) bool {
	return reflect.ValueOf(spec).Pointer() == reflect.ValueOf(AllowAny).Pointer()
}
