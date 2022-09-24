package loon

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestPrinter(t *testing.T) {
	for _, tc := range parserTests {
		t.Run("", func(t *testing.T) {
			got := string(Print(tc.doc))

			if diff := cmp.Diff(strings.TrimSpace(tc.spec), strings.TrimSpace(got)); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPrinterRoundtrip(t *testing.T) {
	for _, tc := range parserTests {
		t.Run("", func(t *testing.T) {
			formatted := bytes.TrimSpace(Print(tc.doc))

			got, err := Parse(formatted)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}

			if diff := cmp.Diff(tc.doc, got, cmpopts.IgnoreFields(Object{}, "Line"), cmpopts.IgnoreFields(Directive{}, "Line")); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
