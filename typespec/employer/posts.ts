import { EmployerPost } from "../common/posts";

export class AddEmployerPostRequest {
  content: string = "";
  tags: string[] = [];

  IsValid(): boolean {
    return (
      this.content.length > 0 &&
      this.content.length <= 4096 &&
      this.tags.length > 0 &&
      this.tags.length <= 3
    );
  }
}

export interface AddEmployerPostResponse {
  post_id: string;
}

export class UpdateEmployerPostRequest {
  post_id: string = "";
  content: string = "";
  tags: string[] = [];

  IsValid(): boolean {
    return (
      this.content.length > 0 &&
      this.content.length <= 4096 &&
      this.tags.length > 0 &&
      this.tags.length <= 3
    );
  }
}

export interface DeleteEmployerPostRequest {
  post_id: string;
}

export class ListEmployerPostsRequest {
  pagination_key?: string;
  limit?: number;

  IsValid(): boolean {
    return this.limit !== undefined && this.limit >= 1 && this.limit <= 40;
  }
}

export interface ListEmployerPostsResponse {
  posts: EmployerPost[];
  pagination_key: string;
}

export class GetEmployerPostRequest {
  post_id: string = "";
}
