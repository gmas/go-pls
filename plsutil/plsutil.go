package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"net/http"
	"io/ioutil"
//	"regexp"
	"sync"

	"github.com/gmas/go-pls"
)

func emit(urls ...string) <- chan string {
	out := make(chan string)
	go func() {
		for _, n := range urls {
			out <- n
			}
		close(out)
	}()
	return out
}

func loadAndParse(in <- chan string) <-chan pls.Playlist {
	out := make(chan pls.Playlist)
	go func() {
		for n := range in {
			plsReader, err := downloadPls(string(n))
			if err != nil {
				log.Printf("WARNING\t %s", err)
				continue
			}
			playlist, err := pls.Parse(plsReader)
			out <- playlist
		}
		close(out)
		}()
		return out
}

func merge(cs ...<-chan pls.Playlist) <-chan pls.Playlist{
    var wg sync.WaitGroup
		out := make(chan pls.Playlist)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan pls.Playlist) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

func main() {
	var playlists []pls.Playlist
	// var wg sync.WaitGroup

	argsLen := len(os.Args)

	if argsLen < 2 {
		os.Exit(1)
	}

	inputs := os.Args[1:argsLen]
	//outFile := os.Args[argsLen-1]

	// for _, inputFile := range inputs {
	// 	var plsReader io.Reader
	// 	var err error

	// 	httpRegexp := regexp.MustCompile(`^http(s*):\/\/`)
	// 	matches := httpRegexp.FindStringSubmatch(inputFile)
	// 	if len(matches) > 0 {
	// 		// plsReader, err = downloadPls(inputFile)
	// 	} else {
	// 		plsReader, err = openPls(inputFile)
	// 	}

	// 	if err != nil {
	// 		log.Printf("WARNING\t %s", err)
	// 		continue
	// 	}
	// 	// log.Printf("opened file\t %s", inputFile)
	// 	// playlist, err := pls.Parse(plsReader)
	// 	//FIXME extract the error handling for warnings/errors
	// 	if err != nil {
	// 		log.Printf("WARNING\t %s", err)
	// 		continue
	// 	}
	// }


	urls := emit(inputs...)

	c1 := loadAndParse(urls)
	c2 := loadAndParse(urls)
	c3 := loadAndParse(urls)
	c4 := loadAndParse(urls)
	
	for playlist := range merge(c1,c2,c3,c4) {
		playlists = append(playlists, playlist)
	}

	pl := pls.Playlist{}
	pl, _ = pl.Merge(playlists...)

	plsContentReader, err := pl.Marshal()
	if err != nil {
		panic(err)
	}
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
