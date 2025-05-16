import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/achievements.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model GetEmployerCountsRequest {
    domain: string;
}

model EmployerCounts {
    active_openings_count: uint32;
    verified_employees_count: uint32;
}

@route("/internal/get-employer-counts")
interface GetEmployerCounts {
    @post
    @tag("Granger")
    getEmployerCounts(@body request: GetEmployerCountsRequest): {
        @statusCode statusCode: 200;
        @body EmployerCounts: EmployerCounts;
    };
}
