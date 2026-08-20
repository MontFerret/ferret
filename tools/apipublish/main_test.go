package main

import (
	"strings"
	"testing"
)

func TestRunValidatesMigrationModeFlags(t *testing.T) {
	for _, test := range []struct {
		name      string
		arguments []string
		want      string
	}{
		{name: "missing pages", arguments: []string{"-migrate-types"}, want: "-pages is required"},
		{name: "publish inputs", arguments: []string{"-migrate-types", "-pages", "pages", "-reference", "api.json"}, want: "cannot be used"},
		{name: "check without migration", arguments: []string{"-check"}, want: "requires -migrate-types"},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := run(test.arguments)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("run error = %v, want %q", err, test.want)
			}
		})
	}
}
