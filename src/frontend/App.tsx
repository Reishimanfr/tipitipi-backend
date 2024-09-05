import { useState } from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Footer from './components/footer'
import Navbar from './components/navbar'
import About from "./pages/about"
import Admin from './pages/admin'
import PostCreating from './pages/admin/dashboardPages/postCreating'
import Blog from "./pages/blog"
import Dashboard from './pages/dashboard'
import Gallery from "./pages/gallery"
import Login from './pages/login'
import Mainpage from "./pages/mainpage"
import Pricing from "./pages/pricing"


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
        <Route path='/admin' element ={<Admin/>}/>
        <Route path="/admin/dashboard" element={<Dashboard  mainpageFirstHeader={mainpageFirstHeader} changeMainpageFirstHeader={changeMainpageFirstHeader}/>}/>
        <Route path="/admin/login" element={<Login/>}/>
        <Route path='/admin/dashboard/create-post' element={<PostCreating/>}/>
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
