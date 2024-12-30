import { useEffect, useState } from "react"
import PostSkeleton from "../components/postSkeletonLoading"
import SlideShow from "../components/slideshow"
import { API_URL } from '../functions/global'
import { GroupInfo } from "../functions/interfaces"
const Gallery = () => {
  const [loading, setLoading] = useState<boolean>(true);
  const [groups, setGroups] = useState<Array<GroupInfo> | null>();

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


  // useEffect(() => {
  //   function handleScroll() {
  //     const scrollTop = document.documentElement.scrollTop;
  //     const scrollHeight = document.documentElement.scrollHeight;
  //     const clientHeight = window.innerHeight;

  //     if (scrollTop + clientHeight >= scrollHeight && isMore) {
  //       if (offset + 6 > groups.length) {
  //         setIsMore(false);
  //       } else {
  //         setOffset((prevOffset) => prevOffset + 6);
  //       }
  //     }
  //   }

  //   window.addEventListener("scroll", handleScroll);
  //   return () => {
  //     window.removeEventListener("scroll", handleScroll);
  //   };
  // });

  if (loading || groups == null) {
    return (
      <div className="globalCss">
        <h1 className="text-3xl mt-5">Galeria</h1>
        <PostSkeleton />
        <PostSkeleton />
        <PostSkeleton />
        <PostSkeleton />
      </div>
    );
  }

  return (
    <div className="bg-gradient-to-r from-tipiOrange to-tipiPink min-h-screen">
      {groups ? (
        groups.map((group) =>
          group.images ? (
            <SlideShow
              images={group.images.map((image) => {
                return `${API_URL}/proxy?key=${image.filename}&type=gallery`;
              })}
            />
          ) : null
        )
      ) : (
        <div></div>
      )}
    </div>
  );
};
export default Gallery;
