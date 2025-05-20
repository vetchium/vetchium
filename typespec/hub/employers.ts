export interface FilterEmployersRequest {
  prefix: string;
}

export interface HubEmployer {
  domain: string;
  name: string;
  ascii_name: string;
}

export interface FilterEmployersResponse {
  employers: HubEmployer[];
}

export interface GetEmployerDetailsRequest {
  domain: string;
}

export interface HubEmployerDetails {
  name: string;
  verified_employees_count: number;
  is_onboarded: boolean;
  active_openings_count: number;
  is_following: boolean;
}
