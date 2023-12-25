package main

import (
	"log"
	"reflect"
	"testing"
)

func TestSliceUnique(t *testing.T) {
	type TestCase struct {
		input    []string
		expected bool
	}

	testCases := []TestCase{
		{
			[]string{"foo", "bar"},
			true,
		},
		{
			[]string{"foo", "foo"},
			false,
		},
		{
			[]string{"foo", "foo", "foo", "bar"},
			false,
		},
		{
			[]string{},
			true,
		},
		{
			[]string{"foo", "bar", "baz", "foobar", "foo"},
			false,
		},
		{
			[]string{"foo", "foo "},
			true,
		},
	}

	for i, tc := range testCases {
		got := SliceUnique(tc.input)
		if !reflect.DeepEqual(tc.expected, got) {
			log.Fatalf("(%d) Expected %v, got %v\n", i, tc.expected, got)
		}
	}
}
