package pls

import (
	"./"
	_ "fmt"
	"reflect"
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
	t.Logf("Expected %d == %d + %d", mergedPlaylist.Length(),
		plst1.Length(), plst2.Length())

	if mergedPlaylist.Length() != plst1.Length()+plst2.Length() {
		t.Error("Merged playlist length != sum of merged playlists")
	}

	mergedEntries := append(entriesPlst1, entriesPlst2...)
	t.Logf("Expected %s \n==\n %s", mergedEntries, mergedPlaylist.Entries)
	if !reflect.DeepEqual(mergedEntries, mergedPlaylist.Entries) {
		t.Error("Merged playlist length != sum of merged playlists")
	}
}

func TestParsePlaylistEntries(t *testing.T) {
	playlistContent := `
[playlist]
numberofentries=3
File1=http://xstream1.somafm.com:8388
Title1=SomaFM: Beat Blender (#1 128k mp3): A late night blend of deep-house and downtempo chill.
Length1=-1
File2=http://xstream1.somafm.com:8384
Title2=SomaFM: Beat Blender (#2 128k mp3): A late night blend of deep-house and downtempo chill.
Length2=-1
File3=http://dcstream1.somafm.com:8384
Title3=SomaFM: Beat Blender (#3 128k mp3): A late night blend of deep-house and downtempo chill.
Length3=-1
`
	parsedPlaylist, err := pls.Parse(playlistContent)
	if err != nil {
		t.Error("Coould not parse playlist")
	}

	testEntries := []pls.PlaylistEntry{
		{File: `http://xstream1.somafm.com:8388`,
			Title: `SomaFM: Beat Blender (#1 128k mp3): A late night blend of deep-house and downtempo chill.`},
		{File: `http://xstream1.somafm.com:8384`,
			Title: `SomaFM: Beat Blender (#2 128k mp3): A late night blend of deep-house and downtempo chill.`},
		{File: `http://dcstream1.somafm.com:8384`,
			Title: `SomaFM: Beat Blender (#3 128k mp3): A late night blend of deep-house and downtempo chill.`},
	}
	if !reflect.DeepEqual(testEntries, parsedPlaylist.Entries) {
		t.Error("Failed to parse playlist entries correctly")
	}
}

func TestParseVersion(t *testing.T) {}

func TestMarshal(t *testing.T) {
}
