import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/education.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model ListHubUserEducationRequest {
    handle: Handle;
}

@route("/employer/list-hub-user-education")
interface EmployerEducation {
    @get
    @useAuth(EmployerAuth)
    listEducation(@query request: ListHubUserEducationRequest): {
        @statusCode statusCode: 200;
        @body Education: Education[];
    };
}
