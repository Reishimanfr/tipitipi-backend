// import { List, X } from "phosphor-react"
// import React, { useEffect, useState } from "react"
// import { Link } from "react-router-dom"
// import logoSmall from "../react.svg"
// import RenderNavbar from "./deprecated_renderNavbar"

// function toggleMenuVisibility(visible: boolean) {
//   const xAndListStyle = "text-white float-right mr-12";
//   const iconSize = 32;

//   if (visible) {
//     return <X className={xAndListStyle} size={iconSize} />;
//   }

//   return <List className={xAndListStyle} size={iconSize} />;
// }

// //ta funkcja decyduje czy generuje komponent navlinks czy alternativenavlinks , jest tak skomplikowana by animacja css dzialala
// function RenderMobileMenu(
//   render: boolean,
//   shouldRender: boolean,
//   setShouldRender: React.Dispatch<React.SetStateAction<boolean>>
// ) {
//   useEffect(() => {
//     if (!render) {
//       // Opóźnij ukrycie komponentu o czas trwania animacji (np. 300ms)
//       const timeoutId = window.setTimeout(() => {
//         setShouldRender(false);
//       }, 500);

//       // Wyczyszczenie timeoutu po odmontowaniu komponentu
//       return () => clearTimeout(timeoutId);
//     } else {
//       setShouldRender(true);
//     }
//   }, [render]);

//   if (shouldRender) {
//     return (
//       <div
//         className={`float-right fixed right-0 top-20 ${
//           render
//             ? "animate-fade-left animate-duration-[350ms] animate-ease-out"
//             : "animate-slide-out"
//         }`}
//       >
//         <div className="opacity-80">
//           <RenderNavbar mobile={true} />
//         </div>
//       </div>
//     );
//   }

//   return (
//     <div className="hidden lg:flex">
//       <RenderNavbar mobile={false} />
//     </div>
//   );
// }

// function Navbar() {
//   const [menuVisible, setMenuVisible] = useState(false);
//   const [shouldRender, setShouldRender] = useState(false);
//   const [lastScrollY, setLastScrollY] = useState(0);
//   const [show, setShow] = useState(true);
//   //u góry mamy state od tego że menu dla urzadzen mobilnych ma byc wyswietlane

//   const handleResize = () => {
//     if (menuVisible) {
//       setMenuVisible(window.innerWidth <= 768);
//     }
//   };

//   const handleScroll = () => {
//     if (typeof window !== 'undefined') {
//       setShow(window.scrollY < lastScrollY)
//       setLastScrollY(window.scrollY);
//     }
//   };

//   useEffect(() => {
//     window.addEventListener("scroll" , handleScroll);
//     return () => {
//       window.removeEventListener("scroll" , handleScroll)
//       }
//   },[lastScrollY]);
//   useEffect(() => {
//     window.addEventListener("resize", handleResize);
//     return () => {
//       window.removeEventListener("resize", handleResize)
//       }
//   });

//   return (
//     <div className={`sticky z-50 top-0 h-20 bg-black w-full transition-transform duration-300 transform ${show ? 'translate-y-0' : '-translate-y-full'}`}>
//       <div>
//         <Link to="/">
//           <img
//             src={logoSmall}
//             alt="logo"
//             className="w-auto h-16 float-left ml-[5%] mt-[1vh] hover:animate-spin"
//           ></img>
//         </Link>
//         <button
//           className="float-right h-20 lg:hidden"
//           onClick={(_) => setMenuVisible(!menuVisible)}
//         >
//           {toggleMenuVisibility(menuVisible)}
//         </button>
//         </div>
//       {RenderMobileMenu(menuVisible, shouldRender, setShouldRender)}
//     </div>
//     //wyświetla navlinks , jeżeli urządzneie jest male pojawia sie przycisk ktory zmienia stan menu , wyswietla lub nie , props definiuje
//     //czy zdjecie ma byc czy nie , zeby nie pojawilo sie 2 razy jest to potrzebne
//   );
// }

// export default Navbar;
