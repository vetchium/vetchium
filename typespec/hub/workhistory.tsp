import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/openings.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model AddWorkHistoryRequest {
    employer_domain: string;
    title: string;
    start_date: string;
    end_date?: string;

    @maxLength(1024)
    description?: string;
}

model AddWorkHistoryResponse {
    work_history_id: string;
}

model WorkHistory {
    id: string;
    employer_domain: string;

    @doc("In case the employer is not yet onboarded, the below employer_* fields will be empty")
    employer_name?: string;

    employer_logo_url?: string;

    @doc("The title of the Job")
    title: string;

    start_date: string;
    end_date?: string;

    @maxLength(1024)
    description?: string;
}

model UpdateWorkHistoryRequest {
    id: string;
    title: string;
    start_date: string;

    @doc("If end date already existed and not passed now, the old one will be removed")
    end_date?: string;

    @doc("If description already existed and not passed now, the old one will be removed")
    @maxLength(1024)
    description?: string;
}

model ListWorkHistoryRequest {
    @doc("The user handle of the user to list the work history for. If not provided, the work history of the current user will be listed.")
    user_handle?: string;
}

model DeleteWorkHistoryRequest {
    id: string;
}

@route("/hub/add-work-history")
interface AddWorkHistory {
    @post
    @useAuth(HubAuth)
    @tag("WorkHistory")
    addWorkHistory(@body request: AddWorkHistoryRequest): {
        @statusCode statusCode: 200;
        @body AddWorkHistoryResponse: AddWorkHistoryResponse;
    } | {
        @doc("Max limit of work history reached")
        @statusCode
        statusCode: 422;
    };
}

@route("/hub/update-work-history")
interface UpdateWorkHistory {
    @put
    @useAuth(HubAuth)
    @tag("WorkHistory")
    updateWorkHistory(@body request: UpdateWorkHistoryRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/delete-work-history")
interface DeleteWorkHistory {
    @post
    @useAuth(HubAuth)
    @tag("WorkHistory")
    deleteWorkHistory(@body request: DeleteWorkHistoryRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/list-work-history")
interface ListWorkHistory {
    @post
    @useAuth(HubAuth)
    @tag("WorkHistory")
    listWorkHistory(@body request: ListWorkHistoryRequest): {
        @statusCode statusCode: 200;
        @body WorkHistory: WorkHistory[];
    };
}
