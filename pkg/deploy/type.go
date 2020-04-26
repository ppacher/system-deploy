package deploy

import "fmt"

// OptionType describes the type of an option. It cannot
// be implemented outside the deploy package.
type OptionType interface {
	option()

	fmt.Stringer
}

// All supported option types.
var (
	StringType      optionType = "string"
	StringSliceType optionType = "[]string"
	BoolType        optionType = "bool"
	IntType         optionType = "int"
	IntSliceType    optionType = "[]int"
	FloatType       optionType = "float"
	FloatSliceType  optionType = "[]float"
)

type optionType string

func (optionType) option() {}

func (o optionType) String() string { return string(o) }
