import { useState } from "react";

interface Props {
  images: string[];
}

const SlideShow = ({ images }: Props) => {
  const [currentIndex, setCurrentIndex] = useState(0);

  // Przejście do wybranego slajdu
  const goToSlide = (index: number) => {
    setCurrentIndex(index);
  };

  return (
    <div className="flex flex-col items-center space-y-4">
      {/* Wybrany slajd powiększony */}
      <div className="w-full max-w-xl">
        {images.length > 0 && (
          <img
            src={images[currentIndex]}
            alt={`Slide ${currentIndex}`}
            className="w-full h-[400px] object-cover rounded-lg shadow-lg transition-transform transform duration-300 scale-105"
          />
        )}
      </div>

      {/* Pasek miniatur */}
      <div className="flex overflow-x-auto space-x-4">
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
              className="w-32 h-32 object-cover rounded-md"
            />
          </div>
        ))}
      </div>
    </div>
  );
};

export default SlideShow;



