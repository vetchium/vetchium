// loadtest/hub_scenario.js
import {
  randomIntBetween,
  randomItem,
} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { check, group, sleep } from "k6";
import http from "k6/http";
import { Trend } from "k6/metrics";

// Helper function to generate a random string
function randomString(length) {
  const charset =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  let result = "";
  for (let i = 0; i < length; i++) {
    result += charset.charAt(Math.floor(Math.random() * charset.length));
  }
  return result;
}

// --- Configuration ---
const API_BASE_URL = __ENV.API_BASE_URL || "http://localhost:8080"; // Your API gateway/service URL
const MAILPIT_URL = __ENV.MAILPIT_URL || "http://localhost:8025"; // Mailpit API URL
const NUM_USERS = parseInt(__ENV.NUM_USERS || "100");
const PASSWORD = "NewPassword123$";
const TEST_DURATION_SECONDS = parseInt(__ENV.TEST_DURATION || "600"); // 10 minutes default

// --- Metrics ---
const followTrend = new Trend("hub_follow_user_duration", true);
const unfollowTrend = new Trend("hub_unfollow_user_duration", true);
const createPostTrend = new Trend("hub_create_post_duration", true);
const timelineReadTrend = new Trend("hub_timeline_read_duration", true);
const postDetailsTrend = new Trend("hub_post_details_duration", true);
const upvoteTrend = new Trend("hub_upvote_duration", true);
const downvoteTrend = new Trend("hub_downvote_duration", true);
const unvoteTrend = new Trend("hub_unvote_duration", true);
const followStatusTrend = new Trend("hub_follow_status_duration", true);

// Default number of users if not provided by environment variable
const DEFAULT_NUM_USERS = 100;

// --- Constants for authentication ---
const MAX_LOGIN_ATTEMPTS = 3;
const MAX_TFA_FETCH_ATTEMPTS = 5;

// --- TFA Code Fetch (Uses Email) ---
function fetchTFACodeForUser(email) {
  let attempts = 0;
  while (attempts < MAX_TFA_FETCH_ATTEMPTS) {
    attempts++;
    console.debug(
      `Attempt ${attempts}/${MAX_TFA_FETCH_ATTEMPTS}: Fetching emails for ${email} from Mailpit...`
    );
    const mailpitRes = http.get(`${MAILPIT_URL}/api/v1/messages?limit=10`);

    if (mailpitRes.status !== 200) {
      console.warn(
        `Mailpit API returned status ${mailpitRes.status}. Waiting before retrying...`
      );
      sleep(5);
      continue;
    }

    try {
      const messages = mailpitRes.json("messages");
      if (!messages || messages.length === 0) {
        console.debug(
          `No messages found in Mailpit yet for ${email}. Waiting...`
        );
        sleep(10);
        continue;
      }

      // Find the most recent email for the target user
      const targetEmail = messages.find(
        (msg) => msg.To && msg.To.length > 0 && msg.To[0].Address === email
      );

      if (targetEmail) {
        // Extract TFA code from the email body (assumes simple text format)
        const codeMatch = targetEmail.Text.match(
          /Your verification code is: (\d{6})/
        );
        if (codeMatch && codeMatch[1]) {
          console.debug(`TFA code found for ${email}: ${codeMatch[1]}`);
          return codeMatch[1];
        }
      }

      // If email found but no code, or email not found yet
      console.debug(
        `Mail for ${email} not found or code extraction failed (attempt ${attempts}). Waiting...`
      );
      sleep(10);
    } catch (e) {
      console.error(`Error processing Mailpit response: ${e}. Waiting...`);
      sleep(10);
    }
  }

  console.error(
    `Failed to fetch TFA code for ${email} after ${MAX_TFA_FETCH_ATTEMPTS} attempts.`
  );
  return null;
}

