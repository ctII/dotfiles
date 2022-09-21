package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var ignoreFiles = [...]string{".git", "README.md", "LICENSE", "go.mod", "go.sum", "install.go"}

func main() {
	log.SetFlags(0)

	verbose := flag.Bool("v", false, "be verbose")
	dryrun := flag.Bool("d", false, "print what we are doing instead of taking the action")

	flag.CommandLine.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "Usage: "+os.Args[0]+" <flags>")
		flag.CommandLine.PrintDefaults()
	}
	flag.Parse()

	if *verbose {
		log.SetFlags(log.Lshortfile)
	}

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get binary path (%v)\n", err)
	}

	var (
		relativePath string
		updateList   []string
	)
	err = filepath.WalkDir(currentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("could not continue walkdir (%w)", err)
		}

		relativePath, err = filepath.Rel(currentDir, path)
		if err != nil {
			return fmt.Errorf("could not make a relative path to current directory (%w)", err)
		}

		if *verbose {
			log.Printf("visiting (%v)\n", relativePath)
		}

		for i := range ignoreFiles {
			if ignoreFiles[i] == relativePath || filepath.SplitList(relativePath)[0] == ignoreFiles[i] {
				if *verbose {
					log.Printf("skipping ignored file (%v)\n", relativePath)
				}
				return fs.SkipDir
			}
		}

		if d.IsDir() {
			if *verbose {
				log.Printf("skipping (%v) as we dont update or symlink directories\n", relativePath)
			}
			return nil
		}

		updateList = append(updateList, relativePath)
		return nil
	})

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get current users home dir due to error (%w)\n", err)
	}

	for _, e := range updateList {
		if *dryrun {
			log.Printf("would have symlinked (%v) -> (%v)\n", filepath.Join(currentDir, e), filepath.Join(homeDir, e))
			log.Printf("would have deleted (%v)\n", filepath.Join(homeDir, e))
			continue
		}
		err = os.Remove(filepath.Join(homeDir, e))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("could not remove file (%v) due to error (%v)\n", filepath.Join(homeDir, e), err)
		}
		err = os.Symlink(filepath.Join(currentDir, e), filepath.Join(homeDir, e))
		if err != nil {
			log.Fatalf("could not symlink (%v) -> (%v)\n error (%v)\n", filepath.Join(currentDir, e), filepath.Join(homeDir, e), err)
		}
		if *verbose {
			log.Printf("symlinked (%v) -> (%v)\n", filepath.Join(currentDir, e), filepath.Join(homeDir, e))
		}
	}
}
