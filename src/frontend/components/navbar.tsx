import { useEffect, useState } from "react";
import { List, X } from "phosphor-react";
import React from "react";
import { Link } from "react-router-dom";
import logoSmall from "../react.svg";
import RenderNavbar from "./renderNavbar";

function toggleMenuVisibility(visible: boolean) {
  const xAndListStyle = "text-white float-right mr-12";
  const iconSize = 32;

  if (visible) {
    return <X className={xAndListStyle} size={iconSize} />;
  }

  return <List className={xAndListStyle} size={iconSize} />;
}

//ta funkcja decyduje czy generuje komponent navlinks czy alternativenavlinks , jest tak skomplikowana by animacja css dzialala
function RenderMobileMenu(
  render: boolean,
  shouldRender: boolean,
  setShouldRender: React.Dispatch<React.SetStateAction<boolean>>
) {
  useEffect(() => {
    if (!render) {
      // Opóźnij ukrycie komponentu o czas trwania animacji (np. 300ms)
      const timeoutId = window.setTimeout(() => {
        setShouldRender(false);
      }, 500);

      // Wyczyszczenie timeoutu po odmontowaniu komponentu
      return () => clearTimeout(timeoutId);
    } else {
      setShouldRender(true);
    }
  }, [render]);

  if (shouldRender) {
    return (
      <div
        className={`float-right fixed right-0 top-20 ${
          render
            ? "animate-slide-in fill-forwards"
            : "animate-slide-out fill-forwards"
        }`}
      >
        <div className="opacity-80">
          <RenderNavbar mobile={true} />
        </div>
      </div>
    );
  }

  return (
    <div className="hidden lg:flex">
      <RenderNavbar mobile={false} />
    </div>
  );
}

function Navbar() {
  const [menuVisible, setMenuVisible] = useState(false);
  const [shouldRender, setShouldRender] = useState(false);
  const [shouldNavbarHide , setShouldNavbarHide] = useState("top-0");
  //u góry mamy state od tego że menu dla urzadzen mobilnych ma byc wyswietlane

  const handleResize = () => {
    console.log(menuVisible);
    if (menuVisible) {
      setMenuVisible(window.innerWidth <= 768);
    }
  };


  useEffect(() => {
    window.addEventListener("resize", handleResize);
    return () => window.removeEventListener("resize", handleResize);
  });

  //ustawia menuVisible na fałsz kiedy okno jest duże , żeby nie było 2 duplikatów i 2 navbarow , po prostu idiotoodpornosc

  return (
    <div className="w-full h-20 bg-black fixed ">
      <div>
        <Link to="/">
          <img
            src={logoSmall}
            alt="logo"
            className="w-auto h-16 float-left ml-[5%] mt-[1vh]"
          ></img>
        </Link>
        <button
          className="float-right h-20 lg:hidden"
          onClick={(_) => setMenuVisible(!menuVisible)}
        >
          {toggleMenuVisibility(menuVisible)}
        </button>
      </div>
      {RenderMobileMenu(menuVisible, shouldRender, setShouldRender)}
    </div>
    //wyświetla navlinks , jeżeli urządzneie jest male pojawia sie przycisk ktory zmienia stan menu , wyswietla lub nie , props definiuje
    //czy zdjecie ma byc czy nie , zeby nie pojawilo sie 2 razy jest to potrzebne
  );
}

export default Navbar;
