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

    @doc("The posts are returned ordered by the newest of created_at timestamp. If two posts are created at the same time, the one with the higher id is returned first. The pagination key is the id of the last post in the previous response. The item corresponding to the passed pagination_key will not be included in the response. In future, once we support updating of existing posts, the pagination will be based on the updated_at timestamp.")
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

model FollowUserRequest {
    @doc("The handle of the user to follow")
    handle: Handle;
}

model UnfollowUserRequest {
    @doc("The handle of the user to unfollow")
    handle: Handle;
}

model GetFollowStatusRequest {
    @doc("The handle of the user to get the follow status of")
    handle: Handle;
}

model FollowStatus {
    @doc("True if the logged in user is following the target user. When checking one's own status, this will always be true.")
    is_following: boolean;

    @doc("True if the logged in user is blocked by the target user. When checking one's own status, this will always be false.")
    is_blocked: boolean;

    @doc("True if the logged in user can follow the target user. Some user accounts may reach a state where they cannot be followed anymore (like tombstoned users). When checking one's own status, this will always be false. When already following a user, this will also be false.")
    can_follow: boolean;
}

model GetMyHomeTimelineRequest {
    @doc("The posts are returned ordered by the newest of update_at timestamp. If two posts are updated at the same time, the one with the higher id is returned first. The pagination key is the id of the last post in the previous response. The item corresponding to the passed pagination_key will not be included in the response.")
    pagination_key?: string;

    @doc("If nothing is passed, 25 Posts are returned")
    @minValue(1)
    @maxValue(40)
    limit?: integer;
}

model MyHomeTimeline {
    posts: Post[];
    pagination_key: string;
}

model GetPostDetailsRequest {
    @doc("The id of the post to get the details of")
    post_id: string;
}

model UpvoteUserPostRequest {
    post_id: string;
}

model DownvoteUserPostRequest {
    post_id: string;
}

model UnvoteUserPostRequest {
    post_id: string;
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

@route("/hub/follow-user")
interface FollowUser {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    followUser(@body request: FollowUserRequest): {
        @doc("The user was followed successfully. If the user was already followed, the request is ignored. If a user attempts to follow themselves, a 200 status is returned without creating any database records. If the target user is not active (e.g., deleted), a 404 status is returned.")
        @statusCode
        statusCode: 200;
    };
}

@route("/hub/unfollow-user")
interface UnfollowUser {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    unfollowUser(@body request: UnfollowUserRequest): {
        @doc("The user was unfollowed successfully. If the user was not followed already, the request is ignored. If a user attempts to unfollow themselves, a 404 status is returned. If the target user is not active (e.g., deleted), a 404 status is returned.")
        @statusCode
        statusCode: 200;
    };
}

@route("/hub/get-follow-status")
interface GetFollowStatus {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    getFollowStatus(@body request: GetFollowStatusRequest): {
        @doc("Returns the follow status between the logged-in user and the requested handle. Special cases: 1) When checking one's own status: isFollowing=true, isBlocked=false, canFollow=false. 2) When already following a user: canFollow=false. If the target user is not active (e.g., deleted), a 404 status is returned.")
        @statusCode
        statusCode: 200;

        @body response: FollowStatus;
    };
}

@route("/hub/get-my-home-timeline")
interface GetMyHomeTimeline {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    getMyHomeTimeline(@body request: GetMyHomeTimelineRequest): {
        @statusCode statusCode: 200;
        @body response: MyHomeTimeline;
    } | {
        @doc("The request cannot be processed. This will most likely happen if the user passed an invalid pagination key. A pagination key could become invalid if the post it refers to is deleted.")
        @statusCode
        statusCode: 422;
    };
}

@route("/hub/get-post-details")
interface GetPostDetails {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    getPostDetails(@body request: GetPostDetailsRequest): {
        @statusCode statusCode: 200;
        @body response: Post;
    };
}

@route("/hub/upvote-user-post")
interface UpvoteUserPost {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    upvoteUserPost(@body request: UpvoteUserPostRequest): {
        @doc("User upvote registered. 200 will be returned even if the user has already upvoted.")
        @statusCode
        statusCode: 200;
    } | {
        @doc("This code will be returned in any of the following scenarios. 1) User has already downvoted the post 2) User is blocked from upvoting this post 3) The Post ID is invalid")
        @statusCode
        statuscode: 422;
    };
}

@route("/hub/downvote-user-post")
interface DownvoteUserPost {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    downvoteUserPost(@body request: DownvoteUserPostRequest): {
        @doc("User downvote registered. 200 will be returned even if the user has already downvoted.")
        @statusCode
        statusCode: 200;
    } | {
        @doc("This code will be returned in any of the following scenarios. 1) User has already upvoted the post 2) User is blocked from downvoting this post 3) The Post ID is invalid")
        @statusCode
        statuscode: 422;
    };
}

@route("/hub/unvote-user-post")
interface UnvoteUserPost {
    @tag("Posts")
    @post
    @useAuth(HubAuth)
    unvoteUserPost(@body request: UnvoteUserPostRequest): {
        @doc("Whatever Upvote or Downvote registered by the logged in user for the Post is removed. No error is returned even if the user has already unvoted or has never voted for the post")
        @statusCode
        statusCode: 200;
    } | {
        @doc("Post ID is invalid")
        @statusCode
        statuscode: 422;
    };
}
