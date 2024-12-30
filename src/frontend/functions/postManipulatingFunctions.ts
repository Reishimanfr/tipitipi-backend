import { toast } from "react-toastify";
export function validateDataForm(title: string, content: string): boolean {
  if (title === "") {
    toast.error("Podano pusty tytuł");

    return false;
  }
  if (content === "<p><br></p>") {
    toast.error("Podano pustą treść");
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

function base64ToBlob(base64: string): Blob | null {
  //mimetype is first part of base64 variable , that matches given regexp in match function
  let mimeType: RegExpMatchArray | null | string = base64
    .split(",")[0]
    .match(/image\/(webp|jpeg|png|gif)/);
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
    /<img\s+src="data:image\/(webp|jpeg|png|gif);base64,([A-Za-z0-9+/=]+)"\s*\/?>/g;

  const contentWithoutImages = content.replace(regexp, (match) => {
    images.push(match);
    i++;
    return `{{${i}}}`;
  });
  return { images: images, contentWithoutImages: contentWithoutImages };
}

export function buildPostMultipart(title: string, content: string) {
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
      title + "." + blob.type.split("/")[1]
    );
  }

  return formData;
}

export function buildGalleryMultipart(images : FileList) : FormData | null{
  const formData = new FormData();
  for(const image of images) {
    const blob = imageToBlob(image) 
    if(!blob) return null
    formData.append("files[]" , blob )
  }
  return formData
}

function imageToBlob (image : File) : Blob | null{
  if (!image.type.match(/image\/(webp|jpeg|png|gif)/)) {
    console.error("File is not a supported image type.");
    return null;
  }

  // Zwracamy nowy Blob z ustawionym typem MIME
  return new Blob([image], { type: image.type });
}

export function getToken(){
    const token = localStorage.getItem("token");
    if (!token) {
      console.error("Token is invalid, redirecting to login page...");
      window.location.href = "/admin/login"
      return false;
    }
    return token
}


