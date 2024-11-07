import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "./common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

union OpeningType {
    FullTime: "FULL_TIME",
    PartTime: "PART_TIME",
    Contract: "CONTRACT",
    Internship: "INTERNSHIP",
    Unspecified: "UNSPECIFIED",
}

union EducationLevel {
    Bachelor: "BACHELOR",
    Master: "MASTER",
    Doctorate: "DOCTORATE",
    NotMatters: "NOT_MATTERS",
    Unspecified: "UNSPECIFIED",
}

model Salary {
    min_amount: integer;
    max_amount: integer;
    currency: Currency;
}

model Opening {
    @minLength(3)
    @maxLength(32)
    title: string;

    @minValue(1)
    @maxValue(20)
    positions: integer;

    @minLength(10)
    @maxLength(1024)
    jd: string;

    @doc("List of recruiters assigned to the opening")
    @minItems(1)
    @maxItems(10)
    recruiters: string[];

    hiring_manager: EmailAddress;

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
    @doc("Years of experience from 0 to 100")
    yoe?: integer;

    min_education_level?: EducationLevel;
    salary?: Salary;
}
