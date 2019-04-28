package gb

import (
	"fmt"
	as "github.com/Adriel-M/gb/assert"
	"github.com/Adriel-M/gb/gb/post"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var postsFolder = filepath.Join("testdata", "testposts")

func TestRetrieveMetaFromFromFolder(t *testing.T) {
	assert := &as.Assert{T: t}
	metaPath := filepath.Join(postsFolder, "001-Post1")
	meta, err := retrieveMetaFromFolder(metaPath)
	if err != nil {
		t.Fatalf("Failed to read %s", metaPath)
	}
	assert.StringEqual(meta.Title, "post 1")
	assert.BoolEqual(meta.Visible, true)
	assert.StringEqual(meta.Path, filepath.Join(postsFolder, "001-Post1", "post1writeup.md"))
	assert.IntEqual(meta.Id, 1)
}

func TestRetrieveMetaFromFromFolderNotExist(t *testing.T) {
	assert := &as.Assert{T: t}
	nonExistantMetaFolderPath := filepath.Join(postsFolder, "nonExistantFolder")
	_, err := retrieveMetaFromFolder(nonExistantMetaFolderPath)
	if err == nil {
		t.Fatalf("%s should not exist", nonExistantMetaFolderPath)
	}
	assert.StringEqual(err.Error(), fmt.Sprintf("open %s: no such file or directory", filepath.Join(nonExistantMetaFolderPath, metaFileName)))
}

func TestRetrieveMetaFromFromFolderPathNotExist(t *testing.T) {
	assert := &as.Assert{T: t}
	folderWithNotBodyPath := filepath.Join(postsFolder, "005-EmptyBodyPath")
	_, err := retrieveMetaFromFolder(folderWithNotBodyPath)
	if err == nil {
		t.Fatalf("body path %s should not exist", folderWithNotBodyPath)
	}
	assert.StringEqual(err.Error(), "path does not exist for this post")
}

func TestRetrieveMetaFromFolderMetaInvalid(t *testing.T) {
	// assert := &as.Assert{T: t}
	pathWithInvalidMeta := filepath.Join(postsFolder, "007-InvalidMeta")
	_, err := retrieveMetaFromFolder(pathWithInvalidMeta)
	if err == nil {
		t.Fatalf("meta should be invalid")
	}
	fmt.Println(err)
}

func TestRetrievePostFromMeta(t *testing.T) {
	assert := &as.Assert{T: t}
	validPostFolder := filepath.Join(postsFolder, "001-Post1")
	validMeta := &post.PostMeta{
		Title:   "some title",
		Visible: true,
		Path:    filepath.Join(validPostFolder, "meta.json"),
		Id:      1,
	}
	post, err := retrievePostFromMeta(validMeta)
	if err != nil {
		t.Fatalf("Valid postMeta, should not fail here")
	}
	assert.StringEqual(validMeta.Title, post.Title)
	assert.IntEqual(validMeta.Id, post.Id)
	assert.PostAddressEqual(nil, post.Next)
	assert.PostAddressEqual(nil, post.Prev)
	body, err := ioutil.ReadFile(validMeta.Path)
	if err != nil {
		t.Fatalf("Valid body path, should not fail here")
	}
	assert.StringEqual(string(body), post.Body)
}

func TestRetrievePostFromMetaBodyNotExist(t *testing.T) {
	assert := &as.Assert{T: t}
	folderWithMissingBodyPath := filepath.Join(postsFolder, "004-MissingBody")
	missingBodyPath := filepath.Join(folderWithMissingBodyPath, "missing.md")
	meta := &post.PostMeta{
		Title:   "some title",
		Visible: true,
		Path:    missingBodyPath,
		Id:      4,
	}
	_, err := retrievePostFromMeta(meta)
	if err == nil {
		t.Fatalf("body file should be missing")
	}
	assert.StringEqual(err.Error(), fmt.Sprintf("open %s: no such file or directory", missingBodyPath))
}

func TestReversePosts(t *testing.T) {
	assert := &as.Assert{T: t}
	numberOfPosts := 10
	posts := make([]*post.Post, numberOfPosts)
	for i := 0; i < numberOfPosts; i++ {
		posts[i] = &post.Post{
			Id: i,
		}
	}
	copyOfPosts := make([]*post.Post, numberOfPosts)
	copy(copyOfPosts, posts)
	reversePosts(posts)
	for i := 0; i < numberOfPosts; i++ {
		reversePost := copyOfPosts[numberOfPosts-1-i]
		forwardPost := posts[i]
		assert.IntEqual(reversePost.Id, forwardPost.Id)
		assert.PostAddressEqual(reversePost, forwardPost)
	}
}

func TestPopulatePrevNext(t *testing.T) {
	assert := &as.Assert{T: t}
	posts := make([]*post.Post, 4)
	for i := 0; i < 4; i++ {
		posts[i] = &post.Post{
			Id: i,
		}
	}
	populatePrevNext(posts)
	for i := 0; i < 4; i++ {
		if i == 0 {
			assert.PostAddressEqual(posts[i].Prev, nil)
		} else {
			assert.PostAddressEqual(posts[i].Prev, posts[i-1])
			assert.IntEqual(posts[i].Prev.Id, i-1)
		}
		if i == 3 {
			assert.PostAddressEqual(posts[i].Next, nil)
		} else {
			assert.PostAddressEqual(posts[i].Next, posts[i+1])
			assert.IntEqual(posts[i].Next.Id, posts[i+1].Id)
		}
	}
}

func TestRetrievePosts(t *testing.T) {
	assert := &as.Assert{T: t}
	postArr, postMap, err := retrievePosts(postsFolder)
	if err != nil {
		t.Fatalf("post location is valid")
	}
	// Should skip metas without valid bodies and not visible
	assert.IntEqual(len(postArr), 3)
	assert.IntEqual(len(postMap), 3)
	expectedTitle := [3]string{"post 1", "another post 3", "valid"}
	expectedIds := [3]int{1, 3, 6}
	for i, post := range postArr {
		assert.StringEqual(post.Title, expectedTitle[i])
		assert.IntEqual(post.Id, expectedIds[i])
	}
	for _, post := range postArr {
		assert.PostAddressEqual(postMap[post.Id], post)
	}
}

func TestRetrievePostsInvalidPath(t *testing.T) {
	assert := &as.Assert{T: t}
	invalidPath := "someInvalid/path"
	_, _, err := retrievePosts(invalidPath)
	if err == nil {
		t.Fatalf("path should be invalid")
	}
	assert.StringEqual(err.Error(), fmt.Sprintf("open %s: no such file or directory", invalidPath))
}
