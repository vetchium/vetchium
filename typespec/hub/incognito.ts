import { VTag, VTagID } from "../common/vtags";

export class IncogntiPost {
  incognito_post_id: string = "";
  content: string = "";
  tags: VTag[] = [];
  created_at: string = "";
  upvotes: number = 0;
  downvotes: number = 0;
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
  in_reply_to?: string = "";
  created_at: Date = new Date();
  upvotes: number = 0;
  downvotes: number = 0;
  is_created_by_me: boolean = false;
  is_deleted: boolean = false;
  my_vote?: "upvote" | "downvote" = undefined;
  depth: number = 0;
}

export class AddIncognitoPostCommentResponse {
  incognito_post_id: string = "";
  comment_id: string = "";
}

export class GetIncognitoPostCommentsRequest {
  incognito_post_id: string = "";
  pagination_key?: string = undefined;
  limit: number = 10;
  parent_comment_id?: string = undefined;
  include_nested_depth?: number = 0;

  IsValid(): boolean {
    return (
      this.incognito_post_id.length > 0 &&
      this.limit >= 1 &&
      this.limit <= 100 &&
      this.include_nested_depth !== undefined &&
      this.include_nested_depth >= 0 &&
      this.include_nested_depth <= 5 &&
      this.parent_comment_id !== undefined &&
      this.parent_comment_id.length > 0
    );
  }
}

export class GetIncognitoPostCommentsResponse {
  comments: IncognitoPostComment[] = [];
  pagination_key?: string = undefined;
  has_more: boolean = false;
  total_count: number = 0;
}

export class DeleteIncognitoPostCommentRequest {
  incognito_post_id: string = "";
  comment_id: string = "";
}

export class VoteIncognitoPostCommentRequest {
  incognito_post_id: string = "";
  comment_id: string = "";
}

export class GetIncognitoPostRequest {
  incognito_post_id: string = "";
}

export class DeleteIncognitoPostRequest {
  incognito_post_id: string = "";
}
