package unit

import (
	"errors"
	"strconv"
	"strings"
)

var (
	// ErrOptionNotSet is returned when a given option is
	// not set.
	ErrOptionNotSet = errors.New("option not set")

	// ErrOptionAllowedOnce is returned if an option is
	// specified multiple times but it's not allowed to.
	ErrOptionAllowedOnce = errors.New("option can only be set once")
)

// IsNotSet returns true if err is ErrOptionNotSet
func IsNotSet(err error) bool {
	if err == ErrOptionNotSet {
		return true
	}
	return errors.Is(err, ErrOptionNotSet)
}

// Options is a convenience type for working with a slice
// of deploy options.
type Options []Option

// GetString returns the value of the option name. It ensures
// that name is only set once. If name can be specified multiple
// time use GetStringSlice
func (opts Options) GetString(name string) (string, error) {
	var found bool
	var s string

	name = strings.ToLower(name)

	for _, opt := range opts {
		if strings.ToLower(opt.Name) == name {
			if found {
				return "", ErrOptionAllowedOnce
			}
			s = opt.Value
			found = true
		}
	}

	if !found {
		return "", ErrOptionNotSet
	}
	return s, nil
}

// GetStringSlice returns a slice of values specified for name.
// If name must be specified at least once, use GetRequiredStringSlice.
func (opts Options) GetStringSlice(name string) []string {
	var s []string
	name = strings.ToLower(name)

	for _, opt := range opts {
		if strings.ToLower(opt.Name) == name {
			s = append(s, opt.Value)
		}
	}

	return s
}

// GetRequiredStringSlice is like GetStringSlice but returns an error if
// name is not specified at least once.
func (opts Options) GetRequiredStringSlice(name string) ([]string, error) {
	s := opts.GetStringSlice(name)
	if len(s) == 0 {
		return nil, ErrOptionNotSet
	}

	return s, nil
}

// GetBool returns the bool set for name or an error if it's not set.
func (opts Options) GetBool(name string) (bool, error) {
	val, err := opts.GetString(name)
	if err != nil {
		return false, err
	}

	return ConvertBool(val)
}

// GetBoolDefault returns the boolean value for name or a default if
// name is not set or cannot be parsed as a bool.
func (opts Options) GetBoolDefault(name string, def bool) bool {
	b, err := opts.GetBool(name)
	if err != nil {
		return def
	}

	return b
}

// GetInt returns the integer value set for name. If name is not set
// or cannot be parsed an error is returned.
func (opts Options) GetInt(name string) (int64, error) {
	val, err := opts.GetString(name)
	if err != nil {
		return 0, err
	}

	return strconv.ParseInt(val, 0, 64)
}

// GetIntDefault is like GetInt but returns a default value if
// name is not set or cannot be parsed
func (opts Options) GetIntDefault(name string, def int64) int64 {
	val, err := opts.GetInt(name)
	if err != nil {
		return def
	}
	return val
}

// GetIntSlice returns a slice of integer values set for name.
func (opts Options) GetIntSlice(name string) ([]int64, error) {
	var sl []int64

	name = strings.ToLower(name)
	for _, opt := range opts {
		if strings.ToLower(opt.Name) == name {
			i, err := strconv.ParseInt(opt.Value, 0, 64)
			if err != nil {
				return nil, err
			}
			sl = append(sl, i)
		}
	}

	return sl, nil
}

// GetRequiredIntSlice is like GetIntSlice but returns an error if
// name is not set at least once.
func (opts Options) GetRequiredIntSlice(name string) ([]int64, error) {
	sl, err := opts.GetIntSlice(name)
	if err != nil {
		return nil, err
	}

	if len(sl) == 0 {
		return nil, ErrOptionNotSet
	}

	return sl, nil
}

// GetFloat returns the float set for name. If name is not set or
// cannot be parsed an error is returned.
func (opts Options) GetFloat(name string) (float64, error) {
	val, err := opts.GetString(name)
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(val, 64)
}

// GetFloatDefault is like GetFloat but returns a default value if
// name is not set or cannot be parsed.
func (opts Options) GetFloatDefault(name string, def float64) float64 {
	val, err := opts.GetFloat(name)
	if err != nil {
		return def
	}

	return val
}

// GetFloatSlice returns a slice of float values set for name. If
// a value cannot be parsed an error is returned.
func (opts Options) GetFloatSlice(name string) ([]float64, error) {
	var fs []float64
	name = strings.ToLower(name)

	for _, opt := range opts {
		if strings.ToLower(opt.Name) == name {
			f, err := strconv.ParseFloat(opt.Value, 64)
			if err != nil {
				return nil, err
			}
			fs = append(fs, f)
		}
	}

	return fs, nil
}

// GetRequiredFloatSlice is like GetFloatSlice but returns an error if
// name is not specified at least once.
func (opts Options) GetRequiredFloatSlice(name string) ([]float64, error) {
	fs, err := opts.GetFloatSlice(name)
	if err != nil {
		return nil, err
	}

	if len(fs) == 0 {
		return nil, ErrOptionNotSet
	}

	return fs, nil
}

// ConvertBool converts the string s into
// a boolean value if it matches one of the supported
// boolean identifiers.
func ConvertBool(s string) (bool, error) {
	switch s {
	case "yes", "Yes", "YES", "on", "ON":
		return true, nil
	case "no", "No", "NO", "off", "OFF":
		return false, nil
	}
	return strconv.ParseBool(s)
}
