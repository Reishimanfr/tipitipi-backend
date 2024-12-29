import { useEffect, useState } from "react";
import validateToken from "../../../functions/validate";
import Unauthorized from "../../errorPages/unauthorized";
import { GroupInfo } from "../../../functions/interfaces";
import { getToken } from "../../../functions/postManipulatingFunctions";
import { toast } from "react-toastify";

const RED_BUTTON_CSS =
  "border w-40 text-white shadow-lg bg-red-500 hover:bg-red-600 hover:duration-300";

async function deleteImages(id: number) {
  if (!window.confirm("Czy napewno chcesz usunąć wszystkie zdjęcia?")) {
    return;
  }
  const token = getToken();

  try {
    const response = await fetch(
      `http://localhost:8080/gallery/groups/${id}/images`,
      {
        method: "DELETE",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      }
    );

    if (response.status >= 200 && response.status < 300) {
      toast.success("Usunięto zdjęcia");
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

async function deleteImage(GroupID : number , imageID:number) {

  if (!window.confirm("Czy napewno chcesz usunąć to zdjęcie?")) {
    return;
  }
  const token = getToken();
  try {
    const response = await fetch(
      `http://localhost:8080/gallery/groups/${GroupID}/images/${imageID}`,
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
      window.location.reload();
    }
  } catch (error) {
    console.error(error);
  }
}

async function deleteGroup(id: number) {
  if (!window.confirm("Czy napewno chcesz usunąć cały album?")) {
    return;
  }
  const token = getToken();

  try {
    const response = await fetch(`http://localhost:8080/gallery/groups/${id}`, {
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

// async function fetchGroups(
//   setGroups: React.Dispatch<React.SetStateAction<GalleryGroup[]>>
// ) {
//   try {
//     const response = await fetch(
//       `http://localhost:8080/gallery/everything`,
//       {
//         method: "GET",
//       }
//     );
//     if (!response.ok) {
//       throw new Error(response.statusText);
//     }

//     const data: Array<GalleryGroup> = await response.json();
//     setGroups(data);
//   } catch (error) {
//     console.error(error);
//   }
// }

const GalleryEdit = () => {
  const [groups, setGroups] = useState<Array<GroupInfo> | null>();
  const [selectedGroup, setSelectedGroup] = useState<GroupInfo | null>();

  useEffect(() => {
    async function fetchPost() {
      try {
        const response = await fetch(
          `http://localhost:8080/gallery/everything`,
          {
            method: "GET",
          }
        );
        if (!response.ok) {
          throw new Error(response.statusText);
        }

        const data = await response.json();
        // setGroups((prevGroups) => prevGroups?.concat(data));
        setGroups(data);
      } catch (error) {
        console.error(error);
      } finally {
        setLoading(false);
      }
    }
    fetchPost();
  }, []);
  // const [selectedGroupImages, setSelectedGroupImages] = useState<
  //   GalleryImage[]
  // >([]);



  // async function getImagesFromOneGroup(id: number) {
  //   try {
  //     const response = await fetch(
  //       `http://localhost:8080/gallery/groups/${id}/images`,
  //       {
  //         method: "GET",
  //       }
  //     );
  //     if (!response.ok) {
  //       throw new Error(response.statusText);
  //     }

  //     const data: GalleryImage[] = await response.json();
  //     setSelectedGroupImages(data);
  //   } catch (error) {
  //     console.error(error);
  //   }
  // }

  // useEffect(() => {
  //   if (selectedGroup != null) {
  //     getImagesFromOneGroup(selectedGroup.id);
  //   }
  // }, [selectedGroup]);

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
        <div>
          <h1 className="text-2xl font-bold">{selectedGroup.name}</h1>
          <br></br>
          {selectedGroup.images ? (
            selectedGroup.images.map((image) => {
              return (
                // <div key={image.id}><img src={`http://localhost:8080/proxy?key=${image.key}`} alt={`${image.alt_text}`}/></div>
                <div className="p-2 border w-1/2 m-2" key={image.id}>
                  {image.id}
                  <img className="max-h-[200px]" src={`http://localhost:8080/proxy?key=${image.filename}&type=gallery`}/>
                  <button
                    className={`ml-[30%] border w-[20%] ${RED_BUTTON_CSS}`}
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
          <button
            className={`mt-[2%] ${RED_BUTTON_CSS} mr-[2%]`}
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
      ) : (
        <div></div>
      )}
    </div>
  );
};

export default GalleryEdit;
