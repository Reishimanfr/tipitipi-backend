import { useEffect, useState } from "react";
import Unauthorized from "../../errorPages/unauthorized";
import validateToken from "../../../functions/validate";
import { validateDataForm , buildMultipart, getToken } from "../../../functions/postManipulatingFunctions";
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


export default function PostCreating() {
  const [title, setTitle] = useState("Tytuł posta");
  const [content, setContent] = useState("Treść posta");

  async function addPost() {
    if (!validateDataForm(title, content)) {
      return;
    }

    const formData = buildMultipart(title, content);

    const token = getToken()

    try{
    const response = await fetch("http://localhost:2333/blog/post/", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`
      },
      body: formData,
    });

    if (response.status === 200) {
      alert("Opublikowano post");
      window.location.reload();
    } else {
      const data: BlogPostDataBodyJson = await response.json();
      alert("Błąd: " + data.error);
    }
  } catch (error){
    console.error(error)
    alert("Wystąpił błąd: " + error)
  }
  }

  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);
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
    <div className="mt-[1%] globalCss">
      <QuillBody title={title} setTitle={setTitle} content={content} setContent={setContent} handlerPost={addPost}/>
    </div>
    
  );
}
