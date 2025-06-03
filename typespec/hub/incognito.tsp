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

model IncognitoPostComment {
    id: string;
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

model UpvoteIncognitoPostCommentRequest {
    comment_id: string;
}

model DownvoteIncognitoPostCommentRequest {
    comment_id: string;
}

model UnvoteIncognitoPostCommentRequest {
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
    @doc("The ID of the comment to delete. Only the comment creator can delete their own comment.")
    comment_id: string;
}

model DeleteIncognitoPostRequest {
    @doc("The ID of the IncognitoPost to delete. Only the IncognitoPost creator can delete their own IncognitoPost. Deleting an IncognitoPost does not delete the comments.")
    incognito_post_id: string;
}
