package envfile

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"text/scanner"

	"github.com/a8m/envsubst/parse"
)

// Config is the configuration for a new envfile parser.
type Config struct {
	Env                map[string]string
	EnableSubstitution bool
}

// Parser is capable of parsing environment files. It supports
// in-line substitution of already defined variables and allows
// C- and Basic-style comments (//, /*, and #). To support existing
// env files, Parser supports variables prefixed with "export".
type Parser struct {
	tokenizer  scanner.Scanner
	substitude bool
	env        map[string]string
}

// New returns a new env-file parser.
func New(fileName string, r io.Reader) *Parser {
	return NewWithConfig(fileName, r, Config{})
}

// NewWithConfig returns a new env-file parser with a config.
func NewWithConfig(fileName string, r io.Reader, cfg Config) *Parser {
	var s scanner.Scanner
	s.Filename = fileName
	s.Mode = scanner.ScanIdents | scanner.ScanStrings | scanner.SkipComments
	s.Init(r)

	env := make(map[string]string)
	for key, value := range cfg.Env {
		env[key] = value
	}

	return &Parser{
		tokenizer:  s,
		env:        env,
		substitude: cfg.EnableSubstitution,
	}
}

// Env returns the parsed environment.
func (lex *Parser) Env() map[string]string {
	return lex.env
}

func (lex *Parser) envSlice() []string {
	sl := make([]string, len(lex.env))
	for key, value := range lex.env {
		sl = append(sl, fmt.Sprintf("%s=%s", key, value))
	}
	return sl
}

func (lex *Parser) replace(val string) (string, error) {
	if !lex.substitude {
		return val, nil
	}

	restrictions := &parse.Restrictions{
		NoEmpty: false,
		NoUnset: true,
	}
	p := parse.New(lex.tokenizer.Pos().Filename, lex.envSlice(), restrictions)
	return p.Parse(val)
}

func (lex *Parser) Parse() error {
	var varName string
	expected := "name"

	lastLine := -1
	inComment := -1

	for tok := lex.tokenizer.Scan(); tok != scanner.EOF; tok = lex.tokenizer.Scan() {
		value := lex.tokenizer.TokenText()

		if inComment > -1 {
			if lex.tokenizer.Pos().Line == inComment {
				// consume the whole line
				continue
			}
			inComment = -1
		}

		if strings.HasPrefix(value, "#") {
			// the rest of the line is a comment and the next
			// line must start with a variable name again
			inComment = lex.tokenizer.Pos().Line

			switch expected {
			case "name": // ok
			case "value":
				lex.env[varName] = "" // empty value
			case "assign":
				return fmt.Errorf("%s: expected assignment but found %s", lex.tokenizer.Pos(), value)
			}

			expected = "name"
			continue
		}

		switch expected {
		case "name":
			if lex.tokenizer.Pos().Line == lastLine {
				return fmt.Errorf("%s: expected newline after value but found %s %s", lex.tokenizer.Pos(), scanner.TokenString(tok), value)
			}

			if tok != scanner.Ident {
				return fmt.Errorf("%s: expected variable name identifier got %s %s", lex.tokenizer.Pos(), scanner.TokenString(tok), value)
			}

			// we support prefixing the variable name with export
			if value == "export" {
				continue
			}

			varName = value
			expected = "assign"

		case "assign":
			if (value != "=" && value != ":") && tok != scanner.Ident {
				return fmt.Errorf("%s: expected assignment ('=' or ':') but found %s %s", lex.tokenizer.Pos(), scanner.TokenString(tok), value)
			}
			expected = "value"

		case "value":
			value, err := unquote(value)
			if err != nil {
				return fmt.Errorf("%s: failed to unquote: %w", lex.tokenizer.Pos(), err)
			}

			value, err = lex.replace(value)
			if err != nil {
				return fmt.Errorf("%s: failed to substitude: %w", lex.tokenizer.Pos(), err)
			}

			lex.env[varName] = value
			expected = "name"
			lastLine = lex.tokenizer.Pos().Line
		}

	}

	return nil
}

func unquote(value string) (string, error) {
	if strings.HasPrefix(value, "\"") || strings.HasPrefix(value, "'") {
		unquoted, err := strconv.Unquote(value)
		if err != nil {
			return "", err
		}
		return unquoted, nil
	}
	return value, nil
}
