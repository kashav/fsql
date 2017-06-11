package query

import "testing"

func TestShouldExclude_ExpectAllExcluded(t *testing.T) {
	type Case struct {
		input    string
		expected bool
	}

	exclusions := []string{".git", ".gitignore"}
	excluder := regexpExclude{exclusions: exclusions}

	cases := []Case{
		Case{input: ".git", expected: true},
		Case{input: ".git/", expected: true},
		Case{input: ".git/some/other/file", expected: true},
		Case{input: ".gitignore", expected: true},
	}

	for _, c := range cases {
		actual := excluder.shouldExclude(c.input)
		if actual != c.expected {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected, actual)
		}
	}
}

func TestShouldExclude_ExpectNotExcluded(t *testing.T) {
	type Case struct {
		input    string
		expected bool
	}

	exclusions := []string{".git"}
	excluder := regexpExclude{exclusions: exclusions}

	cases := []Case{
		Case{input: ".git", expected: true},
		Case{input: ".git/", expected: true},
		Case{input: ".git/some/other/file", expected: true},
		Case{input: ".gitignore", expected: false},
	}

	for _, c := range cases {
		actual := excluder.shouldExclude(c.input)
		if actual != c.expected {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected, actual)
		}
	}
}
