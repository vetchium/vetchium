// distributed_hub_scenario.js - Distributed k6 load test for Vetchium API
// SCRIPT VERSION: 2025-05-04-v1 - Fixed ConfigMap loading and TFA extraction
import {
  randomIntBetween,
  randomItem,
} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { check, group, sleep } from "k6";
import exec from "k6/execution";
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
const TOTAL_USERS = parseInt(__ENV.TOTAL_USERS || "0");
const INSTANCE_INDEX = parseInt(__ENV.INSTANCE_INDEX || "0");
const INSTANCE_COUNT = parseInt(__ENV.INSTANCE_COUNT || "0");
const USERS_PER_INSTANCE = parseInt(__ENV.USERS_PER_INSTANCE || "0");
const SETUP_PARALLELISM = parseInt(__ENV.SETUP_PARALLELISM || "0");
const PASSWORD = "NewPassword123$";
const TEST_DURATION_SECONDS = parseInt(__ENV.TEST_DURATION || "0");

// Validate required environment variables
if (!API_BASE_URL) {
  console.error("API_BASE_URL environment variable is required");
  exec.test.abort();
}

if (!MAILPIT_URL) {
  console.error("MAILPIT_URL environment variable is required");
  exec.test.abort();
}

if (TOTAL_USERS <= 0) {
  console.error("TOTAL_USERS must be a positive number");
  exec.test.abort();
}

if (INSTANCE_COUNT <= 0) {
  console.error("INSTANCE_COUNT must be a positive number");
  exec.test.abort();
}

if (INSTANCE_INDEX < 0 || INSTANCE_INDEX >= INSTANCE_COUNT) {
  console.error(`INSTANCE_INDEX must be between 0 and ${INSTANCE_COUNT - 1}`);
  exec.test.abort();
}

if (USERS_PER_INSTANCE <= 0) {
  console.error("USERS_PER_INSTANCE must be a positive number");
  exec.test.abort();
}

if (SETUP_PARALLELISM <= 0) {
  console.error("SETUP_PARALLELISM must be a positive number");
  exec.test.abort();
}

if (TEST_DURATION_SECONDS <= 0) {
  console.error("TEST_DURATION must be a positive number");
  exec.test.abort();
}

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

// --- IncognitoPosts Metrics ---
const createIncognitoPostTrend = new Trend(
  "hub_create_incognito_post_duration",
  true
);
const getIncognitoPostTrend = new Trend(
  "hub_get_incognito_post_duration",
  true
);
const getIncognitoPostsTrend = new Trend(
  "hub_get_incognito_posts_duration",
  true
);
const getMyIncognitoPostsTrend = new Trend(
  "hub_get_my_incognito_posts_duration",
  true
);
const deleteIncognitoPostTrend = new Trend(
  "hub_delete_incognito_post_duration",
  true
);
const addIncognitoCommentTrend = new Trend(
  "hub_add_incognito_comment_duration",
  true
);
const getIncognitoCommentsTrend = new Trend(
  "hub_get_incognito_comments_duration",
  true
);
const upvoteIncognitoPostTrend = new Trend(
  "hub_upvote_incognito_post_duration",
  true
);
const downvoteIncognitoPostTrend = new Trend(
  "hub_downvote_incognito_post_duration",
  true
);
const unvoteIncognitoPostTrend = new Trend(
  "hub_unvote_incognito_post_duration",
  true
);
const upvoteIncognitoCommentTrend = new Trend(
  "hub_upvote_incognito_comment_duration",
  true
);
const downvoteIncognitoCommentTrend = new Trend(
  "hub_downvote_incognito_comment_duration",
  true
);
const unvoteIncognitoCommentTrend = new Trend(
  "hub_unvote_incognito_comment_duration",
  true
);

// --- Constants for authentication ---
const MAX_LOGIN_ATTEMPTS = 5;
const MAX_TFA_FETCH_ATTEMPTS = 10;

// --- Common tag IDs for incognito posts (from vetchium-tags.json) ---
const AVAILABLE_TAG_IDS = [
  "technology",
  "artificial-intelligence",
  "machine-learning",
  "data-science",
  "software-engineering",
  "web-development",
  "mobile-development",
  "devops",
  "cybersecurity",
  "blockchain",
  "cryptocurrency",
  "cloud-computing",
  "careers",
  "leadership",
  "management",
  "human-resources",
  "entrepreneurship",
  "startups",
  "business-development",
  "marketing",
  "sales",
  "productivity",
  "work-life-balance",
  "remote-work",
  "networking",
  "mentorship",
  "coaching",
  "personal-development",
  "creativity",
  "motivation",
  "writing",
  "content-creation",
  "photography",
  "music",
  "sports",
  "travel",
  "health",
  "finance",
  "real-estate",
];

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

// --- Calculate user range for this instance ---
// Use a more reasonable limit for users per instance based on available memory
const MAX_USERS_PER_INSTANCE = Math.min(5000, USERS_PER_INSTANCE); // Higher limit, but still prevent extreme OOM

const startUserIndex = INSTANCE_INDEX * USERS_PER_INSTANCE + 1;
const endUserIndex = Math.min(
  startUserIndex + MAX_USERS_PER_INSTANCE - 1,
  (INSTANCE_INDEX + 1) * USERS_PER_INSTANCE,
  TOTAL_USERS
);

console.log(
  `Instance ${INSTANCE_INDEX}/${INSTANCE_COUNT} handling users ${startUserIndex} to ${endUserIndex}`
);

