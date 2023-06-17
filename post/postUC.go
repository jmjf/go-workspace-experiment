package post

import (
	"errors"
	"fmt"

	"author"
)

type postUC struct {
	postRepo   PostRepo
	authorRepo author.AuthorRepo
}

// Use case separates business logic from the repo so I don't need to change business logic
// if I change the repo implementation.
// Overkill at this scale, but useful for the goals of the scenario I'm testing.
type PostUC interface {
	Add(authorId int, title string, body string) (Post, error)
	GetById(id int) (Post, error)
	GetByTitle(title string) (Post, error)
}

func NewPostUC(postRepo PostRepo, authorRepo author.AuthorRepo) (PostUC, error) {
	if postRepo == nil {
		return postUC{}, errors.New("postRepo may not be nil")
	}
	if authorRepo == nil {
		return postUC{}, errors.New("authorRepo may not be nil")
	}

	return &postUC{postRepo: postRepo, authorRepo: authorRepo}, nil
}

func isValidTitle(title string) bool {
	// RULE: titles must be at least
	return len(title) >= 5
}

func (uc postUC) fillPost(post repoPost) (Post, error) {
	author, err := uc.authorRepo.GetById(post.AuthorId)
	if err != nil {
		return Post{}, err
	}
	return Post{post: post, author: author}, nil
}

func (uc postUC) Add(authorId int, title string, body string) (Post, error) {
	if !isValidTitle(title) {
		return Post{}, fmt.Errorf("post title not valid |%s|", title)
	}
	// if post.Id < 1, repo will get the next id (max + 1), so pass 0
	repoPost, err := uc.postRepo.Add(repoPost{0, authorId, title, body})
	if err != nil {
		return Post{}, err
	}
	return uc.fillPost(repoPost)
}

func (uc postUC) GetById(id int) (Post, error) {
	// RULE: ids must be at least 1
	if id < 1 {
		return Post{}, fmt.Errorf("post id must be at least 1 -- id: %d", id)
	}

	post, err := uc.postRepo.GetById(id)
	if err != nil {
		return Post{}, err
	}

	return uc.fillPost(post)
}

func (uc postUC) GetByTitle(title string) (Post, error) {
	if !isValidTitle(title) {
		return Post{}, fmt.Errorf("post title not valid |%s|", title)
	}

	post, err := uc.postRepo.GetByTitle(title)
	if err != nil {
		return Post{}, err
	}

	return uc.fillPost(post)
}
