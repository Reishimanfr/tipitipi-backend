interface BlogPostDataBodyJson {
        Content: string
        Created_At: string 
        Edited_At: string
        ID: number
        Images: string
        Title: string
}

const Blog = () => {
        let limit = 6
        let offset = 0

        async function getPosts(amount: number,  sortBy: "likes" | "newest" | "oldest") {
                const request = await fetch("http://localhost:8080/api/blog/posts/1", {
                        method: "GET",
                        body: JSON.stringify({
                                limit: amount,
                                offset: offset,
                                sortBy: sortBy
                        })
                })


                const data: BlogPostDataBodyJson[] = await request.json()

                console.log(data[0].Title)
        }

    return(
        <div>
            <h1>Blog</h1>
            <input onClick={() => getPosts(6, "likes")}>
                
            </input>
        </div>
    )
}

export default Blog