import { Buffer } from "buffer"
import { useEffect, useState } from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Footer from './components/footer'
import Navbar from './components/navbar'
import About from "./pages/about"
import Admin from './pages/admin/admin'
import Dashboard from './pages/admin/dashboard'
import PostCreating from './pages/admin/dashboardPages/postCreating'
import Login from './pages/admin/login'
import Blog from "./pages/blog"
import Unauthorized from './pages/errorPages/unauthorized'
import Gallery from "./pages/gallery"
import Mainpage from "./pages/mainpage"
import Pricing from "./pages/pricing"
const decode = (str: string):string => Buffer.from(str, 'base64').toString('binary');


interface JwtPayload {
    iat: number
    exp: number
    admin: boolean
    user_id: string
} 

function App() {
  const [mainpageFirstHeader,setMainpageFirstHeader] = useState("")
  const changeMainpageFirstHeader = (newMessage : string) => {
    setMainpageFirstHeader(newMessage)
  }

  const [isAuthorized,setIsAuthorized] = useState(false)
  useEffect(() => {
    const token = localStorage.getItem("token")

    if (token === null) {
        console.debug("Token is invalid")
        setIsAuthorized(false)
        return
    }

    // 0: Header | 1: Payload | 2: Signature
    const tokenSplit = token.split(".")

    if (!tokenSplit?.[1]) {
        console.error("Malformed token")
        setIsAuthorized(false)
        return
    }

    const stringPayload = decode(tokenSplit[1])
    const payload: JwtPayload = JSON.parse(stringPayload)
    const now = Date.now() / 1000

    if (now >= payload.exp) {
        console.debug("Token expired, redirecting to login page...")
        setIsAuthorized(false)
        return
    }

    setIsAuthorized(true)
})


 

  return (
    <div className='relative min-h-screen'>
      <BrowserRouter>
      <Navbar/>
      <Routes>
        <Route path='/admin' element ={<Admin isAuthorized={isAuthorized}/>} />
        <Route path="/admin/dashboard" element={isAuthorized ? 
          <Dashboard  mainpageFirstHeader={mainpageFirstHeader} changeMainpageFirstHeader={changeMainpageFirstHeader}/> : 
          <Unauthorized/>} />
        <Route path="/admin/login" element={<Login/>}/>
        {/* <Route path='/admin/dashboard/create-post' element={isAuthorized ? <PostCreating/> : <Unauthorized/>}/> */}
        <Route path='/admin/dashboard/create-post' element={<PostCreating/>}/>
        {/* <Route path="/" element={<Mainpage mainpageFirstHeader={mainpageFirstHeader}/>}/> */}
        <Route path="/" element={<Mainpage/>}/>
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
