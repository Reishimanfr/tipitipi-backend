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
        console.log(login)
        console.log(password)
        const formData = new FormData()
        formData.append("username",login)
        formData.append("password",password)

        const response = await fetch("http://localhost:2333/admin/login", {
            method: "POST",
            body: JSON.stringify({
                username: login,
                password: password
            }) 

            
        })
        // const request = await axios.post("http://localhost:2333/admin/login" , {
        //     "username" : login,
        //     "password" : password
        // })
      
        const data : LoginResponse  = await response.json()    //TODO idk what data type is 

        if(response.status === 200 && data.token != undefined) {
            localStorage.setItem("token",data.token)
            navigate("/admin/dashboard")
        }
        else{
            localStorage.setItem("token", "")
            alert("something went wrong: " +  data.error)
        }
    }

    return(
        <div>
            <h1>Login</h1>
            <form>
                <label htmlFor="login">Podaj login</label>
                <input type="text" name="login" onChange={(event) => setLogin(event.target.value)}/><br></br>

                <label htmlFor="password">Podaj Hasło</label>
                <input type="text" name="password" onChange={(event) => setPassword(event.target.value)}/><br></br>
                
            </form>
            <button onClick={() => loginHandler()}>Zaloguj</button>
        </div>
    )
}
export default Login