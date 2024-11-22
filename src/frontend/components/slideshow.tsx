import { useState } from "react";

interface Props {
  images: string[];
}

const SlideShow = ({ images }: Props) => {
  const [currentIndex, setCurrentIndex] = useState(0);

  const goToSlide = (index: number) => {
    setCurrentIndex(index);
  };

  const goToPrevious = () => {
    setCurrentIndex((prevIndex) =>
      prevIndex === 0 ? images.length - 1 : prevIndex - 1
    );
  };

  const goToNext = () => {
    setCurrentIndex((prevIndex) =>
      prevIndex === images.length - 1 ? 0 : prevIndex + 1
    );
  };

  return (
    <div className="flex flex-col items-center space-y-4 border border-black relative ">
      {/* Wybrany slajd powiększony */}
      <div className="w-full max-w-xl flex justify-center items-center relative h-[500px]">
        {/* Przycisk strzałki w lewo */}
        <button
          onClick={goToPrevious}
          className="absolute left-4 top-1/2 transform -translate-y-1/2 text-white bg-gray-800 p-2 rounded-full shadow-lg hover:bg-gray-600 z-20"
        >
          &#8592;
        </button>

        {images.length > 0 && (
          <img
            src={images[currentIndex]}
            alt={`Slide ${currentIndex}`}
            className="rounded-lg shadow-lg transition-transform transform duration-300 object-contain max-h-full max-w-full"
          />
        )}

        {/* Przycisk strzałki w prawo */}
        <button
          onClick={goToNext}
          className="absolute right-4 top-1/2 transform -translate-y-1/2 text-white bg-gray-800 p-2 rounded-full shadow-lg hover:bg-gray-600 z-20"
        >
          &#8594;
        </button>
      </div>

      {/* Pasek miniatur */}
      <div className="flex overflow-x-auto space-x-4 w-full justify-center px-4 overflow-y-hidden h-[150px]">
        {images.map((image, index) => (
          <div
            key={index}
            className={`${
              index === currentIndex ? "scale-105 " : "opacity-60"
            } cursor-pointer transition-transform transform duration-300`}
            onClick={() => goToSlide(index)}
          >
            <img
              src={image}
              alt={`Thumbnail ${index}`}
              className="h-36 object-cover rounded-md max-h-[144px]"
            />
          </div>
        ))}
      </div>
    </div>
  );
};

export default SlideShow;
