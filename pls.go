package pls

import (
	"bytes"
	"fmt"
	//"io"
	"regexp"
	"strconv"
	"strings"
)

type Playlist struct {
	Entries []PlaylistEntry
}

type PlaylistEntry struct {
	File  string
	Title string
}

func (pl Playlist) Length() int {
	return len(pl.Entries)
}

func (plEntry PlaylistEntry) String() string {
	return fmt.Sprintf("Title: %s\nFile: %s \n", plEntry.Title, plEntry.File)
}

func (plEntry *PlaylistEntry) isEmpty() bool {
	if plEntry.Title == "" && plEntry.File == "" {
		return true
	}
	return false
}

func Parse(contents string) (pl Playlist, err error) {
	//FIXME entries needs to grow itself
	// use a Set so we get no dupes ?
	_entries := make([]PlaylistEntry, 1024)

	fileRegexp := regexp.MustCompile(`File(\d+)=`)
	for _, line := range strings.Split(contents, "\n") {
		fileMatches := fileRegexp.FindStringSubmatch(line)
		if len(fileMatches) > 0 {
			fileId, _ := strconv.ParseInt(fileMatches[1], 10, 64)
			file := strings.TrimPrefix(line, fileMatches[0])
			(_entries[int64(fileId)]).File = file
		}

		titleRegexp := regexp.MustCompile(`Title(\d+)=`)
		titleMatches := titleRegexp.FindStringSubmatch(line)
		if len(titleMatches) > 0 {
			titleId, _ := strconv.ParseInt(titleMatches[1], 10, 64)
			title := strings.TrimPrefix(line, titleMatches[0])
			(_entries[int32(titleId)]).Title = title
		}
	}

	for i := 0; i < len(_entries); i++ {
		if entry := _entries[i]; !entry.isEmpty() {
			pl.Entries = append(pl.Entries, entry)
		}
	}

	return pl, nil
}

func (pl *Playlist) Marshal() (v []byte, err error) {
	var buff bytes.Buffer

	buff.WriteString("[playlist]\n")
	buff.WriteString(fmt.Sprintf("numberofentries=%d\n", len(pl.Entries)))
	for i, _ := range pl.Entries {
		idx := i + 1
		buff.WriteString(fmt.Sprintf(fmt.Sprintf("File%d=%s\n", idx, pl.Entries[i].File)))
		buff.WriteString(fmt.Sprintf(fmt.Sprintf("Title%d=%s\n", idx, pl.Entries[i].Title)))
		buff.WriteString(fmt.Sprintf(fmt.Sprintf("Length%d=-1\n", idx)))
	}
	buff.WriteString("Version=2\n")

	return buff.Bytes(), err
}

func (pl *Playlist) AddEntry(entries ...PlaylistEntry) (int, error) {
	pl.Entries = append(pl.Entries, entries...)
	return len(pl.Entries), nil
}

func (pl *Playlist) Merge(playlists ...Playlist) (Playlist, error) {
	newPl := Playlist{}
	newPl.Entries = append(newPl.Entries, pl.Entries...)
	for i, _ := range playlists {
		newPl.AddEntry(playlists[i].Entries...)
	}
	return newPl, nil
}

//func (pl *Playlist) Read(p []byte) (n int, err error) {
//	p = append(p, "[playlist]\n"...)
//	return len(p), nil
//}
