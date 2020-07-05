package deploy

import "github.com/ppacher/system-conf/conf"

// ApplyDropIns applies all dropins on t. It basically wraps conf.ApplyDropIns and re-creates
// the tasks meta-data.
func ApplyDropIns(t *Task, dropins []*conf.DropIn, spec map[string]map[string]conf.OptionSpec) error {
	if err := conf.ApplyDropIns(t.file, dropins, spec); err != nil {
		return err
	}

	return applyMetaData(t)
}
