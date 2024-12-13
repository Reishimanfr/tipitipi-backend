import { Link } from "react-router-dom";
import x from "../assets/x.svg";
import list from "../assets/list.svg";
import { useEffect, useState } from "react";

const Sidebar = () => {
  const MOBILE_CSS =
    "mt-6 clear-right float-right no-underline text-white text-xl mr-2.5 bg-black-700 font-sans hover:duration-500 hover:bg-gray-300 hover:text-black" as const;
  const [menuVisible, setMenuVisible] = useState(false);
  const [show, setShow] = useState(false);
  // const handleResize = () => {
  //   setMenuVisible(false);
  // };

  useEffect(() => {
    if (menuVisible == false) {
      setTimeout(() => {
        document.body.style.overflow = "auto";
        setShow(false);
      }, 500);
    } else {
      setShow(true);
      document.body.style.overflow = "hidden";
    }
  }, [menuVisible]);

  useEffect(() => {
    window.addEventListener("resize", () => setMenuVisible(false));
    return () => {
      window.removeEventListener("resize", () => setMenuVisible(false));
    };
  });
  useEffect(() => {}, [menuVisible]);
  return (
    <div>
      {/* //burger menu or x  */}
      <div className="float-right mt-4 mr-14">
        {menuVisible ? (
          <button onClick={() => setMenuVisible(false)}>
            <img className="h-12" src={x} />
          </button>
        ) : (
          <button onClick={() => setMenuVisible(true)}>
            <img className="h-12" src={list} />
          </button>
        )}
      </div>

      <div
        className={`float-right absolute right-0 top-20 ${
          menuVisible
            ? "animate-slide-in animate-duration-[350ms] animate-ease-out"
            : "animate-slide-out animate-duration-[350ms] animate-ease-out"
        }`}
      >
        <div
          className={`w-fit p-5 float-right h-screen bg-black opacity-80 ${
            show ? "" : "hidden"
          }`}
        >
          <Link to="/" className={MOBILE_CSS}>
            {" "}
            Strona główna
          </Link>
          <Link to="/blog" className={MOBILE_CSS}>
            {" "}
            Blog{" "}
          </Link>
          <Link to="/gallery" className={MOBILE_CSS}>
            {" "}
            Galeria{" "}
          </Link>
          <Link to="/about" className={MOBILE_CSS}>
            {" "}
            O nas{" "}
          </Link>
          <Link to="/pricing" className={MOBILE_CSS}>
            {" "}
            Cennik{" "}
          </Link>
        </div>
      </div>
    </div>
  );
};

export default Sidebar;
