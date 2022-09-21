package loon

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var parserTests = []struct {
	spec    string
	objects []Object
}{

	{
		spec: `
rule test
	in iron_ore 3
	out iron 1
end
`,

		objects: []Object{
			{
				Type: "rule",
				Name: "test",
				Directives: []Directive{
					{
						Name:    "in",
						Args:    []string{"iron_ore", "3"},
						ArgText: "iron_ore 3",
						Line:    2,
					},
					{
						Name:    "out",
						Args:    []string{"iron", "1"},
						ArgText: "iron 1",
						Line:    3,
					},
				},
				Line: 1,
			},
		},
	},

	{
		spec: `
	rule test
		#comment
		in iron_ore 3
		out iron 3
	end
	#comment
	rule test2
		in iron_ore 1
		out iron 1
	end
	#comment
	`,

		objects: []Object{
			{
				Type: "rule",
				Name: "test",
				Directives: []Directive{
					{
						Name:    "in",
						Args:    []string{"iron_ore", "3"},
						ArgText: "iron_ore 3",
						Line:    3,
					},
					{
						Name:    "out",
						Args:    []string{"iron", "3"},
						ArgText: "iron 3",
						Line:    4,
					},
				},
				Line: 1,
			},
			{
				Type: "rule",
				Name: "test2",
				Directives: []Directive{
					{
						Name:    "in",
						Args:    []string{"iron_ore", "1"},
						ArgText: "iron_ore 1",
						Line:    8,
					},
					{
						Name:    "out",
						Args:    []string{"iron", "1"},
						ArgText: "iron 1",
						Line:    9,
					},
				},
				Line: 7,
			},
		},
	},
}

func TestParser(t *testing.T) {
	for _, tc := range parserTests {
		t.Run("", func(t *testing.T) {
			p := NewParser(strings.NewReader(tc.spec))

			objects := []Object{}

			for p.Next() {
				obj := p.Object()
				if obj == nil {
					t.Errorf("unexpected nil object found")
					return
				}
				objects = append(objects, *obj)
			}

			if p.Err() != nil {
				t.Errorf("unexpected error: %v", p.Err())
				return
			}

			if diff := cmp.Diff(tc.objects, objects); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
