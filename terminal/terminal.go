package terminal

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kshvmdn/fsql"
	"github.com/kshvmdn/fsql/terminal/pager"

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

		// TODO: If the previous character was a paren., bracket, or quote, we
		// don't want to add a space here.
		if query.Len() > 0 {
			query.WriteString(" ")
		}
		query.WriteString(line)

		if strings.HasSuffix(line, ";") {
			query.Truncate(query.Len() - 1)

			if out, err := run(query); err != nil {
				// This error was likely caused by the query itself, so instead of
				// exiting interactive mode, we simply write the error to stdout and
				// proceed.
				term.Write([]byte(err.Error() + "\n"))
			} else if len(out) > 0 {
				b := []byte(out)

				_, h, err := terminal.GetSize(fd)
				if err != nil {
					return err
				}

				// Invoke the pager iff out has more lines than 3/4 of the height of
				// the terminal.
				if float64(strings.Count(out, "\n")) <= 0.75*float64(h) {
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
//
// TODO: We ignore all stderr output, so the terminal breaks whenever anything
// is logged by fsql.Run. We need to find a way of capturing stderr and
// returning the error when it's non-empty.
func run(query bytes.Buffer) (string, error) {
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

	err = fsql.Run(query.String())
	w.Close()
	if err != nil {
		return "", err
	}
	return <-ch, nil
}
