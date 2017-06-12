package prompt

import "testing"

func TestParseLine_ReturnsCorrectValue(t *testing.T) {
	type Case struct {
		line     string
		expected bool
	}

	cases := []Case{
		{line: "select all from .", expected: false},
		{line: "where", expected: false},
		{line: "name like %go", expected: false},
		{line: "select all from where name like %go;", expected: true},
		{line: ";", expected: true},
	}

	for _, c := range cases {
		actual := parseLine([]byte(c.line))
		if c.expected != actual {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected, actual)
		}
	}
}

func TestParseLine_BuffersCorrectQuery(t *testing.T) {
	type Case struct {
		lines    []string
		expected string
	}

	// Note the whitespaces preceding output queries. This happens when the
	// semicolon appears on the next line. This is fine for anything EXCEPT
	// quoted strings (see the last case), since the tokenizer ignores
	// excessive whitespace. Will eventually need to address this, for the time
	// being, we'll need a note in the README highlighting that quoted strings
	// should **not** be spread across multiple lines.
	cases := []Case{
		{
			lines:    []string{"SELECT all", "FROM .;"},
			expected: "SELECT all FROM .",
		},
		{
			lines:    []string{"SELECT all", "FROM .", ";"},
			expected: "SELECT all FROM . ",
		},
		{
			lines:    []string{"SELECT all FROM . WHERE name IN (", "SELECT name FROM .", ");"},
			expected: "SELECT all FROM . WHERE name IN ( SELECT name FROM . )",
		},
		{
			lines:    []string{"SELECT all FROM . WHERE name IN [", "foo, bar, baz", "]", ";"},
			expected: "SELECT all FROM . WHERE name IN [ foo, bar, baz ] ",
		},
		{
			lines:    []string{"SELECT all FROM . WHERE name = \"name with ", "spaces\";"},
			expected: "SELECT all FROM . WHERE name = \"name with  spaces\"",
		},
	}

	for _, c := range cases {
		query.Reset()
		for _, l := range c.lines {
			parseLine([]byte(l))
		}

		actual := query.String()
		if c.expected != actual {
			t.Fatalf("\nExpected %v\n     Got %v", c.expected, actual)
		}
	}
}

func TestReadInput_HasCorrectStatus(t *testing.T) {
	// TODO: Complete this.
	//
	// I think what we can do is invoke the prompt in a pipe and then write some
	// string to stdin from another pipe. The value written to the first pipe
	// should then be either `>>> ` or `... ` depending on if what we wrote ended
	// with a semicolon (I'm not sure how viable this solution is).
}
