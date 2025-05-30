import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model AddEmployerPostRequest {
    @minLength(1)
    @maxLength(4096)
    content: string;

    @maxItems(3)
    tag_ids: string[];
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
    tag_ids: string[];
}

model DeleteEmployerPostRequest {
    @doc("The id of the post to delete")
    post_id: string;
}

model ListEmployerPostsRequest {
    @doc("The pagination key. If nothing is passed, the newest posts are returned")
    pagination_key?: string;

    @doc("If nothing is passed, 5 Posts are returned")
    @minValue(1)
    @maxValue(40)
    limit?: integer;
}

model ListEmployerPostsResponse {
    posts: EmployerPost[];
    pagination_key: string;
}

model GetEmployerPostRequest {
    @doc("The id of the post to get")
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

@route("/employer/list-posts")
interface ListEmployerPosts {
    @tag("Posts")
    @post
    @doc("List employer posts in the descending order of the updated_at field. If there are multiple posts with the same updated_at timestamp, they are further sorted in the order of their post_ids")
    @useAuth(EmployerAuth)
    listEmployerPosts(@body request: ListEmployerPostsRequest): {
        @statusCode statusCode: 200;
        @body response: ListEmployerPostsResponse;
    };
}

@route("/employer/get-post")
interface GetEmployerPost {
    @tag("Posts")
    @post
    @useAuth(EmployerAuth)
    getEmployerPost(@body request: GetEmployerPostRequest): {
        @statusCode statusCode: 200;
        @body response: EmployerPost;
    };
}
