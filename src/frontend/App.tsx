import { BrowserRouter, Route, Routes } from 'react-router-dom';
import Navbar from './components/navbar';
import Mainpage from "./pages/mainpage"
import Blog from "./pages/blog";
import Gallery from "./pages/gallery";
import About from "./pages/about";
import Pricing from "./pages/pricing";

function App() {

  return (
    <>
      <BrowserRouter>
      <Navbar/>
      <Routes>
        <Route path="/" element={<Mainpage/>}/>
        <Route path="/gallery" element={<Gallery/>}/>
        <Route path="/about" element={<About/>}/>
        <Route path="/pricing" element={<Pricing/>}/>
        <Route path="/blog" element={<Blog/>}/>
      </Routes>
      </BrowserRouter>
    </> 
  )
}

export default App
