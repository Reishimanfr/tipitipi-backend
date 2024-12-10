import { useState } from "react"
import { useNavigate } from "react-router-dom"
//import axios from "axios"

interface LoginResponse {
    token?: string
    error?: string
    message?: string
}

const Login = () => {
    const [login,setLogin] = useState("")
    const [password,setPassword] = useState("")

    const navigate = useNavigate()


    async function loginHandler(){
        const formData = new FormData()
        formData.append("username",login)
        formData.append("password",password)

        const response = await fetch("http://localhost:8080/admin/login", {
            method: "POST",
            body: JSON.stringify({
                username: login,
                password: password
            }) 

            
        })
      
        const data : LoginResponse  = await response.json()   

        if(response.status === 200 && data.token != undefined) {
            localStorage.setItem("token",data.token)
            navigate("/admin/dashboard")
        }
        else{
            localStorage.setItem("token", "bad")
            alert("something went wrong: " +  data.error)
        }
    }

    return(
        <div className="m-auto mt-[20vh] border-2 border-gray-800  text-center w-[25%] rounded-lg">
            <form>
                <div className="p-[5%]">
                    <label className="text-xl font-semibold" htmlFor="login">Podaj login: </label>
                    <input className="border-2 w-1/2" type="text" name="login" onChange={(event) => setLogin(event.target.value)}/><br></br>
                </div>
                
                <div className="p-[3%]">
                    <label className="text-xl font-semibold" htmlFor="password">Podaj Hasło: </label>
                    <input className="border-2 w-1/2"  type="password" name="password" onChange={(event) => setPassword(event.target.value)}/><br></br>
                </div>

            </form>
            <button className={"m-[5%] p-[1%] border w-1/2 shadow-lg hover:bg-slate-100 hover:duration-300"} onClick={() => loginHandler()}>Zaloguj</button>
        </div>
    )
}
export default Login