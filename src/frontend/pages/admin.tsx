import { useState } from "react"

const Admin = (props : any) => {
    const BORDER_CSS = "border"

    //edycja tekstu na stronie
    const [mainpageFirstHeader , setMainpageFirstHeader] = useState(props.mainpageFirstHeader)  //tworzymy stan lokalny ktorego poczatkowym stanem jest to co widzą wszyscy , czyli state z app.tsx

    //dodawanie postów
    const [title,setTitle] = useState("")
    const [content,setContent] = useState("")
    const postHandler = () => {
        console.log("Post")
        console.log("tytuł : " + title)
        console.log("treść : " + content)
        console.log("dodano do bazy danych")
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
            <button className={BORDER_CSS +" w-40 ml-5"} onClick={() => postHandler()}>Postuj</button>
         
        </div>
        </>
    )
}
export default Admin