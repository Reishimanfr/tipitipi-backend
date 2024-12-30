import { useState, useEffect } from "react";
import Topbar from "./topbar";
import Sidebar from "./sidebar";
import { Link } from "react-router-dom";
import logoSmall from "../assets/logo.png"

const MainNavbar = () => {
  const [showDesktopNavbar, setShowDesktopNavbar] = useState(
    window.innerWidth >= 768
  );
  const [lastScrollY, setLastScrollY] = useState(0);
  const [show, setShow] = useState(true);

  const handleResize = () => {
    setShowDesktopNavbar(window.innerWidth > 768);
  };

  const handleScroll = () => {
    if (typeof window !== 'undefined') {
      setShow(window.scrollY < lastScrollY)
      setLastScrollY(window.scrollY);
    }
  };

  useEffect(() => {
    window.addEventListener("scroll" , handleScroll);
    return () => {
      window.removeEventListener("scroll" , handleScroll)
      }
  },[lastScrollY]);

  useEffect(() => {
    window.addEventListener("resize", handleResize);
    return () => {
      window.removeEventListener("resize", handleResize)
      }
  });

   return (
    <div className={`sticky z-50 top-0 h-20 bg-[#24252A] w-full transition-transform duration-300 transform ${show ? 'translate-y-0' : '-translate-y-full'}`}>
         <Link to="/">
        <img
          src={logoSmall}
          alt="logo"
          className="w-auto h-16 float-left ml-[5%] mt-[1vh]"
        ></img>
        </Link>
        {showDesktopNavbar ? <Topbar/> : <Sidebar/>}
    </div>
   )
};
export default MainNavbar;

