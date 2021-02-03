package main

import (
	"testing"

	"github.com/jilleJr/flog/pkg/loglevel"
)

func TestGetSkippedLevelsSlice_Empty(t *testing.T) {
	input := map[loglevel.Level]int{}
	want := 0

	got := getSkippedLevelsSlice(input)

	if len(got) != want {
		t.Errorf("slice was not empty: %#v", got)
	}
}

func TestGetSkippedLEvelsSlice_AllInOrder(t *testing.T) {
	input := map[loglevel.Level]int{
		loglevel.Trace:       1,
		loglevel.Debug:       2,
		loglevel.Information: 3,
		loglevel.Warning:     4,
		loglevel.Error:       5,
		loglevel.Critical:    6,
		loglevel.Fatal:       7,
		loglevel.Panic:       8,
	}

	want := []string{
		"1 Trace",
		"2 Debug",
		"3 Information",
		"4 Warning",
		"5 Error",
		"6 Critical",
		"7 Fatal",
		"8 Panic",
	}

	got := getSkippedLevelsSlice(input)

	if len(got) != len(want) {
		t.Fatalf("slice was of wrong length: want %d, got %d: %#v", len(want), len(got), got)
	}

	for i, wantItem := range want {
		gotItem := got[i]
		if gotItem != wantItem {
			t.Errorf("slice was wrong at index %d: want %q, got %q", i, wantItem, gotItem)
		}
	}
}
