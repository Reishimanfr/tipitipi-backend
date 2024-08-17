import { Link } from "react-router-dom"
const Mainpage = (props : any) => {
    return(
        <div>
            <h1>Mainpage</h1>
            <Link to="/admin">Admin</Link>
            <h1>{props.mainpageFirstHeader}</h1>
        </div>
    )
}
export default Mainpage