import { useState, useEffect } from "react";
import validateToken from "../../../functions/validate";
import Unauthorized from "../../errorPages/unauthorized";
import {
  validateDataForm,
  buildMultipart,
  getToken
} from "../../../functions/postManipulatingFunctions";
import QuillBody from "../../../components/quillBody";

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
async function fetchPosts(
  setPosts: React.Dispatch<React.SetStateAction<BlogPostDataBodyJson[]>>
) {
  try {
    const response = await fetch(
      `http://localhost:2333/blog/posts?limit=999&attachments=true`,
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
  }
}




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
        `http://localhost:2333/blog/post/${selectedPost.id}`,
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
    if (title == selectedPost.title && content == selectedPost.content) {
      alert("Nie dokonano żadnych zmian");
      return;
    }
    const formData = buildMultipart(title, content);

    try {
      const response = await fetch(
        `http://localhost:2333/blog/post/${selectedPost.id}`,
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
      setTitle(selectedPost.title);

      if (selectedPost.attachments && selectedPost.content) {
        let tempContent = selectedPost.content
        selectedPost.attachments.forEach((attachment, index) => {
          tempContent = tempContent.replace(
            `{{${index}}}`,
            `<img style="max-height:200px;" src="http://localhost:2333/proxy?key=${attachment.filename}" alt="${attachment.filename}"/>`
          );
          
        });
        setContent(tempContent)
      }
        else {
          if(content) {
            const tempContent = content?.replace(/{{\d+}}/g, "")
            setContent(tempContent)
          }
      
        }
      // setContent(selectedPost.content);
    }
  }, [selectedPost]);


  const test = () => {
    if (selectedPost!.attachments && content) {
      console.log(selectedPost!.attachments)
      let tempContent = content
      selectedPost!.attachments.forEach((attachment, index) => {
        if (!tempContent) {
          return;
        }
        tempContent = tempContent.replace(
          `{{${index}}}`,
          `<img style="max-height:200px;" src="http://localhost:2333/proxy?key=${attachment.filename}" alt="${attachment.filename}"/>`
        );
        
      });
      setContent(tempContent)
    }
    else {
      if(content) {
        const tempContent = content?.replace(/{{\d+}}/g, "")
        setContent(tempContent)
      }
  
    }
  }

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
                <option key={post.id} value={index}>
                  {post.id + " , " + post.title}
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

      <button onClick={test}>test</button>
    </div>
  );
};

export default PostEditing;
