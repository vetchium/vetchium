import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/posts.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model AddPostCommentRequest {
    post_id: string;

    @doc("The comment to add to the post")
    @minLength(1)
    @maxLength(4096)
    content: string;
}

model AddPostCommentResponse {
    post_id: string;
    comment_id: string;
}

model GetPostCommentsRequest {
    post_id: string;

    @doc("The comments are returned ordered by the newest of created_at timestamp. If two comments are created at the same time, the one with the higher id is returned first. The pagination key is the id of the last comment in the previous response. The item corresponding to the passed pagination_key will not be included in the response. If an invalid pagination_key is provided (comment does not exist), the API will return comments from the beginning (first page).")
    pagination_key?: string;

    @doc("If nothing is passed, 10 Comments are returned")
    @minValue(1)
    @maxValue(40)
    limit?: integer;
}

model PostComment {
    id: string;
    content: string;
    author_name: string;
    author_handle: Handle;
    created_at: utcDateTime;
}

model DisablePostCommentsRequest {
    post_id: string;

    @doc("If true, all existing comments for the post will be deleted. If false, new comments will not be allowed, but existing comments will not be deleted.")
    delete_existing_comments: boolean;
}

model EnablePostCommentsRequest {
    post_id: string;
}

model DeletePostCommentRequest {
    post_id: string;
    comment_id: string;
}

model DeleteMyCommentRequest {
    post_id: string;
    comment_id: string;
}

@route("/hub/add-post-comment")
interface AddPostComment {
    @doc("Adds a comment to a Hub User Post. This is not applicable to Employer Posts. We do not support nested comments.")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    addPostComment(@body request: AddPostCommentRequest): {
        @statusCode statusCode: 200;
        @body response: AddPostCommentResponse;
    } | {
        @doc("Comments are disabled for this post")
        @statusCode
        statusCode: 403;
    } | {
        @doc("Post not found")
        @statusCode
        statusCode: 404;
    };
}

@route("/hub/get-post-comments")
interface GetPostComments {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    getPostComments(@body request: GetPostCommentsRequest): {
        @statusCode statusCode: 200;
        @body response: PostComment[];
    } | {
        @doc("Post not found")
        @statusCode
        statusCode: 404;
    };
}

@route("/hub/disable-post-comments")
interface DisablePostComments {
    @doc("Disables the ability to add comments for a post authored by the logged in user. If comments are already disabled, no error will be returned.")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    disablePostComments(@body request: DisablePostCommentsRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Logged in user is not the author of this post")
        @statusCode
        statusCode: 403;
    } | {
        @doc("Post not found")
        @statusCode
        statusCode: 404;
    };
}

@route("/hub/enable-post-comments")
interface EnablePostComments {
    @doc("Enables the ability to add comments for a post authored by the logged in user. If comments are already enabled, no error will be returned.")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    enablePostComments(@body request: EnablePostCommentsRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Logged in user is not the author of this post")
        @statusCode
        statusCode: 403;
    } | {
        @doc("Post not found")
        @statusCode
        statusCode: 404;
    };
}

@route("/hub/delete-post-comment")
interface DeletePostComment {
    @doc("Deletes a comment on the post authored by the logged in user. If the comment is not found, no error will be returned.")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    deletePostComment(@body request: DeletePostCommentRequest): {
        @statusCode statusCode: 200;
    } | {
        @doc("Logged in user is not the author of this post")
        @statusCode
        statusCode: 403;
    } | {
        @doc("Post not found")
        @statusCode
        statusCode: 404;
    };
}

@route("/hub/delete-my-comment")
interface DeleteMyComment {
    @doc("Deletes a comment authored by the logged in user. If the comment is not found, no error will be returned.")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    deleteMyComment(@body request: DeleteMyCommentRequest): {
        @statusCode statusCode: 200;
    };
}
