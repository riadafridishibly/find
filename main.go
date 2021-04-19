package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/riadafridishibly/find/walk"
)

var Usage = func() {
	fmt.Fprintf(os.Stdout, "find [FLAGS/OPTIONS] [<pattern>] [<path>]\n")
	flag.PrintDefaults()
}

func main() {
	var rootDir string = "."
	var searchPattern string = ""

	flag.Parse()

	if flag.NArg() == 1 {
		searchPattern = flag.Arg(0)
	} else if flag.NArg() == 2 {
		searchPattern = flag.Arg(0)
		rootDir = flag.Arg(1)
	}

	if _, err := os.Stat(rootDir); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] No such file or directory\n", rootDir)
		os.Exit(2)
	}

	regex := regexp.MustCompile(searchPattern)

	filterFunc := func(filePath string) bool {
		return regex.MatchString(filePath)
	}

	ch := make(chan string, 1024)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if filterFunc(path) {
			ch <- path
		}
		return nil
	}

	done := walk.Walk(rootDir, walkFn)

	go func() {
		<-done
		close(ch)
	}()

	for filePath := range ch {
		fmt.Println(filePath)
	}
}
