import { useEffect, useState } from "react";
import Unauthorized from "../../errorPages/unauthorized";
import validateToken from "../../../functions/validate";
import {
  validateDataForm,
  buildPostMultipart,
  getToken,
} from "../../../functions/postManipulatingFunctions";
import QuillBody from "../../../components/quillBody";

export default function PostCreating() {
  const [title, setTitle] = useState("Tytuł posta");
  const [content, setContent] = useState("Treść posta");

  async function addPost() {
    if (!validateDataForm(title, content)) {
      return;
    }

    const formData = buildPostMultipart(title, content);

    const token = getToken();

    try {
      const response = await fetch("http://localhost:2333/blog/post/", {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });

      if (response.status >= 200 && response.status < 300) {
        alert("Opublikowano post");
        window.location.reload();
      } 
      // else {
      //   const data: BlogPostDataBodyJson = await response.json();
      //   alert("Błąd: " + data.error);
      // }
      else {
        throw new Error(response.statusText);
      }
    } catch (error) {
      console.error(error);
      alert("Wystąpił błąd: " + error);
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
    <div className="globalCss mt-[1%]">
      <QuillBody
        title={title}
        setTitle={setTitle}
        content={content}
        setContent={setContent}
        handlerPost={addPost}
      />
    </div>
  );
}
