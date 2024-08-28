package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type excludedFiles []string

func (e *excludedFiles) Set(value string) error {
	*e = append(*e, value)
	return nil
}

func (e *excludedFiles) String() string {
	return strings.Join(*e, ",")
}

func Search(dir, search string, ignoreCase, colors bool, e excludedFiles) error {
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() && slices.Contains(e, info.Name()) {
			return filepath.SkipDir
		}
		if err != nil {
			if errors.Is(err, os.ErrPermission) {
				fmt.Println(err) // will exit with success code
			} else {
				return err // will exit with error code
			}
		}
		var idx int

		if ignoreCase {
			idx = strings.Index(strings.ToLower(path), strings.ToLower(search))
		} else {
			idx = strings.Index(path, search)
		}

		// Checking if a match was found
		if idx != -1 {
			// checking if user wants to add colors
			if colors {
				fmt.Println(highlight(path, idx, idx+len(search)))
			} else {
				fmt.Println(path)
			}
		}

		return nil
	})
	return err
}

func highlight(s string, start, end int) string {
	return s[0:start] + "\x1b[31m" + s[start:end] + "\x1b[0m" + s[end:]
}

func main() {
	var dir string
	var excluded excludedFiles
	var ignoreCase bool
	var colors bool
	var search string

	flag.StringVar(&dir, "f", ".", "the directory to search in")
	flag.Var(&excluded, "e", "List of directories to be excluded")
	flag.BoolVar(&ignoreCase, "i", false, "ignore case")
	flag.BoolVar(&colors, "c", false, "highlight text")
	flag.Parse()

	search = flag.Arg(0)
	if search == "" || dir == "" {
		fmt.Printf("Usage: glocate [OPTIONS]... [SEARCH_TEXT]...\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	err := Search(dir, search, ignoreCase, colors, excluded)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
