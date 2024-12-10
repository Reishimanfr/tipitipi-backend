import { Link } from "react-router-dom";
import { BlogAttachments } from "../functions/interfaces";
interface Props {
  id: number;
  title: string;
  content?: string | null;
  date: number;
  attachments?: BlogAttachments[] | null;
  willBeUsedManyTimes: boolean;
  loading?: boolean;
}

const Post = ({
  id,
  title,
  content = null,
  date,
  attachments = null,
  willBeUsedManyTimes,
}: Props) => {
  if (attachments && content && !willBeUsedManyTimes) {
    attachments.forEach((attachment, index) => {
      if (!content) {
        return;
      }
      content = content.replace(
        `{{${index}}}`,
        `<img style="max-height:200px;" src="http://localhost:8080/proxy?key=${attachment.filename}" alt="${attachment.filename}"/>`
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
        {new Date(date * 1000).toLocaleDateString("en-pl").toString()}
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
