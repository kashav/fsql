## fsql

>Search through your file system with SQL-esque queries.

### Demo

[![fsql.gif](./fsql.gif)](https://asciinema.org/a/120534)

### Setup / installation

  - Requires Go to be [installed](https://golang.org/doc/install) and [configured](https://golang.org/doc/install#testing).

  - Install with Go:

    ```sh
    $ go get -v github.com/kshvmdn/fsql/...
    $ which fsql
    $GOPATH/bin/fsql
    ```

  - Or install directly via source:

    ```sh
    $ git clone https://github.com/kshvmdn/fsql.git $GOPATH/src/github.com/kshvmdn/fsql
    $ cd $_ # $GOPATH/src/github.com/kshvmdn/fsql
    $ make install && make
    $ ./fsql
    ```

### Usage

  - Pass the fsql query as a command line argument.

  - Query structure:

    ```sql
    SELECT attribute, ... FROM directory, ... WHERE conditional ...
    ```

    + Attribute can be any of the following: `name`, `size`, `mode`, `time`, or `*` (for all).

    + Directory should be a relative/absolute path to some directory on your file system. Also supports expanding environment variables (e.g. `$GOPATH`) and `~` (for your home directory). You can exclude a directory with a minus sign (`-`), e.g. exclude `.git`: `"... FROM ., -.git/ ..."`.

    + Conditionals:

      * Supported comparators:

        - For numeric comparisons (also applies to time):
          + `>`
          + `>=`
          + `<`
          + `<=`
          + `=`
          + `<>`

        - For string comparisons:
          + `=` - String equality (synonymous to using `LIKE` without any wildcards).
          + `<>` - Synonymous to using `WHERE NOT ... = ...`.
          + `LIKE` - For simple string matching, use `%` to match zero, one, or multiple characters. Check that a string begins with a value using `<value>%`, ends with a value: `%<value>`, or contains a value: `%<value>%`.
          + `RLIKE` - For pattern matching with regular expressions.

      * Use `AND` / `OR` for conditional conjunction/disjunction. Note that precedence is assigned based on order of appearance (i.e. `"WHERE a AND b OR c"` ≠ `"WHERE c OR b AND a"`). Use parentheses to get around this behaviour (`"WHERE a AND b OR c"` = `"WHERE c OR (b AND a)"`).

      * Use `NOT` to negate a conditional statement. This keyword **must** precede the statement (e.g. `"... WHERE NOT name LIKE foo ..."`).

      * If your value contains spaces and/or escaped characters, wrap the value in quotes (either single or double) or backticks.

      * The default unit for size is bytes, to use kilobytes / megabytes / gigabytes, append `kb` / `mb` / `gb` to the size value (e.g. `100kb`).

  - Examples:
    
    - List all files & directories in Desktop and Downloads directory that contain `csc`.

      ```sh
      $ fsql "SELECT name FROM ~/Desktop, ~/Downloads WHERE name LIKE %csc%"
      $ # equivalent to:
      $ fsql "SELECT name FROM ~/Desktop, ~/Downloads WHERE name RLIKE .*csc.*"
      ```

    - List all JavaScript files in the current directory that were modified after April 1st 2017 (try running this on a `node_modules` directory, it's fast :sunglasses:).

      ```sh
      $ fsql "SELECT name, size, time FROM . WHERE name LIKE %.js AND time > 'Apr 01 2017 00 00'"
      ```

    - List all files named `main.go` in `$GOPATH` which are at least 10.5 kilobytes in size or less than 100 bytes in size.

      ```sh
      $ fsql "SELECT * FROM $GOPATH WHERE name = main.go AND (size >= 10.5kb OR size < 100)"
      ```

### Contribute

This project is completely open source, feel free to [open an issue](https://github.com/kshvmdn/fsql/issues) for questions/features/bugs or [submit a pull request](https://github.com/kshvmdn/fsql/pulls).

Use the following to test that your changes comply with [Golint](https://github.com/golang/lint).

  ```sh
  $ make lint
  ```

#### TODO

  - [ ] **Bug**: Exclude skips files with similar names (e.g. excluding `.git` results in `.gitignore` not being listed).
  - [ ] Add unit tests (test files are empty right now).
  - [x] Add support for regex in string comparisons (e.g. `... ENDSWITH jsx?`).
  - [x] Handle errors more gracefully (instead of just panicking everything).
  - [x] Add support for `OR` / `AND`  / `()` (for precedence) in condition statements (lexing is already done for these, just need to add the parsers).
  - [x] Add support for times/dates (to query file creation/modification time).
  - [x] Introduce new attributes to select from (creation/modification time, file mode, _basically whatever else [`os.FileInfo`](https://golang.org/pkg/os/#FileInfo) supports_).
  - [x] **Bug**: Space-separated queries. Currently something like `"... WHERE time > May 1 ..."` is broken since we're splitting conditionals by space. Fix by allowing single quotes and backticks in query strings, so something like `"... WHERE time > 'May 1' ..."` works and evaluates the conditional to have value of `"May 1"`.
  - [x] Add `NOT` operator for negating conditionals.
  - [x] Add support for querying and selecting using other size units (only supports bytes right now, add functionality for KB, MB, and GB as well).
  - [x] **Bug**: Selecting from a directory and it's subdirectory results in duplicates and malformed output.

Lexer & parser are based on the amazing work of [**@JamesOwenHall**](https://github.com/JamesOwenHall) ([json2](https://github.com/JamesOwenHall/json2), [timed](https://github.com/JamesOwenHall/timed)).
