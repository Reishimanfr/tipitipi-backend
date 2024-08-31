import { useNavigate } from "react-router-dom";
const Admin = () => {
    const navigate = useNavigate()
    return(
        <div>
            <h1>Admin</h1>
            <button onClick={() => {navigate("/admin/login")}}>Login</button>
            <button onClick={() => {navigate("/admin/dashboard")}}>Dashboard</button>
        </div>
    )
}
export default Admin