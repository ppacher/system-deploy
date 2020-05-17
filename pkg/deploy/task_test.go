package deploy

import (
	"errors"
	"strings"
	"testing"

	"github.com/ppacher/system-deploy/pkg/unit"
	"github.com/stretchr/testify/assert"
)

func TestDecode(t *testing.T) {
	cases := []struct {
		I string
		T *Task
		E error
	}{
		{
			`
[Task]
StartMasked=yes
Disabled=1

[Section1]
Key1 = Value1

[Section2]
Key2= Value2
			`,
			&Task{
				StartMasked: true,
				Disabled:    true,
				Sections: []unit.Section{
					{
						Name: "Section1",
						Options: []unit.Option{
							{
								Name:  "Key1",
								Value: "Value1",
							},
						},
					},
					{
						Name: "Section2",
						Options: []unit.Option{
							{
								Name:  "Key2",
								Value: "Value2",
							},
						},
					},
				},
			},
			nil,
		},
		{
			"[Task]\nStartMasked=InvalidValue",
			nil,
			ErrInvalidTaskSection,
		},
		{
			"[Task]\n",
			nil,
			ErrNoSections,
		},
	}

	for idx, c := range cases {
		tsk, err := Decode("test-file", strings.NewReader(c.I))
		if tsk != nil {
			// there's not file name in tests.
			tsk.FileName = ""
			tsk.Directory = ""
		}

		if !errors.Is(err, c.E) {
			t.Errorf("case #%d: expected error to be '%v' but got '%v'", idx, c.E, err)
		}

		assert.Equal(t, c.T, tsk)
	}
}
