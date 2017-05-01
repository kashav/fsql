package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kshvmdn/fsql/query"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected query.")
		os.Exit(1)
	}
	input := os.Args[1]

	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	q, err := query.Run(input)

	if err != nil {
		if err == io.ErrUnexpectedEOF {
			log.Fatal("Unexpected end of line.")
		}

		log.Fatal(err)
	}

	var wg sync.WaitGroup

	for _, src := range q.Sources {
		wg.Add(1)
		go func(src string) {
			defer wg.Done()

			if strings.Contains(src, "~/") {
				src = filepath.Join(usr.HomeDir, src[2:])
			}

			filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
				if path == "." || path == ".." {
					return nil
				}

				var isMatch bool

				for _, condition := range q.Conditions {
					switch condition.Attribute {
					case "name":
						isMatch = compareName(condition.Comparator, info.Name(), condition.Value)

					case "size":
						isMatch = compareSize(condition.Comparator, info.Size(), condition.Value)
					}

					if !isMatch {
						return nil
					}
				}

				if isMatch {
					if q.HasAttribute("mode") {
						fmt.Printf("%s\t", info.Mode())
					}

					if q.HasAttribute("size") {
						fmt.Printf("%d\t", info.Size())
					}

					if q.HasAttribute("time") {
						fmt.Printf("%s\t", info.ModTime().Format(time.Stamp))
					}

					if q.HasAttribute("name") {
						if strings.Contains(path, usr.HomeDir) {
							path = filepath.Join("~", path[len(usr.HomeDir):])
						}
						fmt.Printf("%s", path)
					}

					fmt.Printf("\n")
				}
				return nil
			})
		}(src)
	}

	wg.Wait()
}

func compareName(comparator query.TokenType, fileName string, inputFileName string) bool {
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

func compareSize(comparator query.TokenType, fileSize int64, inputSizeStr string) bool {
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
