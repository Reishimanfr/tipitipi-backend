import { useEffect, useState } from "react"
import { toast } from "react-toastify"
import QuillBody from "../../../components/quillBody"
import { API_URL } from '../../../functions/global'
import { BlogPostDataBodyJson } from "../../../functions/interfaces"
import {
        buildPostMultipart,
        getToken,
        validateDataForm,
} from "../../../functions/postManipulatingFunctions"
import validateToken from "../../../functions/validate"
import Unauthorized from "../../errorPages/unauthorized"

async function fetchPosts(
  setPosts: React.Dispatch<React.SetStateAction<BlogPostDataBodyJson[]>>
) {
  try {
    const response = await fetch(
      //TODO niewiem czy bezpieczne / sciagamy wszystkie post yistniejace ze zdjeciami itd // trza zrobic partiala
      `${API_URL}/blog/posts?limit=999&files=true`,
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
      toast.error("Nie znaleziono posta");
      return;
    }
    if (!window.confirm("Czy jesteś pewien że chcesz usunąć ten post?")) {
      return;
    }
    try {
      const response = await fetch(
        `${API_URL}/blog/post/${selectedPost.id}`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.status >= 200 && response.status < 300) {
        alert("Usunięto post");
        window.location.reload();
      } else {
        throw new Error(response.statusText);
      }
    } catch (error) {
      console.error(error);
      toast.error("Wystąpił błąd: " + error);
    }
  }

  async function editPost() {
    if (!validateDataForm(title, content)) {
      return;
    }

    const token = getToken();
    if (!selectedPost) {
      toast.error("Nie znaleziono posta");
      return;
    }
    if (title == selectedPost.title && content == selectedPost.content) {
      toast.warn("Nie dokonano żadnych zmian");
      return;
    }
    const formData = buildPostMultipart(title, content);

    try {
      const response = await fetch(
        `${API_URL}/blog/post/${selectedPost.id}`,
        {
          method: "PATCH",
          headers: {
            Authorization: `Bearer ${token}`,
          },
          body: formData,
        }
      );

      if (response.status >= 200 && response.status < 300) {
        // alert("Edytowano post");
        // window.location.reload();
        toast.success("Edytowano post");
        setSelectedPost(undefined);
        let select = document.getElementById("posts") as HTMLSelectElement;
        if (select) {
          select.selectedIndex = 0;
        }
      } else {
        throw new Error(response.statusText);
      }
    } catch (error) {
      console.error(error);
      toast.error("Wystąpił błąd: " + error);
    }
  }

  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);
  useEffect(() => {
    if (selectedPost) {
      setTitle(selectedPost.title);
      
      if (selectedPost.files && selectedPost.content) {
        let tempContent = selectedPost.content;
        selectedPost.files.forEach((attachment, index) => {
          tempContent = tempContent.replace(
            `{{${index}}}`,
            `<img style="max-height:200px;" src="${API_URL}/proxy?key=${attachment.filename}&type=blog" alt="${attachment.filename}"/>`
          );
        });
        setContent(tempContent);
        
      } else {
        if (selectedPost.content) {
          const tempContent = selectedPost.content?.replace(/{{\d+}}/g, "");
          setContent(tempContent);
          
        }
      }
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
  }, [isAuthorized]);

  if (loading) {
    return <div>Loading</div>;
  }
  if (!isAuthorized) {
    return <Unauthorized />;
  }
  return (
    <div className="globalCss mt-[1%]">
      <div className="my-[1%]">
        <label className="text-xl">Proszę wybrać post</label>
        <br></br>
        <select
          id="posts"
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
    </div>
  );
};

export default PostEditing;
