package gorm

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

type leak string

func (l leak) String() string { panic("yo") }

type printableString string

func (s printableString) String() string { return "custom:" + string(s) }

func TestLogFormatter(t *testing.T) {
	cases := []struct {
		name  string
		value interface{}
		sql   string
	}{
		{"panicking Stringer preserves underlying string", leak("yo"), "SELECT 'yo'"},
		{"ordinary string", "yo", "SELECT 'yo'"},
		{"working Stringer", printableString("yo"), "SELECT 'custom:yo'"},
		{"number", 42, "SELECT 42"},
		{"nil", nil, "SELECT NULL"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := LogFormatter("sql", "formatter.go:1", 123*time.Millisecond, "SELECT ?", []interface{}{tc.value}, int64(12))
			if len(got) != 5 {
				t.Fatalf("log fields = %d, want 5: %v", len(got), got)
			}
			if got[3] != tc.sql {
				t.Errorf("SQL = %q, want %q", got[3], tc.sql)
			}
			for _, expected := range []struct {
				index int
				text  string
			}{{0, "formatter.go:1"}, {2, "123.00ms"}, {4, "12 rows affected or returned"}} {
				if !strings.Contains(fmt.Sprint(got[expected.index]), expected.text) {
					t.Errorf("field %d = %v, want %q", expected.index, got[expected.index], expected.text)
				}
			}
		})
	}
}
