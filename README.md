# Meaning of glyphs
* ⚙️ -> Uses URL queries
* 🔒 -> Requires JWT auth
* 💠 -> Uses JSON for requests
* 📝 -> Uses multipart form for requests

> [!WARNING]
> Endpoints using `images` instead of `attachments` as a query parameter in URLs is subject to change, this is an error

<details>
<summary>Managing blog posts</summary>

## ⚙️ `GET /blog/post/:id`
> Returns the blog post under the specified ID<br>
### Example response
```json
// If something goes wrong:
{
  "message": "An easier to understand error message",
  "error": "The actual error message"
}

// If the request is successful
// /blog/post/1
// (Get the post under ID=1 without attachments)
{
  "ID": 1,
  "Created_At": 123456789,
  "Edited_At": 123456789,
  "Title": "Some title",
  "Content": "<p>Some content as HTML code</p>",
  "Attachments": null
}

// Example response with images=true
// If there are no attachments the array will just be an empty one
// /blog/post/1?images=true
// (Get the post under ID=1 with attachments)
{
  ...
  "Attachments": [
      {
          "ID": 1,
          "Filename": "test-image.png",
          "Path": "path/to/test-image.png",
          "BlogPostID": 1
      }
      ...
  ]
}
```

## ⚙️ `GET /blog/posts`
> Returns multiple blog posts<br>
### Example responses
```json
// If something goes wrong
{
  "message": "An easier to understand error message",
  "error": "The actual error message"
}

// If the request is successful
// /blog/posts?offset=0&limit=3&sort=oldest
// (Get the first 3 posts counting from the first one and sort by oldest to newest)
[
  {
    "ID": 1,
    "Created_At": 1726070770,
    "Edited_At": 1726070770,
    "Title": "Some test title",
    "Content": "<p>Some content as HTML code</p>",
    "Attachments": null
  },
  {
    "ID": 2,
    "Created_At": 1726070834,
    "Edited_At": 1726070834,
    "Title": "Some other test title",
    "Content": "<p>Some content as HTML code</p>",
    "Attachments": null
  },
  ...
]

// If images is set to true the request will look the same as for GET /blog/post/:id just that it's in an array of objects
```

## 🔒 `DELETE /blog/post/:id`
> Deletes a post under some ID<br>
### Example responses
```json
// If something goes wrong during the request
{
  "message": "An easier to understand error message",
  "error": "The actual error message"
}

// If the request is successful
// /blog/post/1
// (Deletes the post under ID=1)
{
  "message": "Post and its attachments deleted successfully",
}
```
> [!TIP]
> The success message will not change. It will always be the one shown here

## 


## 🔒 📝 `PATCH /blog/post/:id`
> Updates a post under some ID with new data<br>
### Example responses
> [!TIP]
> Updating blog posts works the same way as creating them, you just have to pass in the blog post struct with stuff changed
```json
// If something goes wrong during the request
{
  "message": "An easier to understand error message",
  "error": "The actual error message"      
}

// If everything goes well the response will look the same as it does when GETing posts but you don't have to set images=true to get the attachments (if you send any)
// /blog/post/1
```

## 🔒 📝 `POST /blog/post`
> Creates a new blog post<br>
### Example responses
```json
// If something goes wrong during the request
{
  "message": "An easier to understand error message",
  "error": "The actual error message"      
}

// If the request is successful
// /blog/post
{
        "message": "Post added successfully"
}
```
> [!TIP]
> The successful message will also always stay the same here.
</details>