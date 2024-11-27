import { useState } from "react"
import { useNavigate } from "react-router-dom"

const Login = () => {
    const [login,setLogin] = useState("")
    const [password,setPassword] = useState("")

    const navigate = useNavigate()


    async function loginHandler(){
        const formData = new FormData()
        formData.append("username",login)
        formData.append("password",password)
        const response = await fetch("http://localhost:2333/api/admin/login", {
            method: "POST",
            body: formData
        })
      
        const data  = await response.json()       //TODO idk what data type is 
        console.log(data)
        console.log("a")
        console.log(response)

        if(response.ok && data.token != undefined) {
            localStorage.setItem("token",data.token)
            navigate("/admin/dashboard")
        }
        else{
            console.log("something went wrong: " , data.error)
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