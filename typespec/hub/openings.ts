import { CountryCode } from "../common/common";
import {
  EducationLevel,
  OpeningState,
  OpeningType,
  Salary,
} from "../common/openings";
import { VTagID } from "../common/vtags";
export interface ExperienceRange {
  yoe_min: number;
  yoe_max: number;
}

export interface LocationFilter {
  country_code: CountryCode;
  city: string;
}

export interface FindHubOpeningsRequest {
  country_code: CountryCode;
  cities?: string[];
  opening_types?: OpeningType[];
  company_domains?: string[];
  experience_range?: ExperienceRange;
  salary_range?: Salary;
  min_education_level?: EducationLevel;
  tags?: VTagID[];
  terms?: string[];
  pagination_key?: number;
  limit?: number;
}

export interface HubOpening {
  opening_id_within_company: string;
  company_domain: string;
  company_name: string;
  job_title: string;
  jd: string;
  pagination_key: number;
}

export interface GetHubOpeningDetailsRequest {
  opening_id_within_company: string;
  company_domain: string;
}

export interface HubOpeningDetails {
  company_domain: string;
  company_name: string;
  created_at: EpochTimeStamp;
  education_level: EducationLevel;
  hiring_manager_name: string;
  hiring_manager_vetchi_handle?: string;
  is_appliable: boolean;
  jd: string;
  job_title: string;
  opening_id_within_company: string;
  opening_type: OpeningType;
  pagination_key: BigInt;
  recruiter_name: string;
  salary?: Salary;
  state: OpeningState;
  yoe_max: number;
  yoe_min: number;
}

export interface ApplyForOpeningRequest {
  opening_id_within_company: string;
  company_domain: string;
  resume: string;
  filename: string;
  cover_letter?: string;
  endorser_handles?: string[];
}

export interface ApplyForOpeningResponse {
  application_id: string;
}
