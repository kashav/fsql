package transform

import "testing"

func TestFormatName(t *testing.T) {
	type Case struct {
		arg      string
		name     string
		expected string
	}

	cases := []Case{
		{"upper", "foo", "FOO"},
		{"upper", "FOO", "FOO"},
		{"lower", "foo", "foo"},
		{"lower", "FOO", "foo"},
	}

	for _, c := range cases {
		result := formatName(c.arg, c.name)
		if result != c.expected {
			t.Fatalf("\nExpected: %s\n     Got: %s", c.expected, result)
		}
	}
}

func TestUpper(t *testing.T) {
	type Case struct {
		name     string
		expected string
	}

	cases := []Case{
		{"foo", "FOO"},
		{"FOO", "FOO"},
	}

	for _, c := range cases {
		result := upper(c.name)
		if result != c.expected {
			t.Fatalf("\nExpected: %s\n     Got: %s", c.expected, result)
		}
	}
}

func TestLower(t *testing.T) {
	type Case struct {
		name     string
		expected string
	}

	cases := []Case{
		{"foo", "foo"},
		{"FOO", "foo"},
	}

	for _, c := range cases {
		result := lower(c.name)
		if result != c.expected {
			t.Fatalf("\nExpected: %s\n     Got: %s", c.expected, result)
		}
	}
}
