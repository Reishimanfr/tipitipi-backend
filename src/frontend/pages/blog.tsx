import Post from "../components/post";
import { useState, useEffect } from "react";
import { Link } from "react-router-dom";

interface BlogPostDataBodyJson {
  Content: string;
  Created_At: string;
  Edited_At: string;
  ID: number;
  Title: string;
}

const Blog = () => {
  const limit = 6
  const [offset, setOffset] = useState(0);
  const [sortBy,setSortBy] = useState<"newest" | "oldest" | "likes">("newest")

  const [posts, setPosts] = useState<Array<BlogPostDataBodyJson>>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [isMore , setIsMore] = useState(true)

  useEffect(() => {
    async function fetchPost() {
      try {
        const response = await fetch(
          `http://localhost:2333/blog/posts?offset=${offset}&limit=${limit}&sort=${sortBy}`,
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
        alert(error);
      } finally {
        setLoading(false);
      }
    }
    fetchPost();
  }, [offset,sortBy]);


  useEffect(() => {
    function handleScroll() {
      const scrollTop = document.documentElement.scrollTop;
      const scrollHeight = document.documentElement.scrollHeight;
      const clientHeight = window.innerHeight;

      if (scrollTop + clientHeight >= scrollHeight && isMore) {
        if(offset + 6> posts.length) {
          setIsMore(false)
        }
        else {
          setOffset((prevOffset) => prevOffset + 6 );
        }
        
      }
    }

    window.addEventListener("scroll" , handleScroll);
    return () => {
      window.removeEventListener("scroll" , handleScroll)
      }
  });


  if (loading) {
    return <div>Loading</div>;
  }
  return (
    <div className="globalCss">
      <h1 className="text-3xl mt-5">Blog</h1>

      <label htmlFor="sorting">Sortowanie</label>
      <select name="sorts" id="sorting" onChange={(e) => setSortBy(e.target.value as "newest" | "oldest" | "likes")}>
        <option value="newest">Najnowsze</option>
        <option value="oldest">Najstarsze</option>
        <option value="likes">Najwięcej polubień</option>
      </select>
      {posts ? (
        posts.map((post) => {
          return <Link key={post.ID} to={`/blog/${post.ID}`}><Post id={post.ID} title={post.Title} content={post.Content} date={post.Edited_At}/></Link>;
        })
      ) : (
        <div>No post found</div>
      )}
    </div>
  );
};

export default Blog;
