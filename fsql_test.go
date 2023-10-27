package fsql

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
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

func TestRun_All(t *testing.T) {
	type Case struct {
		query    string
		expected string
	}

	cases := []Case{
		{
			query: "SELECT all FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\t-------\tfoo\n",
				strings.Join(GetAttrs("foo", "mode", "size", "time"), "\t")),
		},
		{
			query: "SELECT all FROM ./testdata WHERE name LIKE waldo",
			expected: fmt.Sprintf("%s\twaldo\n",
				strings.Join(GetAttrs("foo/quuz/waldo", "mode", "size", "time", "hash"), "\t")),
		},
		{
			query:    "SELECT all FROM ./testdata WHERE FORMAT(time, 'Jan 02 2006 15:04') > 'Jan 01 2999 00:00'",
			expected: "",
		},
		{
			query: "SELECT all FROM ./testdata WHERE mode IS DIR",
			expected: fmt.Sprintf(
				strings.Repeat("%s\n", 8),
				fmt.Sprintf(
					"%s\t-------\t%-8s",
					strings.Join(GetAttrs(".", "mode", "size", "time"), "\t"),
					"testdata",
				),
				fmt.Sprintf(
					"%s\t-------\t%-8s",
					strings.Join(GetAttrs("bar", "mode", "size", "time"), "\t"),
					"bar",
				),
				fmt.Sprintf(
					"%s\t-------\t%-8s",
					strings.Join(GetAttrs("bar/garply", "mode", "size", "time"), "\t"),
					"garply",
				),
				fmt.Sprintf(
					"%s\t-------\t%-8s",
					strings.Join(GetAttrs("bar/garply/xyzzy", "mode", "size", "time"), "\t"),
					"xyzzy",
				),
				fmt.Sprintf(
					"%s\t-------\t%-8s",
					strings.Join(GetAttrs("bar/garply/xyzzy/thud", "mode", "size", "time"), "\t"),
					"thud",
				),
				fmt.Sprintf(
					"%s\t-------\t%-8s",
					strings.Join(GetAttrs("foo", "mode", "size", "time"), "\t"),
					"foo",
				),
				fmt.Sprintf(
					"%s\t-------\t%-8s",
					strings.Join(GetAttrs("foo/quuz", "mode", "size", "time"), "\t"),
					"quuz",
				),
				fmt.Sprintf(
					"%s\t-------\t%-8s",
					strings.Join(GetAttrs("foo/quuz/fred", "mode", "size", "time"), "\t"),
					"fred",
				),
			),
		},
	}

	for _, c := range cases {
		actual := DoRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("%s\nExpected:\n%v\nGot:\n%v", c.query, c.expected, actual)
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
			expected: fmt.Sprintf("foo\t%s\n", GetAttrs("foo", "size")[0]),
		},
		{
			query:    "SELECT size, name FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\tfoo\n", GetAttrs("foo", "size")[0]),
		},
		{
			query: "SELECT FULLPATH(name), size, time FROM ./testdata/foo",
			expected: fmt.Sprintf(
				strings.Repeat("%s\n", 7),
				fmt.Sprintf(
					"%-31s\t%s",
					"testdata/foo",
					strings.Join(GetAttrs("foo", "size", "time"), "\t"),
				),
				fmt.Sprintf(
					"%-31s\t%s",
					"testdata/foo/quux",
					strings.Join(GetAttrs("foo/quux", "size", "time"), "\t"),
				),
				fmt.Sprintf(
					"%-31s\t%s",
					"testdata/foo/quuz",
					strings.Join(GetAttrs("foo/quuz", "size", "time"), "\t"),
				),
				fmt.Sprintf(
					"%-31s\t%s",
					"testdata/foo/quuz/fred",
					strings.Join(GetAttrs("foo/quuz/fred", "size", "time"), "\t"),
				),
				fmt.Sprintf(
					"%-31s\t%s",
					"testdata/foo/quuz/fred/.gitkeep",
					strings.Join(GetAttrs("foo/quuz/fred/.gitkeep", "size", "time"), "\t"),
				),
				fmt.Sprintf(
					"%-31s\t%s",
					"testdata/foo/quuz/waldo",
					strings.Join(GetAttrs("foo/quuz/waldo", "size", "time"), "\t"),
				),
				fmt.Sprintf(
					"%-31s\t%s",
					"testdata/foo/qux",
					strings.Join(GetAttrs("foo/qux", "size", "time"), "\t"),
				),
			),
		},
	}

	for _, c := range cases {
		actual := DoRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("%s\nExpected:\n%v\nGot:\n%v", c.query, c.expected, actual)
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
			expected: fmt.Sprintf(
				strings.Repeat("%-30s\n", 8),
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
		actual := DoRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("%s\nExpected:\n%v\nGot:\n%v", c.query, c.expected, actual)
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
			expected: fmt.Sprintf("%s\n", GetAttrs("foo", "size")[0]),
		},
		{
			query:    "SELECT FORMAT(size, KB) FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\n", GetAttrs("foo", "size:kb")[0]),
		},
		{
			query:    "SELECT FORMAT(size, MB) FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\n", GetAttrs("foo", "size:mb")[0]),
		},
		{
			query:    "SELECT FORMAT(size, GB) FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\n", GetAttrs("foo", "size:gb")[0]),
		},
		{
			query: "SELECT size FROM ./testdata WHERE name LIKE qu",
			expected: fmt.Sprintf(
				strings.Repeat("%s\n", 3),
				GetAttrs("foo/quux", "size")[0],
				GetAttrs("foo/quuz", "size")[0],
				GetAttrs("foo/qux", "size")[0],
			),
		},
	}

	for _, c := range cases {
		actual := DoRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("%s\nExpected:\n%v\nGot:\n%v", c.query, c.expected, actual)
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
			expected: fmt.Sprintf("%s\n", GetAttrs("baz", "time")[0]),
		},
		{
			query:    "SELECT FORMAT(time, ISO) FROM ./testdata WHERE name = foo",
			expected: fmt.Sprintf("%s\n", GetAttrs("foo", "time:iso")[0]),
		},
		{
			query:    "SELECT FORMAT(time, 2006) FROM ./testdata WHERE NOT name LIKE .%",
			expected: strings.Repeat(fmt.Sprintf("%s\n", GetAttrs(".", "time:year")[0]), 14),
		},
		{
			query: "SELECT time FROM ./testdata/foo/quuz",
			expected: fmt.Sprintf(
				strings.Repeat("%s\n", 4),
				GetAttrs("foo/quuz", "time")[0],
				GetAttrs("foo/quuz/fred", "time")[0],
				GetAttrs("foo/quuz/fred/.gitkeep", "time")[0],
				GetAttrs("foo/quuz/waldo", "time")[0],
			),
		},
	}

	for _, c := range cases {
		actual := DoRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("%s\nExpected:\n%v\nGot:\n%v", c.query, c.expected, actual)
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
			expected: fmt.Sprintf("%s\n", GetAttrs("foo", "mode")[0]),
		},
		{
			query:    "SELECT mode FROM ./testdata WHERE name = baz",
			expected: fmt.Sprintf("%s\n", GetAttrs("baz", "mode")[0]),
		},
		{
			query: "SELECT mode FROM ./testdata WHERE mode IS DIR",
			expected: fmt.Sprintf(strings.Repeat("%s\n", 8),
				GetAttrs("", "mode")[0],
				GetAttrs("foo", "mode")[0],
				GetAttrs("foo/quuz", "mode")[0],
				GetAttrs("foo/quuz/fred", "mode")[0],
				GetAttrs("bar", "mode")[0],
				GetAttrs("bar/garply", "mode")[0],
				GetAttrs("bar/garply/xyzzy", "mode")[0],
				GetAttrs("bar/garply/xyzzy/thud", "mode")[0],
			),
		},
	}

	for _, c := range cases {
		actual := DoRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("%s\nExpected:\n\"%v\"\nGot:\n\"%v\"", c.query, c.expected, actual)
		}
	}
}

