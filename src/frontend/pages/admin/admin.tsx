
import { useNavigate } from "react-router-dom"
import { useEffect } from "react";

const Admin = (props : any) => {
    const navigate = useNavigate()
      console.log("I AM AUTHORIZED ADMIN: " + props.isAuthorized)
    useEffect(() => {
        if(props.isAuthorized){
            navigate("/admin/dashboard")
        }
        else {
            navigate("/admin/login")
        }
    });
   
    
  
    
   return (<></>)
   
}
export default Admin