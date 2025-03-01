export interface GetHubUserBioRequest {
  handle: string;
}

export interface EmployerViewBio {
  handle: string;
  full_name: string;
  short_bio: string;
  long_bio: string;
  verified_mail_domains?: string[];
}
