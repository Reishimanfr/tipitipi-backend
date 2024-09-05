## Public endpoints (no authorization header)
* HEAD `/heartbeat` -> Checks if the server is alive
* POST `/admin/login` -> Logs in and returns a JWT token
* GET  `/blog/post/:id` -> Gets a post at a specified ID
* GET  `/blog/posts` -> Gets posts in batch

## Example requests and responses
### `HEAD /heartbeat`
Request:
```ts
await fetch("http://localhost:2333/heartbeat", {
  method: "HEAD"
})
```
Codes:
| Code | Meaning |
| ---- | -------- |
| 200  | Server is alive |

### `POST /admin/login`
Request:
```ts
await fetch("http://localhost:2333/admin/login", {
  method: "POST",
  body: JSON.stringify({
      username: "admin",
      password: "amazing_password"
  })
})
```
Response:
```json
{
  "token": "JWT token",
  "error?": "Direct error if something went wrong",
  "message?": "A more human readable error"
}
```
Codes:
| Code | Meaning |
| ---- | -------- |
| 200  | Everything went okay |
| 422  | The JSON data was malformed or invalid |
| 401  | Invalid credentials were provided |
| 500  | Something went wrong within the backend server |

### `GET /blog/post/:id`
Request:
```ts
await fetch("http://localhost:2333/blog/posts", {
  method: "GET",
  body: JSON.stringify({
      sortBy: "likes/newest/oldest",
      limit: 5,
      offset: 0
  })
})
```
Response:
```json
[
  {
    "ID": 1,
    "Title": "Some title",
    "Content": "Some content",
    "Created_At": 1234567890, // Unix timestamp
    "Edited_At": 1234567890, // Unix timestamp
    "Images": [
      { "ID": 1, "Path": "/some/path/to/the/image", "Placement": "left/center/right" },
      ...
    ]
  }
  ...
]
```

### `GET /blog/post/:id`
Request:
```ts
await fetch("http://localhost:2333/blog/post/3", {
  method: "GET"
})
```
Response:
```json
{
	"ID": 1,
	"Title": "Some title",
	"Content": "Some content",
	"Created_At": 1725014085,
	"Edited_At": 1725014085,
	"Images": [
      { "ID": 1, "Path": "/some/path/to/the/image", "Placement": "left/center/right" },
      ...
    ]
}
```

## Protected endpoints (required JWT authentication
* DELETE `/blog/post/:id` -> Deletes a selected post
* POST `/blog/post`       -> Creates a new post
* PATCH `/blog/post/:id`  -> Updates an existing post
* PATCH `/admin/account`  -> Changes admin credentials (username and password)

### `DELETE /blog/post/:id`
Request:
```ts
await fetch("http://localhost:2333/blog/post/3", {
  method: "DELETE",
  headers: {
    Authorization: "JWT token here"
  }
})
```
Codes:
| Code | Meaning |
| ---- | -------- |
| 200  | Record deleted successfully |
| 400  | Request was invalid (for example ID was a string) |
| 401  | Authorization failed (likely due to an invalid JWT token) |
| 403  | User doesn't have permission to do this action | 
| 404  | Post record under that ID wasn't found |

### `POST /blog/post`
Request:
```ts
await fetch("http://localhost:2333/blog/post/3", {
  method: "POST",
  headers: {
    Authorization: "JWT token here"
  },
  body: JSON.stringify({
    "ID": 1,
    "Title": "Some title",
    "Content": "Some content",
    "Images": [
      { "ID": 1, "Path": "/some/path/to/the/image", "Placement": "left/center/right" },
      ...
    ]
  })
})
```
Response:
```
(The same data you provided + Created_At and Edited_At in unix)
```
Codes:
| Code | Meaning |
| ---- | -------- |
| 200  | Post was created successfully |
| 400  | Request was invalid (for example title was missing) |
| 401  | Authorization failed (likely due to an invalid JWT token) |
| 403  | User doesn't have permission to do this action | 
| 409  | Post with this title already exists, conflict! |
