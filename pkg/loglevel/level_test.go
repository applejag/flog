package loglevel

import (
	"fmt"
	"testing"
)

func TestShouldIncludeLogInOutput_WhitelistMask(t *testing.T) {
	var testCases = []struct {
		input Level
		want  string
	}{
		{
			input: Warning,
			want:  "Warning",
		},
		{
			input: Information | Panic,
			want:  "Information|Panic",
		},
		{
			input: Undefined,
			want:  "Undefined",
		},
		{
			input: Undefined | Debug,
			want:  "Debug",
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d/input/%v", i, tc.input), func(t *testing.T) {
			got := tc.input.String()
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
