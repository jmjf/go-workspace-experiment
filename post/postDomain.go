package post

import (
	"author"
)

// What we store in the repo (author id, not author details)
type repoPost struct {
	Id       int
	AuthorId int
	Title    string
	Body     string
}

// Get use cases will return entities that make up the aggregate
type Post struct {
	post   repoPost
	author author.Author
}
