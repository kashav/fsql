package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/kshvmdn/fsql/query"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected input string.")
		os.Exit(1)
	}
	input := os.Args[1]

	usr, err := user.Current()
	if err != nil {
		panic(err)
	}

	q, err := query.Run(input)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	var showName bool
	var showSize bool

	for _, attr := range q.Attributes {
		if attr == "name" {
			showName = true
		} else if attr == "size" {
			showSize = true
		}
	}

	for _, src := range q.Sources {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()

			filepath.Walk(filepath.Join(usr.HomeDir, src[2:]), func(path string, info os.FileInfo, err error) error {
				if path == "." {
					return nil
				}

				var isMatch bool

				for _, condition := range q.Conditions {
					switch condition.Attribute {
					case "name":
						isMatch = handleNameComparison(condition.Comparator, info.Name(), condition.Value)

					case "size":
						isMatch = handleSizeComparison(condition.Comparator, info.Size(), condition.Value)
					}

					if !isMatch {
						return nil
					}
				}

				if isMatch {
					if showSize {
						fmt.Printf("%d\t", info.Size())
					}

					if showName {
						fmt.Printf("%s", filepath.Join("~", path[len(usr.HomeDir):]))
					}

					fmt.Printf("\n")
				}
				return nil
			})
		}(src)
	}

	wg.Wait()
}

func handleNameComparison(comparator query.TokenType, fileName string, inputFileName string) bool {
	isMatch := false

	switch comparator {
	case query.BeginsWith:
		isMatch = strings.HasPrefix(fileName, inputFileName)
	case query.EndsWith:
		isMatch = strings.HasSuffix(fileName, inputFileName)
	case query.Is:
		isMatch = fileName == inputFileName
	case query.Contains:
		isMatch = strings.Contains(fileName, inputFileName)
	}

	return isMatch
}

func handleSizeComparison(comparator query.TokenType, fileSize int64, inputSizeStr string) bool {
	isMatch := false

	size, err := strconv.ParseInt(inputSizeStr, 10, 64)
	if err != nil {
		return isMatch
	}

	switch comparator {
	case query.Equals:
		isMatch = fileSize == size
	case query.NotEquals:
		isMatch = fileSize != size
	case query.GreaterThanEquals:
		isMatch = fileSize >= size
	case query.GreaterThan:
		isMatch = fileSize > size
	case query.LessThanEquals:
		isMatch = fileSize <= size
	case query.LessThan:
		isMatch = fileSize < size
	}

	return isMatch
}
