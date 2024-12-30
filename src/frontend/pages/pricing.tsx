import zdj1 from "../assets/cennik1.jpg";
import zdj2 from "../assets/cennik2.jpg";
import zdj3 from "../assets/cennik3.jpg";
import zdj4 from "../assets/cennik4.jpg";
import zdj5 from "../assets/cennik5.jpg";
import zdj6 from "../assets/cennik6.jpg";
import zdj7 from "../assets/cennik7.jpg";
import zdj8 from "../assets/cennik8.jpg";
import { useState } from "react";

const img_css = "max-w-full max-h-screen object-cover  cursor-pointer";
const Pricing = () => {
  const [isOpen, setIsOpen] = useState(false); 
  const [currentImage, setCurrentImage] = useState(""); 

  const openModal = (imgSrc: any) => {
    setCurrentImage(imgSrc);
    setIsOpen(true);
  };

  const closeModal = () => {
    setIsOpen(false);
    setCurrentImage("");
  };
  return (
    <div className="mx-auto grid grid-cols-3 gap-4 p-4">
      <img className={img_css} onClick={() => openModal(zdj1)} src={zdj1} />
      <img className={img_css} onClick={() => openModal(zdj2)}  src={zdj2} />
      <img className={img_css} onClick={() => openModal(zdj3)}  src={zdj3} />
      <img className={img_css} onClick={() => openModal(zdj4)} src={zdj4} />
      <img className={img_css} onClick={() => openModal(zdj5)}  src={zdj5} />
      <img className={img_css} onClick={() => openModal(zdj6)}  src={zdj6} />
      <img className={img_css} onClick={() => openModal(zdj7)}  src={zdj7} />
      <img className={img_css} onClick={() => openModal(zdj8)}  src={zdj8} />

      {isOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-75 flex justify-center items-center z-50"
          onClick={closeModal} //
        >
          <img
            src={currentImage}
            alt="Pełny obraz"
            className="max-w-full max-h-full object-contain cursor-pointer"
          />
        </div>
      )}
    </div>
  );
};

export default Pricing;
