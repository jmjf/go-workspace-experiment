# Notes

## What this is

Some quick experiments with workspaces to see how they work and see if they can solve the problems I see.

## Starting assumptions

* I have a set of code with several executable modules--an HTTP server, a couple of services that listen to Kafka topics, etc.
* Those executable modules want to share data structure definitions.
* The modules are part of a bounded context where read-sharing of a data store is allowed, so data access code is shared.
  * I'm not building a data store here, but may simulate one with a simple array and read/write methods.
* The modules may have other shared code, TBD.
  * For testing purposes, this will probably be a simple function.
* It makes sense to group the code in a single repo because the definitions and code are used within the bounded context only, not to be shared outside it.
  * A utility that's used across several bounded contexts should be in its own repo so it's easier to share. This isn't that case.
  * Or if I want and can justify a corporate monorepo. That isn't my goal, but this approach might work for that case too.

## Questions

Using workspaces, can I:

* Share data structure definitions between modules?
* Share data access code between modules?
* Share functions between modules?
* Build and deploy executables independently?
  * I'm not setting up CI/CD or a deployment platform for this, so build independently is adequate for this test.

## References

* [Golang workspaces tutorial](https://go.dev/doc/tutorial/workspaces)
* [LogRocket overview post](https://blog.logrocket.com/go-workspaces-multi-module-local-development/)
* [Workspace example repo](https://github.com/xmlking/go-workspace) that suggest this is doable and gives some hints
* [Clean architecture example repo](https://github.com/bxcodec/go-clean-arch)

## Setup

Let's build a simple blog with author and posts. Author pieces will be in `./author`, post pieces in `./post`.

Based on the LogRocket post, I need to run the following commands. (The information is also in the official tutorial, but LR's example is more direct.)

```bash
go work init
mkdir cmd && cd cmd
mkdir authorsSvc && cd authorsSvc && go mod init authorsSvc
cd ..
mkdir postsSvc && cd postsSvc && go mod init postsSvc
cd ../..
mkdir author && cd author && go mod init author
cd ..
mkdir post && cd post && go mod init post
cd ..
go work use ./cmd/authorsSvc ./cmd/postsSvc ./author ./post
```

Now I have a workspace (`go.work`) in the project root and three modules (`cmd`, `author`, `post`).

## Can I share strut definitions, data access code, functions?

If the project was more complex, I'd have subdirectories in each module (`domain`, `repo`, `usecase`, etc.). This is a simple test application, so I'll put them in files (`authorDomain`, `authorRepo`, etc.).

After some coding and sorting through syntax, I have code that isn't failing basic checks.

While working, I asked if it was possible to create a repo interface with methods available to the module only so I could protect write methods from other modules. For example, the author repo has a method to add a new author to the database. Only the author module has write permissions to the database. Can I protect it so I don't accidentally add an author in a posts-centric use case but allow posts to read author data?

I found that, if I name private methods with a lowercase first character, the method seems to be protected (error flagged in VS Code). Perfect and consistent with the Golang model of lowercase names are not exported.

```golang
// in authorRepo.go

type AuthorRepo interface {
   add(author Author) (Author, error)
   GetById(id int) (Author, error)
}

// in authorUC.go

func NewAuthorUC(repo AuthorRepo) (AuthorUC, error) {

   // snip
   // following line added for testing
   repo.add(author.Author{}) 
   // no error flagged in VS Code

   return &postUC{postRepo: postRepo, authorRepo: authorRepo}, nil
}

// in postUC.go

import "author"

func NewPostUC(postRepo PostRepo, authorRepo author.AuthorRepo) (PostUC, error) {

   // snip
   // following line added for testing
   authorRepo.add(author.Author{}) 
   // authorRepo.add undefined (type author.AuthorRepo has no field or method add)

   return &postUC{postRepo: postRepo, authorRepo: authorRepo}, nil
}
```

**Result:** Yes.

**COMMIT: DOCS: explore and document how to share code between modules in a workspace**
