import { useState, useEffect } from "react";
import validateToken from "../../../components/validate";
import Unauthorized from "../../errorPages/unauthorized";
import {
  validateDataForm,
  buildMultipart,
  getToken,
  fetchPosts,
} from "./postManipulatingFunctions";
import QuillBody from "../../../components/quillBody";

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
  error?: string;
}

const PostEditing = () => {
  const [selectedPostIndex, setSelectedPostIndex] = useState<number>();
  const [posts, setPosts] = useState<Array<BlogPostDataBodyJson>>([]);
  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");


  //TODO bugged , doesnt delete but it did???? cannot pick newest posts because site crashes
  
  async function deletePost() {
    const token = getToken();
    if (!selectedPostIndex) {
      alert("Nie znaleziono posta");
      return;
    }
    if (!window.confirm("Czy jesteś pewien że chcesz usunąć ten post?")) {
      return;
    }
    const selectedPost = posts[selectedPostIndex];
    const response = await fetch(
      `http://localhost:2333/blog/post/${selectedPost.ID}`,
      {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    if (response.status === 200) {
      alert("Usunięto post");
      window.location.reload();
    } else {
      const data: BlogPostDataBodyJson = await response.json();
      alert("Błąd: " + data.error);
    }
  }
  async function editPost() {
    if (!validateDataForm(title, content)) {
      return;
    }
    const formData = buildMultipart(title, content);
    const token = getToken();
    if (!selectedPostIndex) {
      alert("Nie znaleziono posta");
      return;
    }
    const selectedPost = posts[selectedPostIndex];
    if (title == selectedPost.Title || content == selectedPost.Content) {
      alert("Nie dokonano żadnych zmian");
      return;
    }

    const response = await fetch(
      `http://localhost:2333/blog/post/${selectedPost.ID}`,
      {
        method: "PATCH",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      }
    );

    if (response.status === 200) {
      alert("Edytowano post");
      window.location.reload();
    } else {
      const data: BlogPostDataBodyJson = await response.json();
      alert("Błąd: " + data.error);
    }
  }

  useEffect(() => {
    const fetchPostsEffect = async () => {
      if (isAuthorized && posts.length == 0) {
        await fetchPosts(setPosts);
      }
    };
    fetchPostsEffect();
  }, [isAuthorized]);

  useEffect(() => {
    if (selectedPostIndex) {
      setTitle(posts[selectedPostIndex].Title);
      setContent(posts[selectedPostIndex].Content);
    }
  }, [selectedPostIndex]);

  useEffect(() => {
    const ValidateAuthorization = async () => {
      setIsAuthorized(await validateToken(setLoading));
    };
    ValidateAuthorization();
  }, []);
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
          onChange={(e) => setSelectedPostIndex(parseInt(e.target.value) - 1)}
        >
          <option value="">Wybierz posta</option>
          {posts ? (
            posts.map((post) => {
              return (
                <option key={post.ID} value={post.ID}>
                  {post.ID + " , " + post.Title}
                </option>
              );
            })
          ) : (
            <div>No post found</div>
          )}
        </select>
      </div>

      {selectedPostIndex !== undefined && posts[selectedPostIndex] ? (
        <div>
          <QuillBody
            title={title}
            setTitle={setTitle}
            content={content}
            setContent={setContent}
            handlerPost={editPost}
          />
          <br></br>
          <button
            className="border w-40 bg-red-500 text-white"
            onClick={() => deletePost()}
          >
            Usuń tego posta
          </button>
        </div>
      ) : (
        <div></div>
      )}
    </div>
  );
};

export default PostEditing;
