export interface EmployerPost {
  id: string;
  content: string;
  tags: string[];
  company_domain: string;
  created_at: string;
}

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
