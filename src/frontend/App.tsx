import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Navbar from './components/navbar';
import Mainpage from "./pages/mainpage"
import Blog from "./pages/blog";
import Gallery from "./pages/gallery";
import About from "./pages/about";
import Pricing from "./pages/pricing";
import Footer from './components/footer';
import Admin from './pages/admin';
import { useState } from 'react';

function App() {

  const [mainpageFirstHeader,setMainpageFirstHeader] = useState("")
  const changeMainpageFirstHeader = (newMessage : string) => {
    setMainpageFirstHeader(newMessage)
  }

  return (
    <div className='relative min-h-screen'>
      <BrowserRouter>
      <Navbar/>
      <Routes>
        <Route path="/admin" element={<Admin  mainpageFirstHeader={mainpageFirstHeader} changeMainpageFirstHeader={changeMainpageFirstHeader}/>}/>
        <Route path="/" element={<Mainpage mainpageFirstHeader={mainpageFirstHeader}/>}/>
        <Route path="/gallery" element={<Gallery/>}/>
        <Route path="/about" element={<About/>}/>
        <Route path="/pricing" element={<Pricing/>}/>
        <Route path="/blog" element={<Blog/>}/>
      </Routes>
      <Footer/>
      </BrowserRouter>
    </div> 
  )
}

export default App
