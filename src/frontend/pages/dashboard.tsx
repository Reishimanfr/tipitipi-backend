import { useState } from "react"
import { useNavigate } from "react-router-dom"

interface BlogPostDataBodyJson {
    Content: string
    Created_At: string 
    Edited_At: string
    ID: number
    Images: any[]
    Title: string
    error?:  string
}

const Dashboard = (props : any) => {
    const BORDER_CSS = "border"
    const navigate = useNavigate()

    //edycja tekstu na stronie
    const [mainpageFirstHeader , setMainpageFirstHeader] = useState(props.mainpageFirstHeader)  //tworzymy stan lokalny ktorego poczatkowym stanem jest to co widzą wszyscy , czyli state z app.tsx

    //dodawanie postów
    const [title,setTitle] = useState("")
    const [content,setContent] = useState("")

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
            console.debug("Token is invalid, redirecting to login page...")
            navigate("/admin/login")
            return
        }

        const request = await fetch("http://localhost:2333/api/blog/create", {
            method: "POST",
            body: formData,
            headers: {Authorization: token}
        })
    
    
        if(request.ok){
            alert("Opublikowano post")
            window.location.reload();
            
        }
        else {
            const response: BlogPostDataBodyJson = await request.json()
            alert("Błąd: " + response.error)
        }
    }

    return(
        <>
        <div>
            {/* wyświetlamy to co widzi użytkownik */}
            <h1>Pierwszy nagłówek strony głównej : {props.mainpageFirstHeader}</h1> 

            {/* zmieniamy stan lokalny na to co wpisze admin */}
            <input className={BORDER_CSS} onChange={(event) => setMainpageFirstHeader(event.target.value)}></input><br></br>

            {/* po wcisnieciu przycisku wywolywana jest funkcja podana w props , stan globalny dostaje wartosc lokalnego */}
            <button className={BORDER_CSS+" w-40 ml-5"} onClick={() => props.changeMainpageFirstHeader(mainpageFirstHeader)}>fin</button>
            
            <hr></hr>
        </div>


        <div>
            <h1 className="font-bold text-3xl">Dodawanie postów</h1><br></br>
            <form>
                <label htmlFor="title">Podaj nazwę posta: </label>
                <input type="text" name="title" className={BORDER_CSS} onChange={(event) => setTitle(event.target.value)}/><br></br>

                <label htmlFor="content">Podaj nazwę posta: </label><br/>
                <textarea name="content" className={BORDER_CSS } cols={100} rows={10} onChange={(event) => setContent(event.target.value)}></textarea><br></br>

                <label htmlFor="image">Wybierz zdjęcie: </label>
                <input type="file" name="image" accept="image/*"></input><br></br><br></br>
            </form>
                {/* <input type="submit" value="Postuj" name="submitButton" className={BORDER_CSS +" w-40 ml-5"} onSubmit={() => setTest("xd")}/> */}
            <button className={BORDER_CSS +" w-40 ml-5"} onClick={() => addPost()}>Postuj</button>
         
        </div>
        </>
    )
}
export default Dashboard