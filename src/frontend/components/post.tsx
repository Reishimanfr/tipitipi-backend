import { Link } from "react-router-dom";
interface BlogAttachments {
  ID: number;
  BlogPostID: number;
  Path: string;
  Filename: string;
}
interface Props {
  id: number;
  title: string;
  content?: string | null;
  date: string;
  attachments?: BlogAttachments[] | null;
  willBeUsedManyTimes: boolean;
}

const Post = ({
  id,
  title,
  content = null,
  date,
  attachments = null,
  willBeUsedManyTimes,
}: Props) => {
  if (attachments) {
    attachments.forEach((attachment, index) => {
      if (!content) {
        return;
      }
      content = content.replace(
        `{{${index}}}`,
        `<img src='${attachment.Path}'/>`
      );
    });
  }
  else {
    if(content) {
      content = content?.replace(/{{\d+}}/g, "")
    }

  }
  return (
    <div className=" mb-[1%] border-4 border-gray-800 rounded-lg">
      <h1 className="bg-gray-900 text-white pt-[1%] p-[1%] text-3xl">
        {title + " , " + id}
      </h1>
      {content ? (
        <div
          className={`px-[1%] pt-[1%] text-medium ${
            willBeUsedManyTimes ? "line-clamp-2" : ""
          }`}
          dangerouslySetInnerHTML={{ __html: content }}
        ></div>
      ) : (
        <div></div>
      )}

      <h1 className="pt-[1%] pl-[1%]">
        {new Date(parseInt(date) * 1000).toLocaleDateString("en-pl").toString()}
      </h1>


      {willBeUsedManyTimes ? (
        <Link to={`/blog/${id}`}>
          <button className="border p-[0.5%] ml-[1%] mb-[1%] border-gray-900 hover:bg-gray-900 hover:text-orange-400 hover:duration-300 rounded-md">Zobacz więcej</button>
        </Link>
      ) : (
        <div></div>
      )}
    </div>
  );
};
export default Post;
