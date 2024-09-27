import Post from "../components/post";
import { useState, useEffect } from "react";


interface BlogAttachments {
  ID: number;
  BlogPostID: number;
  Path: string;
  Filename: string;
}

interface BlogPostDataBodyJson {
  Content: string;
  Created_At: string;
  Edited_At: string;
  ID: number;
  Attachments: BlogAttachments[];
  Title: string;
}


const PostPage = () => {
  const [post, setPost] = useState<BlogPostDataBodyJson | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const url = window.location.href.split("/");
    const ID = url[url.length - 1]

    async function fetchPost() {
      try {
        const response = await fetch(`http://localhost:2333/blog/post/${ID}?images=true`, {
          method: "GET"
        });

        if (!response.ok) {
          throw new Error(response.statusText);
        }

        const data: BlogPostDataBodyJson = await response.json();
        setPost(data);
      } catch (error) {
        alert(error);
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
    <div className="globalCss mt-[1%]">
      {post ? (
        <Post title={post.Title} content={post.Content} date={post.Edited_At} id={post.ID} attachments={post.Attachments} willBeUsedManyTimes={false}/>
      ) : (
        <div>No post found</div>
      )}
    </div>
  );
};

export default PostPage;