// --- TFA Code Fetch (Uses Email) ---
async function fetchTFACodeForUser(email) {
  console.log(`Starting TFA code fetch for ${email}`);
  try {
    // Acquire lock before fetching TFA code
    console.log(`Acquiring TFA lock for ${email}`);
    await acquireTfaLock();
    console.log(`TFA lock acquired for ${email}`);

    let attempts = 0;
    let messageId = null;

    // Step 1: Search for the TFA email using the search API
    while (attempts < MAX_TFA_FETCH_ATTEMPTS) {
      attempts++;

      // Build the search URL with query parameters exactly as in the Go code
      // Use a more relaxed search query to find TFA emails
      const searchQuery = `to:${email}`;
      const searchUrl = `${MAILPIT_URL}/api/v1/search?query=${encodeURIComponent(
        searchQuery
      )}`;

      console.log(
        `Attempt ${attempts}/${MAX_TFA_FETCH_ATTEMPTS}: Searching Mailpit at ${searchUrl}`
      );

      // Debug the actual Mailpit URL being used
      console.log(`DEBUG: Full Mailpit URL being accessed: ${MAILPIT_URL}`);
      console.log(`DEBUG: Full search URL: ${searchUrl}`);

      // Add detailed logging for the HTTP request
      console.log(`DEBUG: Making HTTP GET request to Mailpit search API`);
      let searchRes;
      try {
        searchRes = http.get(searchUrl);
        console.log(
          `DEBUG: Mailpit search API response status: ${searchRes.status}`
        );
        console.log(
          `DEBUG: Mailpit search API response size: ${
            searchRes.body ? searchRes.body.length : 0
          } bytes`
        );
      } catch (error) {
        console.error(
          `DEBUG: Exception during Mailpit search API request: ${error}`
        );
        sleep(2);
        continue;
      }

      if (searchRes.status !== 200) {
        console.warn(
          `DEBUG: Mailpit search API returned status ${searchRes.status}. Response body: ${searchRes.body}`
        );
        sleep(2); // Reduced sleep time for faster retries
        continue;
      }

      try {
        console.log(`DEBUG: Attempting to parse JSON response from Mailpit`);
        console.log(
          `DEBUG: First 200 chars of response body: ${searchRes.body.substring(
            0,
            200
          )}...`
        );

        const searchData = JSON.parse(searchRes.body);
        console.log(`DEBUG: JSON parsed successfully`);
        console.log(
          `DEBUG: Search response for ${email}: Found ${
            searchData.messages ? searchData.messages.length : 0
          } messages`
        );

        if (searchData.messages && searchData.messages.length > 0) {
          // Sort messages by date (newest first) to get the most recent TFA email
          searchData.messages.sort(
            (a, b) => new Date(b.Date) - new Date(a.Date)
          );
          messageId = searchData.messages[0].ID;
          console.log(`Found message ID for ${email}: ${messageId}`);
          break;
        }

        console.log(`No matching messages found yet for ${email}. Waiting...`);
        sleep(3); // Increased wait time
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
      releaseTfaLock();
      return null;
    }

    // Step 2: Get the message content using the message ID
    const messageUrl = `${MAILPIT_URL}/api/v1/message/${messageId}`;
    console.log(`DEBUG: Fetching message content from ${messageUrl}`);

    let messageRes;
    try {
      console.log(`DEBUG: Making HTTP GET request to fetch message content`);
      messageRes = http.get(messageUrl);
      console.log(`DEBUG: Message API response status: ${messageRes.status}`);
      console.log(
        `DEBUG: Message API response size: ${
          messageRes.body ? messageRes.body.length : 0
        } bytes`
      );
    } catch (error) {
      console.error(`DEBUG: Exception during message fetch: ${error}`);
      releaseTfaLock();
      return null;
    }

    if (messageRes.status !== 200) {
      console.error(
        `Failed to fetch message content. Status: ${messageRes.status}`
      );
      releaseTfaLock();
      return null;
    }

    try {
      console.log(`DEBUG: Attempting to parse message JSON response`);
      console.log(
        `DEBUG: First 200 chars of message response: ${messageRes.body.substring(
          0,
          200
        )}...`
      );

      const messageData = JSON.parse(messageRes.body);
      console.log(`DEBUG: Message JSON parsed successfully`);

      // Check if we have HTML content
      if (!messageData.HTML) {
        console.log(
          `DEBUG: No HTML content found in message. Available fields: ${Object.keys(
            messageData
          ).join(", ")}`
        );
      }

      const htmlContent = messageData.HTML || "";
      console.log(
        `DEBUG: HTML content length: ${htmlContent.length} characters`
      );

      // Step 3: Extract TFA code using regex pattern
      // Log a portion of the HTML content for debugging
      console.log(
        `DEBUG: Email HTML content for ${email} (first 200 chars): ${htmlContent.substring(
          0,
          200
        )}...`
      );

      // The pattern looks for a 6-digit code that appears after certain text patterns
      // Using the same pattern as in the Go test helpers for consistency
      const tfaCodePattern =
        /verification code is[^\d]*(\d{6})|code:[^\d]*(\d{6})|code is[^\d]*(\d{6})|code.*?(\d{6})|\b(\d{6})\b/i;
      console.log(`DEBUG: Using TFA code regex pattern: ${tfaCodePattern}`);
      const tfaMatch = htmlContent.match(tfaCodePattern);
      console.log(
        `DEBUG: Regex match result: ${tfaMatch ? "Match found" : "No match"}`
      );
      if (tfaMatch) {
        console.log(`DEBUG: Raw TFA match groups: ${JSON.stringify(tfaMatch)}`);
      }

      if (tfaMatch) {
        // The code could be in any of the capture groups, find the non-null one
        const tfaCode =
          tfaMatch[1] ||
          tfaMatch[2] ||
          tfaMatch[3] ||
          tfaMatch[4] ||
          tfaMatch[5];
        console.log(`Successfully extracted TFA code for ${email}: ${tfaCode}`);
        releaseTfaLock();
        return tfaCode;
      } else {
        console.error(
          `TFA code pattern not found in email content for ${email}`
        );
        releaseTfaLock();
        return null;
      }
    } catch (e) {
      console.error(
        `Error parsing message response: ${e}. Body: ${messageRes.body}`
      );
      releaseTfaLock();
      return null;
    }
  } catch (error) {
    console.error(`Unexpected error in fetchTFACodeForUser: ${error}`);
    releaseTfaLock();
    return null;
  }
}

// --- Authentication Function (Accepts user object with email/handle) ---
async function loginAndAuthenticateUser(user) {
  let loginAttempts = 0;
  let authToken = null;

  // Step 1: Login to get TFA token
  while (loginAttempts < MAX_LOGIN_ATTEMPTS && !authToken) {
    loginAttempts++;
    console.debug(
      `Login attempt ${loginAttempts}/${MAX_LOGIN_ATTEMPTS} for ${user.email}`
    );

    // Prepare login payload
    const loginPayload = JSON.stringify({
      email: user.email,
      password: PASSWORD,
    });

    // Make login request
    const loginRes = http.post(`${API_BASE_URL}/hub/login`, loginPayload, {
      headers: { "Content-Type": "application/json" },
      tags: { name: "HubLoginAPI" },
    });

    // Check login response
    if (loginRes.status !== 200) {
      console.error(
        `Login failed for ${user.email}. Status: ${loginRes.status}, Body: ${loginRes.body}`
      );
      continue;
    }

    console.debug(`Login pass for ${user.email} Now for TFA`);

    try {
      const loginData = JSON.parse(loginRes.body);

      // Extract TFA token from response - using 'token' field name as in the Go code
      const tfaToken = loginData.token;
      if (!tfaToken) {
        console.error(
          `TFA token not found in login response for ${user.email}. Response: ${loginRes.body}`
        );
        return;
      }

      console.debug(`Got TFA token for ${user.email}: ${tfaToken}`);

      // Step 2: Fetch TFA code from email
      const tfaCode = await fetchTFACodeForUser(user.email);
      if (!tfaCode) {
        console.error(`Failed to get TFA code for ${user.email}`);
        sleep(2);
        continue;
      }

      // Step 3: Submit TFA code to complete authentication
      // Using the field names from the JSON tags in the Go struct
      const tfaPayload = JSON.stringify({
        tfa_token: tfaToken,
        tfa_code: tfaCode,
      });

      const tfaRes = http.post(`${API_BASE_URL}/hub/tfa`, tfaPayload, {
        headers: { "Content-Type": "application/json" },
        tags: { name: "HubTFAAPI" },
      });

      // Check TFA response
      if (tfaRes.status !== 200) {
        console.error(
          `TFA verification failed for ${user.email}. Status: ${tfaRes.status}, Body: ${tfaRes.body}`
        );
        return;
      }

      try {
        const tfaData = JSON.parse(tfaRes.body);
        // Using session_token as per the JSON tag in the Go struct
        authToken = tfaData.session_token;

        if (!authToken) {
          console.error(
            `Auth token (session_token) not found in TFA response for ${user.email}. Response: ${tfaRes.body}`
          );
          return;
        }

        console.debug(`Successfully authenticated ${user.email}`);
        return authToken;
      } catch (e) {
        console.error(
          `Error parsing TFA response for ${user.email}: ${e}. Body: ${tfaRes.body}`
        );
        return;
      }
    } catch (e) {
      console.error(
        `Error parsing login response for ${user.email}: ${e}. Body: ${loginRes.body}`
      );
      return;
    }
  }

  return null;
}

// --- k6 Setup Function ---
export async function setup() {
  console.error("!!!!!!!!!! MARKER: VERSION 2025-05-04-12:55 !!!!!!!!!!");
  console.log("!!!!!!!!!! MARKER: VERSION 2025-05-04-12:55 !!!!!!!!!!");
  console.log("SETUP FUNCTION CALLED - THIS SHOULD BE VISIBLE");
  console.log(`API_BASE_URL: ${API_BASE_URL}`);
  console.log(`MAILPIT_URL: ${MAILPIT_URL}`);
  console.log(`TOTAL_USERS: ${TOTAL_USERS}`);
  console.log(`INSTANCE_INDEX: ${INSTANCE_INDEX}`);
  console.log(`INSTANCE_COUNT: ${INSTANCE_COUNT}`);
  console.log(`USERS_PER_INSTANCE: ${USERS_PER_INSTANCE}`);
  console.log(`SETUP_PARALLELISM: ${SETUP_PARALLELISM}`);
  console.log(`TEST_DURATION_SECONDS: ${TEST_DURATION_SECONDS}`);

  console.log(
    `Starting setup for instance ${INSTANCE_INDEX} with users ${startUserIndex} to ${endUserIndex}`
  );

  // Test direct Mailpit connectivity
  console.log(`Testing connectivity to Mailpit: ${MAILPIT_URL}`);
  try {
    const mailpitRes = http.get(`${MAILPIT_URL}/api/v1/messages`);
    console.log(`Mailpit connectivity test status: ${mailpitRes.status}`);
  } catch (error) {
    console.error(`Mailpit connectivity test failed: ${error}`);
    return;
  }

  console.log("Delete all existing mails in mailpit");
  try {
    // Use http.del() instead of http.delete() - k6 uses del() for DELETE requests
    const deleteRes = http.del(`${MAILPIT_URL}/api/v1/messages`);
    console.log(`Mailpit delete test status: ${deleteRes.status}`);
    if (deleteRes.status !== 200) {
      console.error(`Mailpit delete failed with status: ${deleteRes.status}`);
      return;
    }
  } catch (error) {
    console.error(`Mailpit delete test failed: ${error}`);
    return;
  }

  // Prepare user data structure
  const users = [];
  const allUserHandles = [];

  // Create user objects for this instance's range
  for (let i = startUserIndex; i <= endUserIndex; i++) {
    const email = `user${i}@example.com`;
    const handle = `user${i}`;
    users.push({ email, handle });
    allUserHandles.push(handle);
  }

  console.log(`Created ${users.length} user objects for authentication`);

  // Authenticate users in parallel batches
  const authenticatedUsers = [];
  // Use a reasonable batch size based on available resources
  const batchSize = Math.min(SETUP_PARALLELISM, 20); // Increased to 20 parallel authentications
  const batches = Math.ceil(users.length / batchSize);

  for (let batchIndex = 0; batchIndex < batches; batchIndex++) {
    const batchStart = batchIndex * batchSize;
    const batchEnd = Math.min((batchIndex + 1) * batchSize, users.length);
    const batchUsers = users.slice(batchStart, batchEnd);

    console.log(
      `Processing authentication batch ${batchIndex + 1}/${batches} (users ${
        batchStart + 1
      } to ${batchEnd})`
    );

    // Process batch in parallel
    const authPromises = batchUsers.map(async (user) => {
      const authToken = await loginAndAuthenticateUser(user);
      if (authToken) {
        return { ...user, authToken };
      }
      return null;
    });

    // Wait for all authentications in this batch to complete
    const batchResults = await Promise.all(authPromises);

    // Filter out failed authentications
    const successfulAuths = batchResults.filter((result) => result !== null);
    authenticatedUsers.push(...successfulAuths);

    console.log(
      `Batch ${batchIndex + 1} complete: ${successfulAuths.length}/${
        batchUsers.length
      } users authenticated`
    );
  }

  console.log(
    `Setup complete: ${authenticatedUsers.length}/${users.length} users authenticated`
  );

  // Return the authenticated users and all handles for the test
  return {
    authenticatedUsers,
    allUserHandles,
  };
}

// --- Social Activity Function ---
function socialActivity(authToken, userHandle, allUserHandles, vuState) {
  // Common auth parameters for all requests
  const authParams = {
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${authToken}`,
    },
  };

  // Define all possible social actions
  const socialActions = [
    "follow_user",
    "unfollow_user",
    "create_post",
    "read_timeline",
    "view_post_details",
    "upvote",
    "downvote",
    "unvote",
    "create_incognito_post",
    "browse_incognito_posts",
    "view_incognito_post_details",
    "my_incognito_posts",
    "add_incognito_comment",
    "upvote_incognito_post",
    "downvote_incognito_post",
    "unvote_incognito_post",
    "upvote_incognito_comment",
    "downvote_incognito_comment",
    "unvote_incognito_comment",
  ];

  // Weight the actions to create a realistic distribution
  // Reading timeline and viewing posts should be more common than posting
  const actionWeights = [
    10, 3, 8, 25, 15, 7, 3, 3, 5, 10, 5, 3, 3, 2, 1, 1, 1, 1, 1,
  ]; // Percentages

  // Calculate cumulative weights
  const cumulativeWeights = [];
  let sum = 0;
  for (const weight of actionWeights) {
    sum += weight;
    cumulativeWeights.push(sum);
  }

  // Function to select a weighted random action
  function selectRandomAction() {
    const randomValue = Math.random() * 100;
    let selectedActionIndex = 0;

    for (let i = 0; i < cumulativeWeights.length; i++) {
      if (randomValue <= cumulativeWeights[i]) {
        selectedActionIndex = i;
        break;
      }
    }
    return socialActions[selectedActionIndex];
  }

  // Function to perform a specific action
  function performAction(action) {
    // Execute the selected social action
    group(`Social Action: ${action}`, () => {
      switch (action) {
        case "follow_user": {
          // Select a random user to follow (not self)
          const otherHandles = allUserHandles.filter((h) => h !== userHandle);
          if (otherHandles.length === 0) {
            console.debug(
              `VU ${__VU} (${userHandle}): No other users to follow.`
            );
            break;
          }

          // Don't follow users we already follow
          const potentialHandlesToFollow = otherHandles.filter(
            (h) => !vuState.followedUsers.includes(h)
          );

          if (potentialHandlesToFollow.length === 0) {
            console.debug(
              `VU ${__VU} (${userHandle}): Already following all available users.`
            );
            break;
          }

          const handleToFollow = randomItem(potentialHandlesToFollow);
          const followPayload = JSON.stringify({ handle: handleToFollow });

          console.debug(
            `VU ${__VU} (${userHandle}): Attempting to follow user ${handleToFollow}`
          );

          const followRes = http.post(
            `${API_BASE_URL}/hub/follow-user`,
            followPayload,
            {
              ...authParams,
              tags: { name: "HubFollowUserAPI" },
            }
          );

          followTrend.add(followRes.timings.duration);

          // Check for success
          check(followRes, {
            "Follow request successful": (r) => r.status === 200,
          });

          // If successful, add to followed users list
          if (followRes.status === 200) {
            vuState.followedUsers.push(handleToFollow);
          }

          // Log unexpected errors
          if (followRes.status !== 200) {
            console.error(
              `VU ${__VU} (${userHandle}): Follow API Error! Handle: ${handleToFollow}, Status: ${followRes.status}, Body: ${followRes.body}`
            );
          }

          sleep(randomIntBetween(1, 3));
          break;
        }

        case "unfollow_user": {
          // Check if we have any users to unfollow
          if (vuState.followedUsers.length === 0) {
            console.debug(`VU ${__VU} (${userHandle}): No users to unfollow.`);
            break;
          }

          const handleToUnfollow = randomItem(vuState.followedUsers);
          const unfollowPayload = JSON.stringify({ handle: handleToUnfollow });

          console.debug(
            `VU ${__VU} (${userHandle}): Attempting to unfollow user ${handleToUnfollow}`
          );

          const unfollowRes = http.post(
            `${API_BASE_URL}/hub/unfollow-user`,
            unfollowPayload,
            {
              ...authParams,
              tags: { name: "HubUnfollowUserAPI" },
            }
          );

          unfollowTrend.add(unfollowRes.timings.duration);

          // Check for success
          check(unfollowRes, {
            "Unfollow request successful": (r) => r.status === 200,
          });

          // If successful, remove from followed users list
          if (unfollowRes.status === 200) {
            vuState.followedUsers = vuState.followedUsers.filter(
              (h) => h !== handleToUnfollow
            );
          }

          // Log unexpected errors
          if (unfollowRes.status !== 200) {
            console.error(
              `VU ${__VU} (${userHandle}): Unfollow API Error! Handle: ${handleToUnfollow}, Status: ${unfollowRes.status}, Body: ${unfollowRes.body}`
            );
          }

          sleep(randomIntBetween(1, 3));
          break;
        }

        case "create_post": {
          // Generate random post content
          const postContent = `Test post from ${userHandle} at ${new Date().toISOString()} - ${randomString(
            20
          )}`;
          const postPayload = JSON.stringify({ content: postContent });

          console.debug(`VU ${__VU} (${userHandle}): Creating new post`);

          const createPostRes = http.post(
            `${API_BASE_URL}/hub/add-post`,
            postPayload,
            {
              ...authParams,
              tags: { name: "HubCreatePostAPI" },
            }
          );

          createPostTrend.add(createPostRes.timings.duration);

          // Check for success
          check(createPostRes, {
            "Create post request successful": (r) => r.status === 200,
          });

          // If successful, extract post ID and add to created posts
          if (createPostRes.status === 200) {
            try {
              const postData = JSON.parse(createPostRes.body);
              if (postData.post_id) {
                vuState.createdPostIds.push(postData.post_id);
                console.debug(
                  `VU ${__VU} (${userHandle}): Created post with ID ${postData.post_id}`
                );
              }
            } catch (e) {
              console.error(
                `VU ${__VU} (${userHandle}): Error parsing create post response: ${e}`
              );
            }
          }

          // Log unexpected errors
          if (createPostRes.status !== 200) {
            console.error(
              `VU ${__VU} (${userHandle}): Create Post API Error! Status: ${createPostRes.status}, Body: ${createPostRes.body}`
            );
          }

          sleep(randomIntBetween(2, 5));
          break;
        }

        case "read_timeline": {
          // Prepare request payload for timeline
          const timelinePayload = JSON.stringify({
            limit: 20,
            pagination_key: vuState.timelineCursor || undefined,
          });

          console.debug(
            `VU ${__VU} (${userHandle}): Reading timeline${
              vuState.timelineCursor ? " with cursor" : ""
            }`
          );

          const timelineRes = http.post(
            `${API_BASE_URL}/hub/get-my-home-timeline`,
            timelinePayload,
            {
              ...authParams,
              tags: { name: "HubTimelineAPI" },
            }
          );

          timelineReadTrend.add(timelineRes.timings.duration);

          // Check for success
          check(timelineRes, {
            "Timeline read successful": (r) => r.status === 200,
          });

          // Process timeline response
          if (timelineRes.status === 200) {
            try {
              const timelineData = JSON.parse(timelineRes.body);

              // Update cursor for next pagination
              if (timelineData.pagination_key) {
                vuState.timelineCursor = timelineData.pagination_key;
              }

              // Store post IDs for potential interactions
              if (timelineData.posts && timelineData.posts.length > 0) {
                // Add new post IDs to the timeline posts collection
                timelineData.posts.forEach((post) => {
                  if (!vuState.timelinePostIds.includes(post.id)) {
                    vuState.timelinePostIds.push(post.id);
                  }
                });

                console.debug(
                  `VU ${__VU} (${userHandle}): Retrieved ${timelineData.posts.length} posts, total in memory: ${vuState.timelinePostIds.length}`
                );
              } else {
                console.debug(
                  `VU ${__VU} (${userHandle}): No posts in timeline response`
                );
              }
            } catch (e) {
              console.error(
                `VU ${__VU} (${userHandle}): Error parsing timeline response: ${e}`
              );
            }
          }

          // Log unexpected errors
          if (timelineRes.status !== 200) {
            console.error(
              `VU ${__VU} (${userHandle}): Timeline API Error! Status: ${timelineRes.status}, Body: ${timelineRes.body}`
            );
          }

          sleep(randomIntBetween(2, 5));
          break;
        }

        case "view_post_details": {
          // Combine timeline posts and created posts for selection
          const availablePosts = [
            ...vuState.timelinePostIds,
            ...vuState.createdPostIds,
          ];

          if (availablePosts.length === 0) {
            console.debug(
              `VU ${__VU} (${userHandle}): No posts available to view details.`
            );
            break;
          }

          const postIdToView = randomItem(availablePosts);
          const postDetailsPayload = JSON.stringify({
            post_id: postIdToView,
          });

          console.debug(
            `VU ${__VU} (${userHandle}): Viewing details for post ${postIdToView}`
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
            "Post details request successful": (r) => r.status === 200,
          });

          // Log unexpected errors
          if (postDetailsRes.status !== 200) {
            console.error(
              `VU ${__VU} (${userHandle}): Post Details API Error! PostID: ${postIdToView}, Status: ${postDetailsRes.status}, Body: ${postDetailsRes.body}`
            );
          }

          sleep(randomIntBetween(1, 4));
          break;
        }

        case "upvote":
        case "downvote": {
          // Combine timeline posts for voting (exclude own posts)
          const availablePosts = vuState.timelinePostIds.filter(
            (id) => !vuState.createdPostIds.includes(id)
          );

          // Also exclude posts we've already voted on
          const votablePosts = availablePosts.filter(
            (id) =>
              !vuState.upvotedPostIds.includes(id) &&
              !vuState.downvotedPostIds.includes(id)
          );

          if (votablePosts.length === 0) {
            console.debug(
              `VU ${__VU} (${userHandle}): No posts available to ${action}.`
            );
            break;
          }

          const postIdToVote = randomItem(votablePosts);
          const votePayload = JSON.stringify({ post_id: postIdToVote });
          const voteUrl = `${API_BASE_URL}/hub/${
            action === "upvote" ? "upvote" : "downvote"
          }-user-post`;

          console.debug(
            `VU ${__VU} (${userHandle}): Attempting to ${action} post ${postIdToVote}`
          );

          // Set appropriate trend and tag based on vote type
          let voteTrend, voteTag;
          if (action === "upvote") {
            voteTrend = upvoteTrend;
            voteTag = "HubUpvoteAPI";
          } else {
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
            [`${action} request successful or expected error`]: (r) =>
              r.status === 200 || r.status === 422,
          });

          // If successful, track the voted post
          if (voteRes.status === 200) {
            if (action === "upvote") {
              vuState.upvotedPostIds.push(postIdToVote);
            } else {
              vuState.downvotedPostIds.push(postIdToVote);
            }
          }

          // Log unexpected errors
          if (voteRes.status !== 200 && voteRes.status !== 422) {
            console.error(
              `VU ${__VU} (${userHandle}): ${action} API Error! PostID: ${postIdToVote}, Status: ${voteRes.status}, Body: ${voteRes.body}`
            );
          }

          sleep(randomIntBetween(1, 4));
          break;
        }

        case "unvote": {
          // Combine upvoted and downvoted posts to pick from
          const votedPosts = [
            ...vuState.upvotedPostIds,
            ...vuState.downvotedPostIds,
          ];

          if (votedPosts.length === 0) {
            console.debug(
              `VU ${__VU} (${userHandle}): No voted posts available to unvote.`
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
            "Unvote request successful or expected error": (r) =>
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
              `VU ${__VU} (${userHandle}): Unvote API Error! PostID: ${postIdToUnvote}, Status: ${unvoteRes.status}, Body: ${unvoteRes.body}`
            );
          }

          sleep(randomIntBetween(1, 4));
          break;
        }

        case "create_incognito_post": {
          // Generate random incognito post content
          const incognitoContent = `Anonymous thought from ${new Date()
            .toISOString()
            .slice(0, 10)} - ${randomString(30)}`;

          // Select 1-3 random tags
          const numTags = randomIntBetween(1, 3);
          const selectedTags = [];
          for (let i = 0; i < numTags; i++) {
            const tag = randomItem(AVAILABLE_TAG_IDS);
            if (!selectedTags.includes(tag)) {
              selectedTags.push(tag);
            }
          }

          const incognitoPostPayload = JSON.stringify({
            content: incognitoContent,
            tag_ids: selectedTags,
          });

          console.debug(
            `VU ${__VU} (${userHandle}): Creating incognito post with tags ${selectedTags.join(
              ", "
            )}`
          );

          const createIncognitoPostRes = http.post(
            `${API_BASE_URL}/hub/add-incognito-post`,
            incognitoPostPayload,
            {
              ...authParams,
              tags: { name: "HubCreateIncognitoPostAPI" },
            }
          );

          createIncognitoPostTrend.add(createIncognitoPostRes.timings.duration);

          // Check for success
          check(createIncognitoPostRes, {
            "Create incognito post request successful": (r) => r.status === 200,
          });

          // If successful, extract incognito post ID and add to created posts
          if (createIncognitoPostRes.status === 200) {
            try {
              const incognitoPostData = JSON.parse(createIncognitoPostRes.body);
              if (incognitoPostData.incognito_post_id) {
                vuState.createdIncognitoPostIds.push(
                  incognitoPostData.incognito_post_id
                );
                console.debug(
                  `VU ${__VU} (${userHandle}): Created incognito post with ID ${incognitoPostData.incognito_post_id}`
                );
              }
            } catch (e) {
              console.error(
                `VU ${__VU} (${userHandle}): Error parsing create incognito post response: ${e}`
              );
            }
          }

          // Log unexpected errors
          if (createIncognitoPostRes.status !== 200) {
            console.error(
              `VU ${__VU} (${userHandle}): Create Incognito Post API Error! Status: ${createIncognitoPostRes.status}, Body: ${createIncognitoPostRes.body}`
            );
          }

          sleep(randomIntBetween(2, 5));
          break;
        }

        case "browse_incognito_posts": {
          // Select a random tag for browsing
          const browseTag = randomItem(AVAILABLE_TAG_IDS);

          const browseIncognitoPayload = JSON.stringify({
            tag_id: browseTag,
            time_filter: "past_24_hours",
            limit: 25,
            pagination_key: vuState.incognitoPostsCursor || undefined,
          });

          console.debug(
            `VU ${__VU} (${userHandle}): Browsing incognito posts with tag ${browseTag}${
              vuState.incognitoPostsCursor ? " with cursor" : ""
            }`
          );

          const browseIncognitoRes = http.get(
            `${API_BASE_URL}/hub/get-incognito-posts`,
            {
              body: browseIncognitoPayload,
              ...authParams,
              tags: { name: "HubBrowseIncognitoPostsAPI" },
            }
          );

          getIncognitoPostsTrend.add(browseIncognitoRes.timings.duration);

          // Check for success
          check(browseIncognitoRes, {
            "Browse incognito posts successful": (r) => r.status === 200,
          });

          // Process browse response
          if (browseIncognitoRes.status === 200) {
            try {
              const browseData = JSON.parse(browseIncognitoRes.body);

              // Update cursor for next pagination
              if (browseData.pagination_key) {
                vuState.incognitoPostsCursor = browseData.pagination_key;
              }

              // Store incognito post IDs for potential interactions
              if (browseData.posts && browseData.posts.length > 0) {
                browseData.posts.forEach((post) => {
                  if (
                    !vuState.browseIncognitoPostIds.includes(
                      post.incognito_post_id
                    )
                  ) {
                    vuState.browseIncognitoPostIds.push(post.incognito_post_id);
                  }
                });

                console.debug(
                  `VU ${__VU} (${userHandle}): Retrieved ${browseData.posts.length} incognito posts, total in memory: ${vuState.browseIncognitoPostIds.length}`
                );
              } else {
                console.debug(
                  `VU ${__VU} (${userHandle}): No incognito posts in browse response`
                );
              }
            } catch (e) {
              console.error(
                `VU ${__VU} (${userHandle}): Error parsing browse incognito posts response: ${e}`
              );
            }
          }

          // Log unexpected errors
          if (browseIncognitoRes.status !== 200) {
            console.error(
              `VU ${__VU} (${userHandle}): Browse Incognito Posts API Error! Status: ${browseIncognitoRes.status}, Body: ${browseIncognitoRes.body}`
            );
          }

          sleep(randomIntBetween(2, 5));
          break;
        }

        case "view_incognito_post_details": {
          // Combine browsed incognito posts and created incognito posts for selection
          const availableIncognitoPosts = [
            ...vuState.browseIncognitoPostIds,
            ...vuState.createdIncognitoPostIds,
          ];

          if (availableIncognitoPosts.length === 0) {
            console.debug(
              `VU ${__VU} (${userHandle}): No incognito posts available to view details.`
            );
            break;
          }

          const incognitoPostIdToView = randomItem(availableIncognitoPosts);
          const incognitoPostDetailsPayload = JSON.stringify({
            incognito_post_id: incognitoPostIdToView,
          });

          console.debug(
            `VU ${__VU} (${userHandle}): Viewing details for incognito post ${incognitoPostIdToView}`
          );

          const incognitoPostDetailsRes = http.get(
            `${API_BASE_URL}/hub/get-incognito-post`,
            {
              body: incognitoPostDetailsPayload,
              ...authParams,
              tags: { name: "HubIncognitoPostDetailsAPI" },
            }
          );

          getIncognitoPostTrend.add(incognitoPostDetailsRes.timings.duration);

          // Check for success
          check(incognitoPostDetailsRes, {
            "Incognito post details request successful": (r) =>
              r.status === 200,
          });

          // Log unexpected errors
          if (incognitoPostDetailsRes.status !== 200) {
            console.error(
              `VU ${__VU} (${userHandle}): Incognito Post Details API Error! PostID: ${incognitoPostIdToView}, Status: ${incognitoPostDetailsRes.status}, Body: ${incognitoPostDetailsRes.body}`
            );
          }

          sleep(randomIntBetween(1, 4));
          break;
        }

        case "my_incognito_posts": {
          // Prepare request payload for my incognito posts
          const myIncognitoPostsPayload = JSON.stringify({
            limit: 25,
            pagination_key: vuState.myIncognitoPostsCursor || undefined,
          });

          console.debug(
            `VU ${__VU} (${userHandle}): Viewing my incognito posts${
              vuState.myIncognitoPostsCursor ? " with cursor" : ""
            }`
          );

          const myIncognitoPostsRes = http.get(
            `${API_BASE_URL}/hub/get-my-incognito-posts`,
            {
              body: myIncognitoPostsPayload,
              ...authParams,
              tags: { name: "HubMyIncognitoPostsAPI" },
            }
          );

          getMyIncognitoPostsTrend.add(myIncognitoPostsRes.timings.duration);

          // Check for success
          check(myIncognitoPostsRes, {
            "My incognito posts read successful": (r) => r.status === 200,
          });

          // Process my incognito posts response
          if (myIncognitoPostsRes.status === 200) {
            try {
              const myIncognitoPostsData = JSON.parse(myIncognitoPostsRes.body);

              // Update cursor for next pagination
              if (myIncognitoPostsData.pagination_key) {
                vuState.myIncognitoPostsCursor =
                  myIncognitoPostsData.pagination_key;
              }

              console.debug(
                `VU ${__VU} (${userHandle}): Retrieved ${
                  myIncognitoPostsData.posts
                    ? myIncognitoPostsData.posts.length
                    : 0
                } of my incognito posts`
              );
            } catch (e) {
              console.error(
                `VU ${__VU} (${userHandle}): Error parsing my incognito posts response: ${e}`
              );
            }
          }

          // Log unexpected errors
          if (myIncognitoPostsRes.status !== 200) {
            console.error(
              `VU ${__VU} (${userHandle}): My Incognito Posts API Error! Status: ${myIncognitoPostsRes.status}, Body: ${myIncognitoPostsRes.body}`
            );
          }

          sleep(randomIntBetween(2, 5));
          break;
        }

        case "add_incognito_comment": {
          // Combine browsed incognito posts for commenting
          const availableIncognitoPosts = [
            ...vuState.browseIncognitoPostIds,
            ...vuState.createdIncognitoPostIds,
          ];

          if (availableIncognitoPosts.length === 0) {
            console.debug(
              `VU ${__VU} (${userHandle}): No incognito posts available to comment on.`
            );
            break;
          }

          const incognitoPostIdToComment = randomItem(availableIncognitoPosts);

          // Generate random comment content
          const commentContent = `Anonymous comment ${randomString(
            20
          )} on ${new Date().toISOString().slice(0, 10)}`;

          const addIncognitoCommentPayload = JSON.stringify({
            incognito_post_id: incognitoPostIdToComment,
            content: commentContent,
            // in_reply_to can be added later for nested comments
          });

          console.debug(
            `VU ${__VU} (${userHandle}): Adding comment to incognito post ${incognitoPostIdToComment}`
          );

          const addIncognitoCommentRes = http.post(
            `${API_BASE_URL}/hub/add-incognito-post-comment`,
            addIncognitoCommentPayload,
            {
              ...authParams,
              tags: { name: "HubAddIncognitoCommentAPI" },
            }
          );

          addIncognitoCommentTrend.add(addIncognitoCommentRes.timings.duration);

          // Check for success
          check(addIncognitoCommentRes, {
            "Add incognito comment request successful": (r) => r.status === 200,
          });

          // If successful, extract comment ID and add to created comments
          if (addIncognitoCommentRes.status === 200) {
            try {
              const commentData = JSON.parse(addIncognitoCommentRes.body);
              if (commentData.comment_id) {
                vuState.incognitoCommentsIds.push(commentData.comment_id);
                console.debug(
                  `VU ${__VU} (${userHandle}): Created incognito comment with ID ${commentData.comment_id}`
                );
              }
            } catch (e) {
              console.error(
                `VU ${__VU} (${userHandle}): Error parsing add incognito comment response: ${e}`
              );
            }
          }

          // Log unexpected errors
          if (addIncognitoCommentRes.status !== 200) {
            console.error(
              `VU ${__VU} (${userHandle}): Add Incognito Comment API Error! PostID: ${incognitoPostIdToComment}, Status: ${addIncognitoCommentRes.status}, Body: ${addIncognitoCommentRes.body}`
            );
          }

          sleep(randomIntBetween(2, 5));
          break;
        }

        case "upvote_incognito_post": {
          // Combine browsed incognito posts for voting (exclude own posts)
          const availableIncognitoPosts = vuState.browseIncognitoPostIds.filter(
            (id) => !vuState.createdIncognitoPostIds.includes(id)
          );

          // Also exclude posts we've already voted on
          const votableIncognitoPosts = availableIncognitoPosts.filter(
            (id) =>
              !vuState.upvotedIncognitoPostIds.includes(id) &&
              !vuState.downvotedIncognitoPostIds.includes(id)
          );

          if (votableIncognitoPosts.length === 0) {
            console.debug(
              `VU ${__VU} (${userHandle}): No incognito posts available to upvote.`
            );
            break;
          }

          const incognitoPostIdToUpvote = randomItem(votableIncognitoPosts);
          const upvoteIncognitoPayload = JSON.stringify({
            incognito_post_id: incognitoPostIdToUpvote,
          });

          console.debug(
            `VU ${__VU} (${userHandle}): Attempting to upvote incognito post ${incognitoPostIdToUpvote}`
          );

          const upvoteIncognitoRes = http.post(
            `${API_BASE_URL}/hub/upvote-incognito-post`,
            upvoteIncognitoPayload,
            {
              ...authParams,
              tags: { name: "HubUpvoteIncognitoPostAPI" },
            }
          );

          upvoteIncognitoPostTrend.add(upvoteIncognitoRes.timings.duration);

          // Check for success or expected error
          check(upvoteIncognitoRes, {
            "Upvote incognito post request successful or expected error": (r) =>
              r.status === 200 || r.status === 422,
          });

          // If successful, track the voted post
          if (upvoteIncognitoRes.status === 200) {
            vuState.upvotedIncognitoPostIds.push(incognitoPostIdToUpvote);
          }

          // Log unexpected errors
          if (
            upvoteIncognitoRes.status !== 200 &&
            upvoteIncognitoRes.status !== 422
          ) {
            console.error(
              `VU ${__VU} (${userHandle}): Upvote Incognito Post API Error! PostID: ${incognitoPostIdToUpvote}, Status: ${upvoteIncognitoRes.status}, Body: ${upvoteIncognitoRes.body}`
            );
          }

          sleep(randomIntBetween(1, 4));
          break;
        }

        case "downvote_incognito_post": {
          // Similar logic to upvote but for downvoting
          const availableIncognitoPosts = vuState.browseIncognitoPostIds.filter(
            (id) => !vuState.createdIncognitoPostIds.includes(id)
          );

          const votableIncognitoPosts = availableIncognitoPosts.filter(
            (id) =>
              !vuState.upvotedIncognitoPostIds.includes(id) &&
              !vuState.downvotedIncognitoPostIds.includes(id)
          );

          if (votableIncognitoPosts.length === 0) {
            console.debug(
              `VU ${__VU} (${userHandle}): No incognito posts available to downvote.`
            );
            break;
          }

          const incognitoPostIdToDownvote = randomItem(votableIncognitoPosts);
          const downvoteIncognitoPayload = JSON.stringify({
            incognito_post_id: incognitoPostIdToDownvote,
          });

          console.debug(
            `VU ${__VU} (${userHandle}): Attempting to downvote incognito post ${incognitoPostIdToDownvote}`
          );

          const downvoteIncognitoRes = http.post(
            `${API_BASE_URL}/hub/downvote-incognito-post`,
            downvoteIncognitoPayload,
            {
              ...authParams,
              tags: { name: "HubDownvoteIncognitoPostAPI" },
            }
          );

          downvoteIncognitoPostTrend.add(downvoteIncognitoRes.timings.duration);

          check(downvoteIncognitoRes, {
            "Downvote incognito post request successful or expected error": (
              r
            ) => r.status === 200 || r.status === 422,
          });

          if (downvoteIncognitoRes.status === 200) {
            vuState.downvotedIncognitoPostIds.push(incognitoPostIdToDownvote);
          }

          if (
            downvoteIncognitoRes.status !== 200 &&
            downvoteIncognitoRes.status !== 422
          ) {
            console.error(
              `VU ${__VU} (${userHandle}): Downvote Incognito Post API Error! PostID: ${incognitoPostIdToDownvote}, Status: ${downvoteIncognitoRes.status}, Body: ${downvoteIncognitoRes.body}`
            );
          }

          sleep(randomIntBetween(1, 4));
          break;
        }

        case "unvote_incognito_post": {
          // Combine upvoted and downvoted incognito posts to pick from
          const votedIncognitoPosts = [
            ...vuState.upvotedIncognitoPostIds,
            ...vuState.downvotedIncognitoPostIds,
          ];

          if (votedIncognitoPosts.length === 0) {
            console.debug(
              `VU ${__VU} (${userHandle}): No voted incognito posts available to unvote.`
            );
            break;
          }

          const incognitoPostIdToUnvote = randomItem(votedIncognitoPosts);
          const unvoteIncognitoPayload = JSON.stringify({
            incognito_post_id: incognitoPostIdToUnvote,
          });

          console.debug(
            `VU ${__VU} (${userHandle}): Attempting to unvote incognito post ${incognitoPostIdToUnvote}`
          );

          const unvoteIncognitoRes = http.post(
            `${API_BASE_URL}/hub/unvote-incognito-post`,
            unvoteIncognitoPayload,
            {
              ...authParams,
              tags: { name: "HubUnvoteIncognitoPostAPI" },
            }
          );

          unvoteIncognitoPostTrend.add(unvoteIncognitoRes.timings.duration);

          check(unvoteIncognitoRes, {
            "Unvote incognito post request successful or expected error": (r) =>
              r.status === 200 || r.status === 422,
          });

          if (unvoteIncognitoRes.status === 200) {
            vuState.upvotedIncognitoPostIds =
              vuState.upvotedIncognitoPostIds.filter(
                (id) => id !== incognitoPostIdToUnvote
              );
            vuState.downvotedIncognitoPostIds =
              vuState.downvotedIncognitoPostIds.filter(
                (id) => id !== incognitoPostIdToUnvote
              );
          }

          if (
            unvoteIncognitoRes.status !== 200 &&
            unvoteIncognitoRes.status !== 422
          ) {
            console.error(
              `VU ${__VU} (${userHandle}): Unvote Incognito Post API Error! PostID: ${incognitoPostIdToUnvote}, Status: ${unvoteIncognitoRes.status}, Body: ${unvoteIncognitoRes.body}`
            );
          }

          sleep(randomIntBetween(1, 4));
          break;
        }

        case "upvote_incognito_comment": {
          // Simplified comment voting - would need comment IDs from get-comments API in practice
          console.debug(
            `VU ${__VU} (${userHandle}): Simulating upvote incognito comment (placeholder)`
          );
          sleep(randomIntBetween(1, 4));
          break;
        }

        case "downvote_incognito_comment": {
          // Simplified comment voting - would need comment IDs from get-comments API in practice
          console.debug(
            `VU ${__VU} (${userHandle}): Simulating downvote incognito comment (placeholder)`
          );
          sleep(randomIntBetween(1, 4));
          break;
        }

        case "unvote_incognito_comment": {
          // Simplified comment voting - would need comment IDs from get-comments API in practice
          console.debug(
            `VU ${__VU} (${userHandle}): Simulating unvote incognito comment (placeholder)`
          );
          sleep(randomIntBetween(1, 4));
          break;
        }
      }

      // Think time between actions
      sleep(randomIntBetween(2, 6));
    });
  }

  // Perform a realistic user session with multiple actions
  // First, always read the timeline as the initial action
  performAction("read_timeline");

  // Then perform 5-10 random actions to simulate a realistic session
  const numberOfActions = Math.floor(Math.random() * 6) + 5; // 5-10 actions
  console.log(
    `VU ${__VU} (${userHandle}): Starting session with ${numberOfActions} actions`
  );

  for (let i = 0; i < numberOfActions; i++) {
    const action = selectRandomAction();
    performAction(action);
  }
}

// --- Main Test Logic (Accepts setup data) ---
export default function (data) {
  // Initialize VU state for tracking social interactions - with smaller arrays to reduce memory usage
  const vuState = {
    // Timeline pagination
    timelineCursor: null,

    // Posts tracking - limit array sizes to reduce memory usage
    timelinePostIds: [],
    createdPostIds: [],
    upvotedPostIds: [],
    downvotedPostIds: [],

    // Social graph tracking - limit array size
    followedUsers: [],

    // Incognito posts tracking
    incognitoPostsCursor: null,
    myIncognitoPostsCursor: null,
    createdIncognitoPostIds: [],
    browseIncognitoPostIds: [],
    upvotedIncognitoPostIds: [],
    downvotedIncognitoPostIds: [],
    incognitoCommentsIds: [],
    upvotedIncognitoCommentIds: [],
    downvotedIncognitoCommentIds: [],
  };

  // Limit array sizes to prevent memory growth
  const MAX_ARRAY_SIZE = 20;

  // Function to trim arrays to prevent memory growth
  function trimArrays() {
    if (vuState.timelinePostIds.length > MAX_ARRAY_SIZE) {
      vuState.timelinePostIds = vuState.timelinePostIds.slice(-MAX_ARRAY_SIZE);
    }
    if (vuState.createdPostIds.length > MAX_ARRAY_SIZE) {
      vuState.createdPostIds = vuState.createdPostIds.slice(-MAX_ARRAY_SIZE);
    }
    if (vuState.upvotedPostIds.length > MAX_ARRAY_SIZE) {
      vuState.upvotedPostIds = vuState.upvotedPostIds.slice(-MAX_ARRAY_SIZE);
    }
    if (vuState.downvotedPostIds.length > MAX_ARRAY_SIZE) {
      vuState.downvotedPostIds = vuState.downvotedPostIds.slice(
        -MAX_ARRAY_SIZE
      );
    }
    if (vuState.followedUsers.length > MAX_ARRAY_SIZE) {
      vuState.followedUsers = vuState.followedUsers.slice(-MAX_ARRAY_SIZE);
    }
    // Trim incognito posts arrays
    if (vuState.createdIncognitoPostIds.length > MAX_ARRAY_SIZE) {
      vuState.createdIncognitoPostIds = vuState.createdIncognitoPostIds.slice(
        -MAX_ARRAY_SIZE
      );
    }
    if (vuState.browseIncognitoPostIds.length > MAX_ARRAY_SIZE) {
      vuState.browseIncognitoPostIds = vuState.browseIncognitoPostIds.slice(
        -MAX_ARRAY_SIZE
      );
    }
    if (vuState.upvotedIncognitoPostIds.length > MAX_ARRAY_SIZE) {
      vuState.upvotedIncognitoPostIds = vuState.upvotedIncognitoPostIds.slice(
        -MAX_ARRAY_SIZE
      );
    }
    if (vuState.downvotedIncognitoPostIds.length > MAX_ARRAY_SIZE) {
      vuState.downvotedIncognitoPostIds =
        vuState.downvotedIncognitoPostIds.slice(-MAX_ARRAY_SIZE);
    }
    if (vuState.incognitoCommentsIds.length > MAX_ARRAY_SIZE) {
      vuState.incognitoCommentsIds = vuState.incognitoCommentsIds.slice(
        -MAX_ARRAY_SIZE
      );
    }
    if (vuState.upvotedIncognitoCommentIds.length > MAX_ARRAY_SIZE) {
      vuState.upvotedIncognitoCommentIds =
        vuState.upvotedIncognitoCommentIds.slice(-MAX_ARRAY_SIZE);
    }
    if (vuState.downvotedIncognitoCommentIds.length > MAX_ARRAY_SIZE) {
      vuState.downvotedIncognitoCommentIds =
        vuState.downvotedIncognitoCommentIds.slice(-MAX_ARRAY_SIZE);
    }
  }

  // Call trimArrays periodically
  setInterval(trimArrays, 5000);

  // Check if we have authenticated users from setup
  if (
    !data ||
    !data.authenticatedUsers ||
    data.authenticatedUsers.length === 0
  ) {
    console.error(
      `No authenticated users available for testing. Setup may have failed.`
    );
    return;
  }

  // Map VU number to a user from our authenticated pool
  const vuIndex = (__VU - 1) % data.authenticatedUsers.length;
  const currentUser = data.authenticatedUsers[vuIndex];

  if (!currentUser || !currentUser.authToken) {
    console.error(
      `No valid user found for VU ${__VU}. Using index ${vuIndex} in a pool of ${data.authenticatedUsers.length} users.`
    );
    return;
  }

  // Pass the specific token, handle, and the list of all handles to socialActivity
  socialActivity(
    currentUser.authToken,
    currentUser.handle,
    data.allUserHandles,
    vuState
  );
}

// --- k6 Test Configuration ---
// Function to generate load test stages based on user count
function generateLoadStages(userCount) {
  console.log(`Generating load stages for ${userCount} users`);

  // Calculate appropriate targets based on user count
  const initialTarget = Math.min(Math.ceil(userCount * 0.1), 100);
  const midTarget = Math.min(
    Math.ceil(userCount * 0.5),
    Math.floor(userCount * 0.5)
  );
  const fullTarget = userCount; // Full load

  // Create stages with appropriate ramp-up periods
  const stages = [
    // Initial ramp-up (10% of users or 100, whichever is smaller)
    { duration: "1m", target: initialTarget },

    // Mid-level ramp-up (50% of users)
    { duration: "2m", target: midTarget },

    // Full load ramp-up (100% of users)
    { duration: "2m", target: fullTarget },

    // Maintain full load for 5 minutes
    { duration: "5m", target: fullTarget },

    // Ramp down
    { duration: "1m", target: 0 },
  ];

  console.log(`Load stages: ${JSON.stringify(stages)}`);
  return stages;
}

export const options = {
  // Calculate setup timeout based on number of users and parallelism - increase timeout significantly
  setupTimeout: `${Math.max(
    300,
    Math.ceil((endUserIndex - startUserIndex + 1) / SETUP_PARALLELISM) * 300
  )}s`,
  // Always run setup - this is critical
  noVUConnectionReuse: true,
  noConnectionReuse: true,
  insecureSkipTLSVerify: true,
  // Debug settings
  discardResponseBodies: false,
  verbose: true,

  scenarios: {
    social_interactions: {
      executor: "ramping-vus",
      startVUs: Math.min(5, Math.ceil(MAX_USERS_PER_INSTANCE * 0.05)),
      stages: generateLoadStages(MAX_USERS_PER_INSTANCE),
      gracefulRampDown: "30s",
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
    // Incognito posts thresholds
    hub_create_incognito_post_duration: ["p(95)<1000"],
    hub_get_incognito_post_duration: ["p(95)<1000"],
    hub_get_incognito_posts_duration: ["p(95)<1000"],
    hub_get_my_incognito_posts_duration: ["p(95)<1000"],
    hub_delete_incognito_post_duration: ["p(95)<1000"],
    hub_add_incognito_comment_duration: ["p(95)<1000"],
    hub_get_incognito_comments_duration: ["p(95)<1000"],
    hub_upvote_incognito_post_duration: ["p(95)<1000"],
    hub_downvote_incognito_post_duration: ["p(95)<1000"],
    hub_unvote_incognito_post_duration: ["p(95)<1000"],
    hub_upvote_incognito_comment_duration: ["p(95)<1000"],
    hub_downvote_incognito_comment_duration: ["p(95)<1000"],
    hub_unvote_incognito_comment_duration: ["p(95)<1000"],
  },
};
