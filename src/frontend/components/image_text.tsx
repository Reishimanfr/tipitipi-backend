const Image_Text = (props: any) => {
  return (
    // TODO justify center nie dziala i blok tekstowy jest za duzy
    <div className="flex items-stretch justify-center pb-[5%]">
      <div id="imageContainer" className="w-[50%] h-fit">
        <img src={props.image} className="float-right"></img>
      </div>
      <div className="bg-blue-500 w-[50%] ">
        <div className="p-[5%]">
          <h1 className="pb-[10%] text-3xl">{props.header}</h1>
          <p className="mr-[15%] text-xl">{props.paragraph}</p>
        </div>
      </div>
    </div>
  );
};

export default Image_Text;
