interface BlogAttachments {
  ID: number;
  BlogPostID: number;
  Path: string;
  Filename: string;
}
interface Props {
  id: number;
  title: string;
  content: string;
  date: string;
  attachments?: BlogAttachments[] | null;
  key?: number;
}

const Post = ({ id, title, content, date, attachments = null }: Props) => {
  if (attachments) {
    attachments.forEach((attachment, index) => {
      content = content.replace(
        `{{${index}}}`,
        `<img src='${attachment.Path}'/>`
      );
    });
  }
  return (
    <div className=" bg-blue-300 mt-[2%] mb-[1%] ">
      <h1>{id}</h1>
      <h1 className="pt-[1%] pl-[1%] text-3xl">{title}</h1>
      <div
        className="p-[1%] text-xl"
        dangerouslySetInnerHTML={{ __html: content }}
      ></div>
      <h1 className="pb-[1%] pl-[1%]">
        {new Date(parseInt(date) * 1000).toLocaleDateString("en-pl").toString()}
      </h1>
    </div>
  );
};
export default Post;
