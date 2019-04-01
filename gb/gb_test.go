package gb

import (
	"fmt"
	"path/filepath"
	"testing"
)

type Assert struct {
	t *testing.T
}

func (a Assert) stringEqual(actual string, expected string) {
	if actual != expected {
		a.t.Fatalf("Expected %s, but got %s", expected, actual)
	}
}

func (a Assert) boolEqual(actual bool, expected bool) {
	if actual != expected {
		a.t.Fatalf("Expected %t, but got %t", expected, actual)
	}
}

func (a Assert) intEqual(actual int, expected int) {
	if actual != expected {
		a.t.Fatalf("Expected %d, but got %d", expected, actual)
	}
}

var postsFolder = filepath.Join("testdata", "testposts")

func TestRetrieveMeta(t *testing.T) {
	assert := &Assert{t: t}
	metaPath := filepath.Join(postsFolder, "001-Post1", metaFileName)
	meta, err := retrieveMeta(metaPath)
	if err != nil {
		t.Fatalf("Failed to read %s", metaPath)
	}
	assert.stringEqual(meta.Title, "post 1")
	assert.boolEqual(meta.Visible, true)
	assert.stringEqual(meta.Path, "post1writeup.md")
	assert.intEqual(meta.Id, 1)
}

func TestMetaNotExist(t *testing.T) {
	assert := &Assert{t: t}
	nonExistantMetaPath := filepath.Join(postsFolder, "nonExistantFolder", metaFileName)
	_, err := retrieveMeta(nonExistantMetaPath)
	if err == nil {
		t.Fatalf("%s should not exist", nonExistantMetaPath)
	}
	assert.stringEqual(err.Error(), fmt.Sprintf("open %s: no such file or directory", nonExistantMetaPath))
}
