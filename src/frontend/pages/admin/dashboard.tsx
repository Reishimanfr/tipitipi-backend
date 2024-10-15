import { useEffect, useState } from "react"
import { Link, useNavigate } from "react-router-dom"
import validateToken from "../../components/validate"
import Unauthorized from "../errorPages/unauthorized"




const Dashboard = () => {
    const navigate = useNavigate()
    //const BORDER_CSS = "border"
    //edycja tekstu na stronie
    //const [mainpageFirstHeader , setMainpageFirstHeader] = useState(props.mainpageFirstHeader)  //tworzymy stan lokalny ktorego poczatkowym stanem jest to co widzą wszyscy , czyli state z app.tsx

    const Logout = () => {
        if(window.confirm("Czy napewno chcesz się wylogować?")) {
            alert("Wylogowano");
            localStorage.setItem("token",'')
            navigate('/admin/login')
        }
    }

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
        {/* <div>
            wyświetlamy to co widzi użytkownik
            <h1>Pierwszy nagłówek strony głównej : {props.mainpageFirstHeader}</h1> 

            zmieniamy stan lokalny na to co wpisze admin
            <input className={BORDER_CSS} onChange={(event) => setMainpageFirstHeader(event.target.value)}></input><br></br>

            po wcisnieciu przycisku wywolywana jest funkcja podana w props , stan globalny dostaje wartosc lokalnego
            <button className={BORDER_CSS+" w-40 ml-5"} onClick={() => props.changeMainpageFirstHeader(mainpageFirstHeader)}>fin</button>
            
            <hr></hr>
        </div> */}


            <div className="mt-[1%]">
                <Link to="/admin/dashboard/create-post"><button className="border p-[0.5%] ml-[1%] mb-[1%] border-gray-900 hover:bg-gray-900 hover:text-white hover:duration-300 rounded-md">Dodawanie postów</button></Link><br></br>
                <Link to="/admin/dashboard/edit-post"><button className="border p-[0.5%] ml-[1%] mb-[1%] border-gray-900 hover:bg-gray-900 hover:text-white hover:duration-300 rounded-md">Edycja postów</button></Link><br></br>
                <Link to="/admin/dashboard/change-credentials"><button className="border p-[0.5%] ml-[1%] mb-[1%] border-gray-900 hover:bg-gray-900 hover:text-white hover:duration-300 rounded-md">Zmień login/hasło</button></Link>  <br></br>
                <button onClick={() => Logout()} className="border p-[0.5%] ml-[1%] mb-[1%] border-gray-900 hover:bg-gray-900 hover:text-white hover:duration-300 rounded-md">Wyloguj się</button>
            </div>
        </>
    )
}
export default Dashboard