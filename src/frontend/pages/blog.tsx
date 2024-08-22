interface BlogPostDataBodyJson {
        Content: string
        Created_At: string 
        Edited_At: string
        ID: number
        Images: string
        Title: string
}

const Blog = () => {
        let posts = 6
        let startIdx = 0

        async function getPosts(amount: number,  sortBy: "likes" | "newest" | "oldest") {
                const request = await fetch("http://localhost:8080/api/blog/getPostById/1", {
                        method: "GET"
                })


                const data: BlogPostDataBodyJson = await request.json()

                console.log(data.Title)
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