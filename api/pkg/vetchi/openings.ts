import { CountryCode, Currency, EmailAddress } from "./common";
import { OpeningType } from "./opening-types";
import { EducationLevel } from "./education-level";
import { CostCenterName } from "./cost-centers";
import { OrgUserShort } from "./org-users";

export type OpeningState =
  | "DRAFT_OPENING"
  | "ACTIVE_OPENING"
  | "SUSPENDED_OPENING"
  | "CLOSED_OPENING";

export interface Salary {
  min_amount: number;
  max_amount: number;
  currency: Currency;
}

export type OpeningID = string;

export interface OpeningInfo {
  id: OpeningID;
  title: string;
  positions: number;
  filled_positions: number;
  recruiter: OrgUserShort;
  hiring_manager: OrgUserShort;
  cost_center_name: CostCenterName;
  opening_type: OpeningType;
  state: OpeningState;
  created_at: Date;
  last_updated_at: Date;
}

export interface Opening {
  id: OpeningID;
  title: string;
  positions: number;
  filled_positions: number;
  jd: string;
  recruiter: OrgUserShort;
  hiring_manager: OrgUserShort;
  hiring_team_members?: OrgUserShort[];
  cost_center_name: CostCenterName;
  employer_notes?: string;
  location_titles?: string[];
  remote_country_codes?: CountryCode[];
  remote_timezones?: string[];
  opening_type: OpeningType;
  yoe_min: number;
  yoe_max: number;
  min_education_level?: EducationLevel;
  salary?: Salary;
  state: OpeningState;
  created_at: Date;
  last_updated_at: Date;
}

export interface CreateOpeningRequest {
  title: string;
  positions: number;
  jd: string;
  recruiter: EmailAddress;
  hiring_manager: EmailAddress;
  hiring_team_members?: EmailAddress[];
  cost_center_name: CostCenterName;
  employer_notes?: string;
  location_titles?: string[];
  remote_country_codes?: CountryCode[];
  remote_timezones?: string[];
  opening_type: OpeningType;
  yoe_min: number;
  yoe_max: number;
  min_education_level?: EducationLevel;
  salary?: Salary;
}

export interface CreateOpeningResponse {
  opening_id: OpeningID;
}

export interface GetOpeningRequest {
  id: OpeningID;
}

export interface FilterOpeningsRequest {
  pagination_key?: OpeningID;
  state?: OpeningState[];
  from_date?: Date;
  to_date?: Date;
  limit?: number;
}

export interface ChangeOpeningStateRequest {
  opening_id: OpeningID;
  from_state: OpeningState;
  to_state: OpeningState;
}

export interface UpdateOpeningRequest {
  id: OpeningID;
}

export interface GetOpeningWatchersRequest {
  opening_id: OpeningID;
}

export interface OpeningWatchers {
  opening_id: OpeningID;
  emails?: EmailAddress[];
}

export interface AddOpeningWatchersRequest {
  opening_id: OpeningID;
  emails: EmailAddress[];
}

export interface RemoveOpeningWatcherRequest {
  opening_id: OpeningID;
  email: EmailAddress;
}
