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
    opening_type?: OpeningType;
    companies?: string[];
    experience?: ExperienceRange;
    salary_range?: SalaryRange;

    @doc("If nothing is passed, all countries are considered. If countries and locations are passed, Openings matching either will be considered.")
    countries?: CountryCode[];

    @doc("If nothing is passed, all jobs in the HubUser's Countries list set in the profile are considered. Remember to set the Countries list during the Hub user signup.")
    locations?: LocationFilter[];

    @doc("If nothing is passed, all education levels are considered")
    min_education_level?: EducationLevel;

    remote_timezones?: TimeZone[];
    created_at: utcDateTime;
    pagination_key?: int64;

    @doc("If nothing is passed, 40 Openings are returned")
    limit?: integer;
}

model Opening {
    opening_id_within_company: string;
    company_domain: string;
    company_name: string;
    logo_url?: url;
    job_title: string;
    locations?: string[];

    /* In future, we will show company avg response time, response rate and some kind of score for how the company treats the applicants */
}

@route("/openings")
interface FindOpenings {
    @tag("Openings")
    @get
    findOpenings(@query query: FindOpeningsRequest): {
        @statusCode statusCode: 200;
    };
}