func TestRun_Hash(t *testing.T) {
	type Case struct {
		query    string
		expected string
	}

	// TODO
	cases := []Case{}

	for _, c := range cases {
		actual := DoRun(c.query)
		if !reflect.DeepEqual(c.expected, actual) {
			t.Fatalf("%s\nExpected:\n%v\nGot:\n%v", c.query, c.expected, actual)
		}
	}
}

func GetAttrs(path string, attrs ...string) []string {
	// If the files map is empty, walk ./testdata and populate it.
	if len(files) == 0 {
		if err := filepath.Walk(
			"./testdata",
			func(path string, info os.FileInfo, err error) error {
				files[filepath.Clean(path)] = &info
				return nil
			},
		); err != nil {
			return []string{}
		}
	}

	path = filepath.Clean(fmt.Sprintf("testdata/%s", path))
	file, ok := files[path]
	if !ok {
		return []string{}
	}

	result := make([]string, len(attrs))
	for i, attr := range attrs {
		// Hard-coding modifiers works for the time being, but we might need a more
		// elegant solution when we introduce new modifiers in the future.
		switch attr {
		case "hash":
			b, err := os.ReadFile(path)
			if err != nil {
				return []string{}
			}
			h := sha1.New()
			if _, err := h.Write(b); err != nil {
				return []string{}
			}
			result[i] = hex.EncodeToString(h.Sum(nil))[:7]
		case "mode":
			result[i] = (*file).Mode().String()
		case "size":
			result[i] = fmt.Sprintf("%d", (*file).Size())
		case "size:kb", "size:mb", "size:gb":
			size := (*file).Size()
			switch attr[len(attr)-2:] {
			case "kb":
				result[i] = fmt.Sprintf("%fkb", float64(size)/(1<<10))
			case "mb":
				result[i] = fmt.Sprintf("%fmb", float64(size)/(1<<20))
			case "gb":
				result[i] = fmt.Sprintf("%fgb", float64(size)/(1<<30))
			}
		case "time":
			result[i] = (*file).ModTime().Format(time.Stamp)
		case "time:iso":
			result[i] = (*file).ModTime().Format(time.RFC3339)
		case "time:year":
			result[i] = (*file).ModTime().Format("2006")
		}
	}
	return result
}

// DoRun executes fsql.Run and returns the output.
func DoRun(query string) string {
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
