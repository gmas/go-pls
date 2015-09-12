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
