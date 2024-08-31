### Get post by ID
```sh
/api/blog/post/:id
```
Returns:
```json
{
  "Title": "Some post title (unique)",
  "Content": "Some content",
  "Created_At": 123456789, // Unix timestamp
  "Edited_at": 123456789, // Unix timestamp
  "Images": [] // TODO
}
```
Auth: `true` - TO BE CHANGED

### Get posts in batch
```sh
/api/blog/posts/
```
Returns: array of posts like in `/api/blog/post/:id`<br>
URL options: 
* `sortBy` -> `likes`, `newest`, `oldest`
* `limit` -> int, amount of posts to get
* `offset` -> int, amount of posts to skip
Auth: `true` - TO BE CHANGED

### Delete post
```sh
/api/blog/delete/:id
```
Returns: `nil` or `error` if something went wrong


### Edit post
```sh
/api/blog/edit/:id
```
Returns: `nil` or `error` if something went wrong
