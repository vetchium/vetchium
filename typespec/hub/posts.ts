import { Post } from "../common/posts";
import { VTagID, VTagName } from "../common/vtags";

export interface AddPostRequest {
  content: string;
  tag_ids: VTagID[];
  new_tags: VTagName[];
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
