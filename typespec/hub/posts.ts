import { Handle } from "../common/common";
import { EmployerPost } from "../common/posts";
import { VTagID } from "../common/vtags";

export class AddFTPostRequest {
  content: string = "";
  tag_ids: VTagID[] = [];

  IsValid(): boolean {
    return (
      this.content.length > 0 &&
      this.content.length <= 255 &&
      this.tag_ids.length <= 3
    );
  }
}

export class AddPostRequest {
  content: string = "";
  tag_ids: VTagID[] = [];

  IsValid(): boolean {
    return (
      this.content.length > 0 &&
      this.content.length <= 4096 &&
      this.tag_ids.length <= 3
    );
  }
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
  me_upvoted: boolean;
  me_downvoted: boolean;
  can_upvote: boolean;
  can_downvote: boolean;
  am_i_author: boolean;
  can_comment: boolean;
  comments_count: number;
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
  employer_posts: EmployerPost[];
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

export interface FollowOrgRequest {
  domain: string;
}

export interface UnfollowOrgRequest {
  domain: string;
}

export interface GetEmployerPostDetailsRequest {
  employer_post_id: string;
}
