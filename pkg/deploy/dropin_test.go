package deploy

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ppacher/system-deploy/pkg/unit"
	"github.com/stretchr/testify/assert"
)

func TestDropInSearchPaths(t *testing.T) {
	paths := DropInSearchPaths("foo-bar-baz.task", "/lib/")

	assert.Equal(t, []string{
		"/lib/task.d",
		"/lib/foo-.task.d",
		"/lib/foo-bar-.task.d",
		"/lib/foo-bar-baz.task.d",
	}, paths)
}

type fakeFileInfo struct {
	name  string
	isDir bool
}

func (t *fakeFileInfo) IsDir() bool {
	return t.isDir
}
func (t *fakeFileInfo) Name() string {
	return t.name
}
func (*fakeFileInfo) ModTime() time.Time { return time.Now() }
func (*fakeFileInfo) Mode() os.FileMode  { return 0600 }
func (*fakeFileInfo) Size() int64        { return 100 }
func (*fakeFileInfo) Sys() interface{}   { return nil }

func fakeFile(name string, dir bool) os.FileInfo {
	return &fakeFileInfo{name, dir}
}

func TestSearchDropinFiles(t *testing.T) {
	// restore readDir after this test case
	defer func() {
		readDir = ioutil.ReadDir
	}()
	readDir = func(path string) ([]os.FileInfo, error) {
		switch {
		case strings.HasPrefix(path, "/lib/task.d"):
			return []os.FileInfo{
				fakeFile("test", false),
				fakeFile("dir.conf", true),
				fakeFile("10-overwrite.conf", false),
				fakeFile("20-task.d.conf", false),
			}, nil
		case strings.HasPrefix(path, "/lib/foo-.task.d"):
			return []os.FileInfo{
				fakeFile("test2", false),
				fakeFile("10-overwrite.conf", false),
				fakeFile("30-foo-task.d.conf", false),
			}, nil
		case strings.HasPrefix(path, "/lib/foo-bar-baz.task.d"):
			return []os.FileInfo{
				fakeFile("10-overwrite.conf", false),
			}, nil
		}

		return nil, os.ErrNotExist
	}

	paths, err := SearchDropinFiles("foo-bar-baz.task", []string{"/lib/"})
	assert.NoError(t, err)
	assert.Equal(t, []string{
		"/lib/foo-bar-baz.task.d/10-overwrite.conf",
		"/lib/task.d/20-task.d.conf",
		"/lib/foo-.task.d/30-foo-task.d.conf",
	}, paths)
}

func TestApplyDropIns(t *testing.T) {
	specs := map[string]map[string]OptionSpec{
		"test": {
			"single": {
				Type: StringType,
			},
			"slice1": {
				Type: StringSliceType,
			},
			"slice2": {
				Type: StringSliceType,
			},
		},
	}
	tsk := &Task{
		Sections: []unit.Section{
			{
				Name: "Test",
				Options: unit.Options{
					{
						Name:  "Single",
						Value: "SV",
					},
					{
						Name:  "Slice1",
						Value: "Value1",
					},
				},
			},
		},
	}
	d1 := &DropIn{
		Task: &Task{
			Sections: []unit.Section{
				{
					Name: "Test",
					Options: unit.Options{
						{
							Name:  "Single",
							Value: "from d1",
						},
						{
							Name:  "Slice2",
							Value: "from d1.1",
						},
						{
							Name:  "Slice2",
							Value: "from d1.2",
						},
						{
							Name:  "Slice1",
							Value: "Value2",
						},
					},
				},
			},
		},
	}

	d2 := &DropIn{
		Meta: &unit.Section{
			Name: "Task",
			Options: unit.Options{
				{
					Name:  "Description",
					Value: "Test Description",
				},
				{
					Name:  "Disabled",
					Value: "true",
				},
			},
		},
		Task: &Task{
			Sections: []unit.Section{
				{
					Name: "Test",
					Options: unit.Options{
						{
							Name:  "Slice1",
							Value: "", // clear all values
						},
						{
							Name:  "Slice1",
							Value: "Value2",
						},
					},
				},
			},
		},
	}

	res, err := ApplyDropIns(tsk, []*DropIn{d1, d2}, specs)
	assert.NoError(t, err)
	assert.Equal(t, &Task{
		Disabled:    true,
		Description: "Test Description",
		Sections: []unit.Section{
			{
				Name: "Test",
				Options: unit.Options{
					{
						Name:  "Single",
						Value: "from d1",
					},
					{
						Name:  "Slice2",
						Value: "from d1.1",
					},
					{
						Name:  "Slice2",
						Value: "from d1.2",
					},
					{
						Name:  "Slice1",
						Value: "Value2",
					},
				},
			},
		},
	}, res)
}

func TestApplyDropInsNotAllowed(t *testing.T) {
	tsk := &Task{
		Sections: []unit.Section{
			{
				Name: "Test",
			},
			{
				Name: "Test",
			},
		},
	}

	tsk, err := ApplyDropIns(tsk,
		[]*DropIn{
			{
				Task: &Task{
					Sections: []unit.Section{
						{ // section Test is not allowed because it's not unique in tsk
							Name: "Test",
						},
					},
				},
			},
		},
		map[string]map[string]OptionSpec{
			"test": nil,
		},
	)

	assert.Nil(t, tsk)
	assert.Error(t, err)
}

func TestApplyDropInsSectionNotExists(t *testing.T) {
	tsk := &Task{
		Sections: []unit.Section{},
	}

	tsk, err := ApplyDropIns(tsk,
		[]*DropIn{
			{
				Task: &Task{
					Sections: []unit.Section{
						{ // section Test is not allowed because it's not unique in tsk
							Name: "Unknown",
						},
					},
				},
			},
		},
		map[string]map[string]OptionSpec{
			"test": nil,
		},
	)

	assert.Nil(t, tsk)
	assert.Error(t, err)
}

func TestApplyDropInsOptionNotExists(t *testing.T) {
	tsk := &Task{
		Sections: []unit.Section{
			{
				Name: "Test",
			},
		},
	}

	tsk, err := ApplyDropIns(tsk,
		[]*DropIn{
			{
				Task: &Task{
					Sections: []unit.Section{
						{ // section Test is not allowed because it's not unique in tsk
							Name: "Test",
							Options: unit.Options{
								{
									Name: "does-not-exist",
								},
							},
						},
					},
				},
			},
		},
		map[string]map[string]OptionSpec{
			"test": nil,
		},
	)

	assert.Nil(t, tsk)
	assert.Error(t, err)
}
