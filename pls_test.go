package pls

import (
	"./"
	"bytes"
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
	plst1 := &pls.Playlist{
		Entries: []pls.PlaylistEntry{
			{`http://localhost:8081`, `test1`},
			{`http://localhost:8082`, `test2`},
			{`http://localhost:8083`, `test3`},
		}}
	plst2 := &pls.Playlist{
		Entries: []pls.PlaylistEntry{
			{`http://localhost:8084`, `test4`},
		}}

	expectedEntries := append(plst1.Entries, plst2.Entries...)
	mergedPlaylist, _ := plst1.Merge(*plst2)

	t.Logf("Expected %d == %d + %d", mergedPlaylist.Length(),
		plst1.Length(), plst2.Length())

	if mergedPlaylist.Length() != plst1.Length()+plst2.Length() {
		t.Error("Merged playlist length != sum of merged playlists")
	}

	t.Logf("Expected %s \n==\n %s", expectedEntries, mergedPlaylist.Entries)
	if !reflect.DeepEqual(expectedEntries, mergedPlaylist.Entries) {
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
	parsedPlaylist, err := pls.Parse(bytes.NewBufferString(playlistContent))
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

func TestParseOutOfOrder(t *testing.T) {
	playlistContent := `
[playlist]
numberofentries=3
File3=http://dcstream1.somafm.com:8384
Title1=SomaFM: Beat Blender (#1 128k mp3): A late night blend of deep-house and downtempo chill.
Length1=-1
Title2=SomaFM: Beat Blender (#2 128k mp3): A late night blend of deep-house and downtempo chill.
Length2=-1
File1=http://xstream1.somafm.com:8388
Title3=SomaFM: Beat Blender (#3 128k mp3): A late night blend of deep-house and downtempo chill.
Length3=-1
File2=http://xstream1.somafm.com:8384
`
	parsedPlaylist, err := pls.Parse(bytes.NewBufferString(playlistContent))
	if err != nil {
		t.Error("Coould not parse playlist")
	}

	expectedEntries := []pls.PlaylistEntry{
		{File: `http://xstream1.somafm.com:8388`,
			Title: `SomaFM: Beat Blender (#1 128k mp3): A late night blend of deep-house and downtempo chill.`},
		{File: `http://xstream1.somafm.com:8384`,
			Title: `SomaFM: Beat Blender (#2 128k mp3): A late night blend of deep-house and downtempo chill.`},
		{File: `http://dcstream1.somafm.com:8384`,
			Title: `SomaFM: Beat Blender (#3 128k mp3): A late night blend of deep-house and downtempo chill.`},
	}
	if !reflect.DeepEqual(expectedEntries, parsedPlaylist.Entries) {
		t.Error("Failed to parse out of order playlist entries")
	}
}

func TestMarshal(t *testing.T) {
	plst := &pls.Playlist{
		Entries: []pls.PlaylistEntry{
			{`http://localhost:8081`, `test1`},
			{`http://localhost:8082`, `test2`},
			{`http://localhost:8083`, `test3`},
		}}

	expected := `[playlist]
numberofentries=3
File1=http://localhost:8081
Title1=test1
Length1=-1
File2=http://localhost:8082
Title2=test2
Length2=-1
File3=http://localhost:8083
Title3=test3
Length3=-1
Version=2
`
	marshalledReader, err := plst.Marshal()
	marshalledBuff := (bytes.Buffer{})
	marshalledBuff.ReadFrom(marshalledReader)
	marshalled := marshalledBuff.String()

	t.Logf("Expected\n%s\nActual\n%s", expected, marshalled)

	if err != nil {
		t.Error("Failed to marshal Playlist")
	}
	if string(marshalled) != expected {
		t.Error("Failed to marshal Playlist")
	}
}
