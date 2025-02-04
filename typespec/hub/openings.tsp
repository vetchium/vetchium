import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model ExperienceRange {
    @minValue(0)
    @maxValue(99)
    yoe_min: integer;

    @minValue(1)
    @maxValue(100)
    yoe_max: integer;
}

model LocationFilter {
    country_code: CountryCode;
    city: string;
}

model FindHubOpeningsRequest {
    country_code: CountryCode;

    @doc("If nothing is passed, all cities under the country_code are considered. No error will be reported if a passed city does not belong to the country_code.")
    @maxItems(4)
    cities?: string[];

    @doc("If nothing is passed, all opening types are considered")
    opening_types?: OpeningType[];

    @doc("If nothing is passed, all companies are considered")
    @maxItems(4)
    company_domains?: string[];

    @doc("If nothing is passed, all experience ranges are considered")
    experience_range?: ExperienceRange;

    @doc("If nothing is passed, all salary ranges are considered")
    salary_range?: Salary;

    @doc("If nothing is passed, all education levels are considered")
    min_education_level?: EducationLevel;

    @doc("If nothing is passed, all tags are considered")
    tags?: OpeningTagID[];

    @doc("The strings passed here will be used to filter on the Opening Titles")
    terms?: string[];

    pagination_key?: int64;

    @doc("If nothing is passed, 40 Openings are returned")
    @minValue(1)
    @maxValue(100)
    limit?: integer;
}

model HubOpening {
    opening_id_within_company: string;
    company_domain: string;
    company_name: string;
    job_title: string;
    jd: string;
    pagination_key: int64;

    /* In future, we will show company's: 
        logo urls
        avg response time
        response rate
        some kind of score for how the company treats the applicants, etc. */
}

model GetHubOpeningDetailsRequest {
    opening_id_within_company: string;
    company_domain: string;
}

model HubOpeningDetails {
    company_domain: string;
    company_name: string;
    created_at: utcDateTime;
    education_level: EducationLevel;
    yoe_min: integer;
    yoe_max: integer;
    hiring_manager_name: string;
    hiring_manager_vetchi_handle?: string; // TODO: Not done yet
    jd: string;
    job_title: string;
    opening_id_within_company: string;
    opening_type: OpeningType;
    pagination_key: int64;
    recruiter_name: string;
    salary?: Salary;
    state: OpeningState;
}

model ApplyForOpeningRequest {
    opening_id_within_company: string;
    company_domain: string;
    resume: string;

    @maxLength(4096)
    cover_letter?: string;

    @maxLength(256)
    filename: string;
}

model ApplyForOpeningResponse {
    application_id: string;
}

@route("/hub/find-openings")
interface FindHubOpenings {
    @tag("Openings")
    @post
    @useAuth(HubAuth)
    findOpenings(request: FindHubOpeningsRequest): {
        @statusCode statusCode: 200;
        @body HubOpening: [];
    };
}

@route("/hub/get-opening-details")
interface GetOpeningDetails {
    @tag("Openings")
    @post
    @useAuth(HubAuth)
    getOpeningDetails(@body request: GetHubOpeningDetailsRequest): {
        @statusCode statusCode: 200;
        @body HubOpeningDetails: HubOpeningDetails;
    };
}

@route("/hub/apply-for-opening")
interface ApplyForOpening {
    @tag("Openings")
    @post
    @useAuth(HubAuth)
    applyForOpening(@body request: ApplyForOpeningRequest): {
        @statusCode statusCode: 200;
        @body ApplyForOpeningResponse: ApplyForOpeningResponse;
    };
}

@route("/hub/filter-opening-tags")
interface FilterHubOpeningTags {
    @tag("Openings")
    @post
    @useAuth(HubAuth)
    filterHubOpeningTags(@body request: FilterOpeningTagsRequest): {
        @statusCode statusCode: 200;
        @body OpeningTag: [];
    };
}
