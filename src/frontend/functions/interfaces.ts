export interface BlogFiles {
  id: number;
  filename: string;
}

export interface BlogPostDataBodyJson {
  content: string;
  created_at: number;
  edited_at: number;
  id: number;
  files: BlogFiles[];
  title: string;
  error?: string;
}

export interface GalleryGroup {
  id: number;
  name: string;
}

// export interface GalleryCreateNewJson {
//   error?: string;
//   message: string;
// }

export interface GalleryImage {
  id: number;
  alt_text: string;
  url: string;
  key: string;
  group_id: number;
}

export interface GroupInfo {
  name: string;
  id: number;
  images: GalleryImages[]
}

//zmienic to images na image tylko tamten interfejs jest stary
interface GalleryImages {
  id: number;
  filename: string;
}