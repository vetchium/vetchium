import { Handle } from "./common";

export interface Post {
  id: string;
  content: string;
  tags: string[];
  author_name: string;
  author_handle: Handle;
  created_at: string;
  upvotes_count: number;
  downvotes_count: number;
  score: number;
}
