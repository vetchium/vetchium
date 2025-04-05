import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

@maxLength(64)
scalar OpeningID extends string;

model OpeningInfo {
    id: OpeningID;
    title: string;
    positions: integer;
    filled_positions: integer;
    recruiter: OrgUserShort;
    hiring_manager: OrgUserShort;
    cost_center_name: CostCenterName;
    opening_type: OpeningType;
    state: OpeningState;
    created_at: utcDateTime;
    last_updated_at: utcDateTime;
}

@doc("The Opening object that should be used within the Employer API. Should not be exposed to the Hub APIs or the Hub UI.")
model Opening {
    id: OpeningID;

    @minLength(3)
    @maxLength(32)
    title: string;

    @minValue(1)
    @maxValue(20)
    positions: integer;

    @minValue(0)
    @maxValue(20)
    filled_positions: integer;

    @minLength(10)
    @maxLength(1024)
    jd: string;

    recruiter: OrgUserShort;
    hiring_manager: OrgUserShort;

    @doc("OrgUser emails of the members of the team for which the opening is created")
    @maxItems(10)
    hiring_team_members?: OrgUserShort[];

    @doc("Cost Center name under which the opening is created")
    cost_center_name: CostCenterName;

    @doc("Privates notes about the opening visible only to the employer")
    @maxLength(1024)
    employer_notes?: string;

    @doc("List of pre-configured Location titles where the opening is available")
    @maxItems(10)
    location_titles?: string[];

    @doc("List of ISO 3166-1 alpha-3 country codes where the opening is available for remote work. Special value ZZG to indicate Global")
    @maxItems(100)
    remote_country_codes?: CountryCode[];

    @doc("List of timezones where the opening is available for remote work")
    @maxItems(200)
    remote_timezones?: TimeZone[];

    opening_type: OpeningType;

    @doc("Minimum years of experience required for the opening")
    @minValue(0)
    @maxValue(100)
    yoe_min: integer;

    @minValue(1)
    @maxValue(100)
    @doc("Maximum years of experience allowed for the opening")
    yoe_max: integer;

    min_education_level: EducationLevel;
    salary?: Salary;

    @doc("Current state of the opening")
    state: OpeningState;

    @doc("List of tags associated with the opening")
    @maxItems(3)
    tags?: OpeningTag[];

    created_at: utcDateTime;
    last_updated_at: utcDateTime;
}

model CreateOpeningRequest {
    @minLength(3)
    @maxLength(32)
    title: string;

    @minValue(1)
    @maxValue(20)
    positions: integer;

    @minLength(10)
    @maxLength(8192)
    jd: string;

    recruiter: EmailAddress;
    hiring_manager: EmailAddress;

    @doc("OrgUser emails of the members of the team for which the opening is created")
    @maxItems(10)
    hiring_team_members?: EmailAddress[];

    @doc("Cost Center name under which the opening is created")
    cost_center_name: CostCenterName;

    @doc("Privates notes about the opening visible only to the employer")
    @maxLength(1024)
    employer_notes?: string;

    @doc("List of pre-configured Location titles where the opening is available")
    @maxItems(10)
    location_titles?: string[];

    @doc("List of ISO 3166-1 alpha-3 country codes where the opening is available for remote work. Special value ZZG to indicate Global")
    @maxItems(100)
    remote_country_codes?: CountryCode[];

    @doc("List of timezones where the opening is available for remote work")
    @maxItems(200)
    remote_timezones?: TimeZone[];

    opening_type: OpeningType;

    @minValue(0)
    @maxValue(100)
    @doc("Minimum years of experience required for the opening")
    yoe_min: integer;

    @minValue(1)
    @maxValue(100)
    @doc("Maximum years of experience allowed for the opening")
    yoe_max: integer;

    min_education_level: EducationLevel;
    salary?: Salary;

    @doc("List of tags to be used for the Opening. Atleast one tag should be specified either here or in new_tags.")
    @maxItems(3)
    tags?: OpeningTagID[];

    @doc("These are list of new tags that should be created and then the Opening should be associated with them. Atleast one tag should be specified either here or in tags.")
    @maxItems(3)
    new_tags?: string[];
}

model CreateOpeningResponse {
    opening_id: OpeningID;
}

model GetOpeningRequest {
    @doc("ID of the Opening to fetch")
    id: OpeningID;
}

model FilterOpeningsRequest {
    @doc("Pagination key to fetch the next page of Openings. The Openings are sorted by the OpeningID in a chronological way")
    pagination_key?: OpeningID;

    @doc("Filter by Opening state. If nothing is provided all Openings in any of DraftOpening, ActiveOpening or SuspendedOpening states will be returned")
    state?: OpeningState[];

    // TODO: Perhaps we should put a limit on the date range ?
    @doc("Filter by Opening creation date. If not provided, (currentdate-30days) will be used")
    from_date?: plainDate;

    @doc("Filter by Opening creation date. If not provided, (currentdate+1day, tomorrow date) will be used")
    to_date?: plainDate;

    @maxValue(40)
    @doc("Number of Openings to return; 40 is the default if not specified")
    limit?: integer;
}

