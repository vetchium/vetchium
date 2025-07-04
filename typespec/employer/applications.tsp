import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

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

model Endorser {
    full_name: string;
    short_bio: string;
    handle: string;
    current_company_domains?: string[];
}

model Application {
    id: string;
    cover_letter?: string;
    created_at: utcDateTime;
    hub_user_handle: string;
    hub_user_name: string;
    hub_user_short_bio: string;
    hub_user_last_employer_domains?: string[];
    state: ApplicationState;
    color_tag?: ApplicationColorTag;

    @doc("List of all the Endorsers who have endorsed the Application")
    endorsers: Endorser[];

    @doc("The scores of the Application by various models")
    scores: ModelScore[];
}

model SetApplicationColorTagRequest {
    application_id: string;
    color_tag: ApplicationColorTag;
}

model RemoveApplicationColorTagRequest {
    application_id: string;
}

model ShortlistApplicationRequest {
    application_id: string;
}

model RejectApplicationRequest {
    application_id: string;
}

model GetResumeRequest {
    application_id: string;
    // TODO: In future, add some kind of versioning here
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

@route("/employer/get-resume")
interface GetResume {
    @tag("Applications")
    @doc("Requires any of ${Admin}, ${ApplicationsCRUD} or ${ApplicationsViewer} roles")
    @post
    getResume(@body request: GetResumeRequest): {
        @statusCode statusCode: 200;
    };
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
