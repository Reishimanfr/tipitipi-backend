import { Link } from "react-router-dom"

const MOBILE_CSS = "mt-6 clear-right float-right no-underline text-white text-xl mr-2.5 bg-black-700 font-sans hover:border border-black-700 hover:bg-white hover:text-black" as const
const DESKTOP_CSS = "float-right no-underline text-white text-3xl mr-2.5 bg-black-700 font-sans hover:duration-500  hover:bg-white hover:text-black" as const

const DIV_MOBILE_CSS = "w-fit  h-screen pt-0 float-right bg-black" as const
const DIV_DESKTOP_CSS = "w-full h-20 flex-wrap bg-black sticky top-0 p-2.5 text-right" as const

function getUrls(mobile: boolean) {
    const css = mobile ? MOBILE_CSS : DESKTOP_CSS
  
    const urls = [
      <Link to='/' className={css}> Strona główna</Link>,
      <Link to='/blog' className={css}> Blog </Link>,
      <Link to='/gallery' className={css}> Galeria </Link>,
      <Link to='/pricing' className={css}> Cennik </Link>,
      <Link to='/about' className={css}> O nas </Link>
    ];
  
    if (mobile) {
      for (let i = 0; i < urls.length; i++) {
        urls.splice(i + 1, 0, <br />)
        i++
      }
    }
  
    return mobile ? urls : urls.reverse();
  }
  
  function RenderNavbar(props: { mobile: boolean }) {
    return (
      <div className={props.mobile ? DIV_MOBILE_CSS : DIV_DESKTOP_CSS}>
        <div>
          {getUrls(props.mobile)}
        </div>
      </div>
    )
  }
  
  export default RenderNavbar