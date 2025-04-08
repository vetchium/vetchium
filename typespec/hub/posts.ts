import { Post } from "../common/posts";

export interface AddPostRequest {
  content: string;
  tags: string[];
}

export interface AddPostResponse {
  post_id: string;
}

export interface GetTimelineRequest {
  timeline_id?: string;
  pagination_key?: string;
  limit?: number;
}

export interface GetTimelineResponse {
  posts: Post[];
  pagination_key: string;
}
