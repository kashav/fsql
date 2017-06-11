# fsql [![Build Status](https://travis-ci.org/kshvmdn/fsql.svg?branch=master)](https://travis-ci.org/kshvmdn/fsql) [![Go Report Card](https://goreportcard.com/badge/github.com/kshvmdn/fsql)](https://goreportcard.com/report/github.com/kshvmdn/fsql)

>Search through your filesystem with SQL-esque queries.

## Contents

- [Demo](#demo)
- [Installation](#setup--installation)
- [Usage](#usage)
- [Query Syntax](#query-syntax)
- [Examples](#usage-examples)
- [Contribute](#contribute)
- [License](#license)

## Demo

[![fsql.gif](./media/fsql.gif)](https://asciinema.org/a/120534)

## Setup / installation

Assumes that Go is [installed](https://golang.org/doc/install) and [setup](https://golang.org/doc/install#testing).

Install with `go get`:

```sh
$ go get -u -v github.com/kshvmdn/fsql/...
$ which fsql
$GOPATH/bin/fsql
```

Install with [Homebrew](https://brew.sh/):

```console
$ brew install fsql
$ which fsql
/usr/local/bin/fsql
```

Install directly via source:

```sh
$ git clone https://github.com/kshvmdn/fsql.git $GOPATH/src/github.com/kshvmdn/fsql
$ cd $_ # $GOPATH/src/github.com/kshvmdn/fsql
$ make install && make
$ ./fsql
```

## Usage

fsql expects a single query via stdin.

You may also choose to use fsql in interactive mode. Note that this mode is currently a work-in-progress, so many common shell features are not implemented yet (e.g. arrow key support, autocomplete, history, etc).

View the usage dialogue with the `-help` flag.

```sh
$ fsql -help
usage: fsql [options] query
  -interactive
      run in interactive mode (Ctrl+D to exit)
  -version
      print version and exit
```

## Query syntax

In general, each query requires a `SELECT` clause (to specify which attributes will be shown), a `FROM` clause (to specify which directories to search), and a `WHERE` clause (to specify conditions to test against).

```console
>>> SELECT attribute, ... FROM source, ... WHERE condition;
```

You may choose to omit the `SELECT` and `WHERE` clause.

If you're providing your query via stdin, quotes are **not** required, however you'll have to escape _reserved_ characters (e.g. `*`, `<`, `>`, etc).

### Attribute

Currently supported attributes include `name`, `size`, `mode`, `time`, and `all`.

If no attribute is provided, `all` is chosen by default.

**Examples**:

Each group features a set of equivalent clauses.

```console
>>> SELECT name, size, time ...
>>> name, size, time ...
```

```console
>>> SELECT all FROM ...
>>> all FROM ...
>>> FROM ...
```

### Source

Each source should be a relative or absolute path to a directory on your machine.

Source paths may include environment variables (e.g. `$GOPATH`) or tildes (`~`). Use a hyphen (`-`) to exclude a directory. Source paths also support usage of [glob patterns](https://en.wikipedia.org/wiki/Glob_(programming)).

In the case that a directory begins with a hyphen (e.g. `-foo`), use the following to include it as a source:

```console
>>> ... FROM ./-foo ...
```

**Examples**:

```console
>>> ... FROM . ...
```

```console
>>> ... FROM ~/Desktop, ./*/**.go ...
```

```console
>>> ... FROM $GOPATH, -.git/ ...
```

### Condition

#### Condition syntax

A single condition is made up of 3 parts: an attribute, an operator, and a value.

- **Attribute**:

  A valid attribute is any of the following: `name`, `size`, `mode`, `time`.

- **Operator**:

  Each attribute has a set of associated operators.

  - `name`:

    | Operator | Description |
    | :---: | --- |
    | `=` | String equality |
    | `<>` / `!=` | Synonymous to using `"NOT ... = ..."` |
    | `IN` | Basic list inclusion |
    | `LIKE` |  Simple pattern matching. Use `%` to match zero, one, or multiple characters. Check that a string begins with a value: `<value>%`, ends with a value: `%<value>`, or contains a value: `%<value>%`. |
    | `RLIKE` | Pattern matching with regular expressions. |

  - `size` / `time`:

    - All basic algebraic operators: `>`, `>=`, `<`, `<=`, `=`, and `<>` / `!=`.

  - `mode`:

    - `IS`

- **Value**:

  If the value contains spaces, wrap the value in quotes (either single or double) or backticks.

  The default unit for `size` is bytes.

  The default format for `time` is `MMM DD YYYY HH MM` (e.g. `"Jan 02 2006 15 04"`).

  Use `mode` to test if a file is regular (`IS REG`) or if it's a directory (`IS DIR`).

#### Conjunction / Disjunction

Use `AND` / `OR` to join conditions. Note that precedence is assigned based on order of appearance.

This means `WHERE a AND b OR c` is **not** the same as `WHERE c OR b AND a`. Use parentheses to get around this behaviour, i.e. `WHERE a AND b OR c` **is** the same as `WHERE c OR (b AND a)`.

**Examples**:

```console
>>> ... WHERE name = main.go OR size = 5 ...
```

```console
>>> ... WHERE name = main.go AND size > 20 ...
```

#### Negation

Use `NOT` to negate a condition. This keyword **must** precede the condition (e.g. `... WHERE NOT a ...`).

Note that negating parenthesized conditions is currently not supported. However, this can easily be resolved by applying [De Morgan's laws](https://en.wikipedia.org/wiki/De_Morgan%27s_laws) to your query. For example, `... WHERE NOT (a AND b) ...` is _logically equivalent_ to `... WHERE NOT a OR NOT b ...` (the latter is actually more optimal, due to [lazy evaluation](https://en.wikipedia.org/wiki/Lazy_evaluation)).

**Examples**:

```console
>>> ... WHERE NOT name = main.go ...
```

### Attribute Modifiers

Attribute modifiers are used to specify how input and output values should be processed. These functions are applied directly to attributes in the `SELECT` and `WHERE` clauses.

The table below lists currently-supported modifiers. Note that the first parameter to `FORMAT` is always the attribute name.

| Attribute | Modifier  | Supported in `SELECT` | Supported in `WHERE` |
| :---: | --- | :---: | :---: |
| `name` | `UPPER` (synonymous to `FORMAT(, UPPER)` | ✔️ | ✔️ |
| | `LOWER` (synonymous to `FORMAT(, LOWER)` | ✔️ | ✔️ |
| | `FULLPATH` | ✔️ |  |
| | `SHORTPATH`  | ✔️ |  |
| `size` | `FORMAT(, unit)` | ✔️ | ✔️ |
| `time` | `FORMAT(, layout)` | ✔️ | ✔️ |


Supported `unit` values: `B` (byte), `KB` (kilobyte), `MB` (megabyte), or `GB` (gigabyte).

Supported `layout` values: [`ISO`](https://en.wikipedia.org/wiki/ISO_8601), [`UNIX`](https://en.wikipedia.org/wiki/Unix_time), or [custom](https://golang.org/pkg/time/#Time.Format). Custom layouts must be provided in reference to the following date: `Mon Jan 2 15:04:05 -0700 MST 2006`.

**Examples**:

```console
>>> ... WHERE UPPER(name) LIKE %.go ...
```

```console
>>> ... WHERE FORMAT(size, GB) > 2 ...
```

```console
>>> ... SELECT FORMAT(time, "Mon Jan 2 2006 15:04:05") ...
```

### Subqueries

Subqueries allow for more complex condition statements. These queries are recursively evaluated while parsing. SELECTing multiple attributes in a subquery is not currently supported; if more than one attribute (or `all`) is provided, only the first attribute is used.

Support for referencing superqueries is not yet implemented, see [#4](https://github.com/kshvmdn/fsql/issues/4) if you'd like to help with this.

**Examples**:

```console
>>> ... WHERE name IN (SELECT name FROM ../foo) ...
```

## Usage Examples

List all attributes of each directory in your home directory (note the escaped `*`):

```console
$ fsql SELECT \* FROM ~ WHERE mode IS DIR
```

List the names of all files in the Desktop and Downloads directory that contain `csc` in the name:

```console
$ fsql "SELECT name FROM ~/Desktop, ~/Downloads WHERE name LIKE %csc%"
```

List all files in the current directory that are also present in some other directory:

```console
$ fsql -interactive
>>> SELECT all FROM . WHERE name IN (
...   SELECT name FROM ~/Desktop/files.bak/
... );
```

Passing queries via stdin without quotes is a bit of a pain, hopefully the next examples highlight that, my suggestion is to use interactive mode or wrap the query in quotes if you're doing anything with subqueries or attribute modifiers.

List all files named `main.go` in `$GOPATH` which are larger than 10.5 kilobytes or smaller than 100 bytes:

```console
$ fsql SELECT all FROM $GOPATH WHERE name = main.go AND \(FORMAT\(size, KB\) \>= 10.5 OR size \< 100\)
$ fsql "SELECT all FROM $GOPATH WHERE name = main.go AND (FORMAT(size, KB) >= 10.5 OR size < 100)"
$ fsql -interactive
>>> SELECT
...   all
... FROM
...   $GOPATH
... WHERE
...   name = main.go
...   AND (
...     FORMAT(size, KB) >= 10.5
...     OR size < 100
...   )
... ;
```

List the name, size, and modification time of JavaScript files in the current directory that were modified after April 1st 2017:

```console
$ fsql SELECT UPPER\(name\), FORMAT\(size, KB\), FORMAT\(time, ISO\) FROM . WHERE name LIKE %.js AND time \> \'Apr 01 2017 00 00\'
$ fsql "SELECT UPPER(name), FORMAT(size, KB), FORMAT(time, ISO) FROM . WHERE name LIKE %.js AND time > 'Apr 01 2017 00 00'"
$ fsql -interactve
>>> SELECT
...   UPPER(name),
...   FORMAT(size, KB),
...   FORMAT(time, ISO)
... FROM
...   .
... WHERE
...   name LIKE %.js
...   AND time > 'Apr 01 2017 00 00'
... ;
```

## Contribute

This project is completely open source, feel free to [open an issue](https://github.com/kshvmdn/fsql/issues) or [submit a pull request](https://github.com/kshvmdn/fsql/pulls).

Before submitting code, please ensure that tests are passing and the linter is happy. The following commands may be of use, refer to the [Makefile](./Makefile) to see what they do.

```sh
$ make fmt
```

```sh
$ make lint
```

```sh
$ make vet
```

```sh
$ make test
```

## License

fsql source code is available under the [MIT license](./LICENSE).
