package gb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

type PostMeta struct {
	Title   string `json:title`
	Visible bool   `json:visible`
	Path    string `json:path`
	Id      int    `json:id`
}

type Post struct {
	Title   string
	Id      int
	Visible bool
	Body    string
	Next    *Post
	Prev    *Post
}

var metaFileName = "meta.json"

func retrieveMeta(path string) (*PostMeta, error) {
	jsonFile, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var postMeta PostMeta
	err = json.Unmarshal(jsonFile, &postMeta)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &postMeta, nil
}

// Path to the meta file
func retrievePost(path string) (*Post, error) {
	postMeta, err := retrieveMeta(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bodyFile, err := ioutil.ReadFile(postMeta.Path)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	postBody := string(bodyFile)

	return &Post{
		Title:   postMeta.Title,
		Visible: postMeta.Visible,
		Body:    postBody,
	}, nil
}

func reversePosts(posts []*Post) {
	for l, r := 0, len(posts)-1; l < r; l, r = l+1, r-1 {
		posts[l], posts[r] = posts[r], posts[l]
	}
}

func retrievePosts(postsLocation string) ([]*Post, map[int]*Post, error) {
	var posts []*Post
	idToPost := make(map[int]*Post)

	folders, err := ioutil.ReadDir(postsLocation)
	if err != nil {
		log.Fatal(err)
		return nil, nil, err
	}

	for _, f := range folders {
		if !f.IsDir() {
			continue
		}
		folderName := f.Name()
		postMetaLocation := filepath.Join(postsLocation, folderName+metaFileName)
		newPost, err := retrievePost(postMetaLocation)
		if err != nil {
			continue
		}
		idToPost[newPost.Id] = newPost
		posts = append(posts, newPost)
	}

	numberOfPosts := len(posts)
	for i, post := range posts {
		if i > 0 {
			post.Prev = posts[i-1]
		}
		if i < numberOfPosts-1 {
			post.Next = posts[i+1]
		}
	}
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
