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
const API_BASE_URL = __ENV.API_BASE_URL;
const MAILPIT_URL = __ENV.MAILPIT_URL;
const NUM_USERS = parseInt(__ENV.NUM_USERS || "100");
const SETUP_PARALLELISM = parseInt(__ENV.SETUP_PARALLELISM || "10"); // Number of users to authenticate in parallel
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

// --- TFA Synchronization ---
let tfaLock = false;
const tfaQueue = [];

// Helper function to acquire TFA lock
function acquireTfaLock() {
  return new Promise((resolve) => {
    if (!tfaLock) {
      tfaLock = true;
      resolve();
    } else {
      tfaQueue.push(resolve);
    }
  });
}

// Helper function to release TFA lock
function releaseTfaLock() {
  if (tfaQueue.length > 0) {
    const nextResolve = tfaQueue.shift();
    nextResolve();
  } else {
    tfaLock = false;
  }
}

// --- TFA Code Fetch (Uses Email) ---
async function fetchTFACodeForUser(email) {
  try {
    // Acquire lock before fetching TFA code
    await acquireTfaLock();

    let attempts = 0;
    let messageId = null;

    // Step 1: Search for the TFA email using the search API
    while (attempts < MAX_TFA_FETCH_ATTEMPTS) {
      attempts++;

      // Build the search URL with query parameters exactly as in the Go code
      const searchQuery = `to:${email} subject:Vetchium Two Factor Authentication`;
      const searchUrl = `${MAILPIT_URL}/api/v1/search?query=${encodeURIComponent(
        searchQuery
      )}`;

      console.debug(
        `Attempt ${attempts}/${MAX_TFA_FETCH_ATTEMPTS}: Searching Mailpit at ${searchUrl}`
      );

      const searchRes = http.get(searchUrl);

      if (searchRes.status !== 200) {
        console.warn(
          `Mailpit search API returned status ${searchRes.status}. Waiting before retrying...`
        );
        sleep(2); // Reduced sleep time for faster retries
        continue;
      }

      try {
        const searchData = JSON.parse(searchRes.body);
        console.debug(
          `Search response: ${JSON.stringify(searchData).substring(0, 100)}...`
        );

        if (searchData.messages && searchData.messages.length > 0) {
          messageId = searchData.messages[0].ID;
          console.debug(`Found message ID: ${messageId}`);
          break;
        }

        console.debug(`No matching messages found yet. Waiting...`);
        sleep(2);
      } catch (e) {
        console.error(
          `Error parsing search response: ${e}. Body: ${searchRes.body}`
        );
        sleep(2);
      }
    }

    if (!messageId) {
      console.error(
        `Failed to find TFA email for ${email} after ${MAX_TFA_FETCH_ATTEMPTS} attempts.`
      );
      return null;
    }

    // Step 2: Get the specific message content using the message ID
    const messageUrl = `${MAILPIT_URL}/api/v1/message/${messageId}`;
    console.debug(`Fetching message content from: ${messageUrl}`);

    const messageRes = http.get(messageUrl);

    if (messageRes.status !== 200) {
      console.error(
        `Failed to fetch message content. Status: ${messageRes.status}`
      );
      return null;
    }

    try {
      const messageData = JSON.parse(messageRes.body);
      console.debug(`Message data retrieved successfully`);

      // Extract the TFA code using the same regex pattern as in the Go code
      const body = messageData.HTML || messageData.Text || "";
      const codeMatch = body.match(
        /Your Two Factor authentication code is:\s*([0-9]+)/
      );

      if (codeMatch && codeMatch[1]) {
        const tfaCode = codeMatch[1];
        console.debug(`TFA code found: ${tfaCode}`);
        return tfaCode;
      } else {
        console.error(`TFA code pattern not found in email body`);
        console.debug(`Email body snippet: ${body.substring(0, 200)}...`);
        return null;
      }
    } catch (e) {
      console.error(
        `Error parsing message response: ${e}. Body: ${messageRes.body}`
      );
      return null;
    }
  } finally {
    // Always release the lock when done
    releaseTfaLock();
  }
}

