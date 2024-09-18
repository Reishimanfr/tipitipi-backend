const Post = (props : any) => {
    return(

        <div className=" bg-blue-300 mt-[2%] mb-[1%] ">
            <h1>{props.id}</h1>
            <h1 className="pt-[1%] pl-[1%] text-3xl">{props.title}</h1>
            <div className="p-[1%] text-xl" dangerouslySetInnerHTML={{__html: props.content}}></div>
            <h1 className="pb-[1%] pl-[1%]">{(new Date(props.date * 1000).toLocaleDateString("en-pl")).toString()}</h1>
        </div>
    )
}
export default Post