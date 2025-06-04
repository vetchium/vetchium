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

model AddIncognitoPostRequest {
    @doc("The content of the IncognitoPost")
    @minLength(1)
    @maxLength(1024)
    content: string;

    @doc("There must be atleast one tag for an IncognitoPost")
    @minItems(1)
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

    @doc("The ID of the parent comment. Null for top-level comments.")
    in_reply_to?: string;

    @doc("The creation timestamp of the comment")
    created_at: utcDateTime;

    upvotes: integer;
    downvotes: integer;

    @doc("Whether the comment was created by the currently logged in user")
    is_created_by_me: boolean;

    @doc("Whether the comment has been deleted by its author")
    is_deleted: boolean;

    @doc("The current user's vote on this comment: 'upvote', 'downvote', or null")
    my_vote?: "upvote" | "downvote";

    @doc("The nesting depth of the comment (0 for top-level, 1 for first reply, etc.)")
    depth: integer;
}

model AddIncognitoPostCommentRequest {
    @doc("The ID of the IncognitoPost to add the comment to")
    incognito_post_id: string;

    @doc("The content of the comment")
    @minLength(1)
    @maxLength(512)
    content: string;

    @doc("The ID of the comment to which this comment is a reply. If not provided, the new comment will be a top-level comment.")
    in_reply_to?: string;
}

model AddIncognitoPostCommentResponse {
    incognito_post_id: string;
    comment_id: string;
}

model GetIncognitoPostCommentsRequest {
    @doc("The ID of the IncognitoPost to get comments for")
    incognito_post_id: string;

    @doc("Pagination key to fetch the next page of comments. Comments are sorted by upvotes (descending), then by creation time (ascending) for ties.")
    pagination_key?: string;

    @doc("Number of comments to fetch per page")
    @minValue(1)
    @maxValue(100)
    limit: integer;

    @doc("The ID of the parent comment to get replies for. If not provided, gets top-level comments.")
    parent_comment_id?: string;

    @doc("Maximum depth of nested comments to include in response. Set to 0 for flat list.")
    @minValue(0)
    @maxValue(5)
    include_nested_depth?: integer = 0;
}

model GetIncognitoPostCommentsResponse {
    @doc("The comments for the IncognitoPost, sorted by upvotes (desc) then creation time (asc)")
    comments: IncognitoPostComment[];

    @doc("Pagination key for fetching the next page of comments")
    pagination_key?: string;

    @doc("Whether there are more comments available")
    has_more: boolean;

    @doc("Total count of comments at the requested level")
    total_count: integer;
}

model DeleteIncognitoPostCommentRequest {
    @doc("The ID of the IncognitoPost")
    incognito_post_id: string;

    @doc("The ID of the comment to delete. Only the comment creator can delete their own comment.")
    comment_id: string;
}

model UpvoteIncognitoPostCommentRequest {
    @doc("The ID of the IncognitoPost")
    incognito_post_id: string;

    @doc("The ID of the comment to vote on")
    comment_id: string;
}

model DownvoteIncognitoPostCommentRequest {
    @doc("The ID of the IncognitoPost")
    incognito_post_id: string;

    @doc("The ID of the comment to vote on")
    comment_id: string;
}

model UnvoteIncognitoPostCommentRequest {
    @doc("The ID of the IncognitoPost")
    incognito_post_id: string;

    @doc("The ID of the comment to unvote on")
    comment_id: string;
}

model GetIncognitoPostRequest {
    @doc("The ID of the IncognitoPost to get")
    incognito_post_id: string;
}

model DeleteIncognitoPostRequest {
    @doc("The ID of the IncognitoPost to delete")
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
    } | {
        @doc("Bad Request: Invalid request data (content too long/short, missing required fields, invalid tag IDs)")
        @statusCode
        statusCode: 400;
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
    } | {
        @doc("Bad Request: Invalid request data")
        @statusCode
        statusCode: 400;
    } | {
        @doc("Not Found: The specified incognito post does not exist")
        @statusCode
        statusCode: 404;
    };
}

