import { useEffect, useState } from "react"
import { Link } from "react-router-dom"
import validateToken from "../../components/validate"
import Unauthorized from "../errorPages/unauthorized"



const Dashboard = (props : any) => {

    const BORDER_CSS = "border"
    //edycja tekstu na stronie
    const [mainpageFirstHeader , setMainpageFirstHeader] = useState(props.mainpageFirstHeader)  //tworzymy stan lokalny ktorego poczatkowym stanem jest to co widzą wszyscy , czyli state z app.tsx


    const [loading ,setLoading] = useState(true)
    const [isAuthorized , setIsAuthorized] = useState(false) 
    useEffect(() => {
        const ValidateAuthorization = async () => {
            setIsAuthorized(await validateToken(setLoading))
        }
        ValidateAuthorization()
    },[])
    if(loading) {
        return(<div>
            Loading
        </div>)
    }
    if(!isAuthorized) {
        return <Unauthorized/>
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
                <Link to="/admin/dashboard/create-post"><button>Dodawanie postów</button></Link>
            </div>
        </>
    )
}
export default Dashboard