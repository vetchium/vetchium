export class AddPostCommentRequest {
  post_id: string = "";
  content: string = "";

  IsValid(): boolean {
    return (
      this.post_id.length > 0 &&
      this.content.length > 0 &&
      this.content.length <= 4096
    );
  }
}

export interface AddPostCommentResponse {
  post_id: string;
  comment_id: string;
}

export class GetPostCommentsRequest {
  post_id: string = "";
  pagination_key?: string;
  limit?: number;
  IsValid(): boolean {
    return (
      this.post_id.length > 0 &&
      (this.limit === undefined || (this.limit >= 0 && this.limit <= 40))
    );
  }
}

export interface PostComment {
  id: string;
  content: string;
  author_name: string;
  author_handle: string;
  created_at: Date;
}

export class DisablePostCommentsRequest {
  post_id: string = "";
  delete_existing_comments: boolean = false;

  IsValid(): boolean {
    return this.post_id.length > 0;
  }
}

export class EnablePostCommentsRequest {
  post_id: string = "";

  IsValid(): boolean {
    return this.post_id.length > 0;
  }
}

export class DeletePostCommentRequest {
  post_id: string = "";
  comment_id: string = "";

  IsValid(): boolean {
    return this.post_id.length > 0 && this.comment_id.length > 0;
  }
}

export class DeleteMyCommentRequest {
  post_id: string = "";
  comment_id: string = "";

  IsValid(): boolean {
    return this.post_id.length > 0 && this.comment_id.length > 0;
  }
}
