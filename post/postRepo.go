package post

import "fmt"

type memoryPostRepo struct {
	// we'll keep a slice of posts
	posts []repoPost
}

type PostRepo interface {
	Add(post repoPost) (repoPost, error)
	GetById(id int) (repoPost, error)
	GetByTitle(title string) (repoPost, error)
}

func NewMemoryPostRepo(posts []repoPost) PostRepo {
	return &memoryPostRepo{posts: posts}
}

func (repo *memoryPostRepo) nextId() int {
	var maxId = 0
	for _, post := range repo.posts {
		if post.Id > maxId {
			maxId = post.Id
		}
	}
	return maxId + 1

}

func (repo *memoryPostRepo) Add(post repoPost) (repoPost, error) {
	if post.Id < 1 {
		post.Id = repo.nextId()
	}
	repo.posts = append(repo.posts, post)
	return post, nil
}

func (repo *memoryPostRepo) GetById(id int) (repoPost, error) {
	for _, post := range repo.posts {
		if post.Id == id {
			return post, nil
		}
	}
	return repoPost{}, fmt.Errorf("could not find id %d", id)
}

func (repo *memoryPostRepo) GetByTitle(title string) (repoPost, error) {
	for _, post := range repo.posts {
		if post.Title == title {
			return post, nil
		}
	}
	return repoPost{}, fmt.Errorf("could not find title %s", title)
}
