
import { useNavigate } from "react-router-dom"
import { useEffect, useState } from "react";

const Admin = (props : any) => {
    const navigate = useNavigate()
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