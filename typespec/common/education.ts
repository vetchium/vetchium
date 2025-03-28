export interface Education {
  id?: string;
  institute_domain: string;
  degree?: string;
  start_date?: string;
  end_date?: string;
  description?: string;
}

export interface Institute {
  domain: string;
  name: string;
}
