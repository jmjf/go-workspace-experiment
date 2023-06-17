package main

import (
	"author"
	"fmt"
)

var authors = []author.Author{
	{Id: 1, Name: "Joe Jones"},
	{Id: 2, Name: "Sue Sutherland"}}

func main() {

	repo := author.NewMemoryAuthorRepo(authors)
	uc, err := author.NewAuthorUC(repo)
	if err != nil {
		fmt.Println("Error creating uc")
		panic(err)
	}

	author, err := uc.Add("Mary Lamb")
	if err != nil {
		fmt.Println("Error creating newAuthor")
		panic(err)
	}

	fmt.Println("New author", author)

	author, err = uc.GetById(1)
	if err != nil {
		fmt.Println("Error getting author1")
		panic(err)
	}
	fmt.Println("Author 1", author)

	author, err = uc.GetById(3)
	if err != nil {
		fmt.Println("Error getting author3")
		panic(err)
	}
	fmt.Println("Author 3", author)
}