// --- Authentication Function (Accepts user object with email/handle) ---
function loginAndAuthenticateUser(user) {
  let loginAttempts = 0;
  let tfaToken = null;
  let handle = user.handle; // Store handle from the input user object

  while (loginAttempts < MAX_LOGIN_ATTEMPTS) {
    loginAttempts++;
    console.debug(
      `Attempt ${loginAttempts}/${MAX_LOGIN_ATTEMPTS}: Logging in ${user.email} (handle: ${user.handle})...`
    );

    // Use Email for login request payload according to LoginRequest model
    const loginPayload = JSON.stringify({
      email: user.email,
      password: PASSWORD,
    });

    const loginRes = http.post(`${API_BASE_URL}/hub/login`, loginPayload, {
      headers: { "Content-Type": "application/json" },
      tags: { name: "HubLoginAPI" },
    });

    // Check if login was successful (status 200)
    if (loginRes.status === 200) {
      // Parse response body to check for TFA requirement or direct token
      try {
        const responseBody = JSON.parse(loginRes.body);

        // Check if TFA is required
        if (responseBody.requires_tfa === true && responseBody.tfa_token) {
          tfaToken = responseBody.tfa_token;
          // Store handle if provided in response
          if (responseBody.handle) {
            handle = responseBody.handle;
          }
          console.debug(
            `Login step 1 successful for ${user.email}. TFA required. Handle: ${handle}`
          );
          break; // Proceed to TFA verification
        }
        // Check if direct login (token provided)
        else if (responseBody.token) {
          const authToken = responseBody.token;
          // Update handle if provided in response
          if (responseBody.handle) {
            handle = responseBody.handle;
          }
          console.debug(
            `Direct login successful for ${user.email}. Handle: ${handle}`
          );
          return { authToken, userHandle: handle }; // Return token and handle
        } else {
          console.warn(
            `Login response for ${user.email} has status 200 but missing expected fields. Body: ${loginRes.body}`
          );
          sleep(randomIntBetween(1, 3));
        }
      } catch (e) {
        console.error(
          `Error parsing login response for ${user.email}: ${e}. Body: ${loginRes.body}`
        );
        sleep(randomIntBetween(1, 3));
      }
    } else if (loginRes.status === 422) {
      console.warn(
        `Login attempt ${loginAttempts} failed for ${user.email} - account not in valid state (422). Body: ${loginRes.body}`
      );
      sleep(randomIntBetween(1, 3));
    } else {
      console.warn(
        `Login attempt ${loginAttempts} failed for ${user.email}. Status: ${loginRes.status}, Body: ${loginRes.body}. Retrying...`
      );
      sleep(randomIntBetween(1, 3));
    }
  }

  if (!tfaToken) {
    console.error(
      `Login failed for ${user.email} after ${MAX_LOGIN_ATTEMPTS} attempts.`
    );
    return null; // Login failure
  }

  // --- TFA Verification Step ---
  const tfaCode = fetchTFACodeForUser(user.email); // Use email to fetch code
  if (!tfaCode) {
    console.error(
      `TFA code retrieval failed for ${user.email}. Cannot complete login.`
    );
    return null; // TFA failure
  }

  // Payload according to HubTFARequest model from hubusers.tsp
  const tfaPayload = JSON.stringify({
    tfa_token: tfaToken,
    tfa_code: tfaCode,
    remember_me: false,
  });

  // Endpoint according to hubusers.tsp
  const tfaVerifyRes = http.post(`${API_BASE_URL}/hub/tfa`, tfaPayload, {
    headers: { "Content-Type": "application/json" },
    tags: { name: "HubVerifyTFA_API" },
  });

  // Process TFA verification response
  if (tfaVerifyRes.status === 200) {
    try {
      const tfaResponseBody = JSON.parse(tfaVerifyRes.body);

      // Extract session_token according to HubTFAResponse model
      if (tfaResponseBody.session_token) {
        const authToken = tfaResponseBody.session_token;
        console.debug(`TFA verification successful for ${user.email}.`);
        return { authToken, userHandle: handle }; // Return token and the confirmed handle
      } else {
        console.error(
          `TFA verification response missing session_token for ${user.email}. Body: ${tfaVerifyRes.body}`
        );
        return null;
      }
    } catch (e) {
      console.error(
        `Error parsing TFA verification response for ${user.email}: ${e}. Body: ${tfaVerifyRes.body}`
      );
      return null;
    }
  } else {
    console.error(
      `TFA verification failed for ${user.email}. Status: ${tfaVerifyRes.status}, Body: ${tfaVerifyRes.body}`
    );
    return null; // TFA verification failure
  }
}

