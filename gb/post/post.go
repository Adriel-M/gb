package post

type PostMeta struct {
	Title   string `json:title`
	Visible bool   `json:visible`
	Path    string `json:path`
	Id      int    `json:id`
}

type Post struct {
	Title string
	Id    int
	Body  string
	Next  *Post
	Prev  *Post
}
