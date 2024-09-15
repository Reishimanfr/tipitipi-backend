import background_example from "../assets/example_background.jpg";
import landscapeImage from "../assets/landscape.jpg"
import Image_Text from "../components/image_text";
//import Text_Component from "../components/text_component";
const Mainpage = () => {
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
      <Image_Text
        image={landscapeImage}
        header="Jakiś nagłówek"
        paragraph="Lorem ipsum dolor sit amet, consectetur adipiscing 
                    elit. Suspendisse tellus lectus, pharetra a aliquet sed, 
                    sagittis vel sapien."
        orientation="right"/>
    </div>
  );
};
export default Mainpage;
