import { useState} from 'react'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import Footer from './components/footer'
import Navbar from './components/navbar'
import About from "./pages/about"
import Admin from './pages/admin/admin'
import PostCreating from './pages/admin/dashboardPages/postCreating'
import Blog from "./pages/blog"
import Dashboard from './pages/admin/dashboard'
import Gallery from "./pages/gallery"
import Login from './pages/admin/login'
import Mainpage from "./pages/mainpage"
import Pricing from "./pages/pricing"
import PostPage from './pages/postPage'
import PostEditing from './pages/admin/dashboardPages/postEditing'
import ChangeCredentials from './pages/admin/dashboardPages/changeCredentials'
import PageNotFound from './pages/errorPages/page_not_found'

function App() {
  const [mainpageFirstHeader,setMainpageFirstHeader] = useState("")
  const changeMainpageFirstHeader = (newMessage : string) => {
    setMainpageFirstHeader(newMessage)
  }
 

  return (
    <div className='relative min-h-screen pb-20'>
      <BrowserRouter>
      <Navbar/>
      <Routes>
        <Route path='/admin' element ={<Admin/>} />
        {/* <Route path="/admin/dashboard" element={isAuthorized ? 
          <Dashboard  mainpageFirstHeader={mainpageFirstHeader} changeMainpageFirstHeader={changeMainpageFirstHeader}/> : 
          <Unauthorized/>} /> */}
        <Route path="/admin/dashboard" element ={<Dashboard mainpageFirstHeader={mainpageFirstHeader} changeMainpageFirstHeader={changeMainpageFirstHeader}/>}/>
        <Route path='/admin/dashboard/create-post' element={<PostCreating/>}/> 
        <Route path='/admin/dashboard/edit-post' element={<PostEditing/>}/>
        <Route path='/admin/dashboard/change-credentials' element={<ChangeCredentials/>}/>
        <Route path="/admin/login" element={<Login/>}/>
        <Route path="/" element={<Mainpage/>}/>
        <Route path="/gallery" element={<Gallery/>}/>
        <Route path="/about" element={<About/>}/>
        <Route path="/pricing" element={<Pricing/>}/>
        <Route path="/blog" element={<Blog/>}/>
        <Route path="/blog/:id" element={<PostPage/>}/>
        <Route path='*' element={<PageNotFound/>}/>
      </Routes>
      <Footer/>
      </BrowserRouter>
    </div> 
  )
}

export default App
