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

model SetApplicationColorTagRequest {
    application_id: string;
    color_tag: ApplicationColorTag;
}

@route("/employer/set-application-color-tag")
interface SetApplicationColorTag {
    @tag("Applications")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    setApplicationColorTag(@body request: SetApplicationColorTagRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Application not found")
        @statusCode
        statusCode: 404;
    } | {
        @doc("Application state is not ApplicationState.Applied")
        @statusCode
        statusCode: 422;
    };
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
    ): {
        @statusCode statusCode: 200;
    } | {
        @doc("Application not found")
        @statusCode
        statusCode: 404;
    } | {
        @doc("Application state is not ApplicationState.Applied")
        @statusCode
        statusCode: 422;
    };
}

model ShortlistApplicationRequest {
    application_id: string;
}

@route("/employer/shortlist-application")
interface ShortlistApplication {
    @tag("Applications")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    shortlistApplication(@body request: ShortlistApplicationRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Application not found")
        @statusCode
        statusCode: 404;
    } | {
        @doc("Application state is not ApplicationState.Applied")
        @statusCode
        statusCode: 422;
    };
}

model RejectApplicationRequest {
    application_id: string;
}

@route("/employer/reject-application")
interface RejectApplication {
    @tag("Applications")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} roles")
    @post
    rejectApplication(@body request: RejectApplicationRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Application not found")
        @statusCode
        statusCode: 404;
    } | {
        @doc("Application state is not ApplicationState.Applied")
        @statusCode
        statusCode: 422;
    };
}
