// import { Link } from "react-router-dom"


// // const mobile_css = "text-white"

// const MOBILE_CSS = "mt-6 clear-right float-right no-underline text-white text-xl mr-2.5 bg-black-700 font-sans hover:duration-500 hover:bg-gray-300 hover:text-black" as const
// const DESKTOP_CSS = "float-right no-underline text-white text-3xl mr-2.5 bg-black-700 font-sans hover:duration-500  hover:bg-white hover:text-black" as const

// const DIV_MOBILE_CSS = "w-fit p-5 float-right bg-black h-screen"  as const
// const DIV_DESKTOP_CSS = "w-full h-20 flex-wrap bg-black top-20 p-5 text-right" as const

// function getUrls(mobile: boolean) {
//     const divCss = mobile ? MOBILE_CSS : DESKTOP_CSS
  
//     const urls = [
//       <Link key={1} to='/' className={divCss}> Strona główna</Link>,
//       <Link key={2} to='/blog' className={divCss}> Blog </Link>,
//       <Link key={3} to='/gallery' className={divCss}> Galeria </Link>,
//       <Link key={4}to='/pricing' className={divCss}> Cennik </Link>,
//       <Link key={5} to='/about' className={divCss}> O nas </Link>
//     ];
  
//     // If we're on mobile add a new line for each menu
//     if (mobile) {
//       for (let i = 0; i < urls.length; i++) {
//         urls.splice(i + 1, 0, <br />)
//         i++
//       }
//     }
  
//     // Menu entries will appear in the wrong order if array isn't reversed
//     return mobile ? urls : urls.reverse();
//   }
  
//   function RenderNavbar(props: { mobile: boolean }) {
//     return (
//       <div className={props.mobile ? DIV_MOBILE_CSS : DIV_DESKTOP_CSS}>
//         <div>
//           {getUrls(props.mobile)}
//         </div>
//       </div>
//     )
//   }
  
//   export default RenderNavbar