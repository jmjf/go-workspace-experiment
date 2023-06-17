package author

import (
	"errors"
	"fmt"
)

type authorUC struct {
	repo AuthorRepo
}

// Use case separates business logic from the repo so I don't need to change business logic
// if I change the repo implementation.
// Overkill at this scale, but useful for the goals of the scenario I'm testing.
type AuthorUC interface {
	Add(name string) (Author, error)
	GetById(id int) (Author, error)
}

func NewAuthorUC(repo AuthorRepo) (AuthorUC, error) {
	if repo == nil {
		return authorUC{}, errors.New("repo may not be nil")
	}

	return &authorUC{repo: repo}, nil
}

func (uc authorUC) Add(name string) (Author, error) {
	// RULE: author names must be at least 3 characters long
	if len(name) < 3 {
		return Author{}, fmt.Errorf("author name must be at least 3 characters long -- name: |%s|", name)
	}
	// if author.Id < 1, repo will get the next id (max + 1), so pass 0
	return uc.repo.add(Author{0, name})
}

func (uc authorUC) GetById(id int) (Author, error) {
	// RULE: ids must be at least 1
	if id < 1 {
		return Author{}, fmt.Errorf("author id must be at least 1 -- id: %d", id)
	}

	return uc.repo.GetById(id)
}
