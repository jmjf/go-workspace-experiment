package author

import "fmt"

// TODO:
// Can I create a public author repo and a private author repo?
// The private repo has write functions available to the author package.
// The public repo is read only.
type AuthorRepo interface {
	add(author Author) (Author, error)
	GetById(id int) (Author, error)
}

type memoryAuthorRepo struct {
	// we'll keep a slice of authors
	authors []Author
}

func NewMemoryAuthorRepo(authors []Author) AuthorRepo {
	return &memoryAuthorRepo{authors: authors}
}

func (repo *memoryAuthorRepo) nextId() int {
	var maxId = 0
	for _, author := range repo.authors {
		if author.Id > maxId {
			maxId = author.Id
		}
	}
	return maxId + 1

}

func (repo *memoryAuthorRepo) add(author Author) (Author, error) {
	if author.Id < 1 {
		author.Id = repo.nextId()
	}
	repo.authors = append(repo.authors, author)
	return author, nil
}

func (repo *memoryAuthorRepo) GetById(id int) (Author, error) {
	for _, author := range repo.authors {
		if author.Id == id {
			return author, nil
		}
	}
	return Author{}, fmt.Errorf("could not find id %d", id)
}
