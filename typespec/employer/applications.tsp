import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model GetApplicationsRequest {
    state: ApplicationState;

    @doc("If passed, the Applicants whose name matches this query will be returned")
    @maxLength(25)
    search_query?: string;

    @doc("The Opening ID for which Applications are to be fetched")
    opening_id: string;

    pagination_key?: string;

    @doc("If nothing is passed a default of 40 will be returned")
    @minValue(1)
    @maxValue(40)
    limit: int64;
}

model Application {
    id: string;
    cover_letter?: string;
    created_at: utcDateTime;
    filename: string;
    hub_user_handle: string;
    hub_user_last_employer_domain?: string;
    resume: string;
    state: ApplicationState;
}

model UpdateApplicationStateRequest {
    id: string;
    from_datestate: ApplicationState;
    to_state: ApplicationState;
}

@route("/employer/get-applications")
interface GetApplications {
    @tag("Applications")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} or ${ApplicationsViewer} roles")
    @post
    getApplications(@body request: GetApplicationsRequest): {
        @statusCode statusCode: 200;
        @body applications: Application[];
    };
}

@route("/employer/update-application-state")
interface UpdateApplicationState {
    @tag("Applications")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    updateApplicationState(@body request: UpdateApplicationStateRequest):
        | {
              @statusCode statusCode: 200;
          }
        | {
              @doc("Invalid state transition requested")
              @statusCode
              statusCode: 422;
          }
        | {
              @doc("Application not found")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("Wrong/Obsolete from_state passed")
              @statusCode
              statusCode: 409;
          };
}
