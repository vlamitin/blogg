# blogg
qraphql tryout 

## Prerequisites
 - go 1.14

## Run
 - `go generate ./...`
 - `go run server.go`
 - open http://localhost:8080/
 ```
mutation createTodo {
  createPost(input:{title:"My Tiyile 1", description:"myDesct"}) {
    id
    title
    publicationDate
    description
  }
}

query getPosts {
  getPosts(limit: 5, offset:0) {
    id
    title
    publicationDate
    description
  }
}
```
