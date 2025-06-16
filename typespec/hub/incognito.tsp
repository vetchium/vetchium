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
    upvotes_count: int32;
    downvotes_count: int32;
    score: int32;

    @doc("Whether the logged in user has upvoted this post")
    me_upvoted: boolean;

    @doc("Whether the logged in user has downvoted this post")
    me_downvoted: boolean;

    @doc("Whether the logged in user can upvote this post. Will be false if the user is blocked from upvoting this post or if the user has already voted for this post.")
    can_upvote: boolean;

    @doc("Whether the logged in user can downvote this post. Will be false if the user is blocked from downvoting this post or if the user has already voted for this post.")
    can_downvote: boolean;

    @doc("Whether this post was created by the currently logged in user")
    is_created_by_me: boolean;

    @doc("Whether this post has been deleted by its author")
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

    upvotes_count: int32;
    downvotes_count: int32;
    score: int32;

    @doc("Whether the logged in user has upvoted this comment")
    me_upvoted: boolean;

    @doc("Whether the logged in user has downvoted this comment")
    me_downvoted: boolean;

    @doc("Whether the logged in user can upvote this comment. Will be false if the user is blocked from upvoting this comment or if the user has already voted for this comment.")
    can_upvote: boolean;

    @doc("Whether the logged in user can downvote this comment. Will be false if the user is blocked from downvoting this comment or if the user has already voted for this comment.")
    can_downvote: boolean;

    @doc("Whether the comment was created by the currently logged in user")
    is_created_by_me: boolean;

    @doc("Whether the comment has been deleted by its author")
    is_deleted: boolean;

    @doc("The nesting depth of the comment (0 for top-level, 1 for first reply, etc.)")
    depth: int32;

    @doc("Total number of direct replies to this comment. This count excludes nested sub-replies and helps frontend developers determine when to show 'Load X more replies' buttons. For example, if a comment has 50 replies but only 5 are included in the response, the frontend can show 'Load 45 more replies'. This count is essential for the progressive loading user experience.")
    replies_count: int32;
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

union IncognitoPostCommentSortBy {
    Top: "top",
    New: "new",
    Old: "old",
}

model GetIncognitoPostCommentsRequest {
    @doc("The ID of the IncognitoPost to get comments for")
    incognito_post_id: string;

    @doc("Sort order for top-level comments. 'top' sorts by score (upvotes - downvotes) in descending order, 'new' sorts by creation time with newest first, 'old' sorts by creation time with oldest first. This sorting only applies to top-level comments (depth=0). All direct replies are sorted by score descending.")
    sort_by: IncognitoPostCommentSortBy = IncognitoPostCommentSortBy.Top;

    @doc("Pagination key to fetch the next page of top-level comments. This API uses cursor-based pagination on top-level comments only.")
    pagination_key?: string;

    @doc("Number of top-level comments to fetch per page.")
    @minValue(1)
    @maxValue(50)
    limit: int32 = 25;

    @doc("Number of direct replies to show for each top-level comment. This loads exactly this many immediate replies (depth+1 only) for each top-level comment, providing a predictable preview. Set to 0 to load only top-level comments without any replies.")
    @minValue(0)
    @maxValue(10)
    direct_replies_per_comment: int32 = 3;
}

model GetIncognitoPostCommentsResponse {
    @doc("Array of comments containing top-level comments with preview replies. Top-level comments (depth=0) are sorted according to the requested sort_by parameter. Each top-level comment includes a limited number of immediate replies (controlled by replies_preview_count) sorted by score descending. The response structure allows frontend developers to render the main discussion threads with popular reply previews. Additional replies can be loaded on-demand using the GetCommentReplies API. Backend developers should implement efficient queries to load top-level comments first, then fetch preview replies for those comments in a separate optimized query to avoid N+1 problems.")
    comments: IncognitoPostComment[];

    @doc("Pagination key for fetching the next page of top-level comments. Will be empty string when no more pages are available. This key should be passed unchanged to subsequent requests to maintain pagination state.")
    pagination_key: string;

    @doc("Total number of top-level comments for this post. This count helps frontend developers show pagination information like 'Showing 1-25 of 847 comments'. The count excludes replies and represents only the main discussion threads.")
    total_comments_count: int32;
}

