package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var Usage = func() {
	fmt.Fprintf(os.Stdout, "find [FLAGS/OPTIONS] [<pattern>] [<path>]\n")
	flag.PrintDefaults()
}

var fileTypeFilter string
var fileExtensionFilter string
var caseSensitive bool

var rootDir string = "."

func searchFunc(regexPatt regexp.Regexp, typeFilter, extensionFilter []string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if regexPatt.MatchString(info.Name()) {
			fmt.Println(path)
		}
		return nil
	}
}

func main() {
	flag.Usage = Usage
	flag.Parse()

	typeFilterSlice := strings.Split(fileTypeFilter, ",")
	extensionFilterSlice := strings.Split(fileExtensionFilter, ",")

	var searchPattern string = ""
	if flag.NArg() == 1 {
		searchPattern = flag.Arg(0)
	} else if flag.NArg() == 2 {
		searchPattern = flag.Arg(0)
		rootDir = flag.Arg(1)
	}

	if !caseSensitive {
		searchPattern = "(?i)" + searchPattern
	}

	if _, err := os.Stat(rootDir); err != nil {
		fmt.Fprintf(os.Stderr, "[%s] No such file or directory\n", rootDir)
		os.Exit(2)
	}

	regex := regexp.MustCompile(searchPattern)

	err := filepath.Walk(rootDir, searchFunc(*regex, typeFilterSlice, extensionFilterSlice))

	if err != nil {
		panic(err)
	}
}

func init() {
	flag.StringVar(&fileTypeFilter, "type", "", "Filter search result by file type")
	flag.StringVar(&fileTypeFilter, "t", "", "Filter search result by file type (shorthand)")
	flag.StringVar(&fileExtensionFilter, "ext", "", "Filter search result by file extension")
	flag.StringVar(&fileExtensionFilter, "e", "", "Filter search result by file extension (shorthand)")
	flag.BoolVar(&caseSensitive, "sensitive-case", false, "Case sensitive search")
	flag.BoolVar(&caseSensitive, "s", false, "Case sensitive search (shorthand)")
}