// --- k6 Setup Function ---
export function setup() {
  console.log("=== Running Setup Phase ===");
  console.log(`Generating and authenticating ${NUM_USERS} users...`);

  const allUsers = [];
  for (let i = 1; i <= NUM_USERS; i++) {
    // Generate user details based on seed_users.sh pattern
    allUsers.push({
      handle: `hubuser${i}`,
      email: `hubuser${i}@example.com`, // Use @example.com as in script
    });
  }

  if (allUsers.length === 0) {
    throw new Error(
      "No users generated. Check NUM_USERS environment variable."
    );
  }

  const authenticatedUsers = [];
  const handles = []; // Collect handles separately for socialActivity

  console.log(`Attempting to authenticate ${allUsers.length} users...`);

  for (let i = 0; i < allUsers.length; i++) {
    const user = allUsers[i];
    console.log(
      `Authenticating user ${i + 1}/${allUsers.length}: ${user.email}`
    );
    const authResult = loginAndAuthenticateUser(user); // Pass the whole user object
    if (authResult && authResult.authToken && authResult.userHandle) {
      authenticatedUsers.push({
        // Store only necessary info for VUs
        authToken: authResult.authToken,
        userHandle: authResult.userHandle,
      });
      handles.push(authResult.userHandle); // Store the handle
      console.log(
        `Authentication successful for: ${user.email} (Handle: ${authResult.userHandle})`
      );
    } else {
      // Fail fast if any user cannot be authenticated
      throw new Error(
        `Setup failed: Could not authenticate user ${user.email}. Halting test.`
      );
    }
    // Optional: Add a small delay between authentications if needed
    // sleep(0.5);
  }

  console.log(
    `=== Setup Phase Complete: ${authenticatedUsers.length} users authenticated ===`
  );
  // Pass authenticated user data and the list of all handles
  return { authenticatedUsers: authenticatedUsers, allUserHandles: handles };
}

// --- Main Test Logic (Accepts setup data) ---
export default function (data) {
  // Initialize VU state for tracking social interactions
  let vuState = {
    // Timeline pagination
    timelineCursor: null,

    // Posts tracking
    fetchedPostIds: [], // Posts seen in timeline
    createdPostIds: [], // Posts created by this user

    // Voting tracking
    upvotedPostIds: [], // Posts this user has upvoted
    downvotedPostIds: [], // Posts this user has downvoted

    // Following tracking
    followingUsers: [], // Users this VU is following

    // Test duration tracking
    startTime: new Date().getTime(),
    endTime: new Date().getTime() + TEST_DURATION_SECONDS * 1000,
  };

  // Get a random authenticated user from the setup data
  const userIndex = __VU % data.authenticatedUsers.length;
  const currentUser = data.authenticatedUsers[userIndex];

  if (!currentUser || !currentUser.authToken) {
    console.error(
      `VU ${__VU}: No authenticated user available. Skipping test execution.`
    );
    return;
  }

  // Pass the specific token, handle, and the list of all handles to socialActivity
  socialActivity(
    currentUser.authToken,
    currentUser.userHandle,
    data.allUserHandles,
    vuState
  );
}

