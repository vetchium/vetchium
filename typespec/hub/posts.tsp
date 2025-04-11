import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/posts.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model AddPostRequest {
    @minLength(1)
    @maxLength(4096)
    content: string;

    @maxItems(3)
    tag_ids: VTagID[];

    @maxItems(3)
    new_tags: VTagName[];

    // TODO: Images
}

model AddPostResponse {
    post_id: string;
}

model GetUserPostsRequest {
    @doc("The handle of the user to get the posts from. If not provided, the posts of the logged in user are returned.")
    handle?: Handle;

    @doc("The pagination key to get the next page of posts. The item corresponding to the passed key will not be included in the response. No guarantees are made about the order of the posts in the response. There may be duplicates in the response and the client is expected to handle it.")
    pagination_key?: string;

    @doc("If nothing is passed, 10 Posts are returned")
    @minValue(1)
    @maxValue(40)
    limit?: integer;
}

model GetUserPostsResponse {
    posts: Post[];
    pagination_key: string;
}

@route("/hub/add-post")
interface AddPost {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    addPost(@body request: AddPostRequest): {
        @statusCode statusCode: 200;
        @body response: AddPostResponse;
    };
}

@route("/hub/get-user-posts")
interface GetUserPosts {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    getUserPosts(@body request: GetUserPostsRequest): {
        @statusCode statusCode: 200;
        @body response: GetUserPostsResponse;
    };
}
