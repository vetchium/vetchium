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

model SalaryRange {
    currency: Currency;
    min: float64;
    max: float64;
}

model LocationFilter {
    country_code: CountryCode;
    city: string;
}

model FindHubOpeningsRequest {
    @doc("If nothing is passed, the HubUser's resident country will be considered.")
    country_code?: CountryCode;

    @doc("If nothing is passed, all cities under the country_code are considered. No error will be reported if a passed city does not belong to the country_code.")
    @maxItems(4)
    cities?: string[];

    @doc("If nothing is passed, all opening types are considered")
    opening_types?: OpeningType[];

    @doc("If nothing is passed, all companies are considered")
    @maxItems(4)
    company_domains?: string[];

    @doc("If nothing is passed, all experience ranges are considered")
    experience?: ExperienceRange;

    @doc("If nothing is passed, all salary ranges are considered")
    salary_range?: SalaryRange;

    @doc("If nothing is passed, all education levels are considered")
    min_education_level?: EducationLevel;

    @doc("If a value is passed, only Openings that either explicitly mention this timezone in the Opening (or) globally remote, are considered. If nothing is passed, all Openings under the country_code irrespective of their remote status are considered.")
    @maxItems(4)
    remote_timezones?: TimeZone[];

    @doc("If a value is passed, either Openings that either explicitly mention this country code in the Opening (or) globally remote, are considered. If nothing is passed, all Openings under the country_code irrespective of their remote status are considered.")
    @maxItems(4)
    remote_country_codes?: CountryCode[];

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

@route("/hub/find-openings")
interface FindHubOpenings {
    @tag("HubOpenings")
    @post
    @useAuth(HubAuth)
    findOpenings(@query query: FindHubOpeningsRequest): {
        @statusCode statusCode: 200;
        @body HubOpening: [];
    };
}
