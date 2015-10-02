package main

import (
	"bytes"
	"fmt"
	"github.com/gmas/go-pls"
	"io"
	"log"
	"os"
)

func main() {
	var playlists []pls.Playlist

	argsLen := len(os.Args)

	if argsLen < 2 {
		os.Exit(1)
	}

	inputs := os.Args[1:argsLen]
	//outFile := os.Args[argsLen-1]

	for _, inputFile := range inputs {

		plsReader, err := openPls(inputFile)
		if err != nil {
			log.Printf("WARNING\t %s", err)
			continue
		}
		log.Printf("opened file\t %s", inputFile)

		playlist, err := pls.Parse(plsReader)
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

	plsContentReader, err := newPls.Marshal()

	// reads contents of io.Reader into a Bytes.buffer
	plsContentBuff := new(bytes.Buffer)
	plsContentBuff.ReadFrom(plsContentReader)

	fmt.Print(plsContentBuff.String())
}

//TODO add option to download from URL
func openPls(inputPls string) (io.Reader, error) {
	reader, err := os.Open(inputPls)

	if err != nil {
		return nil, fmt.Errorf("%s \n", err)
	}

	return reader, nil
}
