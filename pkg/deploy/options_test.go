package deploy

import (
	"errors"
	"testing"

	"github.com/ppacher/system-deploy/pkg/unit"
)

func TestCheckValue(t *testing.T) {
	cases := []struct {
		T OptionType
		V string
		E error
	}{
		{BoolType, "yes", nil},
		{BoolType, "false", nil},
		{BoolType, "0", nil},
		{BoolType, "foo", ErrInvalidBoolean},

		{IntType, "0x10", nil},
		{IntType, "0600", nil},
		{IntSliceType, "0b1100", nil},
		{IntType, "INVALID", ErrInvalidNumber},
		{IntSliceType, "INVALID2", ErrInvalidNumber},

		{FloatType, "0.5", nil},
		{FloatType, ".5", nil},
		{FloatSliceType, "0.1e10", nil},
		{FloatType, ".INVALID", ErrInvalidFloat},
		{FloatSliceType, "0.1eINVALID", ErrInvalidFloat},

		{StringType, "", nil}, // empty strings ARE VALID
	}

	for idx, c := range cases {
		err := checkValue(c.V, c.T)

		if !errors.Is(err, c.E) {
			t.Errorf("case #%d (input=%v) expected error to be %s but got %s", idx, c.V, c.E, err)
		}
	}
}

func TestValidateOption(t *testing.T) {
	cases := []struct {
		I OptionSpec
		V []string
		E error
	}{
		{
			OptionSpec{
				Required: true,
				Type:     BoolType,
			},
			nil,
			ErrOptionRequired,
		},
		{
			OptionSpec{
				Required: true,
				Type:     BoolType,
			},
			[]string{""},
			ErrOptionRequired,
		},
		{
			OptionSpec{
				Required: true,
				Type:     StringSliceType,
			},
			[]string{"value", ""},
			ErrOptionRequired,
		},
		{
			OptionSpec{
				Type: StringType,
			},
			[]string{"one", "two"},
			ErrOptionAllowedOnce,
		},
		{
			OptionSpec{
				Type: IntSliceType,
			},
			[]string{"1", "2", "", "0.5"},
			ErrInvalidNumber,
		},
	}

	for idx, c := range cases {
		err := ValidateOption(c.V, c.I)
		if !errors.Is(err, c.E) {
			t.Errorf("cases #%d (input=%v): expected error to be '%v', got '%v'", idx, c.V, c.E, err)
		}
	}
}

func TestValidate(t *testing.T) {
	cases := []struct {
		I []OptionSpec
		V []unit.Option
		E error
	}{
		{
			[]OptionSpec{
				{Name: "Option1", Type: StringType},
			},
			[]unit.Option{
				{Name: "Option1", Value: "value"},
			},
			nil,
		},
		{
			[]OptionSpec{
				{Name: "Option1", Type: StringSliceType},
			},
			[]unit.Option{
				{Name: "Option1", Value: "value"},
				{Name: "Option1", Value: "value"},
			},
			nil,
		},
		{
			[]OptionSpec{
				{Name: "Option1", Type: StringType},
			},
			[]unit.Option{
				{Name: "Option1", Value: "value"},
				{Name: "Option1", Value: "value"},
			},
			ErrOptionAllowedOnce,
		},
		{
			AllowAny,
			[]unit.Option{
				{Name: "Option1", Value: "value"},
			},
			nil,
		},
		{
			[]OptionSpec{
				{Name: "Option1", Type: StringType},
			},
			[]unit.Option{
				{Name: "Option2", Value: "value"},
			},
			ErrOptionNotExists,
		},
		{
			[]OptionSpec{
				{Name: "Option1", Type: StringType, Required: true},
			},
			nil,
			ErrOptionRequired,
		},
	}

	for idx, c := range cases {
		err := Validate(unit.Section{Options: c.V}, c.I)
		if !errors.Is(err, c.E) {
			t.Errorf("cases #%d (input=%v): expected errot to be '%v' but got '%v'", idx, c.V, c.E, err)
		}
	}
}
