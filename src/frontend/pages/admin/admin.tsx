import { Buffer } from "buffer"
import { useEffect } from "react"
import { useNavigate } from "react-router-dom"

const decode = (str: string):string => Buffer.from(str, 'base64').toString('binary');


interface JwtPayload {
    iat: number
    exp: number
    admin: boolean
    user_id: string
} 

const Admin = () => {
    const navigate = useNavigate()
    useEffect(() => {
        const token = localStorage.getItem("token")

        if (token === null) {
            console.debug("Token is invalid, redirecting to login page...")
            navigate("/admin/login")
            return
        }
    
        // 0: Header | 1: Payload | 2: Signature
        const tokenSplit = token.split(".")
    
        if (!tokenSplit?.[1]) {
            console.error("Malformed token")
            return
        }
    
        const stringPayload = decode(tokenSplit[1])
        const payload: JwtPayload = JSON.parse(stringPayload)
        const now = Date.now() / 1000
    
        if (now >= payload.exp) {
            console.debug("Token expired, redirecting to login page...")
            navigate("/admin/login")
            return
        }

        navigate("/admin/dashboard")
    })
  
    
   return (<></>)
   
}
export default Admin