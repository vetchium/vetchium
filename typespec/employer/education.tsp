import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/education.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model ListHubUserEducationRequest {
    handle: Handle;
}

@route("/employer/list-hub-user-education")
interface EmployerEducation {
    @tag("Education")
    @post
    @useAuth(EmployerAuth)
    listEducation(@body request: ListHubUserEducationRequest): {
        @statusCode statusCode: 200;
        @body Education: Education[];
    };
}
