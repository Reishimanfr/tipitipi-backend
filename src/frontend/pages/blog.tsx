import { useEffect, useState } from "react"
import Post from "../components/post"
import PostSkeleton from "../components/postSkeletonLoading"
import { API_URL } from '../functions/global'
import { BlogPostDataBodyJson } from "../functions/interfaces"

const Blog = () => {
  const limit = 6;
  const [offset, setOffset] = useState(0);
  const [sortBy, setSortBy] = useState<"newest" | "oldest" | "likes">("newest");

  const [posts, setPosts] = useState<Array<BlogPostDataBodyJson>>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [isMore, setIsMore] = useState(true);

  useEffect(() => {
    async function fetchPost() {
      try {
        const response = await fetch(
          `${API_URL}/blog/posts?offset=${offset}&limit=${limit}&sort=${sortBy}`,
          {
            method: "GET",
          }
        );
        if (!response.ok) {
          throw new Error(response.statusText);
        }
  

        const data: Array<BlogPostDataBodyJson> = await response.json();
        setPosts((prevPosts) => prevPosts?.concat(data));
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    }
    fetchPost();
  }, [offset, sortBy]);

  useEffect(() => {
    setPosts([]);
    setOffset(0);
  }, [sortBy]);

  useEffect(() => {
    function handleScroll() {
      const scrollTop = document.documentElement.scrollTop;
      const scrollHeight = document.documentElement.scrollHeight;
      const clientHeight = window.innerHeight;

      if (scrollTop + clientHeight >= scrollHeight && isMore) {
        if (offset + 6 > posts.length) {
          setIsMore(false);
        } else {
          setOffset((prevOffset) => prevOffset + 6);
        }
      }
    }

    window.addEventListener("scroll", handleScroll);
    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  });

  if (loading) {
    return (
      <div className="globalCss">
        <h1 className="text-3xl mt-5">Blog</h1>
        <PostSkeleton />
        <PostSkeleton />
        <PostSkeleton />
        <PostSkeleton />
      </div>
    );
  }
  if(posts.length == 0){
    return (
      <div>
        <h1 className="text-5xl m-12">Brak postów</h1>
      </div>
    )
  }
  return (
    <div className="globalCss">
      <h1 className="text-3xl mt-5">Blog</h1>

      <label htmlFor="sorting" className="mr-2">Sortowanie</label>
      <select
        name="sorts"
        id="sorting"
        onChange={(e) =>
          setSortBy(e.target.value as "newest" | "oldest" | "likes")
        }
      >
        <option value="newest">Najnowsze</option>
        <option value="oldest">Najstarsze</option>
        {/* <option value="likes">Najwięcej polubień</option> */}
      </select>
      {posts ? (
        posts.map((post, index) => {
          return (
            <div key={index} className="mt-12">
              <Post
                id={post.id}
                content={post.content}
                title={post.title}
                date={post.edited_at}
                willBeUsedManyTimes={true}
              />
            </div>
          );
        })
      ) : (
        <div>No post found</div>
      )}
    </div>
  );
};

export default Blog;
