package main

import (
	"reflect"
	"testing"
)

func TestOrgNodeFromString(t *testing.T) {
	type TestCase struct {
		inputVal    string
		expectedVal OrgNode
	}

	testCases := []TestCase{
		{
			"N01.",
			OrgNode{[]int{1}},
		},
		{
			"N01.01.",
			OrgNode{[]int{1, 1}},
		},
		{
			"N01.5",
			OrgNode{[]int{1, 5}},
		},
		{
			"N02.3.5.10.11.20",
			OrgNode{[]int{2, 3, 5, 10, 11, 20}},
		},
		{
			"N01.0.0.3",
			OrgNode{[]int{1, 0, 0, 3}},
		},
	}

	for _, tc := range testCases {
		result, _ := OrgNodeFromString(tc.inputVal, ".")
		if !reflect.DeepEqual(result, tc.expectedVal) {
			t.Fatalf("Expected %v but got %v\n", tc.expectedVal, result)
		}
	}
}
