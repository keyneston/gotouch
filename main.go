package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	g := &GoTouch{}

	flag.BoolVar(&g.SkipTesting, "skip-testing", false, "Skip creating _test.go files")
	flag.BoolVar(&g.Noop, "n", false, "NOOP; Don't do anything just log it")
	flag.BoolVar(&g.Verbose, "v", false, "Verbose")
	flag.Parse()

	if err := g.Touch(flag.Args()...); err != nil {
		log.Fatal(err)
	}
}

type GoTouch struct {
	SkipTesting bool
	Noop        bool
	Verbose     bool
}

func (g *GoTouch) VerboseLog(format string, extra ...interface{}) {
	if !g.Verbose {
		return
	}

	log.Printf(format, extra...)
}

func (g *GoTouch) Touch(files ...string) error {
	fileCreateList := []string{}
	dirCreateList := []string{}

	for _, f := range files {
		dir := filepath.Dir(f)

		fileCreateList = append(fileCreateList, f)
		dirCreateList = append(dirCreateList, dir)

		if t := convertToTest(f); !g.SkipTesting && t != "" {
			fileCreateList = append(fileCreateList, t)
		}
	}

	if err := g.mkdirs(dirCreateList...); err != nil {
		return err
	}
	if err := g.mkfiles(fileCreateList...); err != nil {
		return err
	}

	return nil
}

func (g *GoTouch) mkdirs(dirs ...string) error {
	for _, d := range dirs {
		g.VerboseLog("Creating %q", d)

		if !g.Noop {
			if err := os.MkdirAll(d, 0755); err != nil {
				return fmt.Errorf("Error making %s: %w", d, err)
			}
		}
	}

	return nil
}

func (g *GoTouch) mkfiles(files ...string) error {
	for _, f := range files {
		if err := g.mkfile(f); err != nil {
			return err
		}
	}

	return nil
}

func (g *GoTouch) mkfile(file string) error {
	g.VerboseLog("Creating %q in package %q", file, getPackageName(file))

	if _, err := os.Lstat(file); !errors.Is(err, os.ErrNotExist) {
		g.VerboseLog("File exists continuing: %v", file)
		return nil
	}

	if !g.Noop {
		f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return fmt.Errorf("Error creating %q: %w", file, err)
		}
		defer f.Close()
		fmt.Fprintf(f, "package %s\n", getPackageName(file))
	}

	return nil
}

func getPackageName(input string) string {
	dir := filepath.Dir(input)
	return filepath.Base(dir)
}

func convertToTest(name string) string {
	ext := filepath.Ext(name)
	if ext != ".go" {
		return ""
	}

	base := strings.TrimSuffix(name, ext)
	if strings.HasSuffix(base, "_test") {
		return ""
	}

	return fmt.Sprintf("%s_test%s", base, ext)
}
