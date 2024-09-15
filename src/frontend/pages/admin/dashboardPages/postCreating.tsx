import { useState } from "react";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import { useNavigate } from "react-router-dom";

interface BlogPostDataBodyJson {
  Content: string;
  Created_At: string;
  Edited_At: string;
  ID: number;
  Images: any[];
  Title: string;
  error?: string;
}

export default function PostCreating() {
  const [title, setTitle] = useState("Tytuł posta");
  const [content, setContent] = useState("Treść posta");
  const navigate = useNavigate();

  function validateDataForm() {
    if (title === "") {
      alert("Podano pusty tytuł");

      return false;
    }
    if (content === "") {
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


  async function addPost() {
    if (!validateDataForm()) {
      return;
    }
    const boundary = (Math.random() + 1).toString(36).substring(2)
    const formData = `--${boundary}
Content-Disposition: form-data; name="title"

${title}
--${boundary}
Content-Disposition: form-data; name="content"

${content}
--${boundary}--`;

    const token = localStorage.getItem("token");
    if (!token) {
      alert("Token is invalid, redirecting to login page...");
      navigate("/admin/login");
      return;
    }

    const response = await fetch("http://localhost:2333/blog/post/", {
      method: "POST",
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": `multipart/form-data; boundary=${boundary}`
      },
      body: formData,
    });


    if (response.status === 200) {
      alert("Opublikowano post");
      window.location.reload();
    } else {
      // console.log(response);
      // const data: BlogPostDataBodyJson = await response.json();
      alert("Błąd: ");
    }
  }

  return (
    <div>
      <label htmlFor="title">Podaj nazwę posta: </label>
      <input
        type="text"
        name="title"
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
