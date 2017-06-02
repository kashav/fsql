package prompt

import (
	"bufio"
	"fmt"
	"os"
)

var (
	status int32
	query  string
)

// parseLine appends line to query and returns true iff the last character
// of line is a semicolon.
func parseLine(line []byte) bool {
	if len(line) == 0 {
		return false
	}

	if len(query) == 0 {
		query = string(line)
	} else {
		query = fmt.Sprintf("%s %s", query, string(line))
	}

	// If we reach a semicolon, strip the last character of query and return
	// true (this query is done).
	if line[len(line)-1] == 59 {
		query = query[:len(query)-1]
		return true
	}

	return false
}

// readInput continually reads stdin for input.
func readInput(doneCh, quitCh chan<- bool) {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			switch status {
			case 0:
				fmt.Print(">>> ")
			case 1:
				fmt.Print("... ")
			}

			line, _, err := reader.ReadLine()
			if err != nil {
				quitCh <- true
				break
			}

			if done := parseLine(line); done {
				doneCh <- true
				status = 0
				break
			}
			status = 1
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
			temp := query
			query = ""
			return &temp
		case <-quitCh:
			fmt.Println("\nbye")
			break LOOP
		}
	}

	return nil
}
