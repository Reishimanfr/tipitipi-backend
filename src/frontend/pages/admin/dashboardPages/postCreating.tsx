import { useEffect, useState } from "react";
import Unauthorized from "../../errorPages/unauthorized";
import validateToken from "../../../components/validate";
import { validateDataForm , buildMultipart, getToken } from "./postManipulatingFunctions";
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


export default function PostCreating() {
  const [title, setTitle] = useState("Tytuł posta");
  const [content, setContent] = useState("Treść posta");

  async function addPost() {
    if (!validateDataForm(title, content)) {
      return;
    }

    const formData = buildMultipart(title, content);

    const token = getToken()

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
    <QuillBody title={title} setTitle={setTitle} content={content} setContent={setContent} handlerPost={addPost}/>
  );
}
