package main

import (
	"author"
	"fmt"
	"post"
)

var authors = []author.Author{
	{Id: 1, Name: "Joe Jones"},
	{Id: 2, Name: "Sue Sutherland"}}

var posts = []post.Post{}

func main() {

	postRepo := post.NewMemoryPostRepo(posts)
	authorRepo := author.NewMemoryAuthorRepo(authors)

	uc, err := post.NewPostUC(postRepo, authorRepo)
	if err != nil {
		fmt.Println("Error creating uc")
		panic(err)
	}

	pst, err := uc.Add(1, "Test Post 1.1", "This is the first test post.")
	if err != nil {
		fmt.Println("Error creating new post 1.1")
		panic(err)
	}

	fmt.Println("New post", pst)

	pst, err = uc.Add(1, "Test Post 1.2", "This is the second test post.")
	if err != nil {
		fmt.Println("Error creating new post 1.2")
		panic(err)
	}

	fmt.Println("New post", pst)

	pst, err = uc.Add(2, "Test Post 2.1", "This is the third test post.")
	if err != nil {
		fmt.Println("Error creating new post 2.1")
		panic(err)
	}

	fmt.Println("New post", pst)

	pst, err = uc.GetById(1)
	if err != nil {
		fmt.Println("Error getting post 1")
		panic(err)
	}
	fmt.Println("Post 1", pst)

	pst, err = uc.GetById(3)
	if err != nil {
		fmt.Println("Error getting post 3")
		panic(err)
	}
	fmt.Println("Post 3", pst)

	pst, err = uc.GetByTitle("Test Post 1.2")
	if err != nil {
		fmt.Println("Error getting post by title: Test Post 1.2")
		panic(err)
	}
	fmt.Println("Post by title", pst)
}
