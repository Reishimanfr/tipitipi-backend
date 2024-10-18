import { useState, useEffect } from "react";
import validateToken from "../../../functions/validate";
import Unauthorized from "../../errorPages/unauthorized";

const GalleryAdd = () => {
  const [images, setImages] = useState<FileList | null>();

  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);
  useEffect(() => {
    const ValidateAuthorization = async () => {
      setIsAuthorized(await validateToken(setLoading));
    };
    ValidateAuthorization();
  }, []);
  if (loading) {
    return <div>Loading</div>;
  }
  if (!isAuthorized) {
    return <Unauthorized />;
  }
  return (
    <div className="mt-[1%] globalCss">
      <label className="text-3xl font-bold" htmlFor="image">Dodaj zdjęcia</label>
      <br></br>

      <input
      className="my-[2%]"
        type="file"
        name="image"
        accept="image/*"
        onChange={(e) => {
          setImages(e.target.files);
        }}
        multiple
      /><br></br>

      <button
        className={
          "border w-40 shadow-lg hover:bg-slate-100 hover:duration-300"
        }
        onClick={() => alert("dodano zdjecia ( wcalen ie)")}
      >
        Postuj
      </button>
    </div>
  );
};

export default GalleryAdd;
