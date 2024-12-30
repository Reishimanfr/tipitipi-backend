import { useEffect, useState } from "react"
import { toast } from "react-toastify"
import { API_URL } from '../../../functions/global'
import { GroupInfo } from "../../../functions/interfaces"
import { getToken } from "../../../functions/postManipulatingFunctions"
import validateToken from "../../../functions/validate"
import Unauthorized from "../../errorPages/unauthorized"

const RED_BUTTON_CSS =
  "border w-40 text-white shadow-lg bg-red-500 hover:bg-red-600 hover:duration-300";



async function deleteGroup(id: number) {
  if (!window.confirm("Czy napewno chcesz usunąć cały album?")) {
    return;
  }
  const token = getToken();

  try {
    const response = await fetch(`${API_URL}/gallery/groups/${id}`, {
      method: "DELETE",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (response.status >= 200 && response.status < 300) {
      toast.success("Usunięto album");
      window.location.reload();
    }
    if (!response.ok) {
      throw new Error(response.statusText);
    }
  } catch (error) {
    console.error(error);
    toast.error("Wystąpił błąd: " + error);
  }
}


const GalleryEdit = () => {
  const [groups, setGroups] = useState<Array<GroupInfo> | null>();
  const [selectedGroup, setSelectedGroup] = useState<GroupInfo | null>();

  async function deleteImage(GroupID : number , imageID:number) {

    if (!window.confirm("Czy napewno chcesz usunąć to zdjęcie?")) {
      return;
    }
    const token = getToken();
    try {
      const response = await fetch(
        `${API_URL}/gallery/groups/${GroupID}/images/${imageID}`,
        {
          method: "DELETE",
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );
      if (!response.ok) {
        throw new Error(response.statusText);
      }
  
      if (response.status >= 200 && response.status < 300) {
        toast.success("Usunięto zdjęcie");
        setSelectedGroup((prevGroup) => {
          if (!prevGroup) return null;
  
          return {
            ...prevGroup,
            images: prevGroup.images.filter((image) => image.id !== imageID),
          };
        });
      }
    } catch (error) {
      console.error(error);
    }
  }


  async function deleteImages(id: number) {
  if (!window.confirm("Czy napewno chcesz usunąć wszystkie zdjęcia?")) {
    return;
  }
  const token = getToken();

  try {
    const response = await fetch(
      `${API_URL}/gallery/groups/${id}/images`,
      {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    if (response.status >= 200 && response.status < 300) {
      toast.success("Usunięto zdjęcia");
      setSelectedGroup((prevGroup) => {
        if (!prevGroup) return null;

        return {
          ...prevGroup,
          images: [],
        };
      });
    }
    if (!response.ok) {
      throw new Error(response.statusText);
    }
  } catch (error) {
    console.error(error);
    toast.error("Wystąpił błąd: " + error);
  }
}
  useEffect(() => {
    async function fetchPost() {
      try {
        const response = await fetch(
          `${API_URL}/gallery/everything`,
          {
            method: "GET",
          }
        );
        if (!response.ok) {
          throw new Error(response.statusText);
        }

        const data = await response.json();
        setGroups(data);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    }
    fetchPost();
  }, []);



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
    <div className="globalCss mt-[1%]">
      <label>Wybierz album/grupe którą chcesz edytować</label>
      <br></br>
      <select
        className="mb-8"
        name="groups"
        onChange={(e) => setSelectedGroup(groups![parseInt(e.target.value)])}
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
          <option/>
        )}
      </select>
      <hr></hr>
      {/* ---------------------------------------- */}
      {selectedGroup ? (
        <div >
          <h1 className="text-2xl text-center font-bold">{selectedGroup.name}</h1>
          <br></br>
          {selectedGroup.images ? (
            selectedGroup.images.map((image) => {
              return (
                <div className="p-2 mx-auto bg-white border w-1/2 m-2" key={image.id}>
                  <p className="text-center">Zdjęcie numer: {image.id}</p>
                  <img className="max-h-[200px] mx-auto my-6" src={`${API_URL}/proxy?key=${image.filename}&type=gallery`}/>
                  <button
                    className={`border w-full ${RED_BUTTON_CSS}`}
                    onClick={() => {deleteImage(selectedGroup.id , image.id)}}
                  >
                    Usuń zdjęcie
                  </button>
                </div>
              );
            })
          ) : (
            <div></div>
          )}
          <div className="text-center">

            <button
              className={`${RED_BUTTON_CSS} m-10`}
              onClick={() => deleteImages(selectedGroup.id)}
            >
              Usuń zdjęcia
            </button>
            <button
              className={RED_BUTTON_CSS}
              onClick={() => deleteGroup(selectedGroup.id)}
            >
              Usuń album
            </button>
          </div>
        </div>
      ) : (
        <div></div>
      )}
    </div>
  );
};

export default GalleryEdit;
