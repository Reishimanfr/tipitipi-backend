import { useEffect, useState } from "react";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import { useNavigate } from "react-router-dom";
import Unauthorized from "../../errorPages/unauthorized";
import validateToken from "../../../components/validate";

interface BlogPostDataBodyJson {
  Content: string;
  Created_At: string;
  Edited_At: string;
  ID: number;
  Images: any[];
  Title: string;
  error?: string;
}
function validateDataForm(title: string, content: string): boolean {
  if (title === "") {
    alert("Podano pusty tytuł");

    return false;
  }
  if (content === "<p><br></p>") {
    alert("Podano pustą treść");
    return false;
  }
  const confirm = window.confirm(
    "Czy jesteś pewien że chcesz opublikować ten post?"
  );
  if (!confirm) {
    return false;
  }
  return true;
}
function makeFilename(length: number) {
  let result = "";
  const characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  const charactersLength = characters.length;
  let counter = 0;
  while (counter < length) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
    counter += 1;
  }
  return result;
}

function base64ToBlob(base64: string): Blob | null {
  //mimetype is first part of base64 variable , that matches given regexp in match function
  let mimeType: RegExpMatchArray | null | string = base64
    .split(",")[0]
    .match(/image\/(jpeg|png|gif)/);
  //base64content is second part of base64 variable , that has last two chars trimmed
  let base64content = base64.split(",")[1];
  base64content = base64content.substring(0, base64content.length - 2);

  if (!mimeType) {
    console.error("Mimetype of image is null");
    return null;
  }
  mimeType = mimeType[0];
  const byteString = atob(base64content);

  // Tworzymy tablicę bajtów
  const byteNumbers = new Array(byteString.length);
  for (let i = 0; i < byteString.length; i++) {
    byteNumbers[i] = byteString.charCodeAt(i);
  }
  const byteArray = new Uint8Array(byteNumbers);
  return new Blob([byteArray], { type: mimeType });
}

function extractImagesFromContent(content: string) {
  const images: string[] = [];
  let i = -1;
  const regexp =
    /<img\s+src="data:image\/(jpeg|png|gif);base64,([A-Za-z0-9+/=]+)"\s*\/?>/g;

  const contentWithoutImages = content.replace(regexp, (match) => {
    images.push(match);
    i++;
    return `{{${i}}}`;
  });
  return { images: images, contentWithoutImages: contentWithoutImages };
}

function buildMultipart(title: string, content: string) {
  const extractedData = extractImagesFromContent(content);
  const base64images: string[] = extractedData.images;

  const formData = new FormData()

  formData.append("title", title)
  formData.append("content", extractedData.contentWithoutImages)
  
  for (const image of base64images) {
    const blob = base64ToBlob(image)
    if (!blob) {
      return;
    }
    formData.append("files[]", blob, makeFilename(5) + '.' + blob.type.split("/")[1])
  }

  return formData;
}

export default function PostCreating() {
  const [title, setTitle] = useState("Tytuł posta");
  const [content, setContent] = useState("Treść posta");
  const navigate = useNavigate();

  async function addPost() {
    if (!validateDataForm(title, content)) {
      return;
    }

    const formData = buildMultipart(title, content);

    const token = localStorage.getItem("token");
    if (!token) {
      alert("Token is invalid, redirecting to login page...");
      navigate("/admin/login");
      return;
    }

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
    <div>
      <label htmlFor="title">Podaj nazwę posta: </label>
      <input
        type="text"
        name="title"
        value={title}
        className="border"
        onChange={(event) => setTitle(event.target.value)}
      />
      <br></br>
      <br></br>

      <h1>Podaj treść posta:</h1>
      <ReactQuill
        theme="snow"
        value={content}
        onChange={setContent}
        //style={{ minHeight: "500px" }}
        modules={{
          toolbar: [
            ["bold", "italic", "underline"],
            [{ align: [] }],

            [{ list: "ordered" }, { list: "bullet" }],
            [{ indent: "-1" }, { indent: "+1" }],

            [{ size: ["small", false, "large", "huge"] }],
            [{ header: [1, 2, 3, 4, 5, 6, false] }],
            ["link", "image", "video"],
            [{ color: [] }, { background: [] }],

            ["clean"],
          ],
        }}
      />
      <br></br>
      <button className={"border w-40"} onClick={() => addPost()}>
        Postuj
      </button>
    </div>
  );
}
