import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

import "../common/common.tsp";
import "../common/vtags.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace Vetchium;

model Immi {
    id: string;
    content: string;
    tags: string[];
}

model AddImmiRequest {
    @doc("The content of the Immi")
    @minLength(1)
    @maxLength(4096)
    content: string;

    @maxItems(3)
    tag_ids: VTagID[];
}

model ImmiComment {
    id: string;
    content: string;
    in_reply_to: string;
    created_at: string;
    upvotes: integer;
    downvotes: integer;

    @doc("Whether the comment was created by the currently logged in user")
    is_created_by_me: boolean;
}

model AddImmiCommentRequest {
    @doc("The content of the comment")
    @minLength(1)
    @maxLength(512)
    content: string;

    @doc("The ID of the comment to which this comment is a reply. If nothing is provided, the new comment will be a top-level comment.")
    in_reply_to: string;
}

model UpvoteImmiCommentRequest {
    comment_id: string;
}

model DownvoteImmiCommentRequest {
    comment_id: string;
}

model UnvoteImmiCommentRequest {
    comment_id: string;
}

model GetImmiCommentsRequest {
    @doc("The ID of the Immi to get comments for")
    immi_id: string;

    @doc("Pagination key to fetch the next page of comments. The comments are sorted based on the upvotes of the top-level comments. Use the commentID of the last top-level comment as the pagination_key.")
    pagination_key?: string;

    @doc("Number of top-level comments to fetch per page")
    @minValue(1)
    @maxValue(100)
    limit: integer;
}

model GetImmiCommentsResponse {
    @doc("The comments for the Immi")
    comments: ImmiComment[];

    @doc("Pagination key to fetch the next page of comments. The comments are sorted based on the upvotes of the top-level comments. Use the commentID of the last top-level comment as the pagination_key.")
    pagination_key?: string;
}

model DeleteImmiCommentRequest {
    @doc("The ID of the comment to delete. Only the comment creator can delete their own comment.")
    comment_id: string;
}

model DeleteImmiRequest {
    @doc("The ID of the Immi to delete. Only the Immi creator can delete their own Immi.")
    immi_id: string;
}
