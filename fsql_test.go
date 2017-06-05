package fsql

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

var files = map[string]*os.FileInfo{}

func getAttr(path string, attr string) string {
	// If the files map is empty, walk ./testdata and populate it!
	if len(files) == 0 {
		filepath.Walk("./testdata", func(path string, info os.FileInfo, err error) error {
			files[filepath.Clean(path)] = &info
			return nil
		})
	}

	file, ok := files[filepath.Clean(fmt.Sprintf("testdata/%s", path))]
	if !ok {
		return ""
	}

	// Hard-coding modifiers works for the time being, but we might need a more
	// elegant solution when we introduce new modifiers in the future.
	switch attr {
	case "time":
		return (*file).ModTime().Format(time.Stamp)
	case "time:iso":
		return (*file).ModTime().Format(time.RFC3339)
	case "time:year":
		return (*file).ModTime().Format("2006")

	case "size":
		return fmt.Sprintf("%d", (*file).Size())
	case "size:kb", "size:mb", "size:gb":
		size := (*file).Size()
		unit := attr[len(attr)-2:]
		switch unit {
		case "kb":
			return fmt.Sprintf("%fkb", float64(size)/(1<<10))
		case "mb":
			return fmt.Sprintf("%fmb", float64(size)/(1<<20))
		case "gb":
			return fmt.Sprintf("%fgb", float64(size)/(1<<30))
		}
	}
	return ""
}

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
			query: "SELECT all FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("drwxr-xr-x\t%s\t%s\tfoo\n",
				getAttr("foo", "size"), getAttr("foo", "time")),
		},
		{
			query: "SELECT all FROM ./testdata WHERE name LIKE qu AND size > 0",
			expected: fmt.Sprintf("drwxr-xr-x\t%s\t%s\tquuz\n",
				getAttr("foo/quuz", "size"), getAttr("foo/quuz", "time")),
		},
		{
			query:    "SELECT all FROM ./testdata WHERE FORMAT(time, 'Jan 02 2006 15:04') > 'Jan 01 2999 00:00'",
			expected: "",
		},
		{
			query: "SELECT all FROM ./testdata WHERE mode IS DIR",
			expected: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				fmt.Sprintf("drwxr-xr-x\t%s\t%s\ttestdata",
					getAttr(".", "size"), getAttr(".", "time")),
				fmt.Sprintf("drwxr-xr-x\t%s\t%s\tbar",
					getAttr("bar", "size"), getAttr("bar", "time")),
				fmt.Sprintf("drwxr-xr-x\t%s\t%s\tgarply",
					getAttr("bar/garply", "size"), getAttr("bar/garply", "time")),
				fmt.Sprintf("drwxr-xr-x\t%s\t%s\txyzzy",
					getAttr("bar/garply/xyzzy", "size"), getAttr("bar/garply/xyzzy", "time")),
				fmt.Sprintf("drwxr-xr-x\t%s\t%s\tthud",
					getAttr("bar/garply/xyzzy/thud", "size"), getAttr("bar/garply/xyzzy/thud", "time")),
				fmt.Sprintf("drwxr-xr-x\t%s\t%s\tfoo",
					getAttr("foo", "size"), getAttr("foo", "time")),
				fmt.Sprintf("drwxr-xr-x\t%s\t%s\tquuz",
					getAttr("foo/quuz", "size"), getAttr("foo/quuz", "time")),
				fmt.Sprintf("drwxr-xr-x\t%s\t%s\tfred",
					getAttr("foo/quuz/fred", "size"), getAttr("foo/quuz/fred", "time")),
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
			expected: fmt.Sprintf("%s\tfoo\n", getAttr("foo", "size")),
		},
		{
			query:    "SELECT size, name FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\tfoo\n", getAttr("foo", "size")),
		},
		{
			query: "SELECT size, time, FULLPATH(name) FROM ./testdata/foo",
			expected: fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n",
				fmt.Sprintf("%s\t%s\ttestdata/foo",
					getAttr("foo", "size"), getAttr("foo", "time")),
				fmt.Sprintf("%s\t%s\ttestdata/foo/quux",
					getAttr("foo/quux", "size"), getAttr("foo/quux", "time")),
				fmt.Sprintf("%s\t%s\ttestdata/foo/quuz",
					getAttr("foo/quuz", "size"), getAttr("foo/quuz", "time")),
				fmt.Sprintf("%s\t%s\ttestdata/foo/quuz/fred",
					getAttr("foo/quuz/fred", "size"), getAttr("foo/quuz/fred", "time")),
				fmt.Sprintf("%s\t%s\ttestdata/foo/quuz/fred/.gitkeep",
					getAttr("foo/quuz/fred/.gitkeep", "size"), getAttr("foo/quuz/fred/.gitkeep", "time")),
				fmt.Sprintf("%s\t%s\ttestdata/foo/quuz/waldo",
					getAttr("foo/quuz/waldo", "size"), getAttr("foo/quuz/waldo", "time")),
				fmt.Sprintf("%s\t%s\ttestdata/foo/qux",
					getAttr("foo/qux", "size"), getAttr("foo/qux", "time")),
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
			expected: fmt.Sprintf("%s\n", getAttr("foo", "size")),
		},
		{
			query:    "SELECT FORMAT(size, KB) FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\n", getAttr("foo", "size:kb")),
		},
		{
			query:    "SELECT FORMAT(size, MB) FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\n", getAttr("foo", "size:mb")),
		},
		{
			query:    "SELECT FORMAT(size, GB) FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\n", getAttr("foo", "size:gb")),
		},
		{
			query: "SELECT size FROM ./testdata WHERE name LIKE qu",
			expected: fmt.Sprintf("%s\n%s\n%s\n",
				getAttr("foo/quux", "size"),
				getAttr("foo/quuz", "size"),
				getAttr("foo/qux", "size")),
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
			expected: fmt.Sprintf("%s\n", getAttr("baz", "time")),
		},
		{
			query:    "SELECT FORMAT(time, ISO) FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\n", getAttr("foo", "time:iso")),
		},
		{
			query:    "SELECT FORMAT(time, 2006) FROM ./testdata WHERE NOT name LIKE .%",
			expected: strings.Repeat(fmt.Sprintf("%s\n", getAttr(".", "time:year")), 14),
		},
		{
			query: "SELECT time FROM ./testdata/foo/quuz",
			expected: fmt.Sprintf("%s\n%s\n%s\n%s\n",
				getAttr("foo/quuz", "time"),
				getAttr("foo/quuz/fred", "time"),
				getAttr("foo/quuz/fred/.gitkeep", "time"),
				getAttr("foo/quuz/waldo", "time")),
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
