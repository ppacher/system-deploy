package envfile

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLexer(t *testing.T) {
	r := strings.NewReader(`
# a comment without any values
VAR="value" # some trailing comment
VAR2=value2
BAR=# bar has no value set
export FOO="bar 'test\"' foo" #another trailing comment
	`)
	lex := New("test", r)

	err := lex.Parse()
	require.NoError(t, err)

	env := lex.Env()
	assert.Equal(t, map[string]string{
		"VAR":  "value",
		"VAR2": "value2",
		"BAR":  "",
		"FOO":  "bar 'test\"' foo",
	}, env)
}

func TestLexerSubstitution(t *testing.T) {
	r := strings.NewReader(`
GOPATH="${HOME}/go"
GOSRC="${GOPATH}/src"
GOBIN="${GOBIN:-~/bin}"
PATH="$PATH:$GOBIN"
	`)

	lex := NewWithConfig("test", r, Config{
		Env: map[string]string{
			"HOME":  "/home/user",
			"GOBIN": "/home/user/bin",
			"PATH":  "/usr/local/bin",
		},
		EnableSubstitution: true,
	})

	err := lex.Parse()
	require.NoError(t, err)

	assert.Equal(t,
		map[string]string{
			"HOME":   "/home/user",
			"GOPATH": "/home/user/go",
			"GOSRC":  "/home/user/go/src",
			"GOBIN":  "/home/user/bin",
			"PATH":   "/usr/local/bin:/home/user/bin",
		},
		lex.Env())
}
