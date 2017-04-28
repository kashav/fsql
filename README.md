## fsql

>Search through your file system with SQL-esque queries.

### Demo

![](./fsql.gif)

### Setup / installation

  - Requires Go to be [installed](https://golang.org/doc/install) and [configured](https://golang.org/doc/install#testing).

  - Install with Go tools:

    ```sh
    $ go get -v github.com/kshvmdn/fsql
    $ fsql
    ```

  - Or install directly via source:

    ```sh
    $ git clone https://github.com/kshvmdn/fsql.git $GOPATH/src/github.com/kshvmdn/fsql
    $ cd $_  # $GOPATH/src/github.com/kshvmdn/fsql
    $ make install && make
    $ ./fsql
    ```

### Usage

  - Provide the query as a command line argument via stdin.

  - Query structure:

    ```sql
    SELECT <attribute, ...> FROM <directory, ...> WHERE <condition, ...>
    ```

  - Supported attributes:

    + `name`
    + `size`

  - Supported comparators:

    + For numeric comparisons:

      * `>`
      * `>=`
      * `<`
      * `<=`
      * `=`
      * `<>`

    + For string comparisons:

      * `BEGINSWITH`
      * `ENDSWITH`
      * `IS`
      * `CONTAINS`

  - Examples:
    
    - List all files & directories in Desktop and Downloads directories that start with `csc`:

      ```sh
      $ fsql "SELECT name FROM ~/Desktop, ~/Downloads WHERE name BEGINSWITH csc"
      ```

    - List all JavaScript files that are larger than 300 bytes in the current directory (try running this on a `node_modules` directory, it's fast :sunglasses:).

      ```sh
      $ fsql "SELECT name, size FROM . WHERE size > 300, name ENDSWITH .js"
      ```

### Contribute

This project is completely open source, feel free to [open an issue](https://github.com/kshvmdn/issues) for questions/features/bugs or [submit a pull request](https://github.com/kshvmdn/pulls).

Use the following to test that your changes comply with [Golint](https://github.com/golang/lint).

  ```sh
  $ make lint
  ```

#### __TODO__

  - [ ] Handle errors more gracefully (instead of just panicking everything).
  - [ ] Add support for regex in string comparisons (e.g. `... ENDSWITH jsx?`).
  - [ ] Add support for `OR` / `AND`  / `()` (for precedence) in condition statements (lexing is already done for these, just need to add the parsers).
  - [ ] Add support for times/dates (to query file creation/modification time).
  - [ ] Introduce new attributes to select from (creation/modification time, file mode, _basically whatever else [`os.FileInfo`](https://golang.org/pkg/os/#FileInfo) supports_).
  - [ ] Add unit tests (test files are empty right now).

### Inspirations

Lexer & parser are loosely based on the amazing work of [JamesOwnHall/json2](https://github.com/JamesOwenHall/json2).
