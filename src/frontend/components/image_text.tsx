interface Props {
  image:string;
  header:string;
  paragraph:string;
  leftSide: boolean;
}


const Image_Text = ({image,header,paragraph,leftSide} : Props) => {
  const imageBlock  =  <div id="imageContainer" className="w-[50%]"><img src={image} className="w-full h-full object-cover"></img></div>;
  const textBlock = <div className="bg-gradient-to-r from-tipiOrange to-tipiPink  w-[50%] flex">
  <div className="p-[5%] self-center max-h-[300px] overflow-y-auto">
    <h1 className="pb-[10%] text-4xl">{header}</h1>
    <p className="mr-[15%] text-xl">{paragraph}</p>
  </div>
</div>;

  return (
    <div className="flex items-stretch justify-center pb-[5%] max-h-[400px] md:max-h-[600px] lg:max-h-[800px] overflow-hidden">
      {leftSide ? imageBlock: textBlock}
      {leftSide ? textBlock : imageBlock}
    </div>
  );
};

export default Image_Text;


