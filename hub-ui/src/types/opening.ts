export interface Opening {
  id: string;
  title: string;
  positions: number;
  filled_positions: number;
  jd: string;
  employer_name: string;
  cost_center_name: string;
  location_titles?: string[];
  remote_country_codes?: string[];
  remote_timezones?: string[];
  opening_type: OpeningType;
  yoe_min: number;
  yoe_max: number;
  min_education_level?: EducationLevel;
  salary?: Salary;
  state: OpeningState;
  created_at: string;
  last_updated_at: string;
}

export enum OpeningState {
  ACTIVE_OPENING = 'ACTIVE_OPENING',
  SUSPENDED_OPENING = 'SUSPENDED_OPENING',
  CLOSED_OPENING = 'CLOSED_OPENING',
}

export enum OpeningType {
  FULL_TIME_OPENING = 'FULL_TIME_OPENING',
  PART_TIME_OPENING = 'PART_TIME_OPENING',
  CONTRACT_OPENING = 'CONTRACT_OPENING',
  INTERNSHIP_OPENING = 'INTERNSHIP_OPENING',
  UNSPECIFIED_OPENING = 'UNSPECIFIED_OPENING',
}

export enum EducationLevel {
  BACHELOR_EDUCATION = 'BACHELOR_EDUCATION',
  MASTER_EDUCATION = 'MASTER_EDUCATION',
  DOCTORATE_EDUCATION = 'DOCTORATE_EDUCATION',
  NOT_MATTERS_EDUCATION = 'NOT_MATTERS_EDUCATION',
  UNSPECIFIED_EDUCATION = 'UNSPECIFIED_EDUCATION',
}

export interface Salary {
  min_amount: number;
  max_amount: number;
  currency: string;
}

export interface FilterOpeningsRequest {
  employer_name?: string;
  country_code?: string;
  location_titles?: string[];
  opening_type?: OpeningType[];
  min_education_level?: EducationLevel[];
  yoe_min?: number;
  yoe_max?: number;
  salary_min?: number;
  salary_max?: number;
  salary_currency?: string;
  pagination_key?: string;
  limit?: number;
}

export interface GetOpeningRequest {
  id: string;
} 