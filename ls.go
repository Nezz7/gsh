package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Apostrophe used for names with spaces
const Apostrophe = "'"

// NumberOfNamesPerLine when outputing without list option
const NumberOfNamesPerLine = 4

// DirectorySymbol used at the end of each directory name
const DirectorySymbol = "\\"

var listPtr = flag.Bool("l", false, "Use a long listing format")
var sortPtr = flag.Bool("S", false, "Sort by file size, largest first")
var reversePtr = flag.Bool("r", false, "Reverse order while sorting")
var recursivePtr = flag.Bool("R", false, "List subdirectories recursively")
var allPtr = flag.Bool("a", false, "Do not ignore entries starting with .")

func parseTime(t time.Time) string {
	return fmt.Sprintf("%02d/%02d/%d  %02d:%02d", t.Month(), t.Day(), t.Year(), t.Hour(), t.Minute())
}

func parseName(f os.FileInfo) string {
	var name string
	if strings.Contains(f.Name(), " ") {
		name += Apostrophe + f.Name() + Apostrophe
	} else {
		name += f.Name()
	}
	if f.IsDir() {
		name += DirectorySymbol
	}
	return name
}

func isHidden(filename string) bool {
	if len(filename) == 0 {
		log.Println("empty file name")
	}
	if filename[0:1] == "." && len(filename) > 1 {
		return true
	}
	return false
}

func getDirPath() string {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func getRecursiveFiles(dirPath string) []os.FileInfo {
	files := make([]os.FileInfo, 0)
	err := filepath.Walk(".",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Fatal(err)
				return err
			}

			if isHidden(info.Name()) && info.IsDir() && !*allPtr {
				return filepath.SkipDir
			}
			if !isHidden(info.Name()) || *allPtr {
				files = append(files, info)
			}

			return nil
		})
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func getNonRecursiveFiles(dirPath string) []os.FileInfo {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func getFiles(dirPath string) []os.FileInfo {
	if *recursivePtr {
		return getRecursiveFiles(dirPath)
	}
	return getNonRecursiveFiles(dirPath)
}

func listFiles(files []os.FileInfo, dirPath string) {
	fmt.Printf("\n\tDirectory: %s\n\n\n", dirPath)
	fmt.Printf("Mode \t\t LastWriteTime \t\t Length Name\n")
	fmt.Printf("---- \t\t ------------- \t\t ------ ----\n")
	for _, f := range files {
		fmt.Printf("%s \t%s ", f.Mode(), parseTime(f.ModTime()))
		if f.IsDir() {
			fmt.Printf("\t\t")
		} else {
			fmt.Printf("%13d ", f.Size())
		}
		fmt.Printf("%s\n", parseName(f))
	}
}

func printFiles(files []os.FileInfo) {
	for i, f := range files {
		if i%NumberOfNamesPerLine == 0 {
			fmt.Println()
		}
		fmt.Printf("%40s", parseName(f))
	}
}

func reverse(files []os.FileInfo) []os.FileInfo {
	for left, right := 0, len(files)-1; left < right; left, right = left+1, right-1 {
		files[left], files[right] = files[right], files[left]
	}
	return files
}

func main() {
	flag.Parse()

	dirPath := getDirPath()
	files := getFiles(dirPath)

	if *sortPtr {
		sort.SliceStable(files, func(i, j int) bool {
			return files[i].Size() > files[j].Size()
		})
		if *reversePtr {
			files = reverse(files)
		}
	}

	if *listPtr {
		listFiles(files, dirPath)
	} else {
		printFiles(files)
	}

	fmt.Println()
	fmt.Println()

}
