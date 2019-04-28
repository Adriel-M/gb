package gb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Adriel-M/gb/gb/post"
	"io/ioutil"
	"log"
	"path/filepath"
)

var metaFileName = "meta.json"

func retrieveMetaFromFolder(folderPath string) (*post.PostMeta, error) {
	metaFilePath := filepath.Join(folderPath, metaFileName)
	jsonFile, err := ioutil.ReadFile(metaFilePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var postMeta post.PostMeta
	err = json.Unmarshal(jsonFile, &postMeta)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if postMeta.Path == "" {
		log.Println("Path is not specified!")
		return nil, errors.New("path does not exist for this post")
	}
	pathWithFolder := filepath.Join(folderPath, postMeta.Path)
	postMeta.Path = pathWithFolder
	return &postMeta, nil
}

// Path to the meta file
func retrievePostFromMeta(postMeta *post.PostMeta) (*post.Post, error) {
	bodyFile, err := ioutil.ReadFile(postMeta.Path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	postBody := string(bodyFile)

	return &post.Post{
		Title: postMeta.Title,
		Id:    postMeta.Id,
		Body:  postBody,
	}, nil
}

func reversePosts(posts []*post.Post) {
	for l, r := 0, len(posts)-1; l < r; l, r = l+1, r-1 {
		posts[l], posts[r] = posts[r], posts[l]
	}
}

func populatePrevNext(posts []*post.Post) {
	numOfPosts := len(posts)
	for i, post := range posts {
		if i > 0 {
			post.Prev = posts[i-1]
		}
		if i < numOfPosts-1 {
			post.Next = posts[i+1]
		}
	}
}

func retrievePosts(postsLocation string) ([]*post.Post, map[int]*post.Post, error) {
	var posts []*post.Post
	idToPost := make(map[int]*post.Post)

	folders, err := ioutil.ReadDir(postsLocation)
	if err != nil {
		log.Println(err)
		return nil, nil, err
	}

	for _, f := range folders {
		if !f.IsDir() {
			continue
		}
		folderName := f.Name()
		postMetaFolder := filepath.Join(postsLocation, folderName)
		postMeta, err := retrieveMetaFromFolder(postMetaFolder)
		if err != nil {
			continue
		}
		if !postMeta.Visible {
			continue
		}
		newPost, err := retrievePostFromMeta(postMeta)
		if err != nil {
			continue
		}
		idToPost[newPost.Id] = newPost
		posts = append(posts, newPost)
	}
	populatePrevNext(posts)
	return posts, idToPost, nil
}

type Server struct {
	Path string
	Port int
}

func (s Server) start() {
	postsList, postsMap, err := retrievePosts(s.Path)
	if err != nil {
		// terminate here
	}
	reversePosts(postsList)
	fmt.Println(postsMap)
}
