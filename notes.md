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
mkdir author && cd author && go mod init author
cd ..
mkdir post && cd post && go mod init post
cd ..
go work use ./author ./post
```

Now I have a workspace (`go.work`) in the project root and two modules (`author`, `post`).

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

## Can I build separate executables?

Let's build a couple of simple modules in `cmd`. In the real world, this might be separate REST APIs or message producers/consumers that are part of the bounded context. Here, I just want to setup an `AuthorUC` and a `PostUC` in separate apps, add a couple of items to the repos, and run a `GetById()`

```bash
mkdir cmd && cd cmd
mkdir authorsSvc && cd authorsSvc && go mod init main
cd ..
mkdir postsSvc && cd postsSvc && go mod init main
cd ../..
go work use ./cmd/authorsSvc ./cmd/postsSvc
```

Now both services are set up in the workspace file.

In `authorsSvc`, I build a `main()` function to perform some basic tests.

`cd cmd/authorsSvc && go run main.go` gets the expected output. Note that the `go.mod` calls this module `authorsSvc` but `main.go` is `package main`.

```
New author {3 Mary Lamb}
Author 1 {1 Joe Jones}
Author 3 {3 Mary Lamb}
```

In `postsSvc`, I build a `main()` similar function. I found a few errors in the `post` module along the way and fixed them.

* The use case `Add()` returned a `repoPost` instead of a `Post`.
* The repo `New...()` accepted a `repoPost` instead of a `Post`.

`cd cmd/postSvc && go run main.go` gets the expected output.

```
New post {{1 1 Test Post 1.1 This is the first test post.} {1 Joe Jones}}
New post {{2 1 Test Post 1.2 This is the second test post.} {1 Joe Jones}}
New post {{3 2 Test Post 2.1 This is the third test post.} {2 Sue Sutherland}}
Post 1 {{1 1 Test Post 1.1 This is the first test post.} {1 Joe Jones}}
Post 3 {{3 2 Test Post 2.1 This is the third test post.} {2 Sue Sutherland}}
Post by title {{2 1 Test Post 1.2 This is the second test post.} {1 Joe Jones}}
```

So far, so good. I can run two separate services using shared code and apply the concept of private and public methods.

**COMMIT: FEAT: write runnable code that uses the modules**

## Testing build

I can build each service individually and get the expected results from each.

```bash
cd cmd/authorsSvc && go build 
./authorsSvc
cd ../..
cd cmd/postsSvc && go build 
./postsSvc
cd ../..
```

I can also build from the project root individually and get executables in the project root.

```bash
go build authorsSvc
go build postsSvc
```

But trying to build both at once produces no output. According to `go help build`:

> When compiling multiple packages or a single non-main package, build compiles the packages but discards the resulting object, serving only as a check that the packages can be built.

The help docs say `-o [destination]` will write output. `go build -o . authorsSvc postsSvc` gets two executables in the project root. They produce the same output as the tests above.

I also did a quick test in `postsSvc` to ensure build fails if I try to call `authorRepo.add()`. See comments starting in line 71 of `cmd/postsSvc/main.go`.

And I realized `.gitignore` is ignoring the `go.work` file, so I uncommented it. In theory, you'd exclude `go.work`, but it's needed in this repo for demonstration.

**COMMIT: DOCS: explore and document how build works with workspaces**
