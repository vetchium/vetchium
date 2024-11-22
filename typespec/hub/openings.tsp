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

model FindOpeningsRequest {
    @doc("If nothing is passed, all opening types are considered")
    opening_type?: OpeningType;

    @doc("If nothing is passed, all companies are considered")
    companies?: string[];

    @doc("If nothing is passed, all experience ranges are considered")
    experience?: ExperienceRange;

    @doc("If nothing is passed, all salary ranges are considered")
    salary_range?: SalaryRange;

    @doc("If nothing is passed, all countries are considered. If countries and locations are passed, Openings that match either of them will be considered.")
    countries?: CountryCode[];

    @doc("If nothing is passed, all jobs in the HubUser's Countries list set in the profile are considered. If a city name is passed, only the Openings that explicitly mention this City are considered. TODO: Remember to set the Countries list during the Hub user signup.")
    locations?: LocationFilter[];

    @doc("If nothing is passed, all education levels are considered")
    min_education_level?: EducationLevel;

    @doc("If a value is passed, only Openings that explicitly mention this timezone in the description are considered")
    remote_timezones?: TimeZone[];

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
    logo_url?: url;
    job_title: string;
    locations?: string[];

    /* In future, we will show company's: 
        avg response time
        response rate
        some kind of score for how the company treats the applicants */
}

@route("/find-hub-openings")
interface FindHubOpenings {
    @tag("HubOpenings")
    @post
    @useAuth(HubAuth)
    findOpenings(@query query: FindOpeningsRequest): {
        @statusCode statusCode: 200;
    };
}
