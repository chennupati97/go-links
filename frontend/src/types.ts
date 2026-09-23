export interface Link {
  id: number;
  slug: string;
  url: string;
  createdAt: string;
}

export interface CreateLinkRequest {
  slug: string;
  url: string;
}
