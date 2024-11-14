package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
)

const (
	demoDir   string = "./demo-dir"
	outputDir string = "./output-dir"
	demoDirFileCount int = 1_000_000
)

func main() {
	runtime.MemProfileRate = 1
	// Check demo-dir exists, since we want to walk it.
	_, err := os.Stat(demoDir)
	if err != nil {
		fmt.Println("Demo directory doesn't exist. Create it with the provided script and then run again.")
		log.Fatal(err)
	}

	// Create an output directory for the heap profile
	err = os.Mkdir(outputDir, 0750)
	if err != nil && !os.IsExist(err) {
		log.Fatal(err)
	}

	// Walk the demo directory.
	visitCounter := 0
	err = filepath.Walk(demoDir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			// Stop walking if we get an error.
			return err
		}

		// Write heap profile on first visit
		if visitCounter == 0 {
			runtime.GC() // Comment this line out if you want to see the odd behaviour discussed in the blog post :)
			err := writeHeap(filepath.Join(outputDir, "first-visit.prof"))
			if err != nil {
				fmt.Println("Error creating heap during first visit")
				log.Fatal(err)
			}
		}

		// Write heap profile on last visit
		if visitCounter == demoDirFileCount {
			runtime.GC()
			err := writeHeap(filepath.Join(outputDir, "last-visit.prof"))
			if err != nil {
				fmt.Println("Error creating heap during last visit")
				log.Fatal(err)
			}
		}

		visitCounter++
		// Print to give an indication of progress
		if visitCounter % 10_000 == 0 {
			fmt.Printf("Visited %d files \n", visitCounter)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Encountered error walking demo-dir")
		log.Fatal(err)
	}

	fmt.Println("Finished walking the demo dir")
}

func writeHeap(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)
	if err := pprof.WriteHeapProfile(f); err != nil {
		return err
	}

	return nil
}
