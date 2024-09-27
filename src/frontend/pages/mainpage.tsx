import background_example from "../assets/example_background.jpg";
import landscapeImage from "../assets/landscape.jpg"
import Image_Text from "../components/image_text";
import { useState , useEffect } from "react";
import Post from "../components/post";

interface BlogPostDataBodyJson {
  Content: string;
  Created_At: string;
  Edited_At: string;
  ID: number;
  Title: string;
}
const Mainpage = () => {

  const [posts, setPosts] = useState<Array<BlogPostDataBodyJson>>([]);
  const [loading, setLoading] = useState<boolean>(true);

  useEffect(() => {
    async function fetchPost() {
      try {
        const response = await fetch(
          `http://localhost:2333/blog/posts?limit=3&sort=newest`,
          {
            method: "GET",
          }
        );
        if (!response.ok) {
          throw new Error(response.statusText);
        }

        const data: Array<BlogPostDataBodyJson> = await response.json();
        setPosts((prevPosts) => prevPosts?.concat(data));
      } catch (error) {
        alert(error);
      } finally {
        setLoading(false);
      }
    }
    fetchPost();
  }, []);



  if (loading) {
    return <div>Loading</div>;
  }
  return (
    <div className="globalCss">
      {/* <h1>Mainpage</h1>
            <Link to="/admin">Admin</Link>
            <h1 className="text-center">{props.mainpageFirstHeader}</h1> */}

      <Image_Text
        image={background_example}
        header="Jakiś nagłówek"
        paragraph="Lorem ipsum dolor sit amet, consectetur adipiscing 
                    elit. Suspendisse tellus lectus, pharetra a aliquet sed, 
                    sagittis vel sapien."
        orientation="left"
      />

        <h1 className="text-3xl mb-[1%]">Oto kilka najnowszych postów</h1>
        {posts ? (
        posts.map((post) => {
          return<Post id={post.ID} content={post.Content} title={post.Title} date={post.Edited_At} willBeUsedManyTimes={true}/>;
        })
      ) : (
        <div>No post found</div>
      )}<br></br><br></br>
<br></br>
<br></br>
<br></br>



      <Image_Text
        image={landscapeImage}
        header="Jakiś nagłówek"
        paragraph="Lorem ipsum dolor sit amet, consectetur adipiscing 
                    elit. Suspendisse tellus lectus, pharetra a aliquet sed, 
                    sagittis vel sapien.orem ipsum dolor sit amet, consectetur adipiscing 
                    elit. Suspendisse tellus lectus, pharetra a aliquet sed, 
                    sagittis vel sapien.orem ipsum dolor sit amet, consectetur adipiscing 
                    elit. Suspendisse tellus lectus, pharetra a aliquet sed, 
                    sagittis vel sapien.orem ipsum dolor sit amet, consectetur adipiscing 
                    elit. Suspendisse tellus lectus, pharetra a aliquet sed, 
                    sagittis vel sapien.orem ipsum dolor sit amet, consectetur adipiscing 
                    elit. Suspendisse tellus lectus, pharetra a aliquet sed, 
                    sagittis vel sapien."
        orientation="right"/>
    </div>
  );
};
export default Mainpage;
