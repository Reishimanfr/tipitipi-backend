import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";

interface Props {
    title:string;
    setTitle: React.Dispatch<React.SetStateAction<string>>;
    content: string;
    setContent: React.Dispatch<React.SetStateAction<string>>;
    handlerPost: () => Promise<void>;
}

const QuillBody = ({title,setTitle,content,setContent,handlerPost} : Props) => {
    return(
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
          <button className={"border w-40"} onClick={() => handlerPost()}>
            Postuj
          </button>
        </div>
    )
}

export default QuillBody