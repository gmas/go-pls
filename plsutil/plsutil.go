package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"net/http"
	"io/ioutil"
	"regexp"

	"github.com/gmas/go-pls"
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
		var plsReader io.Reader
		var err error

		httpRegexp := regexp.MustCompile(`^http(s*):\/\/`)
		matches := httpRegexp.FindStringSubmatch(inputFile)
		if len(matches) > 0 {
			plsReader, err = downloadPls(inputFile)
		} else {
			plsReader, err = openPls(inputFile)
		}

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

func downloadPls(plsURL string) (io.Reader, error) {
	response, err := http.Get(plsURL)

	if err != nil {
		return nil, fmt.Errorf("%s \n", err)
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	return bytes.NewReader(body), nil
}