@route("/hub/delete-incognito-post")
interface DeleteIncognitoPost {
    @doc("DeleteIncognitoPost deletes an IncognitoPost by its ID. Only the author of an IncognitoPost can delete it.")
    @tag("Posts")
    @delete
    @useAuth(HubAuth)
    deleteIncognitoPost(@body request: DeleteIncognitoPostRequest):
        | {
              @statusCode statusCode: 200;
          }
        | {
              @doc("Bad Request: Invalid request data")
              @statusCode
              statusCode: 400;
          }
        | {
              @doc("Forbidden: User is not the author of the incognito post")
              @statusCode
              statusCode: 403;
          }
        | {
              @doc("Not Found: The specified incognito post does not exist")
              @statusCode
              statusCode: 404;
          };
}

@route("/hub/add-incognito-post-comment")
interface AddIncognitoPostComment {
    @doc("AddIncognitoPostComment adds a new comment to an IncognitoPost")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    addIncognitoPostComment(@body request: AddIncognitoPostCommentRequest):
        | {
              @statusCode statusCode: 200;
              @body response: AddIncognitoPostCommentResponse;
          }
        | {
              @doc("Bad Request: Invalid request data (content too long/short, missing required fields)")
              @statusCode
              statusCode: 400;
          }
        | {
              @doc("Not Found: The specified incognito post or parent comment does not exist or has been deleted")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("Unprocessable Entity: Cannot add comment due to business rules (e.g., maximum comment depth reached, replying to deleted comment)")
              @statusCode
              statusCode: 422;
          };
}

@route("/hub/get-incognito-post-comments")
interface GetIncognitoPostComments {
    @doc("GetIncognitoPostComments gets the comments for an IncognitoPost. Use parent_comment_id to get replies to a specific comment, or leave it null for top-level comments.")
    @tag("Posts")
    @get
    @useAuth(HubAuth)
    getIncognitoPostComments(@body request: GetIncognitoPostCommentsRequest): {
        @statusCode statusCode: 200;
        @body response: GetIncognitoPostCommentsResponse;
    } | {
        @doc("Bad Request: Invalid request data (invalid limit values)")
        @statusCode
        statusCode: 400;
    } | {
        @doc("Not Found: The specified incognito post does not exist")
        @statusCode
        statusCode: 404;
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
    ):
        | {
              @statusCode statusCode: 200;
          }
        | {
              @doc("Bad Request: Invalid request data")
              @statusCode
              statusCode: 400;
          }
        | {
              @doc("Forbidden: User is not the author of the comment")
              @statusCode
              statusCode: 403;
          }
        | {
              @doc("Not Found: The specified incognito post or comment does not exist")
              @statusCode
              statusCode: 404;
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
        @doc("Comment upvote registered. 200 will be returned even if the user has already upvoted.")
        @statusCode
        statusCode: 200;
    } | {
        @doc("Bad Request: Invalid request data")
        @statusCode
        statusCode: 400;
    } | {
        @doc("This code will be returned in any of the following scenarios: 1) User has already downvoted the comment 2) The comment ID is invalid 3) User is attempting to upvote their own comment")
        @statusCode
        statusCode: 422;
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
        @doc("Comment downvote registered. 200 will be returned even if the user has already downvoted.")
        @statusCode
        statusCode: 200;
    } | {
        @doc("Bad Request: Invalid request data")
        @statusCode
        statusCode: 400;
    } | {
        @doc("This code will be returned in any of the following scenarios: 1) User has already upvoted the comment 2) The comment ID is invalid 3) User is attempting to downvote their own comment")
        @statusCode
        statusCode: 422;
    };
}

@route("/hub/unvote-incognito-post-comment")
interface UnvoteIncognitoPostComment {
    @doc("UnvoteIncognitoPostComment removes the user's vote from a comment on an IncognitoPost")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    unvoteIncognitoPostComment(
        @body request: UnvoteIncognitoPostCommentRequest,
    ): {
        @doc("Whatever upvote or downvote registered by the logged in user for the comment is removed. No error is returned even if the user has already unvoted or has never voted for the comment")
        @statusCode
        statusCode: 200;
    } | {
        @doc("Bad Request: Invalid request data")
        @statusCode
        statusCode: 400;
    } | {
        @doc("This code will be returned in any of the following scenarios: 1) The comment ID is invalid 2) User is attempting to unvote their own comment")
        @statusCode
        statusCode: 422;
    };
}
