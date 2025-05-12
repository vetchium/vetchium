import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model EmployerPost {
    id: string;
    content: string;
    tags: string[];
    company_domain: string;
    created_at: string;
}

model AddEmployerPostRequest {
    @minLength(1)
    @maxLength(4096)
    content: string;

    @maxItems(3)
    tags: string[];
}

model AddEmployerPostResponse {
    post_id: string;
}

model UpdateEmployerPostRequest {
    @doc("The id of the post to update")
    post_id: string;

    @minLength(1)
    @maxLength(4096)
    content: string;

    @maxItems(3)
    tags: string[];
}

model DeleteEmployerPostRequest {
    @doc("The id of the post to delete")
    post_id: string;
}

@route("/employer/add-post")
interface AddEmployerPost {
    @tag("Posts")
    @post
    @useAuth(EmployerAuth)
    addEmployerPost(@body request: AddEmployerPostRequest): {
        @statusCode statusCode: 200;
        @body response: AddEmployerPostResponse;
    };
}

@route("/employer/update-post")
interface UpdateEmployerPost {
    @tag("Posts")
    @post
    @useAuth(EmployerAuth)
    updateEmployerPost(@body request: UpdateEmployerPostRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/employer/delete-post")
interface DeleteEmployerPost {
    @tag("Posts")
    @post
    @useAuth(EmployerAuth)
    deleteEmployerPost(@body request: DeleteEmployerPostRequest): {
        @statusCode statusCode: 200;
    };
}
