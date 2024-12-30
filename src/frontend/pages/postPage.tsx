import { useEffect, useState } from "react"
import Post from "../components/post"
import PostSkeleton from "../components/postSkeletonLoading"
import { API_URL } from '../functions/global'
import { BlogPostDataBodyJson } from "../functions/interfaces"

const PostPage = () => {
  const [post, setPost] = useState<BlogPostDataBodyJson | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    const url = window.location.href.split("/");
    const ID = url[url.length - 1];

    async function fetchPost() {
      try {
        const response = await fetch(
          `${API_URL}/blog/post/${ID}?files=true`,
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
    return <div className="mt-[5%]"><PostSkeleton/></div>
  }
  return (
    <div className="globalCss mt-[1%]">
      {post ? (
        <Post
          title={post.title}
          content={post.content}
          date={post.edited_at}
          id={post.id}
          attachments={post.files}
          willBeUsedManyTimes={false}
        />
      ) : (
        <div>No post found</div>
      )}
    </div>
  );
};

export default PostPage;
