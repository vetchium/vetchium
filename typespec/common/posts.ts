import { Handle } from "./common";

export interface Post {
  id: string;
  content: string;
  tags: string[];
  author_name: string;
  author_handle: Handle;
  created_at: string;
  updated_at: string;
}
