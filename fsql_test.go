package fsql

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

// doRun executes fsql.Run and returns the output.
func doRun(query string) string {
	stdout := os.Stdout
	ch := make(chan string)

	r, w, err := os.Pipe()
	if err != nil {
		return ""
	}
	os.Stdout = w

	if err := Run(query); err != nil {
		return ""
	}

	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		ch <- buf.String()
	}()

	w.Close()
	os.Stdout = stdout
	return <-ch
}

func TestRun_All(t *testing.T) {
	type Case struct {
		query    string
		expected string
	}

	cases := []Case{
		{
			query:    "SELECT all FROM ./testdata WHERE name = foo",
			expected: "drwxr-xr-x\t170\tJun  4 19:24:12\tfoo\n",
		},
		{
			query:    "SELECT all FROM ./testdata WHERE name LIKE qu AND size > 0",
			expected: "drwxr-xr-x\t136\tJun  4 19:25:02\tquuz\n",
		},
		{
			query:    "SELECT all FROM ./testdata WHERE FORMAT(time, 'Jan 02 2006 15:04') > 'Jun 04 2017 23:26' AND NOT name LIKE .%",
			expected: "drwxr-xr-x\t102\tJun  4 19:26:31\tthud\ndrwxr-xr-x\t102\tJun  4 19:26:31\tfred\n",
		},
		{
			query: "SELECT all FROM ./testdata WHERE mode IS DIR",
			expected: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"drwxr-xr-x\t204\tJun  4 19:23:29\ttestdata",
				"drwxr-xr-x\t204\tJun  4 19:25:29\tbar",
				"drwxr-xr-x\t102\tJun  4 19:25:24\tgarply",
				"drwxr-xr-x\t102\tJun  4 19:25:19\txyzzy",
				"drwxr-xr-x\t102\tJun  4 19:26:31\tthud",
				"drwxr-xr-x\t170\tJun  4 19:24:12\tfoo",
				"drwxr-xr-x\t136\tJun  4 19:25:02\tquuz",
				"drwxr-xr-x\t102\tJun  4 19:26:31\tfred",
			),
		},
	}

	for _, c := range cases {
		actual := doRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("\nExpected:\n%v\nGot:\n%v", c.expected, actual)
		}
	}
}

func TestRun_Multiple(t *testing.T) {
	type Case struct {
		query    string
		expected string
	}

	cases := []Case{
		{
			query:    "SELECT name, size FROM ./testdata WHERE name = foo",
			expected: "170\tfoo\n",
		},
		{
			query:    "SELECT size, name FROM ./testdata WHERE name = foo",
			expected: "170\tfoo\n",
		},
		{
			query: "SELECT size, time, FULLPATH(name) FROM ./testdata/foo",
			expected: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"170\tJun  4 19:24:12\ttestdata/foo",
				"0\tJun  4 19:22:01\ttestdata/foo/quux",
				"136\tJun  4 19:25:02\ttestdata/foo/quuz",
				"102\tJun  4 19:26:31\ttestdata/foo/quuz/fred",
				"0\tJun  4 19:26:31\ttestdata/foo/quuz/fred/.gitkeep",
				"0\tJun  4 19:24:56\ttestdata/foo/quuz/waldo",
				"0\tJun  4 19:21:49\ttestdata/foo/qux",
			),
		},
	}

	for _, c := range cases {
		actual := doRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("\nExpected:\n%v\nGot:\n%v", c.expected, actual)
		}
	}
}

func TestRun_Name(t *testing.T) {
	type Case struct {
		query    string
		expected string
	}

	cases := []Case{
		{
			query:    "SELECT name FROM ./testdata WHERE name REGEXP ^g.*",
			expected: "garply\ngrault\n",
		},
		{
			query:    "SELECT FULLPATH(name) FROM ./testdata WHERE name REGEXP ^b.*",
			expected: "testdata/bar\ntestdata/baz\n",
		},
		{
			query: "SELECT UPPER(FULLPATH(name)) FROM ./testdata WHERE mode IS DIR",
			expected: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				"TESTDATA",
				"TESTDATA/BAR",
				"TESTDATA/BAR/GARPLY",
				"TESTDATA/BAR/GARPLY/XYZZY",
				"TESTDATA/BAR/GARPLY/XYZZY/THUD",
				"TESTDATA/FOO",
				"TESTDATA/FOO/QUUZ",
				"TESTDATA/FOO/QUUZ/FRED",
			),
		},
	}

	for _, c := range cases {
		actual := doRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("\nExpected:\n%v\nGot:\n%v", c.expected, actual)
		}
	}
}

func TestRun_Size(t *testing.T) {
	type Case struct {
		query    string
		expected string
	}

	cases := []Case{
		{
			query:    "SELECT size FROM ./testdata WHERE name = foo",
			expected: "170\n",
		},
		{
			query:    "SELECT FORMAT(size, KB) FROM ./testdata WHERE name = baz",
			expected: "0.000000kb\n",
		},
		{
			query:    "SELECT FORMAT(size, MB) FROM ./testdata WHERE name = baz",
			expected: "0.000000mb\n",
		},
		{
			query:    "SELECT FORMAT(size, GB) FROM ./testdata WHERE name = baz",
			expected: "0.000000gb\n",
		},
		{
			query:    "SELECT size FROM ./testdata WHERE name LIKE qu",
			expected: "0\n136\n0\n",
		},
	}

	for _, c := range cases {
		actual := doRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("\nExpected:\n%v\nGot:\n%v", c.expected, actual)
		}
	}
}

func TestRun_Time(t *testing.T) {
	type Case struct {
		query    string
		expected string
	}

	cases := []Case{
		{
			query:    "SELECT time FROM ./testdata WHERE name = baz",
			expected: "Jun  4 19:23:29\n",
		},
		{
			query:    "SELECT FORMAT(time, ISO) FROM ./testdata WHERE name = foo",
			expected: "2017-06-04T19:24:12-04:00\n",
		},
		{
			query:    "SELECT FORMAT(time, 2006) FROM ./testdata WHERE NOT name LIKE .%",
			expected: strings.Repeat("2017\n", 14),
		},
		{
			query:    "SELECT time FROM ./testdata/foo/quuz",
			expected: "Jun  4 19:25:02\nJun  4 19:26:31\nJun  4 19:26:31\nJun  4 19:24:56\n",
		},
	}

	for _, c := range cases {
		actual := doRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("\nExpected:\n%v\nGot:\n%v", c.expected, actual)
		}
	}
}

func TestRun_Mode(t *testing.T) {
	type Case struct {
		query    string
		expected string
	}

	cases := []Case{
		{
			query:    "SELECT mode FROM ./testdata WHERE name = foo",
			expected: "drwxr-xr-x\n",
		},
		{
			query:    "SELECT mode FROM ./testdata WHERE name = baz",
			expected: "-rwxr-xr-x\n",
		},
		{
			query:    "SELECT mode FROM ./testdata WHERE mode IS DIR",
			expected: strings.Repeat("drwxr-xr-x\n", 8),
		},
	}

	for _, c := range cases {
		actual := doRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("\nExpected:\n%v\nGot:\n%v", c.expected, actual)
		}
	}
}

func TestRunInteractive(t *testing.T) {
	// TODO: Complete this.
	//
	// I'm not really sure how to test this. We're already testing the core of
	// this function above, so I guess this test should just ensure that it
	// retrieves input from prompt.Run and runs until the process is killed?
}
