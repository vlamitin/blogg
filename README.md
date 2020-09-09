# blogg
qraphql tryout 

## Prerequisites
 - go 1.14

## Run
 - `go generate ./...`
 - `go run server.go`
 - open http://localhost:8080/
 
## Features:
### Posts crud
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

### WS subscription
(in chrome browser console):
```
let ws = new WebSocket("ws://localhost:8080/subscribe");
ws.onopen = function(evt) {
    console.log("OPEN");
    ws.send("subscribe")
}
ws.onclose = function(evt) {
    console.log("CLOSE");
    ws = null;
}
ws.onmessage = function(evt) {
    console.log("RESPONSE: " + evt.data);
}
ws.onerror = function(evt) {
    console.log("ERROR: " + evt.data);
}
```
