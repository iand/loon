package loon

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func Parse(b []byte) (*Doc, error) {
	p := NewParser(bytes.NewReader(b))
	return p.Parse()
}

type Parser struct {
	scanner *bufio.Scanner
	state   int
}

const (
	readingPreamble = iota
	seekingObject
	readingObject
)

func NewParser(r io.Reader) *Parser {
	return &Parser{
		scanner: bufio.NewScanner(r),
		state:   readingPreamble,
	}
}

func (p *Parser) Parse() (*Doc, error) {
	doc := &Doc{
		Version:          1,
		Comments:         make([]string, 0),
		TrailingComments: make([]string, 0),
	}

	obj := Object{
		Comments:         make([]string, 0),
		TrailingComments: make([]string, 0),
	}
	comments := make([]string, 0)

	linenum := 0
scanloop:
	for p.scanner.Scan() {
		linenum++
		line := strings.TrimSpace(p.scanner.Text())

		switch p.state {
		case readingPreamble, seekingObject:
			// An empty line indicates that the comments don't belong to the first object
			if len(line) == 0 {
				if p.state == readingPreamble {
					doc.Comments = append(doc.Comments, comments...)
					comments = make([]string, 0)
				}
				continue scanloop
			}

			if strings.HasPrefix(line, "#") {
				// gather comment
				c := strings.TrimSpace(line[1:])
				comments = append(comments, c)
				continue scanloop
			} else if strings.HasPrefix(line, "@") {
				if p.state == readingPreamble {
					if line == "@version 1" {
						if linenum != 1 {
							return nil, &ParseError{
								Err:  ErrVersionTagMustBeFirst,
								Line: linenum,
							}
						}
						// no other version exists yet
					} else {
						return nil, &ParseError{
							Err:  ErrInvalidDocumentTag,
							Line: linenum,
						}
					}

					doc.Comments = append(doc.Comments, comments...)
					comments = nil
					continue scanloop
				}
				// doc tags not allowed between objects
				return nil, &ParseError{
					Err:  ErrInvalidDocumentTag,
					Line: linenum,
				}
			}

			// Assume start of object
			keyword, _, args := splitLine(line)
			if len(args) == 0 {
				return nil, &ParseError{
					Err:  ErrMissingObjectName,
					Line: linenum,
				}
			}

			if len(args) > 1 {
				return nil, &ParseError{
					Err:  ErrInvalidObjectName,
					Line: linenum,
				}
			}

			if keyword == "" {
				return nil, &ParseError{
					Err:  ErrInvalidObjectType,
					Line: linenum,
				}
			}

			obj.Type = keyword
			obj.Name = args[0]
			obj.Comments = comments

			comments = make([]string, 0)
			p.state = readingObject

		case readingObject:
			if len(line) == 0 {
				continue scanloop
			}
			if strings.HasPrefix(line, "#") {
				// gather comment
				c := strings.TrimSpace(line[1:])
				comments = append(comments, c)
				continue scanloop
			}

			if line == "end" {
				obj.TrailingComments = comments
				comments = make([]string, 0)
				doc.Objects = append(doc.Objects, obj)
				obj = Object{
					Comments:         make([]string, 0),
					TrailingComments: make([]string, 0),
				}
				p.state = seekingObject
			} else {
				keyword, argtext, args := splitLine(line)
				obj.Directives = append(obj.Directives, Directive{
					Name:     keyword,
					Args:     args,
					ArgText:  argtext,
					Comments: comments,
				})
				comments = make([]string, 0)
			}
		}

	}

	if p.scanner.Err() != nil {
		return nil, p.scanner.Err()
	}

	doc.TrailingComments = comments
	return doc, nil
}

func splitLine(line string) (string, string, []string) {
	keyword := line
	rawargs := ""
	endOfDirective := strings.IndexFunc(line, unicode.IsSpace)
	if endOfDirective != -1 {
		keyword = line[:endOfDirective]
		rawargs = strings.TrimSpace(line[endOfDirective+1:])
	}
	args := strings.Fields(rawargs)
	return keyword, rawargs, args
}

type Doc struct {
	// Version is the version of the loon syntax
	Version          int
	Comments         []string
	Objects          []Object
	TrailingComments []string
}

type Object struct {
	Type             string
	Name             string
	Directives       []Directive
	Line             int
	Comments         []string // contains comment lines that occurred before the object
	TrailingComments []string // contains comment lines that occurred after the last directive in the object
}

type Directive struct {
	Name string
	Args []string

	// Line is the line number of the directive
	Line int

	// ArgText is the original unparsed text of the arguments
	ArgText string

	// Comments that are associated with the directive
	Comments []string
}

var (
	ErrMissingObjectName     = errors.New("missing object name")
	ErrInvalidObjectName     = errors.New("invalid object name")
	ErrInvalidObjectType     = errors.New("invalid object type")
	ErrInvalidDocumentTag    = errors.New("invalid document tag")
	ErrVersionTagMustBeFirst = errors.New("version tag must be on first line")
)

type ParseError struct {
	Err  error
	Line int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error (line:%d): %v", e.Line, e.Err)
}

func (e *ParseError) Unwrap() error {
	return e.Err
}
