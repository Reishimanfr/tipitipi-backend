import Post from "../components/post";
import { useState, useEffect } from "react";
import PostSkeleton from "../components/postSkeletonLoading";

interface BlogAttachments {
  id: number;
  url: string;
  filename: string;
  blog_post_id: number;
}

interface BlogPostDataBodyJson {
  content: string;
  created_at: string;
  edited_at: string;
  id: number;
  attachments: BlogAttachments[];
  title: string;
  error?: string;
}

const PostPage = () => {
  const [post, setPost] = useState<BlogPostDataBodyJson | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const url = window.location.href.split("/");
    const ID = url[url.length - 1];

    async function fetchPost() {
      try {
        const response = await fetch(
          `http://localhost:2333/blog/post/${ID}?attachments=true`,
          {
            method: "GET",
          }
        );

        if (!response.ok) {
          throw new Error(response.statusText);
        }

        const data: BlogPostDataBodyJson = await response.json();
        setPost(data);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    }
    fetchPost();
  }, []);

  if (loading) {
    return <div className="globalCss mt-[1%]"><PostSkeleton/></div>
  }
  return (
    <div className="globalCss mt-[1%]">
      {post ? (
        <Post
          title={post.title}
          content={post.content}
          date={post.edited_at}
          id={post.id}
          attachments={post.attachments}
          willBeUsedManyTimes={false}
        />
      ) : (
        <div>No post found</div>
      )}
    </div>
  );
};

export default PostPage;
