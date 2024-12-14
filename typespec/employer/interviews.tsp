import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/interviews.tsp";
using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

@route("/employer/rsvp-interview")
interface EmployerRSVPInterview {
    @tag("Interviews")
    @doc("The OrgUser doing this must be an Interviewer in the Interview")
    @post
    @useAuth(EmployerAuth)
    employerRSVPInterview(@body request: RSVPInterviewRequest): {
        @statusCode statusCode: 200;
    } | {
        @statusCode statusCode: 400;
        @body error: ValidationErrors;
    } | {
        @doc("The Interview is not in the correct state for this action")
        @statusCode
        statusCode: 422;
    };
}
