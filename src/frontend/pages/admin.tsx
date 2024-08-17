import { useState } from "react"

const Admin = (props : any) => {
    const [mainpageFirstHeader , setMainpageFirstHeader] = useState(props.mainpageFirstHeader)  //tworzymy stan lokalny ktorego poczatkowym stanem jest to co widzą wszyscy , czyli state z app.tsx
    return(
        <div>
            {/* wyświetlamy to co widzi użytkownik */}
            <h1>Pierwszy nagłówek strony głównej : {props.mainpageFirstHeader}</h1> 

            {/* zmieniamy stan lokalny na to co wpisze admin */}
            <input className="border" onChange={(event) => setMainpageFirstHeader(event.target.value)}></input>

            {/* po wcisnieciu przycisku wywolywana jest funkcja podana w props , stan globalny dostaje wartosc lokalnego */}
            <button onClick={() => props.changeMainpageFirstHeader(mainpageFirstHeader)}>fin</button>
        </div>
    )
}
export default Admin