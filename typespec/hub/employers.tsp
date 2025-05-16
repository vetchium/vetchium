import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model FilterEmployersRequest {
    @minLength(2)
    prefix: string;
}

model HubEmployer {
    domain: string;
    name: string;
    ascii_name: string;
}

model FilterEmployersResponse {
    employers: HubEmployer[];
}

model GetEmployerDetailsRequest {
    domain: string;
}

model HubEmployerDetails {
    name: string;
    verified_employees_count: uint32;
    active_openings_count: uint32;
}

@route("/hub/filter-employers")
interface FilterEmployers {
    @post
    @useAuth(HubAuth)
    @tag("HubEmployer")
    filterEmployers(@body request: FilterEmployersRequest): {
        @statusCode statusCode: 200;
        @body FilterEmployersResponse: FilterEmployersResponse;
    };
}

@route("/hub/logo/{domain}/{size}")
interface GetEmployerLogo {
    @get
    @useAuth(HubAuth)
    @tag("HubEmployer")
    getEmployerLogo(@path domain: string, @path size: string): {
        @statusCode statusCode: 200;
        @header contentType: "image/png" | "image/jpeg" | "image/gif";
        @body content: bytes;
    };
}

@route("/hub/get-employer-details")
interface GetEmployerDetails {
    @post
    @useAuth(HubAuth)
    @tag("HubEmployer")
    getEmployerDetails(@body request: GetEmployerDetailsRequest): {
        @statusCode statusCode: 200;
        @body HubEmployerDetails: HubEmployerDetails;
    };
}
