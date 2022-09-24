package loon

import (
	"bytes"
	"fmt"
)

func Print(doc *Doc) []byte {
	pr := NewPrinter()
	pr.Print(doc)
	return pr.Bytes()
}

type Printer struct {
	bytes.Buffer
}

func NewPrinter() *Printer {
	return &Printer{}
}

func (p *Printer) Print(doc *Doc) {
	if doc.Version != 1 {
		fmt.Fprintf(p, "@version %d\n", doc.Version)
	}

	if len(doc.Comments) > 0 {
		for _, c := range doc.Comments {
			fmt.Fprintf(p, "# %s\n", c)
		}
		fmt.Fprintln(p, "")
	}
	for _, o := range doc.Objects {
		for _, c := range o.Comments {
			fmt.Fprintf(p, "# %s\n", c)
		}
		fmt.Fprintf(p, "%s %s\n", o.Type, o.Name)
		for _, d := range o.Directives {
			for _, c := range d.Comments {
				fmt.Fprintf(p, "\t# %s\n", c)
			}
			if d.ArgText != "" {
				fmt.Fprintf(p, "\t%s %s\n", d.Name, d.ArgText)
			}
		}
		for _, c := range o.TrailingComments {
			fmt.Fprintf(p, "\t# %s\n", c)
		}
		fmt.Fprintln(p, "end")
		fmt.Fprintln(p, "")
	}

	for _, c := range doc.TrailingComments {
		fmt.Fprintf(p, "# %s\n", c)
	}
}