// --- Authentication Function (Accepts user object with email/handle) ---
async function loginAndAuthenticateUser(user) {
  let loginAttempts = 0;
  let handle = user.handle;

  while (loginAttempts < MAX_LOGIN_ATTEMPTS) {
    loginAttempts++;
    console.debug(
      `Attempt ${loginAttempts}/${MAX_LOGIN_ATTEMPTS}: Logging in ${user.email} (handle: ${user.handle})...`
    );

    // Step 1: Login
    const loginPayload = JSON.stringify({
      email: user.email,
      password: PASSWORD,
    });

    const loginRes = http.post(`${API_BASE_URL}/hub/login`, loginPayload, {
      headers: { "Content-Type": "application/json" },
      tags: { name: "HubLoginAPI" },
    });

    if (loginRes.status === 200) {
      try {
        const loginResponseBody = JSON.parse(loginRes.body);
        console.debug(`Login response body: ${loginRes.body}`);

        if (!loginResponseBody.token) {
          console.error(`Login response missing token field: ${loginRes.body}`);
          sleep(2);
          continue;
        }

        const tfaToken = loginResponseBody.token;
        console.debug(
          `Login successful, got TFA token: ${tfaToken.substring(0, 10)}...`
        );

        // Step 2: Get TFA code from email (now with synchronization)
        const tfaCode = await fetchTFACodeForUser(user.email);
        if (!tfaCode) {
          console.error(`TFA code retrieval failed for ${user.email}`);
          sleep(2);
          continue;
        }

        // Step 3: Submit TFA code
        const tfaPayload = JSON.stringify({
          tfa_token: tfaToken,
          tfa_code: tfaCode,
          remember_me: true,
        });

        console.debug(`Submitting TFA code for ${user.email}`);

        const tfaVerifyRes = http.post(`${API_BASE_URL}/hub/tfa`, tfaPayload, {
          headers: { "Content-Type": "application/json" },
          tags: { name: "HubVerifyTFA_API" },
        });

        if (tfaVerifyRes.status === 200) {
          try {
            const tfaResponseBody = JSON.parse(tfaVerifyRes.body);

            if (!tfaResponseBody.session_token) {
              console.error(
                `TFA response missing session_token field: ${tfaVerifyRes.body}`
              );
              sleep(2);
              continue;
            }

            const sessionToken = tfaResponseBody.session_token;
            console.debug(
              `TFA verification successful, got session token: ${sessionToken.substring(
                0,
                10
              )}...`
            );

            return { authToken: sessionToken, userHandle: handle };
          } catch (e) {
            console.error(
              `Error parsing TFA response: ${e}. Body: ${tfaVerifyRes.body}`
            );
            sleep(2);
            continue;
          }
        } else {
          console.error(
            `TFA verification failed. Status: ${tfaVerifyRes.status}, Body: ${tfaVerifyRes.body}`
          );
          sleep(2);
          continue;
        }
      } catch (e) {
        console.error(
          `Error parsing login response: ${e}. Body: ${loginRes.body}`
        );
        sleep(2);
        continue;
      }
    } else {
      console.warn(
        `Login failed. Status: ${loginRes.status}, Body: ${loginRes.body}`
      );
      sleep(2);
      continue;
    }
  }

  throw new Error(
    `Authentication failed for ${user.email} after ${MAX_LOGIN_ATTEMPTS} attempts`
  );
}

// --- k6 Setup Function ---
export async function setup() {
  console.log("=== Running Setup Phase ===");

  // Clean up all existing emails from Mailpit
  const cleanupRes = http.del(`${MAILPIT_URL}/api/v1/messages`);
  if (cleanupRes.status !== 200) {
    console.warn(
      `Failed to cleanup Mailpit messages. Status: ${cleanupRes.status}, Body: ${cleanupRes.body}`
    );
  } else {
    console.log("Successfully cleaned up all existing Mailpit messages");
  }

  console.log(
    `Generating and authenticating ${NUM_USERS} users with parallelism of ${SETUP_PARALLELISM}...`
  );

  const allUsers = [];
  for (let i = 1; i <= NUM_USERS; i++) {
    allUsers.push({
      handle: `hubuser${i}`,
      email: `hubuser${i}@example.com`,
    });
  }

  if (allUsers.length === 0) {
    throw new Error(
      "No users generated. Check NUM_USERS environment variable."
    );
  }

  const authenticatedUsers = [];
  const handles = [];

  // Process users in parallel batches
  for (let i = 0; i < allUsers.length; i += SETUP_PARALLELISM) {
    const batch = allUsers.slice(
      i,
      Math.min(i + SETUP_PARALLELISM, allUsers.length)
    );
    console.log(
      `Processing batch ${Math.floor(i / SETUP_PARALLELISM) + 1}/${Math.ceil(
        allUsers.length / SETUP_PARALLELISM
      )} (${batch.length} users)...`
    );

    // Create an array of promises for parallel authentication
    const batchPromises = batch.map((user) => {
      return new Promise(async (resolve, reject) => {
        try {
          console.debug(`Starting authentication for: ${user.email}`);
          const authResult = await loginAndAuthenticateUser(user);
          if (authResult && authResult.authToken && authResult.userHandle) {
            console.debug(
              `Authentication successful for: ${user.email} (Handle: ${authResult.userHandle})`
            );
            resolve({
              authToken: authResult.authToken,
              userHandle: authResult.userHandle,
              handle: user.handle,
            });
          } else {
            reject(new Error(`Authentication failed for user ${user.email}`));
          }
        } catch (error) {
          reject(error);
        }
      });
    });

    // Wait for all authentications in this batch to complete
    try {
      const batchResults = await Promise.all(batchPromises);
      batchResults.forEach((result) => {
        authenticatedUsers.push({
          authToken: result.authToken,
          userHandle: result.userHandle,
        });
        handles.push(result.userHandle);
      });
    } catch (error) {
      throw new Error(`Batch authentication failed: ${error.message}`);
    }
  }

  console.log(
    `=== Setup Phase Complete: ${authenticatedUsers.length} users authenticated ===`
  );
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
  // The API expects the Authorization header to be in the format 'Bearer <token>'
  // Make sure there's a space between 'Bearer' and the token
  // Also ensure the token is trimmed to remove any whitespace
  const cleanToken = authToken.trim();

  // Log the full token for debugging (in production, you would never do this)
  console.debug(`FULL TOKEN BEING USED: ${cleanToken}`);

  const authParams = {
    headers: {
      Authorization: "Bearer " + cleanToken, // Ensure exact format with space after 'Bearer '
      "Content-Type": "application/json",
    },
  };

  // Debug log the exact header being sent
  console.debug(
    `Authorization header: 'Bearer ${cleanToken.substring(0, 10)}...'`
  );
  console.debug(`Full Authorization header: 'Bearer ${cleanToken}'`);

  // Log the token being used (first 10 chars only for security)
  console.debug(
    `VU ${__VU} (${userHandle}): Using token: ${authToken.substring(0, 10)}...`
  );

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
  // Calculate setup timeout based on number of users (10 seconds per user)
  setupTimeout: `${NUM_USERS * 2 * 20}s`,

  scenarios: {
    social_interactions: {
      executor: "shared-iterations",
      vus: `${NUM_USERS}`, // Adjust based on your needs
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
