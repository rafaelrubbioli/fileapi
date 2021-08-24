## File API
This is a graphql api to upload files to cloud storage (only s3 supported for now)

Some examples of queries/mutations on graphql
[explorer](https://rubbioli.com/fileapi/graphql/explorer?query=mutation%20delete%20%7B%0A%20%20delete(id%3A%20%22%22)%0A%7D%0A%0Amutation%20move%20%7B%0A%20%20move(input%3A%20%7Bid%3A%20%22%22%2C%20user%3A%202%2C%20newPath%3A%20%22test%2Facl%2Ffile.txt%22%7D)%20%7B%0A%20%20%20%20id%0A%20%20%7D%0A%7D%0A%0Aquery%20get%20%7B%0A%20%20file(id%3A%20%22Mi90ZXN0L2FjbC9maWxlLnR4dA%3D%3D%22)%20%7B%0A%20%20%20%20id%0A%20%20%20%20name%0A%20%20%20%20path%0A%20%20%20%20user%0A%20%20%20%20fileType%0A%20%20%20%20size%0A%20%20%20%20createdAt%0A%20%20%20%20updatedAt%0A%20%20%20%20downloadURL%0A%20%20%7D%0A%7D%0A%0Aquery%20list%20%7B%0A%20%20listUserFiles(user%3A%201)%20%7B%0A%20%20%20%20id%0A%20%20%20%20name%0A%20%20%20%20path%0A%20%20%20%20user%0A%20%20%20%20fileType%0A%20%20%20%20size%0A%20%20%20%20updatedAt%0A%20%20%20%20downloadURL%0A%20%20%7D%0A%7D%0A&operationName=get)

To upload files use form-files
```
curl -X POST -i https://rubbioli.com/fileapi/graphql \
-F operations='{"query":"mutation($file: Upload!) {  upload(input:{ file: $file   user: 1    path: \"nginx/test/\"  }){    id  }}","variables": { "file": null } }' \  
-F map='{ "0": ["variables.file"] }' \
-F 0=@test.txt
```
