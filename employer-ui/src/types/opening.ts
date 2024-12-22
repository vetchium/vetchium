import { OrgUserShort } from './auth';

export interface Opening {
  id: string;
  title: string;
  positions: number;
  filled_positions: number;
  jd: string;
  recruiter: OrgUserShort;
  hiring_manager: OrgUserShort;
  hiring_team_members?: OrgUserShort[];
  cost_center_name: string;
  employer_notes?: string;
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
  DRAFT_OPENING = 'DRAFT_OPENING',
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

export interface CreateOpeningRequest {
  title: string;
  positions: number;
  jd: string;
  recruiter: string;
  hiring_manager: string;
  hiring_team_members?: string[];
  cost_center_name: string;
  employer_notes?: string;
  location_titles?: string[];
  remote_country_codes?: string[];
  remote_timezones?: string[];
  opening_type: OpeningType;
  yoe_min: number;
  yoe_max: number;
  min_education_level?: EducationLevel;
  salary?: Salary;
}

export interface UpdateOpeningRequest {
  id: string;
  title?: string;
  positions?: number;
  jd?: string;
  recruiter?: string;
  hiring_manager?: string;
  hiring_team_members?: string[];
  cost_center_name?: string;
  employer_notes?: string;
  location_titles?: string[];
  remote_country_codes?: string[];
  remote_timezones?: string[];
  opening_type?: OpeningType;
  yoe_min?: number;
  yoe_max?: number;
  min_education_level?: EducationLevel;
  salary?: Salary;
}

export interface FilterOpeningsRequest {
  state?: OpeningState[];
  from_date?: string;
  to_date?: string;
  pagination_key?: string;
  limit?: number;
}

export interface GetOpeningRequest {
  id: string;
} 