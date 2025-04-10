import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/posts.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model AddPostRequest {
    @maxLength(4096)
    content: string;

    @maxItems(3)
    tags: string[];

    // TODO: Images
}

model AddPostResponse {
    post_id: string;
}

model GetTimelineRequest {
    @doc("The timeline id to get the timeline from. If not provided, the home timeline of the logged in user is returned.")
    timeline_id?: string;

    @doc("The pagination key to get the next page of posts. The item corresponding to the passed key will not be included in the response. No guarantees are made about the order of the posts in the response. There may be duplicates in the response and the client is expected to handle it.")
    pagination_key?: string;

    @doc("If nothing is passed, 40 Posts are returned")
    @minValue(1)
    @maxValue(100)
    limit?: integer;
}

model GetTimelineResponse {
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

@route("/hub/get-timeline")
interface GetTimeline {
    @tag("Posts")
    @get
    getTimeline(@body request: GetTimelineRequest): {
        @statusCode statusCode: 200;
        @body response: GetTimelineResponse;
    };
}
