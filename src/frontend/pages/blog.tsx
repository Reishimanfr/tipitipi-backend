import Post from "../components/post";
import { useState, useEffect } from "react";

interface BlogPostDataBodyJson {
  Content: string;
  Created_At: string;
  Edited_At: string;
  ID: number;
  Images: string;
  Title: string;
}

const Blog = () => {
  let limit = 6;
  let offset = 0;
  const [post, setPost] = useState<BlogPostDataBodyJson | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    async function fetchPost() {
      try {
        const response = await fetch("http://localhost:2333/blog/post/8", {
          method: "GET",
            // body: JSON.stringify({
            //         limit: amount,
            //         offset: offset,
            //         sortBy: sortBy
            // })
        });

        if (!response.ok) {
          throw new Error("Network response was not ok");
        }

        const data: BlogPostDataBodyJson = await response.json();
        setPost(data);
      } catch (error) {
        alert("Błąd: " + error);
      } finally {
        setLoading(false);
      }
    }
    fetchPost();
  }, []);

  if (loading) {
    return <div>Loading</div>;
  }
  return (
    <div>
      <h1>Blog</h1>
      {post ? (
        <Post title={post.Title} content={post.Content} />
      ) : (
        <div>No post found</div>
      )}
    </div>
  );
};

export default Blog;
