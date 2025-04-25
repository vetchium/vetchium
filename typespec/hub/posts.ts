import { Handle } from "../common/common";
import { VTagID, VTagName } from "../common/vtags";

export interface AddPostRequest {
  content: string;
  tag_ids: VTagID[];
  new_tags: VTagName[];
}

export interface AddPostResponse {
  post_id: string;
}

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

export interface GetUserPostsRequest {
  handle?: Handle;
  pagination_key?: string;
  limit?: number;
}

export interface GetUserPostsResponse {
  posts: Post[];
  pagination_key: string;
}

export interface FollowUserRequest {
  handle: Handle;
}

export interface UnfollowUserRequest {
  handle: Handle;
}

export interface GetFollowStatusRequest {
  handle: Handle;
}

export interface FollowStatus {
  is_following: boolean;
  is_blocked: boolean;
  can_follow: boolean;
}

export interface GetMyHomeTimelineRequest {
  pagination_key?: string;
  limit?: number;
}

export interface MyHomeTimeline {
  posts: Post[];
  pagination_key: string;
}

export interface GetPostDetailsRequest {
  post_id: string;
}

export interface UpvoteUserPostRequest {
  post_id: string;
}

export interface DownvoteUserPostRequest {
  post_id: string;
}

export interface UnvoteUserPostRequest {
  post_id: string;
}
