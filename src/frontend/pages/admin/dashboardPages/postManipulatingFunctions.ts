import { useNavigate } from "react-router-dom";
interface BlogAttachments {
    ID: number;
    BlogPostID: number;
    Path: string;
    Filename: string;
  }
  
  interface BlogPostDataBodyJson {
    Content: string;
    Created_At: string;
    Edited_At: string;
    ID: number;
    Attachments: BlogAttachments[];
    Title: string;
    error?: string;
  }

export function validateDataForm(title: string, content: string): boolean {
  if (title === "") {
    alert("Podano pusty tytuł");

    return false;
  }
  if (content === "<p><br></p>") {
    alert("Podano pustą treść");
    return false;
  }
  const confirm = window.confirm(
    "Czy jesteś pewien że chcesz opublikować ten post?"
  );
  if (!confirm) {
    return false;
  }
  return true;
}
function makeFilename(length: number) {
  let result = "";
  const characters =
    "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789";
  const charactersLength = characters.length;
  let counter = 0;
  while (counter < length) {
    result += characters.charAt(Math.floor(Math.random() * charactersLength));
    counter += 1;
  }
  return result;
}

function base64ToBlob(base64: string): Blob | null {
  //mimetype is first part of base64 variable , that matches given regexp in match function
  let mimeType: RegExpMatchArray | null | string = base64
    .split(",")[0]
    .match(/image\/(jpeg|png|gif)/);
  //base64content is second part of base64 variable , that has last two chars trimmed
  let base64content = base64.split(",")[1];
  base64content = base64content.substring(0, base64content.length - 2);

  if (!mimeType) {
    console.error("Mimetype of image is null");
    return null;
  }
  mimeType = mimeType[0];
  const byteString = atob(base64content);

  // Tworzymy tablicę bajtów
  const byteNumbers = new Array(byteString.length);
  for (let i = 0; i < byteString.length; i++) {
    byteNumbers[i] = byteString.charCodeAt(i);
  }
  const byteArray = new Uint8Array(byteNumbers);
  return new Blob([byteArray], { type: mimeType });
}

function extractImagesFromContent(content: string) {
  const images: string[] = [];
  let i = -1;
  const regexp =
    /<img\s+src="data:image\/(jpeg|png|gif);base64,([A-Za-z0-9+/=]+)"\s*\/?>/g;

  const contentWithoutImages = content.replace(regexp, (match) => {
    images.push(match);
    i++;
    return `{{${i}}}`;
  });
  return { images: images, contentWithoutImages: contentWithoutImages };
}

export function buildMultipart(title: string, content: string) {
  const extractedData = extractImagesFromContent(content);
  const base64images: string[] = extractedData.images;

  const formData = new FormData();

  formData.append("title", title);
  formData.append("content", extractedData.contentWithoutImages);

  for (const image of base64images) {
    const blob = base64ToBlob(image);
    if (!blob) {
      return;
    }
    formData.append(
      "files[]",
      blob,
      makeFilename(5) + "." + blob.type.split("/")[1]
    );
  }

  return formData;
}


export function getToken(){
    const token = localStorage.getItem("token");
    if (!token) {
    const navigate = useNavigate();
      alert("Token is invalid, redirecting to login page...");
      navigate("/admin/login");
      return false;
    }
    return token
}

export async function fetchPosts(
    setPosts: React.Dispatch<React.SetStateAction<BlogPostDataBodyJson[]>>
  ) {
    try {
      const response = await fetch(
        `http://localhost:2333/blog/posts?limit=999&images=true`,
        {
          method: "GET",
        }
      );
      if (!response.ok) {
        throw new Error(response.statusText);
      }
  
      const data: Array<BlogPostDataBodyJson> = await response.json();
      setPosts((prevPosts) => prevPosts?.concat(data));
    } catch (error) {
      alert(error);
    }
  }

