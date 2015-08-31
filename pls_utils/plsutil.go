package main

import (
	"./pls"
	"fmt"
	"log"
	"os"
)

func main() {
	var playlists []pls.Playlist

	argsLen := len(os.Args)
	//fmt.Printf("args: %d\n", argsLen)

	//FIXME check arg1 for input, arg2 for output
	if argsLen < 2 {
		os.Exit(1)
	}

	inputs := os.Args[1:argsLen]
	//outFile := os.Args[argsLen-1]

	for _, inputFile := range inputs {
		//TODO make it work with bufio ?
		contents, err := loadPls(inputFile)
		if err != nil {
			log.Printf("WARNING\t %s", err)
			continue
		}

		playlist, err := pls.Parse(contents)
		//FIXME extract the error handling for warnings/errors
		if err != nil {
			log.Printf("WARNING\t %s", err)
			continue
		}

		playlists = append(playlists, playlist)
	}

	pl := pls.Playlist{}
	newPls, err := pl.Merge(playlists...)
	if err != nil {
		panic(err)
	}
	fmt.Print(newPls.ToPls())
}

//TODO add option to download from URL
func loadPls(inputPls string) (_ string, _ error) {
	contents, _err := readLocalFile(inputPls)
	if _err != nil {
		return "", fmt.Errorf("%s \n", _err)
	}
	return contents, nil
}

func readLocalFile(inputFile string) (fileContents string, err error) {
	f, err := os.Open(inputFile)
	if err != nil {
		//log.Printf("missing input file: %v \n", inputFile)
		return "", fmt.Errorf("skipping missing file %s ", inputFile)
	} else {
		defer f.Close()
		//FIXME check count, return nil if 0?
		data := make([]byte, 2048)
		f.Read(data)
		return string(data), nil
	}
	return
}
