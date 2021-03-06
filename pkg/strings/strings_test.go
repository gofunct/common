package strings

import (
	"testing"
)

func TestSplitQualifiedName(t *testing.T) {
	testCases := []struct {
		input  string
		output []string
	}{
		{"kubernetes.io/blah", []string{"kubernetes.io", "blah"}},
		{"blah", []string{"", "blah"}},
		{"kubernetes.io/blah/blah", []string{"kubernetes.io", "blah"}},
	}
	for i, tc := range testCases {
		namespace, name := SplitQualifiedName(tc.input)
		if namespace != tc.output[0] || name != tc.output[1] {
			t.Errorf("case[%d]: expected (%q, %q), got (%q, %q)", i, tc.output[0], tc.output[1], namespace, name)
		}
	}
}

func TestJoinQualifiedName(t *testing.T) {
	testCases := []struct {
		input  []string
		output string
	}{
		{[]string{"kubernetes.io", "blah"}, "kubernetes.io/blah"},
		{[]string{"blah", ""}, "blah"},
		{[]string{"kubernetes.io", "blah"}, "kubernetes.io/blah"},
	}
	for i, tc := range testCases {
		res := JoinQualifiedName(tc.input[0], tc.input[1])
		if res != tc.output {
			t.Errorf("case[%d]: expected %q, got %q", i, tc.output, res)
		}
	}
}

func TestShortenString(t *testing.T) {
	testCases := []struct {
		input  string
		outLen int
		output string
	}{
		{"kubernetes.io", 5, "kuber"},
		{"blah", 34, "blah"},
		{"kubernetes.io", 13, "kubernetes.io"},
	}
	for i, tc := range testCases {
		res := ShortenString(tc.input, tc.outLen)
		if res != tc.output {
			t.Errorf("case[%d]: expected %q, got %q", i, tc.output, res)
		}
	}
}
