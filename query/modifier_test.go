package query

import "testing"

func TestModifierString(t *testing.T) {
	type Case struct {
		input    Modifier
		expected string
	}

	cases := []Case{
		{Modifier{"upper", []string{}}, "upper()"},
		{Modifier{"format", []string{"upper"}}, "format(upper)"},
	}

	for _, c := range cases {
		result := c.input.String()
		if result != c.expected {
			t.Fatalf("\nExpected: %s\n     Got: %s", c.expected, result)
		}
	}
}

func TestApplyModifiers(t *testing.T) {
	// TODO
}
