const Image_Text = (props: any) => {
  return (
    <div className="flex items-stretch justify-center pb-[5%]">
      <div id="imageContainer" className="w-[50%]">
        <img src={props.image} className="w-full h-full object-cover"></img>
      </div>
      <div className="bg-blue-500 w-[50%] flex">
        <div className="p-[5%] self-center">
          <h1 className="pb-[10%] text-5xl">{props.header}</h1>
          <p className="mr-[15%] text-xl">{props.paragraph}</p>
        </div>
      </div>
    </div>
  );
};

export default Image_Text;


