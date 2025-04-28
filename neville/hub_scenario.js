// loadtest/hub_scenario.js
import {
  randomIntBetween,
  randomItem,
} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { check, group, sleep } from "k6";
import { SharedArray } from "k6/data";
import http from "k6/http";
import { Trend } from "k6/metrics";

// --- Configuration ---
const API_BASE_URL = __ENV.API_BASE_URL || "http://localhost:8080"; // Your API gateway/service URL
const MAILPIT_URL = __ENV.MAILPIT_URL || "http://localhost:8025"; // Mailpit API URL
const NUM_USERS = parseInt(__ENV.NUM_USERS || "1000");
const PASSWORD = "NewPassword123$";
const LOGIN_BATCH_SIZE = parseInt(__ENV.LOGIN_BATCH_SIZE || "50");
const LOGIN_BATCH_INTERVAL_S = parseInt(__ENV.LOGIN_BATCH_INTERVAL_S || "30"); // Seconds between login batches

// --- Metrics ---
const loginTrend = new Trend("hub_login_duration", true);
const tfaFetchTrend = new Trend("hub_tfa_fetch_duration", true);
const tfaVerifyTrend = new Trend("hub_tfa_verify_duration", true);
const followTrend = new Trend("hub_follow_user_duration", true); // Updated metric name
const createPostTrend = new Trend("hub_create_post_duration", true);
const timelineReadTrend = new Trend("hub_timeline_read_duration", true);
const upvoteTrend = new Trend("hub_upvote_duration", true);
const downvoteTrend = new Trend("hub_downvote_duration", true);
const unfollowTrend = new Trend("hub_unfollow_user_duration", true); // Added trend for unfollow

// --- Data ---
// Assumes users hubuser1, hubuser2, ..., hubuserN exist
const usernames = new SharedArray("usernames", function () {
  const users = [];
  for (let i = 1; i <= NUM_USERS; i++) {
    users.push(`hubuser${i}`);
  }
  return users;
});

// Store user handles (assuming they are the primary identifier for actions like follow)
const userHandles = usernames; // Alias for clarity, assuming usernames are the handles

// Store tokens obtained after successful login+TFA
// Each VU manages its own state (token, handle, maybe fetched post IDs)
let vuState = {
  authToken: null,
  userHandle: null,
  fetchedPostIds: [],
  timelineCursor: null, // For pagination
  tfaToken: null, // Store intermediate TFA token from login response
};

