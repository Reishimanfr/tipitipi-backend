import { useState, useEffect } from "react";
import validateToken from "../../../components/validate";
import Unauthorized from "../../errorPages/unauthorized";
import Post from "../../../components/post";

interface BlogPostDataBodyJson {
  Content: string;
  Created_At: string;
  Edited_At: string;
  ID: number;
  Title: string;
}

async function fetchPosts(
  setPosts: React.Dispatch<React.SetStateAction<BlogPostDataBodyJson[]>>
) {
  try {
    const response = await fetch(`http://localhost:2333/blog/posts?offset=0`, {
      method: "GET",
    });
    if (!response.ok) {
      throw new Error(response.statusText);
    }

    const data: Array<BlogPostDataBodyJson> = await response.json();
    setPosts((prevPosts) => prevPosts?.concat(data));
  } catch (error) {
    alert(error);
  }
}

const PostEditing = () => {
  const [selectedPost, setSelectedPost] = useState<number>();
  const [posts, setPosts] = useState<Array<BlogPostDataBodyJson>>([]);
  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);

  useEffect(() => {
    const ValidateAuthorization = async () => {
      setIsAuthorized(await validateToken(setLoading));
    };

    ValidateAuthorization();
  }, []);

  useEffect(() => {
    const fetchPostsEffect = async () => {
      if (isAuthorized) {
        await fetchPosts(await setPosts);
      }
    };
    fetchPostsEffect();
  }, [isAuthorized]);

  if (loading) {
    return <div>Loading</div>;
  }
  if (!isAuthorized) {
    return <Unauthorized />;
  }
  return (
    <div className="globalCss">
      <div>
        <label>Proszę wybrać post</label>
        <br></br>
        <select
          name="posts"
          onChange={(e) => setSelectedPost(parseInt(e.target.value)-1)}
        >
          <option value="">Wybierz posta</option>
          {posts ? (
            posts.map((post) => {
              return (
                <option key={post.ID} value={post.ID}>
                  {post.Title}
                </option>
              );
            })
          ) : (
            <div>No post found</div>
          )}
        </select>
      </div>

      {selectedPost !== undefined && posts[selectedPost] ? (
        <div>
          <Post id={posts[selectedPost].ID} title={posts[selectedPost].Title} content={posts[selectedPost].Content} date={posts[selectedPost].Created_At}/>
        </div>
      ) : (
        <div></div>
      )}
    </div>
  );
};

export default PostEditing;
