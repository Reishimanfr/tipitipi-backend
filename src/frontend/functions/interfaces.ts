export interface BlogAttachments {
  id: number;
  url: string;
  filename: string;
}

export interface BlogPostDataBodyJson {
  content: string;
  created_at: number;
  edited_at: number;
  id: number;
  attachments: BlogAttachments[];
  title: string;
  error?: string;
}

export interface GalleryGroup {
  id: number;
  name: string;
}

export interface GalleryCreateNewJson {
  error?: string;
  message: string;
}