model GetCommentRepliesRequest {
    @doc("The ID of the IncognitoPost containing the parent comment.")
    incognito_post_id: string;

    @doc("The ID of the parent comment whose replies should be loaded.")
    parent_comment_id: string;

    @doc("Pagination key for loading more replies to the specified parent comment.")
    pagination_key?: string;

    @doc("Number of replies to fetch for the parent comment.")
    @minValue(1)
    @maxValue(100)
    limit: int32 = 50;

    @doc("Whether to load only direct replies (depth+1) or allow nested replies. When true, loads only immediate children of the parent comment, making the 'Load X more replies' buttons predictable. When false, loads nested reply trees up to max_depth for 'continue this thread' functionality.")
    direct_only: boolean = true;

    @doc("Maximum depth of nested replies to include when direct_only=false. Only used for thread continuation. Ignored when direct_only=true.")
    @minValue(1)
    @maxValue(3)
    max_depth: int32 = 2;
}

model GetCommentRepliesResponse {
    @doc("Array of replies to the specified parent comment, sorted by score descending with creation time as tiebreaker. The depth field indicates nesting level relative to the original post. Replies include nested sub-replies up to the requested max_depth. Frontend developers should render these as nested conversations under the parent comment. Backend developers should ensure efficient loading of reply trees while respecting depth limits to prevent infinite recursion.")
    replies: IncognitoPostComment[];

    @doc("Pagination key for loading more replies to the same parent comment. Empty string indicates no more replies are available. This pagination is scoped specifically to the parent comment's replies.")
    pagination_key: string;

    @doc("Total number of direct replies to the parent comment. This count excludes nested sub-replies and helps frontend show information like 'Showing 1-50 of 234 replies'. Useful for implementing 'Load more replies' buttons with context.")
    total_replies_count: int32;

    @doc("The ID of the parent comment whose replies are being returned. This helps frontend developers associate the response with the correct parent comment, especially when multiple reply requests are in flight simultaneously.")
    parent_comment_id: string;
}

model GetIncognitoPostCommentPermalinkRequest {
    @doc("The ID of the IncognitoPost containing the target comment. This validates that the comment belongs to the specified post for security purposes.")
    incognito_post_id: string;

    @doc("The ID of the specific comment to retrieve with context. This comment becomes the focal point of the response and will be included along with contextual comments to provide a meaningful conversation view.")
    comment_id: string;

    @doc("Number of sibling comments to include at the same level as the target comment. This provides additional context about the broader conversation at that level. Siblings are sorted by score and limited to prevent excessive response sizes.")
    @minValue(0)
    @maxValue(10)
    context_siblings_count: int32 = 3;

    @doc("Number of replies to include under the target comment. This shows the immediate responses to the target comment, helping users understand its impact and the ensuing discussion. Replies are sorted by score descending.")
    @minValue(0)
    @maxValue(20)
    context_replies_count: int32 = 10;
}

