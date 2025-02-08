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
