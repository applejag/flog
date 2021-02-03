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
		loglevel.Debug: 1,
	}
	want := 0

	got := getSkippedLevelsSlice(input)

	if len(got) != want {
		t.Errorf("slice was not empty: %#v", got)
	}
}