// --- Stages ---
// Stage 1: Login batches
// Stage 2: Ramp up social activity
// Stage 3: Sustained social activity
export const options = {
  scenarios: {
    hub_load_test: {
      executor: "ramping-arrival-rate",
      // Define maxVUs at the scenario level, not per stage
      maxVUs: NUM_USERS + 50, // Keep buffer
      // Define preAllocatedVUs at the scenario level
      preAllocatedVUs: NUM_USERS, // Pre-allocate based on expected user count
      // Define timeUnit at the scenario level
      timeUnit: "30s", // Revert to 30s to allow integer targets
      stages: [
        // Stage 1: Login users in batches (gradual ramp-up)
        {
          duration: "5m", // Ramp up over 5 minutes
          target: 10, // Target 10 iterations per 30s
        }, // Ramp up to 10 iterations/30s over 5 mins

        // Stage 2: Full Load Social Activity (e.g., for 10 minutes)
        {
          duration: "10m",
          target: 3000, // Target 3000 iterations per 30s (equivalent to 100/sec)
        }, // Sustain 3000 iterations/30s

        // Stage 3: Ramp-down (optional)
        {
          duration: "2m",
          target: 0,
        },
      ],
      exec: "socialActivity", // Function executed by VUs
      gracefulStop: "30s",
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"], // http errors should be less than 1%
    http_req_duration: ["p(95)<1000"], // 95% of requests should be below 1000ms
    "group_duration{group:::Hub Login}": ["p(95)<1500"],

    // This is disabled because we do not need to profile TFA as mailpit is known to be slow
    // "group_duration{group:::Hub TFA Fetch and Verify}": ["p(95)<3000"],

    "group_duration{group:::Hub Social Interaction}": ["p(95)<800"],
    // Add specific thresholds for trends if needed
    [loginTrend.name]: ["p(95)<1500"],

    // The two below are disabled because we do not need to profile TFA as mailpit is known to be slow
    // [tfaFetchTrend.name]: ["p(95)<2000"],
    // [tfaVerifyTrend.name]: ["p(95)<1000"],

    [createPostTrend.name]: ["p(95)<500"],
    [timelineReadTrend.name]: ["p(95)<700"],
    [followTrend.name]: ["p(95)<400"],
    [upvoteTrend.name]: ["p(95)<300"],
    [downvoteTrend.name]: ["p(95)<300"],
    [unfollowTrend.name]: ["p(95)<400"], // Added threshold for unfollow
  },
  // Increase default HTTP timeout if needed
  // httpTimeout: '60s',
};

// --- Helper Functions ---

// Function to fetch TFA code from Mailpit
function fetchTFACodeFromMailpit(username) {
  let mailId = null;
  let tfaCode = null;
  let attempts = 0;
  const maxAttempts = 6;

  // Construct the search query using the full email address
  const fullEmail = `${username}@example.com`;
  const searchQuery = encodeURIComponent(
    `to:${fullEmail} subject:"Vetchium Two Factor Authentication"` // Use full email here
  );
  const searchUrl = `${MAILPIT_URL}/api/v1/search?query=${searchQuery}`;
  console.debug(`VU ${__VU}: Using Mailpit search URL: ${searchUrl}`); // Log the search URL using INFO level

  // Poll Mailpit for the email
  while (attempts < maxAttempts && !mailId) {
    const searchRes = http.get(searchUrl, { tags: { name: "MailpitSearch" } });

    // Check status and try to parse the full JSON body
    if (searchRes.status === 200) {
      try {
        const responseBody = searchRes.json(); // Get the whole JSON object

        // Check if messages array exists and is not empty
        if (
          responseBody &&
          Array.isArray(responseBody.messages) &&
          responseBody.messages.length > 0
        ) {
          const firstMessage = responseBody.messages[0];
          console.debug(
            `VU ${__VU}: First message object: ${JSON.stringify(firstMessage)}`
          );

          if (firstMessage && firstMessage.ID) {
            mailId = firstMessage.ID;
            console.debug(
              `VU ${__VU}: Found mail ID ${mailId} for ${username}`
            );
            // Break the loop once found
            break;
          } else {
            console.warn(
              `VU ${__VU}: First message found but missing ID: ${JSON.stringify(
                firstMessage
              )}`
            );
          }
        } else {
          // Response might be ok but no messages array or empty
          console.debug(
            `VU ${__VU}: Mailpit search OK, but no messages found yet. Body: ${searchRes.body.substring(
              0,
              100
            )}...`
          );
        }
      } catch (e) {
        console.error(
          `VU ${__VU}: Failed to parse Mailpit search JSON response. Error: ${e}. Body: ${searchRes.body.substring(
            0,
            500
          )}...`
        );
        // Don't retry immediately on parse error, wait
      }
    }

    // If mailId wasn't found in this attempt, increment and wait
    if (!mailId) {
      attempts++;
      console.debug(
        `VU ${__VU}: Mail for ${username} not found yet (attempt ${attempts}). Waiting...`
      );
      sleep(2); // Reduced wait time from 10s
    }
  }

  // --- Fetch the specific email content using the found mailId ---
  if (!mailId) {
    console.error(
      `VU ${__VU}: Failed to find Mailpit email for ${username} after ${maxAttempts} attempts.`
    );
    return { tfaCode: null, mailId: null };
  }

  // Fetch the email content
  const messageUrl = `${MAILPIT_URL}/api/v1/message/${mailId}`;
  const msgRes = http.get(messageUrl, { tags: { name: "MailpitGetMessage" } });
  if (msgRes.status === 200 && msgRes.body) {
    // Use regex found in Go test helper
    const match = msgRes.body.match(
      /Your Two Factor authentication code is:\s*(\d+)/
    );
    if (match && match[1]) {
      tfaCode = match[1];
      console.debug(
        `VU ${__VU}: Extracted TFA code ${tfaCode} for ${username}`
      );
    } else {
      console.error(
        `VU ${__VU}: Could not extract TFA code from email body for ${username}. Mail Body: ${msgRes.body.substring(
          0,
          500
        )}...`
      );
    }
  } else {
    console.error(
      `VU ${__VU}: Failed to fetch Mailpit message ${mailId} for ${username}. Status: ${msgRes.status}`
    );
  }

  return { tfaCode, mailId };
}

// Function to delete email from Mailpit
function deleteMailpitEmail(mailId) {
  if (!mailId) return;
  // *** PLACEHOLDER - CONFIRM MAILPIT DELETE API ***
  const deleteUrl = `${MAILPIT_URL}/api/v1/messages`;
  const res = http.del(deleteUrl, JSON.stringify({ ids: [mailId] }), {
    headers: { "Content-Type": "application/json" },
    tags: { name: "MailpitDelete" },
  });
  console.debug(
    `VU ${__VU}: Attempted Mailpit delete for ${mailId}. Status: ${res.status}`
  );
  check(res, { "Mailpit email deleted": (r) => r.status === 200 });
}

// --- Setup function (runs once before test) ---
// We could potentially pre-fetch user handles/IDs here if needed,
// but the SharedArray approach is generally better for large datasets.
export function setup() {
  console.log(`Starting test with NUM_USERS=${NUM_USERS}`);
  // Initialize state for VUs if necessary (though VUs have independent state)
  return {}; // Pass data to VUs if needed
}

// --- Main Execution Logic --- (Executed by each VU)

// Function to perform login and TFA for a VU
function loginAndAuthenticate(username) {
  let success = false;
  let loginTfaToken = null; // Store the intermediate token from login response
  const email = `${username}@example.com`;

  group("Hub Login", function () {
    const loginPayload = JSON.stringify({
      // Use 'Email' field and construct value by appending domain to username
      Email: email, // Construct email address here
      password: PASSWORD,
    });
    const loginParams = {
      headers: { "Content-Type": "application/json" },
      tags: { name: "HubLoginAPI" },
    };
    const loginRes = http.post(
      `${API_BASE_URL}/hub/login`,
      loginPayload,
      loginParams
    );
    const loginTime = loginRes.timings.duration;
    loginTrend.add(loginTime);

    check(loginRes, {
      "Login indicates success/TFA needed (status 200/202)": (r) =>
        r.status === 200 || r.status === 202,
    });

    if (loginRes.status === 200 || loginRes.status === 202) {
      // Extract the intermediate TFA token from the response
      if (loginRes.json("token")) {
        loginTfaToken = loginRes.json("token");
        console.debug(
          `VU ${__VU} (${email}): Login successful (Status: ${loginRes.status}). Got TFA token.`
        );
        success = true;
      } else {
        console.error(
          `VU ${__VU} (${email}): Login response OK/Accepted but missing 'token'. Body: ${loginRes.body}`
        );
        success = false;
      }
    } else {
      console.error(
        `VU ${__VU} (${username}): Login failed! Status: ${loginRes.status}, Body: ${loginRes.body}`
      );
      success = false;
    }
  });

  if (!success) return false; // Don't proceed if login failed

  group("Hub TFA Fetch and Verify", function () {
    // 1. Fetch TFA Code from Mailpit
    const mailpitStartTime = Date.now();
    const { tfaCode, mailId } = fetchTFACodeFromMailpit(username);
    const tfaFetchTime = Date.now() - mailpitStartTime;
    tfaFetchTrend.add(tfaFetchTime);

    if (tfaCode) {
      // 2. Verify TFA Code
      const tfaPayload = JSON.stringify({
        tfa_token: loginTfaToken,
        tfa_code: tfaCode,
        remember_me: true,
      });
      const tfaParams = {
        headers: { "Content-Type": "application/json" },
        tags: { name: "HubTfaVerifyAPI" },
      };
      const tfaRes = http.post(
        `${API_BASE_URL}/hub/tfa`,
        tfaPayload,
        tfaParams
      );
      const tfaVerifyTime = tfaRes.timings.duration;
      tfaVerifyTrend.add(tfaVerifyTime);

      check(tfaRes, {
        "TFA verification successful (status 200)": (r) => r.status === 200,
        // Expect 'sessionToken' based on Go test
        "TFA response contains session_token": (r) =>
          r.body &&
          r.json("session_token") !== null &&
          r.json("session_token") !== undefined,
      });

      if (tfaRes.status === 200 && tfaRes.json("session_token")) {
        // Store the final session token
        vuState.authToken = tfaRes.json("session_token");
        // *** PLACEHOLDER: Extract logged-in user ID/handle if available ***
        // vuState.userHandle = tfaRes.json('handle'); // Or 'userId', 'username' etc.
        vuState.userHandle = username; // Assume handle is the username for now
        console.debug(
          `VU ${__VU} (${vuState.userHandle}): TFA successful. Token obtained.`
        );
        // 3. Clean up Mailpit
        deleteMailpitEmail(mailId);
        success = true;
      } else {
        console.error(
          `VU ${__VU} (${username}): TFA verification failed! Status: ${tfaRes.status}, Body: ${tfaRes.body}`
        );
        success = false;
      }
    } else {
      console.error(`VU ${__VU} (${username}): Could not retrieve TFA code.`);
      success = false;
    }
  });
  return success;
}

// Function for the main social activity loop, executed by each VU
export function socialActivity() {
  // Pick a user for this VU based on its ID
  // Ensure VU index is within bounds of the usernames array
  const userIndex = (__VU - 1) % usernames.length;
  const username = usernames[userIndex];

  // --- Authentication Check ---
  // Attempt login ONLY if the VU doesn't have a token yet for this session
  if (!vuState.authToken) {
    console.debug(
      `VU ${__VU} (${username}): Not authenticated. Attempting login and TFA...`
    );
    const loginSuccess = loginAndAuthenticate(username);

    if (!loginSuccess) {
      console.error(
        `VU ${__VU} (${username}): Authentication failed. Skipping social actions for this iteration.`
      );
      // Wait a bit longer after a failed login attempt before the VU retries
      sleep(randomIntBetween(5, 10));
      return; // Exit this iteration
    }
    // If login succeeded, vuState.authToken is now set by loginAndAuthenticate
    console.debug(`VU ${__VU} (${username}): Authentication successful.`);
  }

  // If we reach here, the VU is authenticated (either previously or just now)
  // Proceed with social actions...

  // Define possible actions and their weights (adjust as needed)
  // Example: More reading/voting than posting/following
  const possibleActions = [
    "follow",
    "unfollow", // Added unfollow action
    "createPost",
    "readTimeline",
    "readTimeline", // Read timeline more often
    "readTimeline",
    "vote",
    "vote", // Vote more often
    "vote",
  ];
  const action = randomItem(possibleActions);

  // Group related actions for better reporting
  group("Hub Social Interaction", function () {
    const authParams = {
      headers: {
        Authorization: `Bearer ${vuState.authToken}`,
        "Content-Type": "application/json",
      },
    };

    switch (action) {
      case "follow":
        // Select a random user to follow (ensure it's not the current VU)
        let userToFollow;
        do {
          userToFollow = randomItem(userHandles);
        } while (userToFollow === vuState.userHandle);

        console.debug(
          `VU ${__VU} (${vuState.userHandle}): Attempting to follow ${userToFollow}`
        );
        const followPayload = JSON.stringify({ handle: userToFollow });
        const followRes = http.post(
          `${API_BASE_URL}/hub/follow-user`,
          followPayload,
          { ...authParams, tags: { name: "HubFollowAPI" } }
        );
        followTrend.add(followRes.timings.duration);
        // Check for success (200) or already following (e.g., 409 Conflict or similar - adjust based on actual API)
        // Assuming 200 for success, 4xx might indicate already followed or other client error.
        check(followRes, {
          "Follow request successful or already following (status 200/4xx)": (
            r
          ) => r.status === 200 || (r.status >= 400 && r.status < 500),
        });
        if (followRes.status !== 200 && followRes.status < 400) {
          // Log unexpected non-200/4xx statuses
          console.warn(
            `VU ${__VU} (${vuState.userHandle}): Follow request unusual status. Status: ${followRes.status}, Body: ${followRes.body}`
          );
        }
        sleep(randomIntBetween(1, 3)); // Pause after action
        break;

      case "unfollow":
        // Select a random user to unfollow (ensure it's not the current VU)
        let userToUnfollow;
        do {
          userToUnfollow = randomItem(userHandles);
        } while (userToUnfollow === vuState.userHandle);

        console.debug(
          `VU ${__VU} (${vuState.userHandle}): Attempting to unfollow ${userToUnfollow}`
        );
        const unfollowPayload = JSON.stringify({ handle: userToUnfollow });
        const unfollowRes = http.post(
          `${API_BASE_URL}/hub/unfollow-user`,
          unfollowPayload,
          { ...authParams, tags: { name: "HubUnfollowAPI" } }
        );
        unfollowTrend.add(unfollowRes.timings.duration);
        // Check for success (200) or potentially errors like not following (e.g., 4xx) - adjust based on actual API.
        // Assuming 200 for success, 4xx might indicate not following or other client error.
        check(unfollowRes, {
          "Unfollow request successful or user not followed (status 200/4xx)": (
            r
          ) => r.status === 200 || (r.status >= 400 && r.status < 500),
        });
        if (unfollowRes.status !== 200 && unfollowRes.status < 400) {
          // Log unexpected non-200/4xx statuses
          console.warn(
            `VU ${__VU} (${vuState.userHandle}): Unfollow request unusual status. Status: ${unfollowRes.status}, Body: ${unfollowRes.body}`
          );
        }
        sleep(randomIntBetween(1, 3)); // Pause after action
        break;

      case "createPost":
        console.debug(
          `VU ${__VU} (${vuState.userHandle}): Attempting to create post`
        );
        const postPayload = JSON.stringify({
          content: `This is a test post from ${
            vuState.userHandle
          } VU ${__VU} at ${new Date().toISOString()}`,
          // Add other fields like visibility if required by CreateUserPostRequest
        });
        const postRes = http.post(`${API_BASE_URL}/hub/add-post`, postPayload, {
          ...authParams,
          tags: { name: "HubCreatePostAPI" },
        });
        createPostTrend.add(postRes.timings.duration);
        check(postRes, {
          "Create post successful (status 200)": (r) => r.status === 200,
          "Create post response has post_id": (r) =>
            r.body &&
            r.json("post_id") !== null &&
            r.json("post_id") !== undefined,
        });
        if (postRes.status === 200 && postRes.json("post_id")) {
          const newPostId = postRes.json("post_id");
          // Optionally add to VU's list of posts for potential voting later
          if (vuState.fetchedPostIds.length < 50) {
            vuState.fetchedPostIds.push(newPostId);
          } else {
            // Simple FIFO replacement if list is full
            vuState.fetchedPostIds.shift();
            vuState.fetchedPostIds.push(newPostId);
          }
          console.debug(
            `VU ${__VU} (${vuState.userHandle}): Created post ${newPostId}`
          );
        } else {
          console.error(
            `VU ${__VU} (${vuState.userHandle}): Failed to create post. Status: ${postRes.status}, Body: ${postRes.body}`
          );
        }
        sleep(randomIntBetween(2, 5)); // Pause after action
        break;

      case "readTimeline":
        console.debug(
          `VU ${__VU} (${vuState.userHandle}): Attempting to read timeline (cursor: ${vuState.timelineCursor})`
        );
        const timelinePayload = JSON.stringify({
          limit: 10, // Request 10 posts per page
          cursor: vuState.timelineCursor, // Use stored cursor for pagination
        });
        const timelineRes = http.post(
          `${API_BASE_URL}/hub/get-my-home-timeline`,
          timelinePayload,
          { ...authParams, tags: { name: "HubTimelineAPI" } }
        );
        timelineReadTrend.add(timelineRes.timings.duration);
        check(timelineRes, {
          "Read timeline successful (status 200)": (r) => r.status === 200,
        });

        if (timelineRes.status === 200 && timelineRes.body) {
          const timelineData = timelineRes.json();
          if (timelineData.posts && timelineData.posts.length > 0) {
            // Extract post IDs from the response for potential voting
            const postIds = timelineData.posts
              .map((post) => post.postId)
              .filter((id) => id);
            if (postIds.length > 0) {
              // Add new, unique post IDs to the VU's state
              const uniqueNewIds = postIds.filter(
                (id) => !vuState.fetchedPostIds.includes(id)
              );
              vuState.fetchedPostIds.push(...uniqueNewIds);
              // Keep the list size manageable
              if (vuState.fetchedPostIds.length > 50) {
                vuState.fetchedPostIds = vuState.fetchedPostIds.slice(
                  vuState.fetchedPostIds.length - 50
                );
              }
              console.debug(
                `VU ${__VU} (${vuState.userHandle}): Fetched ${postIds.length} posts. Total known IDs: ${vuState.fetchedPostIds.length}`
              );
            }
          }
          // Update cursor for next timeline request (assuming response provides 'nextCursor')
          vuState.timelineCursor = timelineData.nextCursor || null;
        } else {
          // Reset cursor on error or empty response?
          // vuState.timelineCursor = null;
        }
        sleep(randomIntBetween(1, 4)); // Pause after action
        break;

      case "vote":
        // Only vote if we have some post IDs fetched previously
        if (vuState.fetchedPostIds.length > 0) {
          const postIdToVote = randomItem(vuState.fetchedPostIds);
          const voteType = randomItem(["upvote", "downvote"]);
          console.debug(
            `VU ${__VU} (${vuState.userHandle}): Attempting to ${voteType} post ${postIdToVote}`
          );

          const votePayload = JSON.stringify({
            post_id: postIdToVote,
          });
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
          // Note: API spec says 200 is returned even if already voted/downvoted, or 422 for specific errors.
          check(voteRes, {
            [`${voteType} request sent (status 200/422)`]: (r) =>
              r.status === 200 || r.status === 422,
          });
          if (voteRes.status !== 200 && voteRes.status !== 422) {
            console.warn(
              `VU ${__VU} (${vuState.userHandle}): Vote request failed unexpectedly. Status: ${voteRes.status}, Body: ${voteRes.body}`
            );
          }
          sleep(randomIntBetween(1, 2)); // Pause after action
        } else {
          console.debug(
            `VU ${__VU} (${vuState.userHandle}): No post IDs available to vote on yet.`
          );
          sleep(1); // Short pause if no voting possible
        }
        break;
    }

    // Think time between actions in the main loop
    sleep(randomIntBetween(2, 6));
  }); // End 'Hub Social Interaction' group
}

// --- Teardown function (runs once after test) ---
export function teardown(data) {
  console.log("Test finished.");
  // Clean up resources if needed
}
