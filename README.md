# API entry URL
```sh
http://localhost:8080/api
```

# Blog operations
## Create blog post
```sh
/api/blog/create
```
Example request body:
```json
{
  "title": "Some title for the post",
  "content": "Some content of the blog post",
  "images": [
      { "id": 1, "path": "path/to/image", "settings": [] }
  ]
}
```

## Edit blog post
```sh
/api/blog/edit/:id
```
Example request body:
```json
{
  "title": "Some title for the post",
  "content": "Some content of the blog post",
  "images": [
      { "id": 1, "path": "path/to/image", "settings": [] }
  ]
}
```
## Delete blog post
```sh
/api/blog/delete/:id
```

# Server heartbeat
```sh
/api/heartbeat
```
