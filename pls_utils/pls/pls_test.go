package pls

import (
	"./"
	_ "fmt"
	"testing"
)

func TestAddEntry(t *testing.T) {
	playlist := pls.Playlist{}
	entry := pls.PlaylistEntry{File: `http://soma.fm/dronezone.pls`, Title: `DroneZone`}

	playlist.AddEntry(entry)
	t.Logf("Expected %v == %v", 1, len(playlist.Entries))
	if entriesLen := len(playlist.Entries); entriesLen < 1 {
		t.Error("Expected %d did not match value %d", 1, entriesLen)
	}

}

func TestMerge(t *testing.T) {
	entriesPlst1 := []pls.PlaylistEntry{
		{`http://localhost:8081`, `test1`},
		{`http://localhost:8082`, `test2`},
		{`http://localhost:8083`, `test3`},
	}
	entriesPlst2 := []pls.PlaylistEntry{
		{`http://localhost:8084`, `test4`},
	}
	plst1 := &pls.Playlist{Entries: entriesPlst1}
	plst2 := &pls.Playlist{Entries: entriesPlst2}

	mergedPlaylist, _ := plst1.Merge(*plst2)
	t.Logf("Expected %d == %d + %d", mergedPlaylist.Length(), plst1.Length(), plst2.Length())

	if mergedPlaylist.Length() != plst1.Length()+plst2.Length() {
		t.Error("Merged playlist length != sum of merged playlists")
	}
}

func TestParse(t *testing.T) {
}

func TestMarshal(t *testing.T) {
}
