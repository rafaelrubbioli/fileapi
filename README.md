## File API
This is a graphql api to upload and retrieve files on cloud storage (only s3 supported for now)

Some examples of queries/mutations on graphql
[explorer](https://rubbioli.com/fileapi/graphql/explorer?query=mutation%20delete%20%7B%0A%20%20delete(id%3A%20%22%22)%0A%7D%0A%0Amutation%20move%20%7B%0A%20%20move(input%3A%20%7Bid%3A%20%22%22%2C%20user%3A%202%2C%20newPath%3A%20%22test%2Facl%2Ffile.txt%22%7D)%20%7B%0A%20%20%20%20id%0A%20%20%7D%0A%7D%0A%0Aquery%20get%20%7B%0A%20%20file(id%3A%20%22Mi90ZXN0L2FjbC9maWxlLnR4dA%3D%3D%22)%20%7B%0A%20%20%20%20id%0A%20%20%20%20name%0A%20%20%20%20path%0A%20%20%20%20user%0A%20%20%20%20fileType%0A%20%20%20%20size%0A%20%20%20%20createdAt%0A%20%20%20%20updatedAt%0A%20%20%20%20downloadURL%0A%20%20%7D%0A%7D%0A%0Aquery%20list%20%7B%0A%20%20listUserFiles(user%3A%201)%20%7B%0A%20%20%20%20id%0A%20%20%20%20name%0A%20%20%20%20path%0A%20%20%20%20user%0A%20%20%20%20fileType%0A%20%20%20%20size%0A%20%20%20%20updatedAt%0A%20%20%20%20downloadURL%0A%20%20%7D%0A%7D%0A&operationName=get)


## Usage
### Upload
To upload files use form-files as shown
```
curl -X POST -i https://rubbioli.com/fileapi/graphql \
-F operations='{"query":"mutation($file: Upload!) {  upload(input:{ file: $file   user: 1    path: \"nginx/test/\"  }){    id  }}","variables": { "file": null } }' \  
-F map='{ "0": ["variables.file"] }' \
-F 0=@test.txt
```

### Get
It is possible to get a file by `id`. The `id` is a unique string given to every file after the upload.
```graphql
query get {
  file(id: "") {
    id
    name
    path
    user
    fileType
    size
    createdAt
    updatedAt
    downloadURL
  }
}
```

### Delete
Delete takes the file `id` and removes it from the storage permanently.
```graphql
mutation delete {
  delete(id: "")
}
```

### Move
Move takes the file `id` and a new path and moves it, returning the new file.
```graphql
mutation move {
  move(input: {id: "", user: 0, newPath: "test/acl/file.txt"}) {
    id
  }
}
```

### List files
List takes and user and a path prefix (optional) and returns the list of al user files under that path.
```graphql
query list {
  listUserFiles(user: 1) {
    id
    name
    path
    user
    fileType
    size
    updatedAt
    downloadURL
  }
}
```

## Roadmap
- [ ] Parse s3 custom errors (such as not found, bad request)
- [ ] List file tree
- [ ] cmd/worker to process jobs asynchronously with retry (such as deleting a file)
- [ ] Finish unit tests