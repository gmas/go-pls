package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type playlistEntry struct {
	File  string
	Title string
	//Length  uint32
}

func (p *playlistEntry) String() string {
	return fmt.Sprintf("Title: %s, File: %s \n", p.Title, p.File)
}

func (p *playlistEntry) Empty() bool {
	if p.Title == "" && p.File == "" {
		return true
	}
	return false
}

func main() {
	var entries []playlistEntry
	argsLen := len(os.Args)
	//fmt.Printf("args: %d\n", argsLen)

	//FIXME check arg1 for input, arg2 for output
	if argsLen < 2 {
		os.Exit(1)
	}

	inputs := os.Args[1 : argsLen-1]
	outFile := os.Args[argsLen-1]

	for _, inputFile := range inputs {
		readEntries, err := parsePls(inputFile)
		if err != nil {
			log.Printf("%s \n", err)
			continue
		}
		for i := 0; i < len(readEntries); i++ {
			if entry := readEntries[i]; !entry.Empty() {
				entries = append(entries, entry)
			}
		}
	}
	writePls(entries, outFile)
}

func writePls(entries []playlistEntry, fileName string) bool {
	file, err := os.Create(fileName)
	writer := bufio.NewWriter(file)
	defer func() {
		if err := writer.Flush(); err != nil {
			panic(err)
		}
		file.Close()
	}()

	if err != nil {
		panic(err)
	}

	fmt.Fprint(writer, "[playlist]\n")
	fmt.Fprintf(writer, fmt.Sprintf("numberofentries=%d\n", len(entries)))
	for i, _ := range entries {
		idx := i + 1
		fmt.Fprintf(writer, fmt.Sprintf("File%d=%s\n", idx, entries[i].File))
		fmt.Fprintf(writer, fmt.Sprintf("Title%d=%s\n", idx, entries[i].Title))
		fmt.Fprintf(writer, fmt.Sprintf("Length%d=-1\n", idx))
	}
	fmt.Fprintln(writer, "Version=2")

	return true
}

func parsePls(inputPlsFile string) (parsedEntries []playlistEntry, err error) {
	log.Printf("Reading %s, ", inputPlsFile)
	contents, err := readLocalFile(inputPlsFile)
	if err != nil {
		return nil, fmt.Errorf("%s \n", err)
	}

	//FIXME entries needs to grow itself
	entries := make([]playlistEntry, 1024)

	fileRegexp := regexp.MustCompile(`File(\d+)=`)
	for _, line := range strings.Split(contents, "\n") {
		fileMatches := fileRegexp.FindStringSubmatch(line)
		if len(fileMatches) > 0 {
			fileId, _ := strconv.ParseInt(fileMatches[1], 10, 64)
			file := strings.TrimPrefix(line, fileMatches[0])
			(entries[int32(fileId)]).File = file
		}

		titleRegexp := regexp.MustCompile(`Title(\d+)=`)
		titleMatches := titleRegexp.FindStringSubmatch(line)
		if len(titleMatches) > 0 {
			titleId, _ := strconv.ParseInt(titleMatches[1], 10, 64)
			title := strings.TrimPrefix(line, titleMatches[0])
			(entries[int32(titleId)]).Title = title
		}
	}
	//... lets you pass multiple arguments to a variadic function from a slice
	return append(parsedEntries, entries...), nil
}

func readLocalFile(inputFile string) (fileContents string, err error) {
	f, err := os.Open(inputFile)
	defer f.Close()
	if err != nil {
		//log.Printf("missing input file: %v \n", inputFile)
		return "", fmt.Errorf("skipping missing file %s ", inputFile)
	} else {
		//FIXME check count, return nil if 0 maybe?
		data := make([]byte, 2048)
		f.Read(data)
		return string(data), nil
	}
	return
}
