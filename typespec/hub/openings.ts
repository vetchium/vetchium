import { CountryCode, Currency, TimeZone } from "../common/common";
import { EducationLevel, OpeningType } from "../common/openings";

export interface ExperienceRange {
  yoe_min: number;
  yoe_max: number;
}

export interface SalaryRange {
  currency: Currency;
  min: number;
  max: number;
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
  salary_range?: SalaryRange;
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

export interface ApplyForOpeningRequest {
  opening_id_within_company: string;
  company_domain: string;
  resume: string;
  filename: string;
  cover_letter?: string;
}
