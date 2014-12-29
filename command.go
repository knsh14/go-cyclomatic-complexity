package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	path  string
	limit int
)

func init() {
	flag.StringVar(&path, "path", ".", "root path to measure")
	flag.IntVar(&limit, "limit", 12, "max complexity to allow")
}

func main() {
	flag.Parse()
	filepath.Walk(path, walk)
}

func walk(path string, info os.FileInfo, err error) error {
	if strings.HasPrefix(path, ".") {
		return nil
	}
	fmt.Println(path)
	return nil
}
