import { useState } from "react"

const Login = () => {
    const [login,setLogin] = useState("")
    const [password,setPassword] = useState("")

    console.log("Trying to log in")

    async function loginHandler(){
        const formData = new FormData()
        formData.append("username",login)
        formData.append("password",password)
        const request = await fetch("http://localhost:2333/api/admin/login", {
            method: "POST",
            body: formData
        })

        const data = await request.json()       //TODO idk what data type is 
        console.log(data.Code)
    }

    return(
        <div>
            <h1>Login</h1>
            <form>
                <label htmlFor="login">Podaj login</label>
                <input type="text" name="login" onChange={(event) => setLogin(event.target.value)}/><br></br>

                <label htmlFor="password">Podaj login</label>
                <input type="text" name="password" onChange={(event) => setPassword(event.target.value)}/><br></br>
                
            </form>
            <button onClick={() => loginHandler()}>Zaloguj</button>
        </div>
    )
}
export default Login