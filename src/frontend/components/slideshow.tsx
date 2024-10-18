import{ useState } from "react";

interface Props {
  images: string[]; 
}

const SlideShow = ({ images } : Props) => {
  const [currentIndex, setCurrentIndex] = useState(0);

  const nextSlide = () => {
    setCurrentIndex((prevIndex) =>
      prevIndex === images.length - 1 ? 0 : prevIndex + 1
    );
  };

  const prevSlide = () => {
    setCurrentIndex((prevIndex) =>
      prevIndex === 0 ? images.length - 1 : prevIndex - 1
    );
  };

  const goToSlide = (index: number) => {
    setCurrentIndex(index);
  };

  return (
    <div className="slideshow-container relative">
      {/* Render the current image */}
      {images.length > 0 && (
        <img
          src={images[currentIndex]}
          alt={`Slide ${currentIndex}`}
          className="slideshow-image center mx-auto max-h-[200px]"
        />
      )}

      {/* Previous button */}
      <button
        className="prev-button absolute left-0 top-1/2 transform -translate-y-1/2 bg-black text-white p-2"
        onClick={prevSlide}
      >
        &#10094;
      </button>

      {/* Next button */}
      <button
        className="next-button absolute right-0 top-1/2 transform -translate-y-1/2 bg-black text-white p-2"
        onClick={nextSlide}
      >
        &#10095;
      </button>

      {/* Image indicators */}
      <div className="indicators flex justify-center mt-2">
        {images.map((_, index) => (
          <span
            key={index}
            className={`indicator cursor-pointer w-3 h-3 mx-1 rounded-full ${
              index === currentIndex ? "bg-blue-600" : "bg-gray-300"
            }`}
            onClick={() => goToSlide(index)}
          />
        ))}
      </div>
    </div>
  );
};

export default SlideShow;
