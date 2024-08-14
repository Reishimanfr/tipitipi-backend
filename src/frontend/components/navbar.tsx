import { Link } from "react-router-dom"
import logoSmall from "../react.svg"
const Navbar = () => {
    const CSS = "float-right mt-[1%] text-white text-3xl mr-2.5 bg-black-700 font-sans hover:duration-500  hover:bg-white hover:text-black"
    return(
        <div className="w-full h-20 bg-black sticky top-0">
            <div>
            <Link to='/'>
                <img src={logoSmall} alt="logo" className="h-16 float-left ml-[5%] mt-[1vh]">
                </img>
            </Link>
            <Link to='/about' className={CSS}> O nas </Link>,
            <Link to='/pricing' className={CSS}>Pricing</Link>,
            <Link to='/gallery' className={CSS}> Gallery </Link>,
            <Link to='/blog' className={CSS}> Blog </Link>,
            <Link to='/' className={CSS}> Strona główna</Link>
            </div>
        </div>
    )
}

export default Navbar