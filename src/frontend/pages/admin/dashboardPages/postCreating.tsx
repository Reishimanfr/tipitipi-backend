import { useState } from "react";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import { useNavigate } from "react-router-dom";
import axios from "axios";



interface BlogPostDataBodyJson {
  Content: string
  Created_At: string 
  Edited_At: string
  ID: number
  Images: any[]
  Title: string
  error?:  string
}


export default function PostCreating() {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("Treść posta");
  const navigate = useNavigate()

  function validateDataForm() {
    if(title === "") {
        alert("Podano pusty tytuł")

        return false
    }
    if(content === "") {
        alert("Podano pustą treść")
        return false
    }
    const confirm = window.confirm("Czy jesteś pewien że chcesz opublikować ten post?")
    if(!confirm){
        return false;
    }
    return true
}

  async function addPost() {
    if (!validateDataForm()) {return}

    const formData = new FormData()
    formData.append("title",title)
    formData.append("content",content)
    //formData.append("images","")

    const token = localStorage.getItem("token")
    if (!token) {
        alert("Token is invalid, redirecting to login page...")
        navigate("/admin/login")
        return
    }

    const request = await fetch("http://localhost:2333/blog/post", {
        method: "POST",
        headers: {Authorization: token},
        body: JSON.stringify({
          "Title": title,
          "Content": content,
          "Images": []
        })

    })
    // const request = await axios.post("http://localhost:2333/blog/post",{
    //     headers: {Authorization: token},
    //     body: {
    //       "Title": title,
    //       "Content": content,
    //       "Images": []
    //     }
    // })


    if(request.ok){
        alert("Opublikowano post")
        window.location.reload();
        
    }
    else {
        const response: BlogPostDataBodyJson = await request.json()
        alert("Błąd: " + response.error)
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
      /><br></br>
      <button className={"border w-40"} onClick={() => addPost()}>Postuj</button>
    </div>
  );
}
