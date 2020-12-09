package main

import (
	"strings"
	"testing"
)

type testCase struct {
	name     string
	args     []string
	input    string
	exitCode int
}

var prog = "shamir"

func TestRun(t *testing.T) {
	testCases := []testCase{
		{
			name:     "no args",
			args:     []string{prog},
			exitCode: 1,
		},
		{
			name:     "invalid command foo",
			args:     []string{prog, "foo"},
			exitCode: 1,
		},
		{
			name:     "split with no args",
			args:     []string{prog, "split"},
			exitCode: 1,
		},
		{
			name:     "combine with no args",
			args:     []string{prog, "combine"},
			exitCode: 1,
		},
		{
			name:     "split with invalid threshold and parts",
			args:     []string{prog, "split", "--threshold=4", "--parts=3"},
			exitCode: 1,
		},
		{
			name:     "valid split",
			args:     []string{prog, "split", "--parts=5", "--threshold=2"},
			input:    "foobar\n",
			exitCode: 0,
		},
		{
			name:     "valid combine",
			args:     []string{prog, "combine", "--parts=2"},
			input:    "pEjEriAxfkdtkgU=\nM7fP0GdOU+dx4Aw=\n",
			exitCode: 0,
		},
		{
			name:     "valid combine with retries",
			args:     []string{prog, "combine", "--parts=2"},
			input:    "PuHK1xUwLTsjyM=\nePuHK1xUwLTsjyM=\nM7fP0GdOU+\nM7fP0GdOU+dx4Aw=\n",
			exitCode: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			exitCode := run(tc.args, strings.NewReader(tc.input))
			if exitCode != tc.exitCode {
				t.Fatalf("expected: exitCode=%d, got: exitCode=%d", tc.exitCode, exitCode)
			}
		})
	}
}
