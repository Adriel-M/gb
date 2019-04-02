package gb

import (
	"fmt"
	"io/ioutil"
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

func (a Assert) postAddressEqual(actual *Post, expected *Post) {
	if actual != expected {
		a.t.Fatalf("Expected %#v, but got %#v", expected, actual)
	}
}

var postsFolder = filepath.Join("testdata", "testposts")

func TestRetrieveMetaFromFromFolder(t *testing.T) {
	assert := &Assert{t: t}
	metaPath := filepath.Join(postsFolder, "001-Post1")
	meta, err := retrieveMetaFromFolder(metaPath)
	if err != nil {
		t.Fatalf("Failed to read %s", metaPath)
	}
	assert.stringEqual(meta.Title, "post 1")
	assert.boolEqual(meta.Visible, true)
	assert.stringEqual(meta.Path, filepath.Join(postsFolder, "001-Post1", "post1writeup.md"))
	assert.intEqual(meta.Id, 1)
}

func TestRetrieveMetaFromFromFolderNotExist(t *testing.T) {
	assert := &Assert{t: t}
	nonExistantMetaFolderPath := filepath.Join(postsFolder, "nonExistantFolder")
	_, err := retrieveMetaFromFolder(nonExistantMetaFolderPath)
	if err == nil {
		t.Fatalf("%s should not exist", nonExistantMetaFolderPath)
	}
	assert.stringEqual(err.Error(), fmt.Sprintf("open %s: no such file or directory", filepath.Join(nonExistantMetaFolderPath, metaFileName)))
}

func TestRetrieveMetaFromFromFolderPathNotExist(t *testing.T) {
	assert := &Assert{t: t}
	folderWithNotBodyPath := filepath.Join(postsFolder, "005-EmptyBodyPath")
	_, err := retrieveMetaFromFolder(folderWithNotBodyPath)
	if err == nil {
		t.Fatalf("body path %s should not exist", folderWithNotBodyPath)
	}
	assert.stringEqual(err.Error(), "path does not exist for this post")
}

func TestRetrievePostFromMeta(t *testing.T) {
	assert := &Assert{t: t}
	validPostFolder := filepath.Join(postsFolder, "001-Post1")
	validMeta := &PostMeta{
		Title:   "some title",
		Visible: true,
		Path:    filepath.Join(validPostFolder, "meta.json"),
		Id:      1,
	}
	post, err := retrievePostFromMeta(validMeta)
	if err != nil {
		t.Fatalf("Valid postMeta, should not fail here")
	}
	assert.stringEqual(validMeta.Title, post.Title)
	assert.intEqual(validMeta.Id, post.Id)
	assert.boolEqual(validMeta.Visible, post.Visible)
	assert.postAddressEqual(nil, post.Next)
	assert.postAddressEqual(nil, post.Prev)
	body, err := ioutil.ReadFile(validMeta.Path)
	if err != nil {
		t.Fatalf("Valid body path, should not fail here")
	}
	assert.stringEqual(string(body), post.Body)
}

func TestRetrievePostFromMetaBodyNotExist(t *testing.T) {
	assert := &Assert{t: t}
	folderWithMissingBodyPath := filepath.Join(postsFolder, "004-MissingBody")
	missingBodyPath := filepath.Join(folderWithMissingBodyPath, "missing.md")
	meta := &PostMeta{
		Title:   "some title",
		Visible: true,
		Path:    missingBodyPath,
		Id:      4,
	}
	_, err := retrievePostFromMeta(meta)
	if err == nil {
		t.Fatalf("body file should be missing")
	}
	assert.stringEqual(err.Error(), fmt.Sprintf("open %s: no such file or directory", missingBodyPath))
}

func TestReversePosts(t *testing.T) {
	assert := &Assert{t: t}
	numberOfPosts := 10
	posts := make([]*Post, numberOfPosts)
	for i := 0; i < numberOfPosts; i++ {
		posts[i] = &Post{
			Id: i,
		}
	}
	copyOfPosts := make([]*Post, numberOfPosts)
	copy(copyOfPosts, posts)
	reversePosts(posts)
	for i := 0; i < numberOfPosts; i++ {
		reversePost := copyOfPosts[numberOfPosts-1-i]
		forwardPost := posts[i]
		assert.intEqual(reversePost.Id, forwardPost.Id)
		assert.postAddressEqual(reversePost, forwardPost)
	}
}
