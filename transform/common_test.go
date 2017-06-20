package transform

import (
	"crypto/sha1"
	"hash"
	"os"
	"reflect"
	"strings"
	"testing"
)

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

func TestCommon_Truncate(t *testing.T) {
	input := "foo-bar-baz"

	type Case struct {
		n        int
		expected string
	}

	cases := []Case{
		{n: 3, expected: "foo"},
		{n: 7, expected: "foo-bar"},
		{n: 100, expected: input},
		{n: -1, expected: input},
		{n: len(input), expected: "foo-bar-baz"},
		{n: len(input) - 1, expected: "foo-bar-ba"},
	}

	for _, c := range cases {
		actual := truncate(input, c.n)
		if c.expected != actual {
			t.Fatalf("\nExpected: %s\n     Got: %s", c.expected, actual)
		}
	}
}

func TestCommon_FindHash(t *testing.T) {
	type Case struct {
		name     string
		expected hash.Hash
	}

	cases := []Case{
		{name: "SHA1", expected: sha1.New()},
		{name: "FOO", expected: nil},
	}

	for _, c := range cases {
		actual := FindHash(c.name)
		if actual == nil {
			if c.expected != nil {
				t.Fatalf("\nExpected: %s\n     Got: nil", c.expected)
			}
		} else if h := actual(); !reflect.DeepEqual(c.expected, h) {
			t.Fatalf("\nExpected: %s\n     Got: %s", c.expected, h)
		}
	}
}

func TestCommon_ComputeHash(t *testing.T) {
	type Case struct {
		path     string
		expected string
	}

	cases := []Case{
		{path: "../testdata/foo", expected: strings.Repeat("-", 40)},
		{path: "../testdata/baz", expected: "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
	}

	for _, c := range cases {
		info, err := os.Stat(c.path)
		if err != nil {
			t.Fatalf("\nExpected no error\n     Got: %s", err.Error())
		}
		actual, err := ComputeHash(info, c.path, sha1.New())
		if err != nil {
			t.Fatalf("\nExpected no error\n     Got: %s", err.Error())
		}
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("\nExpected: %s\n     Got: %s", c.expected, actual)
		}
	}
}