// --- Social Activity Function ---
export function socialActivity(authToken, userHandle, allUserHandles, vuState) {
  // Check if test duration has been reached
  if (new Date().getTime() > vuState.endTime) {
    console.debug(
      `VU ${__VU} (${userHandle}): Test duration reached, ending test.`
    );
    return;
  }

  // Prepare auth params using the passed token
  const authParams = {
    headers: {
      Authorization: `Bearer ${authToken}`,
      "Content-Type": "application/json",
    },
  };

  // Select a random user handle to interact with (exclude self)
  let handlesToInteractWith = allUserHandles.filter((h) => h !== userHandle);
  if (handlesToInteractWith.length === 0 && allUserHandles.length > 0) {
    // Fallback if only one user exists or filtering failed
    handlesToInteractWith = allUserHandles;
  } else if (handlesToInteractWith.length === 0) {
    console.warn(
      `VU ${__VU} (${userHandle}): No other handles available to interact with.`
    );
  }

  console.debug(
    `VU ${__VU} (${userHandle}): Starting social activity iteration.`
  );

  group("Social Interaction Loop", function () {
    // --- Action Selection ---
    // Weighted action selection based on realistic user behavior
    let availableActions = [
      "follow",
      "createPost",
      "createPost", // Higher weight for creating posts
      "readTimeline",
      "readTimeline", // Higher weight for reading timeline
      "readTimeline",
      "vote",
      "vote", // Higher weight for voting
      "vote",
      "getPostDetails",
      "getFollowStatus",
    ];

    // Only add unfollow if the user is following someone
    if (vuState.followingUsers.length > 0) {
      availableActions.push("unfollow");
    }

    // Only add unvote if the user has voted on something
    if (
      vuState.upvotedPostIds.length > 0 ||
      vuState.downvotedPostIds.length > 0
    ) {
      availableActions.push("unvote");
    }

    const action = randomItem(availableActions);

    switch (action) {
      case "follow":
        if (!handlesToInteractWith.length) break;
        const userToFollowHandle = randomItem(handlesToInteractWith);
        const followPayload = JSON.stringify({ handle: userToFollowHandle });

        console.debug(
          `VU ${__VU} (${userHandle}): Attempting to follow ${userToFollowHandle}`
        );

        const followRes = http.post(
          `${API_BASE_URL}/hub/follow-user`,
          followPayload,
          { ...authParams, tags: { name: "HubFollowAPI" } }
        );

        followTrend.add(followRes.timings.duration);

        // Check for success
        check(followRes, {
          "Follow request successful or expected client error (status 200/4xx)":
            (r) => r.status === 200 || (r.status >= 400 && r.status < 500),
        });

        // If successful, add to following list
        if (followRes.status === 200) {
          if (!vuState.followingUsers.includes(userToFollowHandle)) {
            vuState.followingUsers.push(userToFollowHandle);
          }
        }

        // Log unexpected errors
        if (
          followRes.status !== 200 &&
          !(followRes.status >= 400 && followRes.status < 500)
        ) {
          console.error(
            `VU ${__VU} (${userHandle}): Follow API Unexpected Error! Status: ${followRes.status}, Body: ${followRes.body}`
          );
        }

        sleep(randomIntBetween(1, 3));
        break;

      case "unfollow":
        if (vuState.followingUsers.length === 0) break;

        // Pick a random user to unfollow from those we're following
        const userToUnfollowHandle = randomItem(vuState.followingUsers);
        const unfollowPayload = JSON.stringify({
          handle: userToUnfollowHandle,
        });

        console.debug(
          `VU ${__VU} (${userHandle}): Attempting to unfollow ${userToUnfollowHandle}`
        );

        const unfollowRes = http.post(
          `${API_BASE_URL}/hub/unfollow-user`,
          unfollowPayload,
          { ...authParams, tags: { name: "HubUnfollowAPI" } }
        );

        unfollowTrend.add(unfollowRes.timings.duration);

        // Check for success
        check(unfollowRes, {
          "Unfollow request successful or expected client error (status 200/4xx)":
            (r) => r.status === 200 || (r.status >= 400 && r.status < 500),
        });

        // If successful, remove from following list
        if (unfollowRes.status === 200) {
          vuState.followingUsers = vuState.followingUsers.filter(
            (h) => h !== userToUnfollowHandle
          );
        }

        // Log unexpected errors
        if (
          unfollowRes.status !== 200 &&
          !(unfollowRes.status >= 400 && unfollowRes.status < 500)
        ) {
          console.error(
            `VU ${__VU} (${userHandle}): Unfollow API Unexpected Error! Status: ${unfollowRes.status}, Body: ${unfollowRes.body}`
          );
        }

        sleep(randomIntBetween(1, 3));
        break;

      case "getFollowStatus":
        if (!handlesToInteractWith.length) break;

        // Pick a random user to check follow status
        const userToCheckHandle = randomItem(handlesToInteractWith);
        const followStatusPayload = JSON.stringify({
          handle: userToCheckHandle,
        });

        console.debug(
          `VU ${__VU} (${userHandle}): Checking follow status with ${userToCheckHandle}`
        );

        const followStatusRes = http.post(
          `${API_BASE_URL}/hub/get-follow-status`,
          followStatusPayload,
          { ...authParams, tags: { name: "HubFollowStatusAPI" } }
        );

        followStatusTrend.add(followStatusRes.timings.duration);

        // Check for success
        check(followStatusRes, {
          "Get follow status successful (status 200)": (r) => r.status === 200,
        });

        // Update following status if response is valid
        if (followStatusRes.status === 200) {
          try {
            const responseBody = JSON.parse(followStatusRes.body);
            if (responseBody.is_following) {
              if (!vuState.followingUsers.includes(userToCheckHandle)) {
                vuState.followingUsers.push(userToCheckHandle);
              }
            } else {
              vuState.followingUsers = vuState.followingUsers.filter(
                (h) => h !== userToCheckHandle
              );
            }
          } catch (e) {
            console.error(
              `VU ${__VU} (${userHandle}): Error parsing follow status response: ${e}`
            );
          }
        }

        // Log errors
        if (followStatusRes.status !== 200) {
          console.error(
            `VU ${__VU} (${userHandle}): Get Follow Status API Error! Status: ${followStatusRes.status}, Body: ${followStatusRes.body}`
          );
        }

        sleep(randomIntBetween(1, 3));
        break;

      case "createPost":
        console.debug(`VU ${__VU} (${userHandle}): Attempting to create post`);

        const postContent = `Post content from VU ${__VU} at ${new Date().toISOString()}: ${randomString(
          50
        )}`;

        const numTags = randomIntBetween(0, 3);
        const postTags = [];
        for (let i = 0; i < numTags; i++) {
          postTags.push(`tag_${randomString(5)}`);
        }

        // Use new_tags as per TypeSpec
        const postPayload = JSON.stringify({
          content: postContent,
          new_tags: postTags,
        });

        console.debug(
          `VU ${__VU} (${userHandle}): Attempting to create post. Tags: ${postTags.join(
            ", "
          )}`
        );

        const postRes = http.post(`${API_BASE_URL}/hub/add-post`, postPayload, {
          ...authParams,
          tags: { name: "HubCreatePostAPI" },
        });

        createPostTrend.add(postRes.timings.duration);

        // Check for success and post_id
        check(postRes, {
          "Create post successful (status 200)": (r) => r.status === 200,
          "Create post response has post_id": (r) =>
            r.body &&
            r.json("post_id") !== null &&
            r.json("post_id") !== undefined,
        });

        // Store created post ID if successful
        if (postRes.status === 200) {
          try {
            const postId = postRes.json("post_id");
            if (postId) {
              vuState.createdPostIds.push(postId);
              // Also add to fetched posts so we can interact with it
              if (!vuState.fetchedPostIds.includes(postId)) {
                vuState.fetchedPostIds.push(postId);
              }
            }
          } catch (e) {
            console.error(
              `VU ${__VU} (${userHandle}): Error extracting post_id: ${e}`
            );
          }
        }

        // Log errors
        if (postRes.status !== 200) {
          console.error(
            `VU ${__VU} (${userHandle}): Create Post API Error! Status: ${postRes.status}, Body: ${postRes.body}`
          );
        } else if (!postRes.json("post_id")) {
          console.error(
            `VU ${__VU} (${userHandle}): Create Post API Error! Status 200 but missing post_id. Body: ${postRes.body}`
          );
        }

        sleep(randomIntBetween(2, 5));
        break;

      case "readTimeline":
        let timelineUrl = `${API_BASE_URL}/hub/get-my-home-timeline`;

        // Prepare body for POST request
        let timelinePayload = {};
        if (vuState.timelineCursor) {
          timelinePayload = { pagination_key: vuState.timelineCursor };
          console.debug(
            `VU ${__VU} (${userHandle}): Reading timeline with pagination_key: ${vuState.timelineCursor}`
          );
        } else {
          console.debug(
            `VU ${__VU} (${userHandle}): Reading timeline (first page). Sending empty body {}`
          );
        }

        // Send request
        const timelineRes = http.post(
          timelineUrl,
          JSON.stringify(timelinePayload),
          {
            ...authParams,
            tags: { name: "HubTimelineReadAPI" },
          }
        );

        timelineReadTrend.add(timelineRes.timings.duration);

        // Check for success
        check(timelineRes, {
          "Read timeline successful (status 200)": (r) => r.status === 200,
        });

        // Process timeline data if successful
        if (timelineRes.status === 200) {
          try {
            const timelineData = JSON.parse(timelineRes.body);

            // Update pagination cursor
            if (timelineData.pagination_key) {
              vuState.timelineCursor = timelineData.pagination_key;
            }

            // Store post IDs for future interactions
            if (timelineData.posts && timelineData.posts.length > 0) {
              timelineData.posts.forEach((post) => {
                if (post.id && !vuState.fetchedPostIds.includes(post.id)) {
                  vuState.fetchedPostIds.push(post.id);
                }
              });

              console.debug(
                `VU ${__VU} (${userHandle}): Fetched ${timelineData.posts.length} posts. Total known posts: ${vuState.fetchedPostIds.length}`
              );
            } else {
              console.debug(
                `VU ${__VU} (${userHandle}): Timeline empty or no new posts.`
              );
            }
          } catch (e) {
            console.error(
              `VU ${__VU} (${userHandle}): Error processing timeline response: ${e}`
            );
          }
        } else {
          console.error(
            `VU ${__VU} (${userHandle}): Read Timeline API Error! Status: ${timelineRes.status}, Body: ${timelineRes.body}`
          );
        }

        sleep(randomIntBetween(3, 7));
        break;

      case "getPostDetails":
        if (vuState.fetchedPostIds.length === 0) {
          console.debug(
            `VU ${__VU} (${userHandle}): No posts fetched yet, skipping get post details.`
          );
          break;
        }

        const postIdToView = randomItem(vuState.fetchedPostIds);
        const postDetailsPayload = JSON.stringify({ post_id: postIdToView });

        console.debug(
          `VU ${__VU} (${userHandle}): Getting details for post ${postIdToView}`
        );

        const postDetailsRes = http.post(
          `${API_BASE_URL}/hub/get-post-details`,
          postDetailsPayload,
          {
            ...authParams,
            tags: { name: "HubPostDetailsAPI" },
          }
        );

        postDetailsTrend.add(postDetailsRes.timings.duration);

        // Check for success
        check(postDetailsRes, {
          "Get post details successful (status 200)": (r) => r.status === 200,
        });

        // Log errors
        if (postDetailsRes.status !== 200) {
          console.error(
            `VU ${__VU} (${userHandle}): Get Post Details API Error! PostID: ${postIdToView}, Status: ${postDetailsRes.status}, Body: ${postDetailsRes.body}`
          );
        }

        sleep(randomIntBetween(2, 5));
        break;

      case "vote":
        if (vuState.fetchedPostIds.length === 0) {
          console.debug(
            `VU ${__VU} (${userHandle}): No posts fetched yet, skipping vote.`
          );
          break;
        }

        // Filter out posts that are created by this user (can't vote on own posts)
        // and posts that have already been voted on
        const votablePosts = vuState.fetchedPostIds.filter(
          (id) =>
            !vuState.createdPostIds.includes(id) &&
            !vuState.upvotedPostIds.includes(id) &&
            !vuState.downvotedPostIds.includes(id)
        );

        if (votablePosts.length === 0) {
          console.debug(
            `VU ${__VU} (${userHandle}): No votable posts available, skipping vote.`
          );
          break;
        }

        const postIdToVote = randomItem(votablePosts);
        const voteType = randomItem(["upvote", "downvote"]);

        console.debug(
          `VU ${__VU} (${userHandle}): Attempting to ${voteType} post ${postIdToVote}`
        );

        const votePayload = JSON.stringify({ post_id: postIdToVote });

        let voteUrl;
        let voteTrend;
        let voteTag;
        if (voteType === "upvote") {
          voteUrl = `${API_BASE_URL}/hub/upvote-user-post`;
          voteTrend = upvoteTrend;
          voteTag = "HubUpvoteAPI";
        } else {
          voteUrl = `${API_BASE_URL}/hub/downvote-user-post`;
          voteTrend = downvoteTrend;
          voteTag = "HubDownvoteAPI";
        }

        const voteRes = http.post(voteUrl, votePayload, {
          ...authParams,
          tags: { name: voteTag },
        });

        voteTrend.add(voteRes.timings.duration);

        // Check for success or expected error
        check(voteRes, {
          [`${voteType} request successful or expected error (status 200/422)`]:
            (r) => r.status === 200 || r.status === 422,
        });

        // If successful, track the voted post
        if (voteRes.status === 200) {
          if (voteType === "upvote") {
            vuState.upvotedPostIds.push(postIdToVote);
          } else {
            vuState.downvotedPostIds.push(postIdToVote);
          }
        }

        // Log unexpected errors
        if (voteRes.status !== 200 && voteRes.status !== 422) {
          console.error(
            `VU ${__VU} (${userHandle}): ${voteType} API Unexpected Error! PostID: ${postIdToVote}, Status: ${voteRes.status}, Body: ${voteRes.body}`
          );
        }

        sleep(randomIntBetween(1, 4));
        break;

      case "unvote":
        // Combine upvoted and downvoted posts to pick from
        const votedPosts = [
          ...vuState.upvotedPostIds,
          ...vuState.downvotedPostIds,
        ];

        if (votedPosts.length === 0) {
          console.debug(
            `VU ${__VU} (${userHandle}): No voted posts available, skipping unvote.`
          );
          break;
        }

        const postIdToUnvote = randomItem(votedPosts);
        const unvotePayload = JSON.stringify({ post_id: postIdToUnvote });

        console.debug(
          `VU ${__VU} (${userHandle}): Attempting to unvote post ${postIdToUnvote}`
        );

        const unvoteRes = http.post(
          `${API_BASE_URL}/hub/unvote-user-post`,
          unvotePayload,
          {
            ...authParams,
            tags: { name: "HubUnvoteAPI" },
          }
        );

        unvoteTrend.add(unvoteRes.timings.duration);

        // Check for success or expected error
        check(unvoteRes, {
          "Unvote request successful or expected error (status 200/422)": (r) =>
            r.status === 200 || r.status === 422,
        });

        // If successful, remove from voted lists
        if (unvoteRes.status === 200) {
          vuState.upvotedPostIds = vuState.upvotedPostIds.filter(
            (id) => id !== postIdToUnvote
          );
          vuState.downvotedPostIds = vuState.downvotedPostIds.filter(
            (id) => id !== postIdToUnvote
          );
        }

        // Log unexpected errors
        if (unvoteRes.status !== 200 && unvoteRes.status !== 422) {
          console.error(
            `VU ${__VU} (${userHandle}): Unvote API Unexpected Error! PostID: ${postIdToUnvote}, Status: ${unvoteRes.status}, Body: ${unvoteRes.body}`
          );
        }

        sleep(randomIntBetween(1, 4));
        break;
    }

    // Think time between actions in the main loop
    sleep(randomIntBetween(2, 6));
  });

  // Recursively call socialActivity to continue the test until duration is reached
  socialActivity(authToken, userHandle, allUserHandles, vuState);
}
// --- k6 Test Configuration ---
export const options = {
  scenarios: {
    social_interactions: {
      executor: "shared-iterations",
      vus: 10, // Adjust based on your needs
      iterations: 100, // This is per VU
      maxDuration: `${TEST_DURATION_SECONDS}s`,
    },
  },
  thresholds: {
    hub_follow_user_duration: ["p(95)<1000"], // 95% of requests should be under 1s
    hub_unfollow_user_duration: ["p(95)<1000"],
    hub_create_post_duration: ["p(95)<1000"],
    hub_timeline_read_duration: ["p(95)<1000"],
    hub_post_details_duration: ["p(95)<1000"],
    hub_upvote_duration: ["p(95)<1000"],
    hub_downvote_duration: ["p(95)<1000"],
    hub_unvote_duration: ["p(95)<1000"],
    hub_follow_status_duration: ["p(95)<1000"],
  },
};
