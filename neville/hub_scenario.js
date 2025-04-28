// loadtest/hub_scenario.js
import {
  randomIntBetween,
  randomItem,
} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { check, group, sleep } from "k6";
import { SharedArray } from "k6/data";
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
        if (!userHandles.length) break;
        const userToFollow = randomItem(userHandles);
        const followPayload = JSON.stringify({ handle: userToFollow });
        console.debug(
          `VU ${__VU} (${vuState.userHandle}): Attempting to follow ${userToFollow}`
        );
        const followRes = http.post(
          `${API_BASE_URL}/hub/follow-user`,
          followPayload,
          { ...authParams, tags: { name: "HubFollowAPI" } }
        );
        followTrend.add(followRes.timings.duration);
        // Reverted check: Allow 200 or 4xx
        check(followRes, {
          "Follow request successful or expected client error (status 200/4xx)":
            (r) => r.status === 200 || (r.status >= 400 && r.status < 500),
        });
        // Log only truly unexpected errors (not 200 and not 4xx)
        if (
          followRes.status !== 200 &&
          !(followRes.status >= 400 && followRes.status < 500)
        ) {
          console.error(
            `VU ${__VU} (${vuState.userHandle}): Follow API Unexpected Error! Status: ${followRes.status}, Body: ${followRes.body}`
          );
        }
        sleep(randomIntBetween(1, 3)); // Pause after action
        break;

      case "unfollow":
        if (!userHandles.length) break;
        // Prefer unfollowing someone recently followed by this VU, otherwise random
        // Simplified: Just pick a random user for unfollow attempt in load test
        const userToUnfollow = randomItem(userHandles);
        const unfollowPayload = JSON.stringify({ handle: userToUnfollow });
        console.debug(
          `VU ${__VU} (${vuState.userHandle}): Attempting to unfollow ${userToUnfollow}`
        );
        const unfollowRes = http.post(
          `${API_BASE_URL}/hub/unfollow-user`,
          unfollowPayload,
          { ...authParams, tags: { name: "HubUnfollowAPI" } }
        );
        unfollowTrend.add(unfollowRes.timings.duration);
        // Reverted check: Allow 200 or 4xx
        check(unfollowRes, {
          "Unfollow request successful or expected client error (status 200/4xx)":
            (r) => r.status === 200 || (r.status >= 400 && r.status < 500),
        });
        // Log only truly unexpected errors (not 200 and not 4xx)
        if (
          unfollowRes.status !== 200 &&
          !(unfollowRes.status >= 400 && unfollowRes.status < 500)
        ) {
          console.error(
            `VU ${__VU} (${vuState.userHandle}): Unfollow API Unexpected Error! Status: ${unfollowRes.status}, Body: ${unfollowRes.body}`
          );
        }
        sleep(randomIntBetween(1, 3)); // Pause after action
        break;

      case "createPost":
        console.debug(
          `VU ${__VU} (${vuState.userHandle}): Attempting to create post`
        );
        const postContent = `Post content from VU ${__VU} at ${new Date().toISOString()}: ${randomString(
          50
        )}`;
        const numTags = randomIntBetween(0, 3);
        const postTags = [];
        for (let i = 0; i < numTags; i++) {
          postTags.push(`tag_${randomString(5)}`);
        }
        // Use new_tags as per TypeSpec, omit tag_ids for simplicity
        const postPayload = JSON.stringify({
          content: postContent,
          new_tags: postTags,
        });
        console.debug(
          `VU ${__VU} (${
            vuState.userHandle
          }): Attempting to create post. Tags: ${postTags.join(", ")}`
        );
        const postRes = http.post(`${API_BASE_URL}/hub/add-post`, postPayload, {
          ...authParams,
          tags: { name: "HubCreatePostAPI" },
        });
        createPostTrend.add(postRes.timings.duration);
        // Keep strict check for 200 and post_id
        check(postRes, {
          "Create post successful (status 200)": (r) => r.status === 200,
          "Create post response has post_id": (r) =>
            r.body &&
            r.json("post_id") !== null &&
            r.json("post_id") !== undefined,
        });
        // Log any non-200 response or missing post_id
        if (postRes.status !== 200) {
          console.error(
            `VU ${__VU} (${vuState.userHandle}): Create Post API Error! Status: ${postRes.status}, Body: ${postRes.body}`
          );
        } else if (!postRes.json("post_id")) {
          console.error(
            `VU ${__VU} (${vuState.userHandle}): Create Post API Error! Status 200 but missing post_id. Body: ${postRes.body}`
          );
        }
        // ... (post processing remains the same)
        sleep(randomIntBetween(2, 5)); // Pause after action
        break;

      case "readTimeline":
        let timelineUrl = `${API_BASE_URL}/hub/get-my-home-timeline`;
        // Prepare body for POST request
        let timelinePayload = null;
        if (vuState.timelineCursor) {
          timelinePayload = JSON.stringify({
            pagination_key: vuState.timelineCursor,
          });
          console.debug(
            `VU ${__VU} (${vuState.userHandle}): Reading timeline with pagination_key: ${vuState.timelineCursor}`
          );
        } else {
          console.debug(
            `VU ${__VU} (${vuState.userHandle}): Reading timeline (first page).`
          );
        }
        // Change to POST and send payload in body
        const timelineRes = http.post(timelineUrl, timelinePayload, {
          ...authParams,
          tags: { name: "HubTimelineReadAPI" },
        });
        timelineReadTrend.add(timelineRes.timings.duration);
        // Keep strict check for 200
        check(timelineRes, {
          "Read timeline successful (status 200)": (r) => r.status === 200,
        });
        // Log any non-200 response
        if (timelineRes.status !== 200) {
          console.error(
            `VU ${__VU} (${vuState.userHandle}): Read Timeline API Error! Status: ${timelineRes.status}, Body: ${timelineRes.body}`
          );
        } else {
          // ... (timeline processing remains the same)
        }
        sleep(randomIntBetween(3, 7)); // Pause after reading
        break;

      case "vote":
        if (vuState.fetchedPostIds.length === 0) {
          console.debug(
            `VU ${__VU} (${vuState.userHandle}): No posts fetched yet, skipping vote.`
          );
          break; // Cannot vote if no posts are known
        }
        const postIdToVote = randomItem(vuState.fetchedPostIds);
        const voteType = randomItem(["upvote", "downvote"]);

        console.debug(
          `VU ${__VU} (${vuState.userHandle}): Attempting to ${voteType} post ${postIdToVote}`
        );

        // Use post_id in JSON body
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
        // Reverted check: Allow 200 or 422
        check(voteRes, {
          [`${voteType} request successful or expected error (status 200/422)`]:
            (r) => r.status === 200 || r.status === 422,
        });
        // Log only truly unexpected errors (not 200 and not 422)
        if (voteRes.status !== 200 && voteRes.status !== 422) {
          console.error(
            `VU ${__VU} (${vuState.userHandle}): ${voteType} API Unexpected Error! PostID: ${postIdToVote}, Status: ${voteRes.status}, Body: ${voteRes.body}`
          );
        }
        sleep(randomIntBetween(1, 4)); // Pause after voting
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
