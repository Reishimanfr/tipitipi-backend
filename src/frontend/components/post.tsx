const Post = (props : any) => {
    return(
        <div>
            <h1>{props.title}</h1>
            <div dangerouslySetInnerHTML={{__html: props.content}}></div>
        </div>
    )
}
export default Post