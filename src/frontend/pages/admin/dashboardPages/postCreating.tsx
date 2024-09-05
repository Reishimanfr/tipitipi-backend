import { useState } from "react"
import ReactQuill from "react-quill"
import 'react-quill/dist/quill.snow.css'

export default function PostCreating() {
  const [convertedText, setConvertedText] = useState("Some default content");
  return (
    <div>
      <ReactQuill
        theme='snow'
        value={convertedText}
        onChange={setConvertedText}
        style={{minHeight: '300px'}}
        modules={{toolbar: [
            ['bold', 'italic', 'underline'],
            [{ align: [] }],
        
            [{ list: 'ordered'}, { list: 'bullet' }],
            [{ indent: '-1'}, { indent: '+1' }],
        
            [{ size: ['small', false, 'large', 'huge'] }],
            [{ header: [1, 2, 3, 4, 5, 6, false] }],
            ['link', 'image', 'video'],
            [{ color: [] }, { background: [] }],
        
            ['clean'],
          ]}}
        
      />
      <button onClick={()=>{console.log(convertedText)}}>przeslij</button>
    </div>
  );
}

