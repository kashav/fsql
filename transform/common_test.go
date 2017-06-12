package transform

import "testing"

func TestCommon_FormatName(t *testing.T) {
	type Case struct {
		arg      string
		name     string
		expected string
	}

	cases := []Case{
		{arg: "upper", name: "foo", expected: "FOO"},
		{arg: "upper", name: "FOO", expected: "FOO"},
		{arg: "lower", name: "foo", expected: "foo"},
		{arg: "lower", name: "FOO", expected: "foo"},
	}

	for _, c := range cases {
		result := formatName(c.arg, c.name)
		if result != c.expected {
			t.Fatalf("\nExpected: %s\n     Got: %s", c.expected, result)
		}
	}
}

func TestCommon_Upper(t *testing.T) {
	type Case struct {
		name     string
		expected string
	}

	cases := []Case{
		{name: "foo", expected: "FOO"},
		{name: "FOO", expected: "FOO"},
	}

	for _, c := range cases {
		result := upper(c.name)
		if result != c.expected {
			t.Fatalf("\nExpected: %s\n     Got: %s", c.expected, result)
		}
	}
}

func TestCommon_Lower(t *testing.T) {
	type Case struct {
		name     string
		expected string
	}

	cases := []Case{
		{name: "foo", expected: "foo"},
		{name: "FOO", expected: "foo"},
	}

	for _, c := range cases {
		result := lower(c.name)
		if result != c.expected {
			t.Fatalf("\nExpected: %s\n     Got: %s", c.expected, result)
		}
	}
}
