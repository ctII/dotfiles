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

	ignoreFiles := [...]string{".git", "README.md", "LICENSE", "go.mod", "go.sum", "install.go"}

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
	if err != nil {
		log.Fatal(err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("could not get current users home dir due to error (%v)\n", err)
	}

	for _, e := range updateList {
		currentDirLink := filepath.Join(currentDir, e)
		homeLink := filepath.Join(homeDir, e)

		if *dryrun {
			log.Printf("would have symlinked (%v) -> (%v)\n", currentDirLink, homeLink)
			log.Printf("would have deleted (%v)\n", homeLink)
			continue
		}

		err = os.Remove(homeLink)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("could not remove file (%v) due to error (%v)\n", homeLink, err)
		}

		err = os.MkdirAll(filepath.Dir(homeLink), 0o700)
		if err != nil {
			log.Fatalf("could not mkdir -p (%v) error (%v)", filepath.Dir(homeLink), err)
		}

		if *verbose {
			log.Printf("did equivalent of mkdir -p %v", filepath.Dir(homeLink))
		}

		err = os.Symlink(currentDirLink, homeLink)
		if err != nil {
			log.Fatalf("could not symlink (%v) -> (%v)\n error (%v)\n", currentDirLink, homeLink, err)
		}

		if *verbose {
			log.Printf("symlinked (%v) -> (%v)\n", currentDirLink, homeLink)
		}
	}
}
