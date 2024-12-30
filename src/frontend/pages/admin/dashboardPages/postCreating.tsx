import { useEffect, useState } from "react"
import { toast } from "react-toastify"
import QuillBody from "../../../components/quillBody"
import { API_URL } from '../../../functions/global'
import {
        buildPostMultipart,
        getToken,
        validateDataForm,
} from "../../../functions/postManipulatingFunctions"
import validateToken from "../../../functions/validate"
import Unauthorized from "../../errorPages/unauthorized"

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
      const response = await fetch(`${API_URL}/blog/post/`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: formData,
      });

      if (response.status >= 200 && response.status < 300) {
        toast.success("Opublikowano post");
        setTitle("Tytuł posta")
        setContent("Treść posta")
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
