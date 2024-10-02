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

// async function fetchFileAsBlob(path: string): Promise<Blob> {
//   const response = await fetch(path);
  
//   if (!response.ok) {
//     throw new Error("Błąd podczas pobierania pliku");
//   }

//   const blob = await response.blob();
//   return blob;
// }

// async function getBase64(path : string) {
//   const file = await fetchFileAsBlob(path)
//   var reader = new FileReader();
//   reader.readAsDataURL(file);
//   reader.onload = function () {
//     console.log(reader.result);
//   };
//   reader.onerror = function (error) {
//     console.log('Error: ', error);
//   };
// }

const PostEditing = () => {
  const [selectedPost, setSelectedPost] = useState<BlogPostDataBodyJson>();
  const [posts, setPosts] = useState<Array<BlogPostDataBodyJson>>([]);
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");

  async function deletePost() {
    const token = getToken();
    if (!selectedPost) {
      alert("Nie znaleziono posta");
      return;
    }
    if (!window.confirm("Czy jesteś pewien że chcesz usunąć ten post?")) {
      return;
    }
    try {
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
    } catch (error) {
      console.error(error);
      alert("Wystąpił błąd: " + error);
    }
  }
  async function editPost() {
    if (!validateDataForm(title, content)) {
      return;
    }
  
    const token = getToken();
    if (!selectedPost) {
      alert("Nie znaleziono posta");
      return;
    }
    if (title == selectedPost.Title && content == selectedPost.Content) {
      alert("Nie dokonano żadnych zmian");
      return;
    }
    const formData = buildMultipart(title, content);

    try {
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
    } catch (error) {
      console.error(error);
      alert("Wystąpił błąd: " + error);
    }
  }

  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);
;

  useEffect(() => {
    if (selectedPost) {
      setTitle(selectedPost.Title);
      setContent(selectedPost.Content);
    }
  }, [selectedPost]);




  useEffect(() => {
    const ValidateAuthorization = async () => {
      setIsAuthorized(await validateToken(setLoading));
    };
    ValidateAuthorization();
  }, []);
  useEffect(() => {
    const fetchPostsEffect = async () => {
      if (isAuthorized && posts.length == 0) {
        await fetchPosts(setPosts);
      }
    };
    fetchPostsEffect();
  }, [isAuthorized])
  if (loading) {
    return <div>Loading</div>;
  }
  if (!isAuthorized) {
    return <Unauthorized />;
  }
  return (
    <div className="globalCss">
      <div className="my-[1%]">
        <label className="text-xl">Proszę wybrać post</label>
        <br></br>
        <select
          name="posts"
          onChange={(e) => setSelectedPost(posts[parseInt(e.target.value)])}
        >
          <option value="">--post--</option>
          {posts ? (
            posts.map((post, index) => {
              return (
                <option key={post.ID} value={index}>
                  {post.ID + " , " + post.Title}
                </option>
              );
            })
          ) : (
            <div>No post found</div>
          )}
        </select>
      </div>

      {selectedPost !== undefined ? (
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
            className="border w-40 text-white shadow-lg bg-red-500 hover:bg-red-600 hover:duration-300"
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
