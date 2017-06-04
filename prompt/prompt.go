package prompt

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

type state int8

const (
	initial state = iota
	second
)

func (s state) String() string {
	switch s {
	case initial:
		return ">>>"
	case second:
		return "..."
	default:
		return ""
	}
}

var query bytes.Buffer

// parseLine appends line to query and returns true iff the last character
// of line is a semicolon.
func parseLine(line []byte) bool {
	if len(line) == 0 {
		return false
	}

	if query.Len() > 0 && !bytes.ContainsAny(query.Bytes(), "([") {
		query.WriteString(" ")
	}
	query.WriteString(string(line))

	// If we reach a semicolon, strip the last character of query and return
	// true (this query is done).
	if line[len(line)-1] == ';' {
		query.Truncate(query.Len() - 1)
		return true
	}

	return false
}

// readInput continually reads stdin for input.
func readInput(doneCh, quitCh chan<- bool) {
	var s state

	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("%s ", s.String())

			line, _, err := reader.ReadLine()
			if err != nil {
				quitCh <- true
				break
			}

			if done := parseLine(line); done {
				doneCh <- true
				s = initial
				break
			}
			s = second
		}
	}()
}

// Run reads input via stdin and returns the string upon reading a semicolon.
func Run() *string {
	doneCh := make(chan bool)
	quitCh := make(chan bool)

LOOP:
	for {
		readInput(doneCh, quitCh)

		select {
		case <-doneCh:
			temp := query.String()
			query.Reset()
			return &temp
		case <-quitCh:
			fmt.Println("\nbye")
			break LOOP
		}
	}

	return nil
}
