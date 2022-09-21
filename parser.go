package loon

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"unicode"
)

type Parser struct {
	scanner *bufio.Scanner
	err     error
	linenum int
	object  *Object
}

func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
		linenum: -1,
	}
}

func (p *Parser) Next() bool {
	p.object = nil

	if p.err != nil {
		return false
	}

	for p.scanner.Scan() {
		p.linenum++

		line := strings.TrimSpace(p.scanner.Text())
		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		directive := line
		data := ""

		endOfDirective := strings.IndexFunc(directive, unicode.IsSpace)
		if endOfDirective != -1 {
			directive = line[:endOfDirective]
			data = strings.TrimSpace(line[endOfDirective+1:])
		}
		args := strings.Fields(data)

		if p.object == nil {
			if len(args) == 0 {
				p.err = ErrMissingObjectName
				return false
			}

			if len(args) > 1 {
				p.err = ErrInvalidObjectName
				return false
			}

			// TODO: check if args len > 0
			p.object = &Object{
				Line: p.linenum,
				Type: directive,
				Name: args[0],
			}
		} else {
			if directive == "end" {
				return true
			} else {
				p.object.Directives = append(p.object.Directives, Directive{
					Line:    p.linenum,
					Name:    directive,
					Args:    args,
					ArgText: data,
				})
			}
		}
	}

	if p.scanner.Err() != nil {
		p.err = p.scanner.Err()
	}

	return false
}

func (p *Parser) Object() *Object {
	return p.object
}

func (p *Parser) Err() error {
	return p.err
}

type Object struct {
	Type       string
	Name       string
	Directives []Directive
	Line       int
}

type Directive struct {
	Name string
	Args []string
	// ArgText is the original unparsed text of the arguments
	ArgText string
	Line    int
}

var (
	ErrMissingObjectName = errors.New("missing object name")
	ErrInvalidObjectName = errors.New("invalid object name")
)
