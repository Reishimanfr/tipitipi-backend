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
                const request = await fetch("http://localhost:2333/blog/posts/2", {
                        method: "GET"
                        // body: JSON.stringify({
                        //         limit: amount,
                        //         offset: offset,
                        //         sortBy: sortBy
                        // })
                })


                const response: BlogPostDataBodyJson = await request.json()

                console.log(response)
                console.log(request.ok)
        }

    return(
        <div>
            <h1>Blog</h1>
            <button onClick={() => getPosts(6, "likes")}>
                asfasfs
            </button>
        </div>
    )
}

export default Blog