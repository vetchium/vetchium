import { VTag, VTagID } from "../common/vtags";

export class IncognitoPost {
  incognito_post_id: string = "";
  content: string = "";
  tags: VTag[] = [];
  created_at: string = "";
  upvotes_count: number = 0;
  downvotes_count: number = 0;
  score: number = 0;
  me_upvoted: boolean = false;
  me_downvoted: boolean = false;
  can_upvote: boolean = false;
  can_downvote: boolean = false;
  is_created_by_me: boolean = false;
  is_deleted: boolean = false;
}

export class AddIncognitoPostRequest {
  content: string = "";
  tag_ids: VTagID[] = [];

  IsValid(): boolean {
    return (
      this.content.length > 0 &&
      this.content.length <= 1024 &&
      this.tag_ids.length >= 1 &&
      this.tag_ids.length <= 3
    );
  }
}

export class AddIncognitoPostResponse {
  incognito_post_id: string = "";
}

export class IncognitoPostComment {
  comment_id: string = "";
  content: string = "";
  in_reply_to?: string = undefined;
  created_at: string = "";
  upvotes_count: number = 0;
  downvotes_count: number = 0;
  score: number = 0;
  me_upvoted: boolean = false;
  me_downvoted: boolean = false;
  can_upvote: boolean = false;
  can_downvote: boolean = false;
  is_created_by_me: boolean = false;
  is_deleted: boolean = false;
  depth: number = 0;
}

export class AddIncognitoPostCommentRequest {
  incognito_post_id: string = "";
  content: string = "";
  in_reply_to?: string = undefined;

  IsValid(): boolean {
    return (
      this.incognito_post_id.length > 0 &&
      this.content.length > 0 &&
      this.content.length <= 512
    );
  }
}

export class AddIncognitoPostCommentResponse {
  incognito_post_id: string = "";
  comment_id: string = "";
}

export class GetIncognitoPostCommentsRequest {
  incognito_post_id: string = "";
  pagination_key?: string = undefined;
  limit: number = 25;
  parent_comment_id?: string = undefined;
  include_nested_depth?: number = 0;

  IsValid(): boolean {
    return (
      this.incognito_post_id.length > 0 &&
      this.limit >= 1 &&
      this.limit <= 100 &&
      (this.include_nested_depth === undefined ||
        (this.include_nested_depth >= 0 && this.include_nested_depth <= 5))
    );
  }
}

export class GetIncognitoPostCommentsResponse {
  comments: IncognitoPostComment[] = [];
  pagination_key: string = "";
}

export class DeleteIncognitoPostCommentRequest {
  incognito_post_id: string = "";
  comment_id: string = "";
}

export class UpvoteIncognitoPostCommentRequest {
  incognito_post_id: string = "";
  comment_id: string = "";
}

export class DownvoteIncognitoPostCommentRequest {
  incognito_post_id: string = "";
  comment_id: string = "";
}

export class UnvoteIncognitoPostCommentRequest {
  incognito_post_id: string = "";
  comment_id: string = "";
}

export class GetIncognitoPostRequest {
  incognito_post_id: string = "";
}

export class DeleteIncognitoPostRequest {
  incognito_post_id: string = "";
}

export class UpvoteIncognitoPostRequest {
  incognito_post_id: string = "";
}

export class DownvoteIncognitoPostRequest {
  incognito_post_id: string = "";
}

export class UnvoteIncognitoPostRequest {
  incognito_post_id: string = "";
}

export class MyIncognitoPostComment {
  comment_id: string = "";
  content: string = "";
  in_reply_to?: string = undefined;
  created_at: string = "";
  upvotes_count: number = 0;
  downvotes_count: number = 0;
  score: number = 0;
  me_upvoted: boolean = false;
  me_downvoted: boolean = false;
  is_deleted: boolean = false;
  depth: number = 0;
  incognito_post_id: string = "";
  post_content_preview: string = "";
  post_tags: VTag[] = [];
}

export class GetMyIncognitoPostCommentsRequest {
  pagination_key?: string = undefined;
  limit: number = 25;

  IsValid(): boolean {
    return this.limit >= 1 && this.limit <= 40;
  }
}

export class GetMyIncognitoPostCommentsResponse {
  comments: MyIncognitoPostComment[] = [];
  pagination_key: string = "";
}

export enum IncognitoPostTimeFilter {
  Past24Hours = "past_24_hours",
  PastWeek = "past_week",
  PastMonth = "past_month",
  PastYear = "past_year",
}

export class GetIncognitoPostsRequest {
  tag_id: VTagID = "";
  time_filter?: IncognitoPostTimeFilter = IncognitoPostTimeFilter.Past24Hours;
  limit: number = 25;
  pagination_key?: string = undefined;

  IsValid(): boolean {
    return this.tag_id.length > 0 && this.limit >= 1 && this.limit <= 100;
  }
}

export class IncognitoPostSummary {
  incognito_post_id: string = "";
  content: string = "";
  tags: VTag[] = [];
  created_at: string = "";
  upvotes_count: number = 0;
  downvotes_count: number = 0;
  score: number = 0;
  me_upvoted: boolean = false;
  me_downvoted: boolean = false;
  can_upvote: boolean = false;
  can_downvote: boolean = false;
  comments_count: number = 0;
  is_created_by_me: boolean = false;
  is_deleted: boolean = false;
}

export class GetIncognitoPostsResponse {
  posts: IncognitoPostSummary[] = [];
  pagination_key: string = "";
}

export class GetMyIncognitoPostsRequest {
  pagination_key?: string = undefined;
  limit: number = 25;

  IsValid(): boolean {
    return this.limit >= 1 && this.limit <= 40;
  }
}

export class GetMyIncognitoPostsResponse {
  posts: IncognitoPostSummary[] = [];
  pagination_key: string = "";
}
