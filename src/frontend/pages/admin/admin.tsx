
import { useNavigate } from "react-router-dom"

const Admin = (props : any) => {
    const navigate = useNavigate()
    if(props.isAuthorized){
        navigate("/admin/dashboard")
    }
    else {
        navigate("/admin/login")
    }
    
  
    
   return (<></>)
   
}
export default Admin