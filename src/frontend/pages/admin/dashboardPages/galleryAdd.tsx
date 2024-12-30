import { useEffect, useState } from "react"
import { GalleryGroup } from "../../../functions/interfaces"
import {
        buildGalleryMultipart,
        getToken,
} from "../../../functions/postManipulatingFunctions"
import validateToken from "../../../functions/validate"
import Unauthorized from "../../errorPages/unauthorized"

import { toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'
import { API_URL } from '../../../functions/global'


async function addNewGroup(name: string , setNewGroupName: React.Dispatch<React.SetStateAction<string>>) {
  const token = getToken();
  if (name == "") {
    toast.warn("Nie podano nazwy nowego albumu");
    return;
  }
  if (!window.confirm("Czy napewno chcesz dodać nowy album?")) {
    return;
  }
  try {
    const response = await fetch(
      `${API_URL}/gallery/groups/new/${name}`,
      {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    if (response.status >= 200 && response.status < 300) {
      setNewGroupName("")
      toast.success("Dodano album")
      window.location.reload()
    } 
    else{
      throw new Error(response.statusText);
    }
  } catch (error) {
    console.error(error);
    toast.error("Wystąpił błąd: " + error);
  }
}

const GalleryAdd = () => {
  const [newGroupName, setNewGroupName] = useState("");
  const [groups, setGroups] = useState<GalleryGroup[]>([]);
  const [selectedGroup, setSelectedGroup] = useState<GalleryGroup | null>();
  const [images, setImages] = useState<FileList | null>();

  async function fetchGroups() {
    try {
      const response = await fetch(
        `${API_URL}/gallery/groups/all/info`,
        {
          method: "GET",
        }
      );
      if (!response.ok) {
        throw new Error(response.statusText);
      }

      const data: Array<GalleryGroup> = await response.json();
      setGroups((prevGroups) => prevGroups?.concat(data));
    } catch (error) {
      console.error(error);
    }
  }

  async function addImages() {
    if (selectedGroup == null) {
      toast.warn("Nie wybrano do którego albumu docelowego");
      return;
    }
    if (images?.length == null) {
      toast.warn("Nie wybrano zdjęć");
      return;
    }
    if (!window.confirm("Czy napewno chcesz dodać zdjęcia?")) {
      return;
    }
    const token = getToken();
    const formData = buildGalleryMultipart(images);
    try {
      const response = await fetch(
        `${API_URL}/gallery/groups/${selectedGroup.id}/images`,
        {
          method: "POST",
          headers: {
            Authorization: `Bearer ${token}`,
          },
          body: formData,
        }
      );

      if (response.status >= 200 && response.status < 300) {
        toast.success("Dodano zdjęcia");
        setImages(null)
      } else{
        throw new Error(response.statusText);
      }
    } catch (error) {
      console.error(error);
      
      toast.error("Wystąpił błąd: " + error);
    }
  }

  const [loading, setLoading] = useState(true);
  const [isAuthorized, setIsAuthorized] = useState(false);
  useEffect(() => {
    const ValidateAuthorization = async () => {
      setIsAuthorized(await validateToken(setLoading));
    };
    ValidateAuthorization();
  }, []);
  useEffect(() => {
    const fetchGroupsEffect = async () => {
      if (isAuthorized && groups.length == 0) {
        await fetchGroups();
      }
    };
    fetchGroupsEffect();
  }, [isAuthorized]);
  if (loading) {
    return <div>Loading</div>;
  }
  if (!isAuthorized) {
    return <Unauthorized />;
  }
  return (
    <div className="globalCss mt-[1%]">
      <div>
        <h1 className="text-3xl font-bold mb-6">Tworzenie albumów/grup</h1>
        <label className="text-xl" htmlFor="newAlbum">
          Podaj nazwę nowego albumu:{" "}
        </label>
        <input
          className="border-2 "
          type="text"
          name="newAlbum"
          value={newGroupName}
          onChange={(e) => setNewGroupName(e.target.value)}
        />
        <br></br>
        <button
          className={
            "border w-40 bg-white shadow-lg hover:bg-slate-100 hover:duration-300 mt-6"
          }
          onClick={() => addNewGroup(newGroupName,setNewGroupName)}
        >
          Stwórz nowy album
        </button>
      </div>

      <br></br>
      <hr></hr>
      <br></br>

      <div>
        <h1 className="text-3xl font-bold mb-6">Dodawanie zdjęć</h1>

        <label htmlFor="groups" className="text-2xl">
          Do której grupy chcesz dodać zdjęcia?
        </label>
        <br></br>
        <select
          className="mb-8"
          name="groups"
          onChange={(e) => setSelectedGroup(groups[parseInt(e.target.value)])}
        >
          <option value="">--albumy--</option>
          {groups ? (
            groups.map((group, index) => {
              return (
                <option key={group.id} value={index}>
                  {group.id + " , " + group.name}
                </option>
              );
            })
          ) : (
            <div>No group found</div>
          )}
        </select>
        <br></br>

        <label className="text-2xl" htmlFor="image">
          Dodaj zdjęcia
        </label>
        <br></br>
        <input
          className="mb-8"
          type="file"
          name="image"
          accept="image/*"
          value={images ? undefined : ""}
          onChange={(e) => {
            setImages(e.target.files);
          }}
          multiple
        />
        <br></br>

        <button
          className={
            "border w-40 bg-white shadow-lg hover:bg-slate-100 hover:duration-300"
          }
          onClick={() => addImages()}
        >
          Dodaj
        </button>
      </div>
    </div>
  );
};

export default GalleryAdd;
