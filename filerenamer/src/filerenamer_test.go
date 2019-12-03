package filerenamer

import (
	"os"
	"testing"
)

func TestRenameFileName(t *testing.T) {
	fileName := "abc"
	textToRemove := "c"

	got,_ := renameFileName(fileName, textToRemove)
	if got != "ab" {
		t.Errorf("RenameFileName(%s) = %s", fileName, "ab")
	}
}

func TestRenameFileName_With_Does_Not_Contain(t *testing.T) {
	fileName := "abc"
	textToRemove := "x"

	got,_ := renameFileName(fileName, textToRemove)
	if got != "abc" {
		t.Errorf("RenameFileName(%s) = %s", fileName, "abc")
	}
}

func TestGetFiles(t *testing.T) {
	folderToRead := "."
	got, _ := getFiles(folderToRead)
	if len(got) != 0 {
		t.Errorf("GetFiles should return empty results")
	}
}

func TestGetFiles_Should_Return_Something(t *testing.T) {
	folderToRead := os.TempDir()
	got, _ := getFiles(folderToRead)
	if len(got) < 0 {
		t.Errorf("GetFiles should return files")
	}
}
