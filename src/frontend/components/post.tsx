import "react-quill/dist/quill.snow.css"; // Załaduj style Quill, np. z wersji 'snow' (lub 'bubble', jeśli używasz innego stylu)
import { Link } from "react-router-dom"
import { API_URL } from '../functions/global'
import { BlogFiles } from "../functions/interfaces"

interface Props {
  id: number;
  title: string;
  content?: string | null;
  date: number;
  attachments?: BlogFiles[] | null;
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
        `<img style="max-height:200px;" src="${API_URL}/proxy?key=${attachment.filename}&type=blog" alt="${attachment.filename}"/>`
      );
    });
  } else {
    if (content) {
      content = content?.replace(/{{\d+}}/g, "");
    }
  }
  return (
    <div className=" mb-6 border-4 bg-white border-gray-800 rounded-lg">
      <div className="bg-gray-900 text-white pt-4 p-4 text-xl md:text-3xl">
        <div className="float-left ">{title}</div>
        <div className="float-right">
          {" "}
          {new Date(date * 1000).toLocaleDateString("en-GB") +
            " " +
            new Date(date * 1000).toLocaleTimeString("en-GB", {
              hour: "2-digit",
              minute: "2-digit",
            })}
        </div>
        <div className="clear-both"></div>
      </div>
      {content ? (
        <div
          className={`ql-editor  px-4 pt-4 text-medium ${
            willBeUsedManyTimes ? "line-clamp-2" : ""
          }`}
          dangerouslySetInnerHTML={{ __html: content }}
        ></div>
      ) : (
        <div></div>
      )}

      {/* <h1 className="pt-4 pl-4">
        {new Date(date * 1000).toLocaleDateString("en-pl").toString()}
      </h1> */}

      {willBeUsedManyTimes ? (
        <Link to={`/blog/${id}`}>
          <button className="border p-2 ml-4 mb-4 border-gray-900 hover:bg-gray-900 hover:text-orange-400 hover:duration-300 rounded-md">
            Zobacz więcej
          </button>
        </Link>
      ) : (
        <div></div>
      )}
    </div>
  );
};
export default Post;
