package loon

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var parserTests = []struct {
	spec string
	doc  *Doc
}{
	{
		spec: `
rule test
	in iron_ore 3
	out iron 1
end
`,
		doc: &Doc{
			Version: 1,
			Objects: []Object{
				{
					Type: "rule",
					Name: "test",
					Directives: []Directive{
						{
							Name:     "in",
							Args:     []string{"iron_ore", "3"},
							ArgText:  "iron_ore 3",
							Comments: []string{},
						},
						{
							Name:     "out",
							Args:     []string{"iron", "1"},
							ArgText:  "iron 1",
							Comments: []string{},
						},
					},
					Comments:         []string{},
					TrailingComments: []string{},
				},
			},
			Comments:         []string{},
			TrailingComments: []string{},
		},
	},

	{
		spec: `
rule test
	# directive comment
	in iron_ore 3
	out iron 3
	# trailing directive comment
end

# test2 comment
rule test2
	in iron_ore 1
	out iron 1
end
	`,
		doc: &Doc{
			Version: 1,
			Objects: []Object{
				{
					Type: "rule",
					Name: "test",
					Directives: []Directive{
						{
							Name:    "in",
							Args:    []string{"iron_ore", "3"},
							ArgText: "iron_ore 3",
							Comments: []string{
								"directive comment",
							},
						},
						{
							Name:     "out",
							Args:     []string{"iron", "3"},
							ArgText:  "iron 3",
							Comments: []string{},
						},
					},
					Comments: []string{},
					TrailingComments: []string{
						"trailing directive comment",
					},
				},
				{
					Type: "rule",
					Name: "test2",
					Directives: []Directive{
						{
							Name:     "in",
							Args:     []string{"iron_ore", "1"},
							ArgText:  "iron_ore 1",
							Comments: []string{},
						},
						{
							Name:     "out",
							Args:     []string{"iron", "1"},
							ArgText:  "iron 1",
							Comments: []string{},
						},
					},
					Comments: []string{
						"test2 comment",
					},
					TrailingComments: []string{},
				},
			},
			Comments:         []string{},
			TrailingComments: []string{},
		},
	},

	{
		spec: `
# This is a comment about the document
# Since there is an empty line between it and
# the object comment below

# Define some resources
resource iron_ore
	name Iron Ore
end

resource iron_ingot
	name Iron Ingot
end

# Define some production rules
rule smelter
	in iron_ore 3
	out iron_ingot 3
end

rule blacksmith
	in iron_ingot 1
	out horseshoe 1
end

# a trailing document comment
`,
		doc: &Doc{
			Version: 1,
			Comments: []string{
				"This is a comment about the document",
				"Since there is an empty line between it and",
				"the object comment below",
			},
			Objects: []Object{
				{
					Comments: []string{
						"Define some resources",
					},
					Type: "resource",
					Name: "iron_ore",
					Directives: []Directive{
						{Name: "name", Args: []string{"Iron", "Ore"}, ArgText: "Iron Ore", Comments: []string{}},
					},
					TrailingComments: []string{},
				},
				{
					Comments: []string{},
					Type:     "resource",
					Name:     "iron_ingot",
					Directives: []Directive{
						{Name: "name", Args: []string{"Iron", "Ingot"}, ArgText: "Iron Ingot", Comments: []string{}},
					},
					TrailingComments: []string{},
				},
				{
					Comments: []string{
						"Define some production rules",
					},
					Type: "rule",
					Name: "smelter",
					Directives: []Directive{
						{Name: "in", Args: []string{"iron_ore", "3"}, ArgText: "iron_ore 3", Comments: []string{}},
						{Name: "out", Args: []string{"iron_ingot", "3"}, ArgText: "iron_ingot 3", Comments: []string{}},
					},
					TrailingComments: []string{},
				},
				{
					Comments: []string{},
					Type:     "rule",
					Name:     "blacksmith",
					Directives: []Directive{
						{Name: "in", Args: []string{"iron_ingot", "1"}, ArgText: "iron_ingot 1", Comments: []string{}},
						{Name: "out", Args: []string{"horseshoe", "1"}, ArgText: "horseshoe 1", Comments: []string{}},
					},
					TrailingComments: []string{},
				},
			},
			TrailingComments: []string{
				"a trailing document comment",
			},
		},
	},
}

func TestParser(t *testing.T) {
	for _, tc := range parserTests {
		t.Run("", func(t *testing.T) {
			p := NewParser(strings.NewReader(tc.spec))
			doc, err := p.Parse()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if diff := cmp.Diff(tc.doc, doc, cmpopts.IgnoreFields(Object{}, "Line"), cmpopts.IgnoreFields(Directive{}, "Line")); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
