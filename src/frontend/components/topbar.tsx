import { Link } from "react-router-dom";
const Topbar = () => {
    const DESKTOP_CSS = "float-right no-underline text-white text-3xl mr-2.5 bg-black-700 font-sans hover:duration-500  hover:bg-white hover:text-black" as const
    return (
        <div>
          <div className="w-full h-20 flex-wrap top-20 p-5 text-right">
            <Link to='/pricing' className={DESKTOP_CSS}> Cennik </Link>
            <Link to='/about' className={DESKTOP_CSS}> O nas </Link>
            <Link to='/gallery' className={DESKTOP_CSS}> Galeria </Link>
            <Link to='/blog' className={DESKTOP_CSS}> Blog </Link>
            <Link to='/' className={DESKTOP_CSS}> Strona główna</Link>
          </div>
         
      </div>
    )
}

export default Topbar