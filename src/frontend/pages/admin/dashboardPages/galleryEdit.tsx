import { useEffect, useState } from "react";
import validateToken from "../../../functions/validate";
import Unauthorized from "../../errorPages/unauthorized";
import { GalleryGroup , GalleryCreateNewJson } from "../../../functions/interfaces";
import { getToken } from "../../../functions/postManipulatingFunctions";

async function deleteImages(id : number) {
    if (!window.confirm("Czy napewno chcesz usunąć wszystkie zdjęcia?")) {
        return;
      }
      const token = getToken();

    try {
        const response = await fetch(
          `http://localhost:2333/gallery/groups/${id}/images`,
          {
            method: "DELETE",
            headers: {
              Authorization: `Bearer ${token}`,
            }
          }
        );
  
        if (response.status >= 200 && response.status < 300) {
          alert("Usunięto zdjęcia");
          window.location.reload();
        } else {
          const data: GalleryCreateNewJson = await response.json();
          alert("Błąd: " + data.error);
        }
      } catch (error) {
        console.error(error);
        alert("Wystąpił błąd: " + error);
      }
}
async function deleteGroup(id : number) {
    if (!window.confirm("Czy napewno chcesz usunąć cały album?")) {
        return;
      }
      const token = getToken();

    try {
        const response = await fetch(
          `http://localhost:2333/gallery/groups/${id}`,
          {
            method: "DELETE",
            headers: {
              Authorization: `Bearer ${token}`,
            }
          }
        );
  
        if (response.status >= 200 && response.status < 300) {
          alert("Usunięto album");
          window.location.reload();
        } else {
          const data: GalleryCreateNewJson = await response.json();
          alert("Błąd: " + data.error);
        }
      } catch (error) {
        console.error(error);
        alert("Wystąpił błąd: " + error);
      }
}

async function fetchGroups(setGroups : React.Dispatch<React.SetStateAction<GalleryGroup[]>>) {
    try {
      const response = await fetch(
        `http://localhost:2333/gallery/groups/all/info`,
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


  async function getInfo(id:number) {
    try {
        const response = await fetch(
          `http://localhost:2333/gallery/groups/${id}/info`,
          {
            method: "GET",
          }
        );
        if (!response.ok) {
          throw new Error(response.statusText);
        }
  
        const data: GalleryGroup = await response.json();
        console.log(data)
      } catch (error) {
        console.error(error);
      }
  }


const GalleryEdit = () => {
  const [groups, setGroups] = useState<GalleryGroup[]>([]);
  const [selectedGroup, setSelectedGroup] = useState<GalleryGroup | null>();

  useEffect(() =>{
    if(selectedGroup != null) {
     getInfo(selectedGroup.id)
    }
  },[selectedGroup])

  //public.GET("/gallery/groups/:groupId/info", s.GalleryGetGroupOne)
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
        await fetchGroups(setGroups);
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
      <label>Wybierz album/grupe którą chcesz edytować</label>
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
      <hr></hr>
      {/* ---------------------------------------- */}
      {selectedGroup ? (
        <div>
          <h1 className="text-2xl font-bold">{selectedGroup.name}</h1><br></br>
          <button className="border w-40 text-white shadow-lg bg-red-500 hover:bg-red-600 hover:duration-300 mr-[2%]" onClick={() => deleteImages(selectedGroup.id)}>Usuń zdjęcia</button>
          <button className="border w-40 text-white shadow-lg bg-red-500 hover:bg-red-600 hover:duration-300" onClick={() => deleteGroup(selectedGroup.id)}>Usuń album</button>
        </div>
      ) : (
        <div></div>
      )}
    </div>
  );
};

export default GalleryEdit;
