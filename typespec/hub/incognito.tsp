import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/vtags.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model IncognitoPost {
    incognito_post_id: string;
    content: string;
    tags: VTag[];
    created_at: utcDateTime;
    upvotes: integer;
    downvotes: integer;
    is_created_by_me: boolean;
    is_deleted: boolean;
}

model GetIncognitoPostRequest {
    @doc("The ID of the IncognitoPost to get")
    incognito_post_id: string;
}

model AddIncognitoPostRequest {
    @doc("The content of the IncognitoPost")
    @minLength(1)
    @maxLength(4096)
    content: string;

    @maxItems(3)
    tag_ids: VTagID[];
}

model AddIncognitoPostResponse {
    @doc("The ID of the IncognitoPost")
    incognito_post_id: string;
}

model IncognitoPostComment {
    comment_id: string;
    content: string;
    in_reply_to: string;
    created_at: string;
    upvotes: integer;
    downvotes: integer;

    @doc("Whether the comment was created by the currently logged in user")
    is_created_by_me: boolean;
}

model AddIncognitoPostCommentRequest {
    @doc("The content of the comment")
    @minLength(1)
    @maxLength(512)
    content: string;

    @doc("The ID of the comment to which this comment is a reply. If nothing is provided, the new comment will be a top-level comment.")
    in_reply_to: string;
}

model AddIncognitoPostCommentResponse {
    incognito_post_id: string;
    comment_id: string;
}

model UpvoteIncognitoPostCommentRequest {
    incognito_post_id: string;
    comment_id: string;
}

model DownvoteIncognitoPostCommentRequest {
    incognito_post_id: string;
    comment_id: string;
}

model UnvoteIncognitoPostCommentRequest {
    incognito_post_id: string;
    comment_id: string;
}

model GetIncognitoPostCommentsRequest {
    @doc("The ID of the IncognitoPost to get comments for")
    incognito_post_id: string;

    @doc("Pagination key to fetch the next page of comments. The comments are sorted based on the upvotes of the top-level comments. Use the commentID of the last top-level comment as the pagination_key.")
    pagination_key?: string;

    @doc("Number of top-level comments to fetch per page")
    @minValue(1)
    @maxValue(100)
    limit: integer;
}

model GetIncognitoPostCommentsResponse {
    @doc("The comments for the IncognitoPost")
    comments: IncognitoPostComment[];

    @doc("Pagination key to fetch the next page of comments. The comments are sorted based on the upvotes of the top-level comments. Use the commentID of the last top-level comment as the pagination_key.")
    pagination_key?: string;
}

model DeleteIncognitoPostCommentRequest {
    @doc("The ID of the IncognitoPost to which the comment belongs")
    incognito_post_id: string;

    @doc("The ID of the comment to delete. Only the comment creator can delete their own comment.")
    comment_id: string;
}

model DeleteIncognitoPostRequest {
    @doc("The ID of the IncognitoPost to delete. Only the IncognitoPost creator can delete their own IncognitoPost. Deleting an IncognitoPost does not delete the comments.")
    incognito_post_id: string;
}

@route("/hub/add-incognito-post")
interface AddIncognitoPost {
    @doc("AddIncognitoPost adds a new IncognitoPost. It is available for both free and paid tier users.")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    addIncognitoPost(@body request: AddIncognitoPostRequest): {
        @statusCode statusCode: 200;
        @body response: AddIncognitoPostResponse;
    };
}

@route("/hub/get-incognito-post")
interface GetIncognitoPost {
    @doc("GetIncognitoPost gets an IncognitoPost by its ID")
    @tag("Posts")
    @get
    @useAuth(HubAuth)
    getIncognitoPost(@body request: GetIncognitoPostRequest): {
        @statusCode statusCode: 200;
        @body response: IncognitoPost;
    };
}

@route("/hub/delete-incognito-post")
interface DeleteIncognitoPost {
    @doc("DeleteIncognitoPost deletes an IncognitoPost by its ID. Only the author of an IncognitoPost can delete it.")
    @tag("Posts")
    @delete
    @useAuth(HubAuth)
    deleteIncognitoPost(@path request: DeleteIncognitoPostRequest): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/add-incognito-post-comment")
interface AddIncognitoPostComment {
    @doc("AddIncognitoPostComment adds a new comment to an IncognitoPost")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    addIncognitoPostComment(@body request: AddIncognitoPostCommentRequest): {
        @statusCode statusCode: 200;
        @body response: AddIncognitoPostCommentResponse;
    };
}

@route("/hub/get-incognito-post-comments")
interface GetIncognitoPostComments {
    @doc("GetIncognitoPostComments gets the comments for an IncognitoPost")
    @tag("Posts")
    @get
    @useAuth(HubAuth)
    getIncognitoPostComments(@body request: GetIncognitoPostCommentsRequest): {
        @statusCode statusCode: 200;
        @body response: GetIncognitoPostCommentsResponse;
    };
}

@route("/hub/delete-incognito-post-comment")
interface DeleteIncognitoPostComment {
    @doc("DeleteIncognitoPostComment deletes a comment from an IncognitoPost. Only the comment author can delete their own comment.")
    @tag("Posts")
    @delete
    @useAuth(HubAuth)
    deleteIncognitoPostComment(
        @body request: DeleteIncognitoPostCommentRequest,
    ): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/upvote-incognito-post-comment")
interface UpvoteIncognitoPostComment {
    @doc("UpvoteIncognitoPostComment upvotes a comment on an IncognitoPost")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    upvoteIncognitoPostComment(
        @body request: UpvoteIncognitoPostCommentRequest,
    ): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/downvote-incognito-post-comment")
interface DownvoteIncognitoPostComment {
    @doc("DownvoteIncognitoPostComment downvotes a comment on an IncognitoPost")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    downvoteIncognitoPostComment(
        @body request: DownvoteIncognitoPostCommentRequest,
    ): {
        @statusCode statusCode: 200;
    };
}

@route("/hub/unvote-incognito-post-comment")
interface UnvoteIncognitoPostComment {
    @doc("UnvoteIncognitoPostComment unvotes a comment on an IncognitoPost")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    unvoteIncognitoPostComment(
        @body request: UnvoteIncognitoPostCommentRequest,
    ): {
        @statusCode statusCode: 200;
    };
}
