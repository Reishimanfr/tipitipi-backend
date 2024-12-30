import { Link } from "react-router-dom";
const Topbar = () => {
    const DESKTOP_CSS = "float-right no-underline font-custom p-2 text-white text-[24px] mr-5 bg-black-700 font-sans hover:duration-500 rounded-xl hover:bg-white hover:text-black" as const
    return (
        <div>
          <div className="mt-[14px] mr-20">
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