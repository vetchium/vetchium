export interface AddEducationRequest {
  institute_domain: string;
  degree?: string;
  start_date?: string;
  end_date?: string;
  description?: string;
}

export interface AddEducationResponse {
  education_id: string;
}

export interface Institute {
  domain: string;
  name: string;
}

export interface FilterInstitutesRequest {
  prefix: string;
}

export interface DeleteEducationRequest {
  education_id: string;
}

export interface Education {
  id?: string;
  institute_domain: string;
  degree?: string;
  start_date?: string;
  end_date?: string;
  description?: string;
}

export interface ListEducationRequest {
  user_handle?: string;
}

export interface FilterInstitutesResponse {
  // Maximum 10 institutes will be returned in random order.
  institutes: Institute[];
}