model GetIncognitoPostCommentPermalinkResponse {
    @doc("Array of comments providing context around the target comment. Includes the complete parent chain (ancestors) from the target comment up to the top-level comment, the target comment itself, contextual sibling comments at the same level, and a limited number of replies to the target comment. The response is bounded to prevent performance issues while providing meaningful conversation context. Frontend developers should highlight the target comment and provide clear visual hierarchy showing the conversation flow. Backend developers should optimize queries to fetch the parent chain, then gather contextual siblings and replies efficiently.")
    comments: IncognitoPostComment[];

    @doc("The ID of the target comment that was requested. Frontend should use this to identify which comment to highlight, focus, or scroll to in the rendered conversation.")
    target_comment_id: string;

    @doc("Array of comment IDs representing the path from the top-level comment down to the target comment. This breadcrumb trail helps frontend developers build navigation UI and shows users where they are in the conversation hierarchy. The array is ordered from root to target comment.")
    breadcrumb_path: string[];
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

model UpvoteIncognitoPostRequest {
    @doc("The ID of the IncognitoPost to vote on")
    incognito_post_id: string;
}

model DownvoteIncognitoPostRequest {
    @doc("The ID of the IncognitoPost to vote on")
    incognito_post_id: string;
}

model UnvoteIncognitoPostRequest {
    @doc("The ID of the IncognitoPost to unvote on")
    incognito_post_id: string;
}

union IncognitoPostTimeFilter {
    Past24Hours: "past_24_hours",
    PastWeek: "past_week",
    PastMonth: "past_month",
    PastYear: "past_year",
}

model GetIncognitoPostsRequest {
    @doc("Tag ID to filter posts by. Exactly one tag must be specified to ensure focused content discovery.")
    tag_id: VTagID;

    @doc("Time period to filter posts by. Controls the age of posts to include in results.")
    time_filter?: IncognitoPostTimeFilter = IncognitoPostTimeFilter.Past24Hours;

    @doc("Number of posts to fetch per page")
    @minValue(1)
    @maxValue(100)
    limit: int32 = 25;

    @doc("Pagination key to fetch the next page of posts.")
    pagination_key?: string;
}

model IncognitoPostSummary {
    incognito_post_id: string;
    content: string;
    tags: VTag[];
    created_at: utcDateTime;
    upvotes_count: int32;
    downvotes_count: int32;
    score: int32;

    @doc("Whether the logged in user has upvoted this post")
    me_upvoted: boolean;

    @doc("Whether the logged in user has downvoted this post")
    me_downvoted: boolean;

    @doc("Whether the logged in user can upvote this post. Will be false if the user is blocked from upvoting this post or if the user has already voted for this post.")
    can_upvote: boolean;

    @doc("Whether the logged in user can downvote this post. Will be false if the user is blocked from downvoting this post or if the user has already voted for this post.")
    can_downvote: boolean;

    @doc("Total number of comments on this post (includes nested comments)")
    comments_count: int32;

    @doc("Whether this post was created by the currently logged in user")
    is_created_by_me: boolean;

    @doc("Whether this post has been deleted by its author")
    is_deleted: boolean;
}

model GetIncognitoPostsResponse {
    @doc("The list of incognito posts matching the specified criteria and filters. Posts are sorted by upvotes_count (descending), then by creation date (descending, newest first) for posts with the same upvote count. Deleted posts (where is_deleted=true) are excluded from results.")
    posts: IncognitoPostSummary[];

    @doc("Pagination key for fetching the next page of posts. Use this in subsequent requests. If this is empty, there are no more posts to fetch.")
    pagination_key: string;
}

model GetMyIncognitoPostsRequest {
    @doc("The posts are returned ordered by the newest created_at timestamp. If two posts are created at the same time, the one with the higher id is returned first. The pagination key is the id of the last post in the previous response. The item corresponding to the passed pagination_key will not be included in the response.")
    pagination_key?: string;

    @doc("Number of posts to fetch per page")
    @minValue(1)
    @maxValue(40)
    limit: int32 = 25;
}

model GetMyIncognitoPostsResponse {
    @doc("The list of incognito posts created by the logged-in user. Posts are sorted by creation date (descending, newest first). Includes both active and deleted posts.")
    posts: IncognitoPostSummary[];

    @doc("Pagination key for fetching the next page of posts. Use this in subsequent requests. If this is empty, there are no more posts to fetch.")
    pagination_key: string;
}

model MyIncognitoPostComment {
    comment_id: string;
    content: string;

    @doc("The ID of the parent comment. Null for top-level comments.")
    in_reply_to?: string;

    @doc("The creation timestamp of the comment")
    created_at: utcDateTime;

    upvotes_count: int32;
    downvotes_count: int32;
    score: int32;

    @doc("Whether the logged in user has upvoted this comment")
    me_upvoted: boolean;

    @doc("Whether the logged in user has downvoted this comment")
    me_downvoted: boolean;

    @doc("Whether the comment has been deleted by its author")
    is_deleted: boolean;

    @doc("The nesting depth of the comment (0 for top-level, 1 for first reply, etc.)")
    depth: int32;

    @doc("The ID of the incognito post this comment belongs to")
    incognito_post_id: string;

    @doc("Preview of the incognito post content (first 100 characters) to provide context")
    @maxLength(100)
    post_content_preview: string;

    @doc("Tags associated with the incognito post for additional context")
    post_tags: VTag[];
}

model GetMyIncognitoPostCommentsRequest {
    @doc("The comments are returned ordered by the newest created_at timestamp. If two comments are created at the same time, the one with the higher id is returned first. The pagination key is the id of the last comment in the previous response. The item corresponding to the passed pagination_key will not be included in the response.")
    pagination_key?: string;

    @doc("Number of comments to fetch per page")
    @minValue(1)
    @maxValue(40)
    limit: int32 = 25;
}

model GetMyIncognitoPostCommentsResponse {
    @doc("The list of incognito post comments created by the logged-in user across all incognito posts. Comments are sorted by creation date (descending, newest first). Includes both active and deleted comments.")
    comments: MyIncognitoPostComment[];

    @doc("Pagination key for fetching the next page of comments. Use this in subsequent requests. If this is empty, there are no more comments to fetch.")
    pagination_key: string;
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

@route("/hub/get-incognito-posts")
interface GetIncognitoPosts {
    @doc("GetIncognitoPosts gets a paginated list of incognito posts filtered by a single tag within a specified time period. Posts are sorted by upvotes_count (descending), then by creation date (descending, newest first) for ties. Deleted posts are automatically excluded from results. This API provides discovery of anonymous content within specific topic areas.")
    @tag("Posts")
    @get
    @useAuth(HubAuth)
    getIncognitoPosts(@body request: GetIncognitoPostsRequest): {
        @statusCode statusCode: 200;
        @body response: GetIncognitoPostsResponse;
    } | {
        @doc("Bad Request: Invalid request data (invalid tag ID, invalid limit values, invalid time_filter values)")
        @statusCode
        statusCode: 400;
    };
}

@route("/hub/get-my-incognito-posts")
interface GetMyIncognitoPosts {
    @doc("GetMyIncognitoPosts gets a paginated list of incognito posts created by the logged-in user. Posts are sorted by creation date (descending, newest first). This includes both active and deleted posts, allowing users to manage their own incognito content.")
    @tag("Posts")
    @get
    @useAuth(HubAuth)
    getMyIncognitoPosts(@body request: GetMyIncognitoPostsRequest): {
        @statusCode statusCode: 200;
        @body response: GetMyIncognitoPostsResponse;
    } | {
        @doc("Bad Request: Invalid request data (invalid limit values)")
        @statusCode
        statusCode: 400;
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
    @doc("GetIncognitoPostComments retrieves top-level comments with predictable direct reply previews. This API loads a specified number of top-level comments, and for each top-level comment, loads exactly the requested number of direct replies (immediate children only). This provides a Reddit-style experience where users can see the main discussion threads with a small preview of responses, then load more direct replies predictably using the GetCommentReplies API. The system ensures that 'Load X more replies' buttons load exactly X direct replies, never surprising users with nested content.")
    @tag("Posts")
    @get
    @useAuth(HubAuth)
    getIncognitoPostComments(@body request: GetIncognitoPostCommentsRequest): {
        @statusCode statusCode: 200;
        @body response: GetIncognitoPostCommentsResponse;
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

@route("/hub/get-comment-replies")
interface GetCommentReplies {
    @doc("GetCommentReplies loads replies for a specific comment with two modes: predictable direct reply loading and thread continuation. When direct_only=true (default), it loads exactly the requested number of immediate replies to the parent comment, making 'Load X more replies' buttons work predictably. When direct_only=false, it loads nested reply trees for 'continue this thread' functionality. This dual approach enables both Reddit-style incremental loading and deep thread exploration.")
    @tag("Posts")
    @get
    @useAuth(HubAuth)
    getCommentReplies(@body request: GetCommentRepliesRequest): {
        @statusCode statusCode: 200;
        @body response: GetCommentRepliesResponse;
    } | {
        @doc("Bad Request: Invalid request data")
        @statusCode
        statusCode: 400;
    } | {
        @doc("Not Found: The specified incognito post or parent comment does not exist")
        @statusCode
        statusCode: 404;
    };
}

@route("/hub/get-incognito-post-comment-permalink")
interface GetIncognitoPostCommentPermalink {
    @doc("GetIncognitoPostCommentPermalink retrieves a specific comment with bounded conversation context for permalink and direct linking functionality. This API is designed to provide meaningful context around any comment while maintaining performance bounds, unlike loading entire conversation trees which could overwhelm the system for viral comments. The response includes the complete parent chain (showing how the conversation led to this comment), the target comment itself, contextual sibling comments at the same level (showing alternative responses), and a preview of replies to the target comment (showing its impact). This bounded approach ensures consistent performance regardless of conversation size while providing sufficient context for users to understand and engage with the linked comment. Backend implementation should: 1) Validate the target comment exists and belongs to the specified post, 2) Traverse up the parent chain to collect all ancestors, 3) Query sibling comments at the target's level sorted by score, 4) Fetch top replies to the target comment, 5) Construct breadcrumb navigation path. Frontend implementation should: 1) Visually highlight the target comment distinctly, 2) Render the parent chain with clear conversation flow, 3) Show contextual siblings to provide conversation alternatives, 4) Display reply previews with 'Load more replies' if needed, 5) Implement breadcrumb navigation for conversation hierarchy.")
    @tag("Posts")
    @get
    @useAuth(HubAuth)
    getIncognitoPostCommentPermalink(
        @body request: GetIncognitoPostCommentPermalinkRequest,
    ): {
        @statusCode statusCode: 200;
        @body response: GetIncognitoPostCommentPermalinkResponse;
    } | {
        @doc("Bad Request: Invalid request data (malformed comment_id, invalid context counts)")
        @statusCode
        statusCode: 400;
    } | {
        @doc("Not Found: The specified incognito post or comment does not exist")
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
    ):
        | {
              @doc("Comment upvote registered. 200 will be returned even if the user has already upvoted.")
              @statusCode
              statusCode: 200;
          }
        | {
              @doc("Bad Request: Invalid request data")
              @statusCode
              statusCode: 400;
          }
        | {
              @doc("Not Found: The specified incognito post or comment does not exist or has been deleted")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("This code will be returned in any of the following scenarios: 1) User has already downvoted the comment 2) User is attempting to upvote their own comment")
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
    ):
        | {
              @doc("Comment downvote registered. 200 will be returned even if the user has already downvoted.")
              @statusCode
              statusCode: 200;
          }
        | {
              @doc("Bad Request: Invalid request data")
              @statusCode
              statusCode: 400;
          }
        | {
              @doc("Not Found: The specified incognito post or comment does not exist or has been deleted")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("This code will be returned in any of the following scenarios: 1) User has already upvoted the comment 2) User is attempting to downvote their own comment")
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
    ):
        | {
              @doc("Whatever upvote or downvote registered by the logged in user for the comment is removed. No error is returned even if the user has already unvoted or has never voted for the comment")
              @statusCode
              statusCode: 200;
          }
        | {
              @doc("Bad Request: Invalid request data")
              @statusCode
              statusCode: 400;
          }
        | {
              @doc("Not Found: The specified incognito post or comment does not exist or has been deleted")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("This code will be returned in any of the following scenarios: 1) User is attempting to unvote their own comment")
              @statusCode
              statusCode: 422;
          };
}

@route("/hub/upvote-incognito-post")
interface UpvoteIncognitoPost {
    @doc("UpvoteIncognitoPost upvotes an IncognitoPost")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    upvoteIncognitoPost(@body request: UpvoteIncognitoPostRequest):
        | {
              @doc("Post upvote registered. 200 will be returned even if the user has already upvoted.")
              @statusCode
              statusCode: 200;
          }
        | {
              @doc("Bad Request: Invalid request data")
              @statusCode
              statusCode: 400;
          }
        | {
              @doc("Not Found: The specified incognito post does not exist or has been deleted")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("This code will be returned in any of the following scenarios: 1) User has already downvoted the post 2) User is attempting to upvote their own post")
              @statusCode
              statusCode: 422;
          };
}

@route("/hub/downvote-incognito-post")
interface DownvoteIncognitoPost {
    @doc("DownvoteIncognitoPost downvotes an IncognitoPost")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    downvoteIncognitoPost(@body request: DownvoteIncognitoPostRequest):
        | {
              @doc("Post downvote registered. 200 will be returned even if the user has already downvoted.")
              @statusCode
              statusCode: 200;
          }
        | {
              @doc("Bad Request: Invalid request data")
              @statusCode
              statusCode: 400;
          }
        | {
              @doc("Not Found: The specified incognito post does not exist or has been deleted")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("This code will be returned in any of the following scenarios: 1) User has already upvoted the post 2) User is attempting to downvote their own post")
              @statusCode
              statusCode: 422;
          };
}

@route("/hub/unvote-incognito-post")
interface UnvoteIncognitoPost {
    @doc("UnvoteIncognitoPost removes the user's vote from an IncognitoPost")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    unvoteIncognitoPost(@body request: UnvoteIncognitoPostRequest):
        | {
              @doc("Whatever upvote or downvote registered by the logged in user for the post is removed. No error is returned even if the user has already unvoted or has never voted for the post")
              @statusCode
              statusCode: 200;
          }
        | {
              @doc("Bad Request: Invalid request data")
              @statusCode
              statusCode: 400;
          }
        | {
              @doc("Not Found: The specified incognito post does not exist or has been deleted")
              @statusCode
              statusCode: 404;
          }
        | {
              @doc("This code will be returned in any of the following scenarios: 1) User is attempting to unvote their own post")
              @statusCode
              statusCode: 422;
          };
}

@route("/hub/get-my-incognito-post-comments")
interface GetMyIncognitoPostComments {
    @doc("GetMyIncognitoPostComments gets a paginated list of incognito post comments created by the logged-in user across all incognito posts. Comments are sorted by creation date (descending, newest first). This includes both active and deleted comments, allowing users to manage their own anonymous comment activity.")
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    getMyIncognitoPostComments(
        @body request: GetMyIncognitoPostCommentsRequest,
    ): {
        @statusCode statusCode: 200;
        @body response: GetMyIncognitoPostCommentsResponse;
    } | {
        @doc("Bad Request: Invalid request data (invalid limit values)")
        @statusCode
        statusCode: 400;
    };
}
