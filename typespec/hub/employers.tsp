import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

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
