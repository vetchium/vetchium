export interface GetHubUserBioRequest {
  handle: string;
}

export interface EmployerWorkHistory {
  id: string;
  employer_domain: string;
  employer_name?: string;
  title: string;
  start_date: Date;
  end_date?: Date;
  description?: string;
}

export interface EmployerViewBio {
  handle: string;
  full_name: string;
  short_bio: string;
  long_bio: string;
  verified_mail_domains?: string[];
  work_history: EmployerWorkHistory[];
}
