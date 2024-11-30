import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

union ApplicationColorTag {
    Green: "GREEN",
    Yellow: "YELLOW",
    Red: "RED",
}

model GetApplicationsRequest {
    state: ApplicationState;

    @doc("If passed, the Applicants whose name matches this query will be returned")
    @maxLength(25)
    search_query?: string;

    @doc("If passed, only Applications of the passed color tag will be returned")
    color_tag_filter?: ApplicationColorTag;

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
    from_state: ApplicationState;
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

model AddApplicationColorTagRequest {
    application_id: string;
    color_tag: ApplicationColorTag;
}

@route("/employer/add-application-color-tag")
interface AddApplicationColorTag {
    @tag("Applications")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    addApplicationColorTag(@body request: AddApplicationColorTagRequest): {};
}

model RemoveApplicationColorTagRequest {
    application_id: string;
}

@route("/employer/remove-application-color-tag")
interface RemoveApplicationColorTag {
    @tag("Applications")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    removeApplicationColorTag(
        @body request: RemoveApplicationColorTagRequest,
    ): {};
}

model ShortlistApplicationRequest {
    application_id: string;
}

@route("/employer/shortlist-application")
interface ShortlistApplication {
    @tag("Applications")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    shortlistApplication(@body request: ShortlistApplicationRequest): {};
}

model RejectApplicationRequest {
    application_id: string;
}

@route("/employer/reject-application")
interface RejectApplication {
    @tag("Applications")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    rejectApplication(@body request: RejectApplicationRequest): {};
}
