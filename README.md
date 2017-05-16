# fsql

>Search through your file system with SQL-esque queries.

- [Demo](#demo)
- [Installation](#setup-installation)
- [Usage](#usage)
  + [Query](#query-syntax)
    * [Attribute](#attribute)
    * [Source](#source)
    * [Conditon](#condition)
  + [Examples](#examples-3)
- [Contribute](#contribute)
- [Credits](#credits)
- [License](#license)

## Demo

[![fsql.gif](./fsql.gif)](https://asciinema.org/a/120534)

## Setup / installation

Requires Go to be [installed](https://golang.org/doc/install) and [configured](https://golang.org/doc/install#testing).

Install with `go get`:

```console
$ go get -v github.com/kshvmdn/fsql/...
$ which fsql
$GOPATH/bin/fsql
```

Or, install directly via source:

```console
$ git clone https://github.com/kshvmdn/fsql.git $GOPATH/src/github.com/kshvmdn/fsql
$ cd $_ # $GOPATH/src/github.com/kshvmdn/fsql
$ make install && make
$ ./fsql
```

## Usage

fsql expects the query as a command line argument.

### Query syntax

In general, each query requires a `SELECT` clause (to specify which attributes should be shown), a `FROM` clause (to specify the directories to search in), and a `WHERE` clause (to specify conditions for the files).

```sql
SELECT attribute, ... FROM source, ... WHERE conditon
```

You may omit the `SELECT` clause, as well as the `WHERE` clause.

Quotes are **not** required, however you'll have to escape reserved characters (e.g. `<`, `>`, etc).

#### Attribute

Currently supported attributes include `name`, `size`, `mode`, `time`, or `all` / `*`.

If no attribute is provided, `all` is chosen by default.

##### Examples

Each group features a set of equivalent clauses.

```sh
$ fsql SELECT name, size, time FROM ...
$ fsql name, size, time FROM ...
```

```sh
$ fsql SELECT all FROM ...
$ fsql all FROM ...
$ fsql FROM ...
```

#### Source

Each source should be a relative or absolute path to some directory on your machine. You can also use environment variables (e.g. `$GOPATH`) or `~` (for your home directory).

Use `-` to exclude a directory. For example, to exclude `.git`: `"... FROM ., -.git/ ..."`.

##### Examples

```sh
$ fsql ... FROM . WHERE ...
```

```sh
$ fsql ... FROM ~/Desktop WHERE ...
```

```sh
$ fsql ... FROM ~/Desktop, $GOPATH WHERE ...
```

#### Condition

##### Conjunction/Disjunction

Use `AND` / `OR` to join conditions. Note that precedence is assigned based on order of appearance.

This means `WHERE a AND b OR c` is **not** the same as `WHERE c OR b AND a`. Use parentheses to get around this behaviour, `WHERE a AND b OR c` **is** the same as `WHERE c OR (b AND a)`.

##### Negation

Use `NOT` to negate a condition. This keyword **must** precede the condition (e.g. `... WHERE NOT a ...`).

Note that wrapping parentheses with `NOT` is currently not supported. This can easily be resolved with [De Morgan's laws](https://en.wikipedia.org/wiki/De_Morgan%27s_laws). For example, `... WHERE NOT (a AND b) ...` is the same as `... WHERE NOT a OR NOT b ...`.

##### Condition Syntax

A single condition is made up of 3 parts: attribute, comparator, and value.

###### attribute

A valid attribute is any of the following: `name`, `size`, `file`, `time`.

###### comparator

Comparators depend on the attribute.

For `name`:

  - `=` - Strings that are an exact match.
  - `<>` - Synonymous to using `WHERE NOT ... = ...`.
  - `LIKE` - For simple pattern matching. Use `%` to match zero, one, or multiple characters. Check that a string begins with a value: `<value>%`, ends with a value: `%<value>`, or contains a value: `<value>`.
  - `RLIKE` - For pattern matching with regular expressions.

For `size` and `time`:

  - `>`
  - `>=`
  - `<`
  - `<=`
  - `=`
  - `<>`

And, for `file`:

  - `IS`

###### value

If the value contains spaces and/or escaped characters, wrap the value in quotes (either single or double) or backticks.

The default unit for `size` is bytes. To use kilobytes / megabytes / gigabytes, append `kb` / `mb` / `gb` to the size value (e.g. `100kb` for 100 kilobytes).

Attribute `file` only has 2 supported values: `dir` (to check that the file is a directory) and `reg` (to check that the file is regular).

Use the following format for `time` values: `MMM DD YYYY HH MM` (eg. `Jan 02 2006 15 04`).

##### Examples

See the next section for examples.

### Examples

List all files & directories in Desktop and Downloads that contain `csc` in the name:

```console
$ fsql SELECT name FROM ~/Desktop, ~/Downloads WHERE name LIKE %csc%
$ # this is equivalent to:
$ fsql SELECT name FROM ~/Desktop, ~/Downloads WHERE name RLIKE .*csc.*
```

List all JavaScript files in the current directory that were modified after April 1st 2017 (try running this on a `node_modules` directory, it's fast :sunglasses:).

```sh
$ fsql name, size, time FROM . WHERE name LIKE %.js AND time > 'Apr 01 2017 00 00'
```

List all files named `main.go` in `$GOPATH` which are larger than 10.5 kilobytes or smaller than 100 bytes (note the escaped redirection symbols, to avoid this wrap the query in quotes).

```sh
$ fsql FROM $GOPATH WHERE name = main.go AND (size \>= 10.5kb OR size \< 100)
$ fsql "FROM $GOPATH WHERE name = main.go AND (size >= 10.5kb OR size < 100)"
```

## Contribute

This project is completely open source, feel free to [open an issue](https://github.com/kshvmdn/fsql/issues) or [submit a pull request](https://github.com/kshvmdn/fsql/pulls).

Before submitting code, please ensure your changes comply with [Golint](https://github.com/golang/lint). Use `make lint` to test this.

## Credits

Lexer & parser are based on the work of [JamesOwenHall](https://github.com/JamesOwenHall) ([json2](https://github.com/JamesOwenHall/json2), [timed](https://github.com/JamesOwenHall/timed)).

## License

fsql source code is available under the [MIT license](./LICENSE).
