import { CountryCode, EmailAddress, TimeZone } from "../common/common";
import {
  EducationLevel,
  OpeningState,
  OpeningType,
  Salary,
} from "../common/openings";
import { VTag, VTagID } from "../common/vtags";
import type { CostCenterName } from "../employer/costcenters";
import type { OrgUserShort } from "../employer/orgusers";

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
  hiring_team?: OrgUserShort[];
  cost_center_name: CostCenterName;
  location_titles?: string[];
  remote_country_codes?: CountryCode[];
  remote_timezones?: TimeZone[];
  opening_type: OpeningType;
  yoe_min: number;
  yoe_max: number;
  state: OpeningState;
  created_at: Date;
  last_updated_at: Date;
  employer_notes?: string;
  min_education_level: EducationLevel;
  salary?: Salary;
  tags?: VTag[];
}

export interface CreateOpeningRequest {
  title: string;
  positions: number;
  jd: string;
  recruiter: EmailAddress;
  hiring_manager: EmailAddress;
  hiring_team?: EmailAddress[];
  cost_center_name: CostCenterName;
  location_titles?: string[];
  remote_country_codes?: CountryCode[];
  remote_timezones?: TimeZone[];
  opening_type: OpeningType;
  yoe_min: number;
  yoe_max: number;
  employer_notes?: string;
  min_education_level: EducationLevel;
  salary?: Salary;

  // Should be minimum 1 and maximum 3
  tags: VTagID[];

  // Can be maximum 3
  new_tags?: string[];
}

export interface CreateOpeningResponse {
  opening_id: OpeningID;
}

export interface GetOpeningRequest {
  id: OpeningID;
}

export interface FilterOpeningsRequest {
  state?: OpeningState[];
  from_date?: Date;
  to_date?: Date;
  pagination_key?: string;
  limit?: number;
}

export interface ChangeOpeningStateRequest {
  opening_id: OpeningID;
  from_state: OpeningState;
  to_state: OpeningState;
}

export interface UpdateOpeningRequest {
  opening_id: OpeningID;
  // TODO: Decide what fields are allowed to be updated
}

export interface GetOpeningWatchersRequest {
  opening_id: OpeningID;
}

export interface AddOpeningWatchersRequest {
  opening_id: OpeningID;
  emails: EmailAddress[];
}

export interface RemoveOpeningWatcherRequest {
  opening_id: OpeningID;
  email: EmailAddress;
}