model ChangeOpeningStateRequest {
    opening_id: OpeningID;
    from_state: OpeningState;
    to_state: OpeningState;
}

model UpdateOpeningRequest {
    id: OpeningID;
    // TODO: Decide what fields are allowed to be updated
}

model GetOpeningWatchersRequest {
    opening_id: OpeningID;
}

model AddOpeningWatchersRequest {
    opening_id: OpeningID;
    emails: EmailAddress[];
}

model RemoveOpeningWatcherRequest {
    opening_id: OpeningID;
    email: EmailAddress;
}

@route("/employer/create-opening")
interface CreateOpening {
    @tag("Openings")
    @doc("Requires any of ${Admin}, ${OpeningsCRUD} roles")
    @post
    @useAuth(EmployerAuth)
    createOpening(@body createOpeningRequest: CreateOpeningRequest): {
        @statusCode statusCode: 200;
    } | {
        @statusCode statusCode: 400;
        @body error: ValidationErrors;
    } | {
        @doc("One or more of the provided values for locations, recruiters, hiring_team are invalid")
        @statusCode
        statusCode: 422;
    };
}

@route("/employer/get-opening")
interface GetOpening {
    @tag("Openings")
    @doc("Requires any of ${Admin}, $${OpeningsCRUD} {OpeningsViewer} roles")
    @post
    @useAuth(EmployerAuth)
    getOpening(@body getOpeningRequest: GetOpeningRequest): {
        @statusCode statusCode: 200;
        @body opening: Opening;
    };
}

@route("/employer/filter-openings")
interface FilterOpenings {
    @tag("Openings")
    @doc("Requires any of ${Admin}, ${OpeningsCRUD} or ${OpeningsViewer} roles")
    @post
    @useAuth(EmployerAuth)
    filterOpenings(@body filterOpeningsRequest: FilterOpeningsRequest): {
        @statusCode statusCode: 200;
        @body opening_infos: OpeningInfo[];
    };
}

@route("/employer/update-opening")
interface UpdateOpening {
    @tag("Openings")
    @doc("Requires any of ${Admin}, ${OpeningsCRUD} roles")
    @post
    @useAuth(EmployerAuth)
    updateOpening(@body updateOpeningRequest: UpdateOpeningRequest): {
        @statusCode statusCode: 200;
    } | {
        @statusCode statusCode: 400;
        @body error: ValidationErrors;
    };
}

@route("/employer/get-opening-watchers")
interface GetOpeningWatchers {
    @tag("Openings")
    @doc("Requires any of ${Admin}, ${OpeningsCRUD} or ${OpeningsViewer} roles")
    @post
    @useAuth(EmployerAuth)
    getOpeningWatchers(
        @body getOpeningWatchersRequest: GetOpeningWatchersRequest,
    ): {
        @statusCode statusCode: 200;
        @body org_users: OrgUserShort[];
    };
}

@route("/employer/add-opening-watchers")
interface AddOpeningWatchers {
    @tag("Openings")
    @doc("Requires any of ${Admin}, ${OpeningsCRUD} roles")
    @post
    @useAuth(EmployerAuth)
    addOpeningWatchers(@body addWatchersRequest: AddOpeningWatchersRequest):
        | {
              @statusCode statusCode: 200;
          }
        | {
              @doc("invalid opening_id")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("invalid or non-existent email addresses passed")
              @statusCode
              statusCode: 400;
          }
        | {
              @doc("adding watchers would exceed 25 watchers")
              @statusCode
              statusCode: 422;
          };
}

@route("/employer/remove-opening-watcher")
interface RemoveOpeningWatcher {
    @tag("Openings")
    @doc("Requires any of ${Admin}, ${OpeningsCRUD} roles")
    @post
    @useAuth(EmployerAuth)
    removeOpeningWatcher(
        @body removeWatcherRequest: RemoveOpeningWatcherRequest,
    ): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/change-opening-state")
interface ChangeOpeningState {
    @tag("Openings")
    @doc("Requires any of ${Admin}, ${OpeningsCRUD} roles")
    @post
    @useAuth(EmployerAuth)
    changeOpeningState(
        @body changeOpeningStateRequest: ChangeOpeningStateRequest,
    ): {
        @statusCode statusCode: 200;
    } | {
        @doc("invalid from_state")
        @statusCode
        statusCode: 409;
    } | {
        @doc("invalid transition from from_state to to_state")
        @statusCode
        statusCode: 422;
    };
}

@route("/employer/filter-opening-tags")
interface FilterOpeningTags {
    @tag("Openings")
    @post
    @useAuth(EmployerAuth)
    filterOpeningTags(
        @body filterOpeningTagsRequest: FilterOpeningTagsRequest,
    ): {
        @statusCode statusCode: 200;
        @body opening_tags: OpeningTag[];
    };
}
