package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type node struct {
	parent *node
	files map[string]int
	subdirs map[string]*node
	size int
}

func main() {
	file, err := os.Open(os.Args[len(os.Args)-1])
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	root := node{
		files: map[string]int{},
		subdirs: map[string]*node{},
	}
	slash := &node{
		parent: &root,
		files: map[string]int{},
		subdirs: map[string]*node{},
	}
	root.subdirs["/"] = slash

	cur := &root

	for scanner.Scan() {
		lineElems := strings.Split(scanner.Text(), " ")
		fmt.Printf("%+v\n", lineElems)
		switch lineElems[0] {
		case "$":
			switch lineElems[1] {
			case "cd":
				switch lineElems[2] {
				case "/":
					cur = root.subdirs["/"]
					fmt.Printf("switched to /: %+v\n", cur)
				case "..":
					cur = cur.parent
					fmt.Printf("switched to parent: %+v\n", cur)
				default:
					cur = cur.subdirs[lineElems[2]]
					fmt.Printf("switched to %s: %+v\n", lineElems[2], cur)
				}
			case "ls":
				// nothing to do
			}
		case "dir":
			newDir := &node{
				parent: cur,
				files: map[string]int{},
				subdirs: map[string]*node{},
			}

			cur.subdirs[lineElems[1]] = newDir
		default:
			fileSize, _ := strconv.Atoi(lineElems[0])
			cur.files[lineElems[1]] = fileSize
		}
	}

	fmt.Printf("finished reading file; root: %+v\n", root)
	printTree(&root)
	
	totalSize := getAndSetSize(&root)
	fmt.Printf("total tree size: %d\n", totalSize)

	totalSizeUnder100k := getSizeUnder100k(&root)
	fmt.Printf("total size <100k: %d\n", totalSizeUnder100k)

	totalSpace, targetSpace := 70000000, 30000000
	neededSpace := totalSize - (totalSpace-targetSpace)
	fmt.Printf("need to free up: %d\n", neededSpace)
	fmt.Printf("size of smallest viable dir: %d\n", getSmallestSubdirOver(&root, neededSpace))
}

func getAndSetSize(dir *node) int {
	if dir.size == 0 {
		filesSize, subdirsSize := 0, 0
		for _, file := range dir.files {
			filesSize += file
		}
		for _, sub := range dir.subdirs {
			subdirsSize += getAndSetSize(sub)
		}
		dir.size = filesSize + subdirsSize
	}
	return dir.size
}

func getSizeUnder100k(dir * node) int {
	totalSize := 0
	if dir.size < 100000 {
		totalSize += dir.size
	}
	totalSize += getSizeUnder100kHelper(dir)

	return totalSize
}

func getSizeUnder100kHelper(dir *node) int {
	totalSize := 0
	for _, subdir := range dir.subdirs {
		if subdir.size < 100000 {
			totalSize += subdir.size
		}
		totalSize += getSizeUnder100kHelper(subdir)
	}

	return totalSize
}

func getSmallestSubdirOver(dir *node, minimum int) int {
	return getSmallestSubdirOverHelper(dir, minimum, -1)
}

func getSmallestSubdirOverHelper(dir *node, minimum, guess int) int {
	if dir.size < minimum {
		return guess
	}
	localGuess := -1
	if dir.size < guess ||  guess == -1 {
		localGuess = dir.size
	}

	for _, subdir := range dir.subdirs {
		subDirGuess := getSmallestSubdirOverHelper(subdir, minimum, guess)
		if subDirGuess > minimum && (subDirGuess < localGuess || localGuess < 0) {
			localGuess = subDirGuess
		}
	}
	return localGuess
}

func printTree(dir *node) {
	printTreeHelper(dir, 0)
}

func printTreeHelper(dir *node, level int) {
	prefix := ""
	for i := 0; i < level; i++ {
		prefix += "  "
	}
	for name, size := range dir.files {
		fmt.Printf("%s%s (%d)\n", prefix, name, size)
	}
	for name, subdir := range dir.subdirs {
		fmt.Printf("%s%s (%d) ->\n", prefix, name, subdir.size)
		printTreeHelper(subdir, level+1)
	}
}