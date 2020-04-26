// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package unit provides a systemd unit file lexer based on
// github.com/coreos/go-systemd with slight modifications.
package unit

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	// SystemdLineMax mimics the maximum line length that systemd can use.
	// On typical systemd platforms (i.e. modern Linux), this will most
	// commonly be 2048, so let's use that as a sanity check.
	// Technically, we should probably pull this at runtime:
	//    SystemdLineMax = int(C.sysconf(C.__SC_LINE_MAX))
	// but this would introduce an (unfortunate) dependency on cgo
	SystemdLineMax = 2048

	// SystemdNewline defines characters that systemd considers indicators
	// for a newline.
	SystemdNewline = "\r\n"
)

var (
	// ErrLineTooLong gets returned when a line is too long for systemd to handle.
	ErrLineTooLong = fmt.Errorf("line too long (max %d bytes)", SystemdLineMax)
)

type (
	Option struct {
		Name  string
		Value string
	}

	Section struct {
		Options
		Name string
	}
)

// Deserialize parses a systemd unit file into a list of UnitOption objects.
func Deserialize(f io.Reader) (sections []Section, err error) {
	lexer, secchan, errchan := newLexer(f)
	go lexer.lex()

	for sec := range secchan {
		sections = append(sections, *sec)
	}

	err = <-errchan
	return sections, err
}

func newLexer(f io.Reader) (*lexer, <-chan *Section, <-chan error) {
	secchan := make(chan *Section)
	errchan := make(chan error, 1)
	buf := bufio.NewReader(f)

	return &lexer{buf, secchan, errchan, nil}, secchan, errchan
}

type lexer struct {
	buf     *bufio.Reader
	secchan chan *Section
	errchan chan error
	section *Section
}

func (l *lexer) lex() {
	defer func() {
		close(l.secchan)
		close(l.errchan)
	}()
	next := l.lexNextSection
	for next != nil {
		if l.buf.Buffered() >= SystemdLineMax {
			// systemd truncates lines longer than LINE_MAX
			// https://bugs.freedesktop.org/show_bug.cgi?id=85308
			// Rather than allowing this to pass silently, let's
			// explicitly gate people from encountering this
			line, err := l.buf.Peek(SystemdLineMax)
			if err != nil {
				l.errchan <- err
				return
			}
			if !bytes.ContainsAny(line, SystemdNewline) {
				l.errchan <- ErrLineTooLong
				return
			}
		}

		var err error
		next, err = next()
		if err != nil {
			l.errchan <- err
			return
		}
	}

	if l.section != nil {
		l.secchan <- l.section
	}
}

type lexStep func() (lexStep, error)

func (l *lexer) lexSectionName() (lexStep, error) {
	sec, err := l.buf.ReadBytes(']')
	if err != nil {
		return nil, errors.New("unable to find end of section")
	}

	sectionName := string(sec[:len(sec)-1])

	if l.section != nil {
		l.secchan <- l.section
	}

	l.section = &Section{
		Name: sectionName,
	}

	return l.lexSectionSuffixFunc(), nil
}

func (l *lexer) lexSectionSuffixFunc() lexStep {
	return func() (lexStep, error) {
		garbage, _, err := l.toEOL()
		if err != nil {
			return nil, err
		}

		garbage = bytes.TrimSpace(garbage)
		if len(garbage) > 0 {
			return nil, fmt.Errorf("found garbage after section name %s: %v", l.section, garbage)
		}

		return l.lexNextSectionOrOptionFunc(), nil
	}
}

func (l *lexer) ignoreLineFunc(next lexStep) lexStep {
	return func() (lexStep, error) {
		for {
			line, _, err := l.toEOL()
			if err != nil {
				return nil, err
			}

			line = bytes.TrimSuffix(line, []byte{' '})

			// lack of continuation means this line has been exhausted
			if !bytes.HasSuffix(line, []byte{'\\'}) {
				break
			}
		}

		// reached end of buffer, safe to exit
		return next, nil
	}
}

func (l *lexer) lexNextSection() (lexStep, error) {
	r, _, err := l.buf.ReadRune()
	if err != nil {
		if err == io.EOF {
			err = nil
		}
		return nil, err
	}

	if r == '[' {
		return l.lexSectionName, nil
	} else if isComment(r) {
		return l.ignoreLineFunc(l.lexNextSection), nil
	}

	return l.lexNextSection, nil
}

func (l *lexer) lexNextSectionOrOptionFunc() lexStep {
	return func() (lexStep, error) {
		r, _, err := l.buf.ReadRune()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return nil, err
		}

		if unicode.IsSpace(r) {
			return l.lexNextSectionOrOptionFunc(), nil
		} else if r == '[' {
			return l.lexSectionName, nil
		} else if isComment(r) {
			return l.ignoreLineFunc(l.lexNextSectionOrOptionFunc()), nil
		}

		l.buf.UnreadRune()
		return l.lexOptionNameFunc(), nil
	}
}

func (l *lexer) lexOptionNameFunc() lexStep {
	return func() (lexStep, error) {
		var partial bytes.Buffer
		for {
			r, _, err := l.buf.ReadRune()
			if err != nil {
				return nil, err
			}

			if r == '\n' || r == '\r' {
				return nil, errors.New("unexpected newline encountered while parsing option name")
			}

			if r == '=' {
				break
			}

			partial.WriteRune(r)
		}

		name := strings.TrimSpace(partial.String())
		return l.lexOptionValueFunc(name, bytes.Buffer{}), nil
	}
}

func (l *lexer) lexOptionValueFunc(name string, partial bytes.Buffer) lexStep {
	return func() (lexStep, error) {
		for {
			line, eof, err := l.toEOL()
			if err != nil {
				return nil, err
			}

			if len(bytes.TrimSpace(line)) == 0 {
				break
			}

			partial.Write(line)

			// lack of continuation means this value has been exhausted
			idx := bytes.LastIndex(line, []byte{'\\'})
			if idx == -1 || idx != (len(line)-1) {
				break
			}

			if !eof {
				partial.WriteRune('\n')
			}

			return l.lexOptionValueFunc(name, partial), nil
		}

		val := partial.String()
		if strings.HasSuffix(val, "\n") {
			// A newline was added to the end, so the file didn't end with a backslash.
			// => Keep the newline
			val = strings.TrimSpace(val) + "\n"
		} else {
			val = strings.TrimSpace(val)
		}

		if l.section == nil {
			return nil, fmt.Errorf("found option outside of section")
		}

		l.section.Options = append(l.section.Options, Option{Name: name, Value: val})

		return l.lexNextSectionOrOptionFunc(), nil
	}
}

// toEOL reads until the end-of-line or end-of-file.
// Returns (data, EOFfound, error)
func (l *lexer) toEOL() ([]byte, bool, error) {
	line, err := l.buf.ReadBytes('\n')
	// ignore EOF here since it's roughly equivalent to EOL
	if err != nil && err != io.EOF {
		return nil, false, err
	}

	line = bytes.TrimSuffix(line, []byte{'\r'})
	line = bytes.TrimSuffix(line, []byte{'\n'})

	return line, err == io.EOF, nil
}

func isComment(r rune) bool {
	return r == '#' || r == ';'
}
