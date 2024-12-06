import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/applications.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchi;

model MyApplicationsRequest {
    @doc("Filter by Application state. If nothing is provided all Applications in the Hub are fetched")
    state?: ApplicationState;

    @doc("Pagination key to fetch the next page of Applications. The Applications are sorted by the created_at timestamp in descending order with newest Applications first and if there are multiple Applications with the same created_at timestamp, they are sorted by the ApplicationID in ascending order. Use the ApplicationID as the pagination_key")
    pagination_key?: string;

    @doc("Number of Applications to fetch per page")
    @minValue(1)
    @maxValue(40)
    limit: integer;
}

model HubApplication {
    application_id: string;
    state: ApplicationState;
    opening_id: string;
    opening_title: string;
    employer_name: string;
    employer_domain: string;
    created_at: string;
}

@route("/hub/my-applications")
interface MyApplications {
    @tag("Applications")
    @post
    @useAuth(HubAuth)
    myApplications(@body request: MyApplicationsRequest): {
        @statusCode statusCode: 200;
        @body applications: HubApplication[];
    };
}

model WithdrawApplicationRequest {
    application_id: string;
}

@route("/hub/withdraw-application")
interface WithdrawApplication {
    @tag("Applications")
    @post
    @useAuth(HubAuth)
    withdrawApplication(@body request: WithdrawApplicationRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Application not found")
        @statusCode
        statusCode: 404;
    } | {
        @doc("Application is not in a state that allows for withdrawal")
        @statusCode
        statusCode: 422;
    };
}
