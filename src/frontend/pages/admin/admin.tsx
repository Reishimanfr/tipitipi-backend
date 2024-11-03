
import { useNavigate } from "react-router-dom"
import { useEffect, useState } from "react";
import validateToken from "../../functions/validate";
const Admin = () => {
    const navigate = useNavigate()

    const [loading ,setLoading] = useState(true)
    const [isAuthorized , setIsAuthorized] = useState(false) 
    useEffect(() => {
        const ValidateAuthorization = async () => {
            setIsAuthorized(await validateToken(setLoading))
        }
        ValidateAuthorization()
    },[])
    useEffect(() => {
        if(loading == false){
            if(isAuthorized){
                navigate("/admin/dashboard")
            }
            else {
                navigate("/admin/login")
            }
        }
    },[loading]);


   return (<></>)
   
}
export default Admin