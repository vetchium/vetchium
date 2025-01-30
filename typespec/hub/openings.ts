import { CountryCode } from "../common/common";
import {
  EducationLevel,
  OpeningState,
  OpeningType,
  Salary,
} from "../common/openings";

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
  yoe_min: number;
  yoe_max: number;
  hiring_manager_name: string;
  hiring_manager_vetchi_handle?: string;
  jd: string;
  job_title: string;
  opening_id_within_company: string;
  opening_type: OpeningType;
  pagination_key: BigInt;
  recruiter_name: string;
  salary?: Salary;
  state: OpeningState;
}

export interface ApplyForOpeningRequest {
  opening_id_within_company: string;
  company_domain: string;
  resume: string;
  filename: string;
  cover_letter?: string;
}

export interface ApplyForOpeningResponse {
  application_id: string;
}
