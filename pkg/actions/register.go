package actions

import (
	"errors"
	"strings"

	"github.com/ppacher/system-deploy/pkg/deploy"
	"github.com/ppacher/system-deploy/pkg/unit"
)

// ErrNoSetupFunc is returned when a plugin without a Setup
// function is registered.
var ErrNoSetupFunc = errors.New("no setup function defined")

// ErrInvalidAction is returned when a SetupFunc returns an invalid
// (nil) action
var ErrInvalidAction = errors.New("invalid (nil) action returned")

// Register registers a new action fn under Rname.
func Register(plg Plugin) error {
	actionsLock.Lock()
	defer actionsLock.Unlock()

	if plg.Setup == nil {
		return ErrNoSetupFunc
	}

	key := strings.ToLower(plg.Name)

	if actions == nil {
		actions = make(map[string]*Plugin)
	}

	if _, ok := actions[key]; ok {
		return errors.New("action exists")
	}
	actions[key] = &plg

	return nil
}

// MustRegister is like Register but panics on error.
func MustRegister(plg Plugin) {
	if err := Register(plg); err != nil {
		panic(err)
	}
}

// GetPlugin returns the plugin by name.
func GetPlugin(name string) (Plugin, bool) {
	actionsLock.RLock()
	defer actionsLock.RUnlock()

	plg, ok := actions[strings.ToLower(name)]
	if !ok {
		return Plugin{}, false
	}

	return *plg, true
}

// ListActions returns a list of all section handlers.
func ListActions() []string {
	actionsLock.RLock()
	defer actionsLock.RUnlock()

	names := []string{}

	for key := range actions {
		names = append(names, key)
	}

	return names
}

// Setup returns the action function for name.
func Setup(name string, log Logger, task deploy.Task, section unit.Section) (Action, error) {
	actionsLock.RLock()
	defer actionsLock.RUnlock()

	key := strings.ToLower(name)

	plg, ok := actions[key]
	if !ok {
		return nil, errors.New("unknown action")
	}

	// validate all section options before calling Setup()
	if err := deploy.Validate(section, plg.Options); err != nil {
		return nil, err
	}

	act, err := plg.Setup(task, section)
	if err != nil {
		return nil, err
	}

	if act == nil {
		return nil, ErrInvalidAction
	}

	act.SetLogger(log)

	return act, nil
}
