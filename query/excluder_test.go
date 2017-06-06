package query

import "testing"

type ExcluderCase struct {
	input    string
	expected bool
}

func TestShouldExclude_ExpectAllExcluded(t *testing.T) {
	exclusions := []string{".git", ".gitignore"}
	excluder := regexpExclude{exclusions: exclusions}
	cases := []ExcluderCase{
		{".git", true},
		{".git/", true},
		{".git/some/other/file", true},
		{".gitignore", true},
	}

	for _, c := range cases {
		actual := excluder.shouldExclude(c.input)
		if actual != c.expected {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected, actual)
		}
	}
}

func TestShouldExclude_ExpectNotExcluded(t *testing.T) {
	exclusions := []string{".git"}
	excluder := regexpExclude{exclusions: exclusions}
	cases := []ExcluderCase{
		{".git", true},
		{".git/", true},
		{".git/some/other/file", true},
		{".gitignore", false},
	}

	for _, c := range cases {
		actual := excluder.shouldExclude(c.input)
		if actual != c.expected {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected, actual)
		}
	}
}
