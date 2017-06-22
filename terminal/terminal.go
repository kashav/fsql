package terminal

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kshvmdn/fsql"

	"golang.org/x/crypto/ssh/terminal"
)

// Start listens for queries via stdin and invokes fsql.Run whenever a
// semicolon is read.
func Start() error {
	state, err := terminal.MakeRaw(0)
	if err != nil {
		return err
	}
	defer terminal.Restore(0, state)

	prompt := ">>> "
	term := terminal.NewTerminal(os.Stdin, prompt)

	var query bytes.Buffer

	for {
		line, err := term.ReadLine()
		if err == io.EOF {
			fmt.Print("\r\nbye\r\n")
			break
		}
		if err != nil {
			return err
		}

		// TODO: If the previous character was a paren, bracket, or quote, we don't
		// want to add a space after it.
		if query.Len() > 0 {
			query.WriteString(" ")
		}
		query.WriteString(line)

		if strings.HasSuffix(line, ";") {
			// Remove trailing semicolon.
			query.Truncate(query.Len() - 1)

			var buf bytes.Buffer
			if out, err := run(query); err != nil {
				// This error was likely caused by the query itself, so instead of
				// exiting interactive mode, we simply write the error to stdout and
				// continue.
				buf.WriteString(err.Error() + "\n")
			} else {
				buf.WriteString(out)
			}
			term.Write(buf.Bytes())

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
// is logged. We need to find a way of capturing stderr and returning an error
// when it's non-empty.
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
