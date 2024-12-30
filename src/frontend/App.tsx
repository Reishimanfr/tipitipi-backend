import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { ToastContainer } from 'react-toastify'
import Footer from './components/footer'
import MainNavbar from './components/mainNavbar'
import About from "./pages/about"
import Admin from './pages/admin/admin'
import Dashboard from './pages/admin/dashboard'
import ChangeCredentials from './pages/admin/dashboardPages/changeCredentials'
import GalleryAdd from './pages/admin/dashboardPages/galleryAdd'
import GalleryEdit from './pages/admin/dashboardPages/galleryEdit'
import PostCreating from './pages/admin/dashboardPages/postCreating'
import PostEditing from './pages/admin/dashboardPages/postEditing'
import Login from './pages/admin/login'
import Blog from "./pages/blog"
import PageNotFound from './pages/errorPages/page_not_found'
import Gallery from "./pages/gallery"
import Mainpage from "./pages/mainpage"
import PostPage from './pages/postPage'
import Pricing from "./pages/pricing"



// import NewNavbar from './components/newNavbar'

function App() {
  // const [mainpageFirstHeader,setMainpageFirstHeader] = useState("")
  // const changeMainpageFirstHeader = (newMessage : string) => {
  //   setMainpageFirstHeader(newMessage)
  // }

  return (
    <div className="flex flex-col min-h-screen">
      {/* <Monitoring apiKey="AnqmUAZGoTDxYC7R9b3aZIGoEn8NoHS_" params={{}} path="/" url="https://monitoring.react-scan.com/api/v1/ingest"/> */}
      <ToastContainer
        position="top-center"
        autoClose={2000}
        hideProgressBar={false}
        closeOnClick
        rtl={false}
        pauseOnFocusLoss
        draggable
        pauseOnHover
        theme="colored"/>
      <BrowserRouter>
      <MainNavbar/>
      <div className='flex-grow'>
      <Routes>
        {/* admin pages */}
        <Route path='/admin' element ={<Admin/>} />
        {/* <Route path="/admin/dashboard" element={isAuthorized ? 
          <Dashboard  mainpageFirstHeader={mainpageFirstHeader} changeMainpageFirstHeader={changeMainpageFirstHeader}/> : 
          <Unauthorized/>} /> */}
        <Route path="/admin/dashboard" element ={<Dashboard />}/>
        <Route path='/admin/dashboard/create-post' element={<PostCreating/>}/> 
        <Route path='/admin/dashboard/edit-post' element={<PostEditing/>}/>
        <Route path='/admin/dashboard/gallery-add' element={<GalleryAdd/>}/>
        <Route path='/admin/dashboard/gallery-edit' element={<GalleryEdit/>}/>
        <Route path='/admin/dashboard/change-credentials' element={<ChangeCredentials/>}/>
        <Route path="/admin/login" element={<Login/>}/>


        {/* user pages */}
        <Route path="/" element={<Mainpage/>}/>
        <Route path="/gallery" element={<Gallery/>}/>
        <Route path="/about" element={<About/>}/>
        <Route path="/pricing" element={<Pricing/>}/>
        <Route path="/blog" element={<Blog/>}/>
        <Route path="/blog/:id" element={<PostPage/>}/>
        <Route path='*' element={<PageNotFound/>}/>
      </Routes>
      </div>

      <Footer/>
      </BrowserRouter>
    </div> 
  )
}

export default App
