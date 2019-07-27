package terminal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kashav/fsql"
	"github.com/kashav/fsql/terminal/pager"

	"golang.org/x/crypto/ssh/terminal"
)

var fd = int(os.Stdin.Fd())
var query bytes.Buffer

// Start listens for queries via stdin and invokes fsql.Run whenever a
// semicolon is read.
func Start() error {
	if !terminal.IsTerminal(fd) {
		return errors.New("not a terminal")
	}

	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, state)

	prompt := ">>> "
	term := terminal.NewTerminal(os.Stdin, prompt)

	// Listen for queries and invoke run whenever a semicolon is read. Continues
	// until receiving an EOF (Ctrl-D) or _fatal_ error (i.e. anything not
	// caused by the query itself).
	for {
		line, err := term.ReadLine()
		if err == io.EOF {
			fmt.Print("\r\nbye\r\n")
			break
		}
		if err != nil {
			return err
		}

		if line == "exit" {
			fmt.Print("bye\r\n")
			break
		}

		// TODO: If the previous character was a paren., bracket, or quote, we
		// don't want to add a space here (although not necessary, since the
		// tokenizer handles excess whitespace).
		if query.Len() > 0 {
			query.WriteString(" ")
		}
		query.WriteString(line)

		if strings.HasSuffix(line, ";") {
			query.Truncate(query.Len() - 1)

			b := []byte{}
			if out, err := run(query.String()); err != nil {
				// This error likely corresponds to the query, so instead of exiting
				// interactive mode, we simply write the error to stdout and proceed.
				b = append(b, []byte(err.Error())...)
				b = append(b, '\a', '\n')
				term.Write(b)
			} else if len(out) > 0 {
				_, h, err := terminal.GetSize(fd)
				if err != nil {
					return err
				}

				b = append(b, []byte(out)...)

				// Write to stdout if out is less than 3/4 of the height of the
				// window OR the `less` command doesn't exist; otherwise, invoke the
				// pager.
				if float64(strings.Count(out, "\n")) <= 0.75*float64(h) ||
					!pager.CommandExists() {
					term.Write(b)
				} else if err = pager.New(b); err != nil {
					return err
				}
			}

			query.Reset()
		}

		prompt = "... "
		if query.Len() == 0 {
			prompt = ">>> "
		}
		term.SetPrompt(prompt)
	}
	return nil
}

// run invokes fsql.Run with the provided query string.
func run(query string) (out string, err error) {
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return "", err
	}
	os.Stdout = w
	defer func() { os.Stdout = stdout }()

	ch := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		ch <- buf.String()
	}()

	err = fsql.Run(query)
	// Must happen after the function call and before we try to read from ch.
	if closeErr := w.Close(); closeErr != nil {
		return "", closeErr
	}
	out = <-ch
	return
}
