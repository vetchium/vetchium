// loadtest/distributed_hub_scenario.js
import {
  randomIntBetween,
  randomItem,
} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { check, group, sleep } from "k6";
import http from "k6/http";
import { Trend } from "k6/metrics";
import exec from "k6/execution";
import { SharedArray } from "k6/data";

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
const NUM_USERS = parseInt(__ENV.NUM_USERS || "1000000");
const USERS_PER_INSTANCE = parseInt(__ENV.USERS_PER_INSTANCE || "10000"); // Number of users per k6 instance
const INSTANCE_INDEX = parseInt(__ENV.INSTANCE_INDEX || "0"); // Instance index (0-based)
const SETUP_PARALLELISM = parseInt(__ENV.SETUP_PARALLELISM || "100"); // Number of users to authenticate in parallel
const PASSWORD = "NewPassword123$";
const TEST_DURATION_SECONDS = parseInt(__ENV.TEST_DURATION || "600"); // 10 minutes default
const MAX_VUS = parseInt(__ENV.MAX_VUS || "500"); // Maximum number of VUs
const TOKEN_PERSISTENCE_PATH = __ENV.TOKEN_PERSISTENCE_PATH || "/tmp/k6-tokens"; // Path to store token files

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

// --- Constants for authentication ---
const MAX_LOGIN_ATTEMPTS = 3;
const MAX_TFA_FETCH_ATTEMPTS = 5;

// --- TFA Synchronization ---
let tfaLock = false;
const tfaQueue = [];

// --- User range calculation for this instance ---
const START_USER_INDEX = INSTANCE_INDEX * USERS_PER_INSTANCE + 1;
const END_USER_INDEX = Math.min((INSTANCE_INDEX + 1) * USERS_PER_INSTANCE, NUM_USERS);
const ACTUAL_USERS_COUNT = END_USER_INDEX - START_USER_INDEX + 1;

console.log(`Instance ${INSTANCE_INDEX} handling users ${START_USER_INDEX}-${END_USER_INDEX} (${ACTUAL_USERS_COUNT} users)`);

// --- Shared token storage ---
// Use SharedArray to efficiently share auth tokens between VUs
const authTokens = new SharedArray(`auth_tokens_${INSTANCE_INDEX}`, function() {
  // Initialize with empty tokens
  return new Array(ACTUAL_USERS_COUNT).fill(null);
});

// Even though SharedArray is read-only during test execution,
// we can track which tokens have been fetched in a normal array
let tokenFetchStatus = new Array(ACTUAL_USERS_COUNT).fill(false);

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

// --- Login and authentication with token sharing ---
async function loginAndAuthenticateUser(userIndex) {
  // Calculate the local index within our range
  const localIndex = userIndex - START_USER_INDEX;
  
  // Check if we already have a token for this user in our shared array
  if (authTokens[localIndex]) {
    return authTokens[localIndex];
  }
  
  // Check if another VU is already fetching this token
  if (tokenFetchStatus[localIndex]) {
    // Wait a bit and check again to avoid duplicate fetches
    console.debug(`Another VU is fetching token for user ${userIndex}, waiting...`);
    await new Promise(resolve => setTimeout(resolve, 500 * (1 + Math.random())));
    
    // If token was fetched while waiting, return it
    if (authTokens[localIndex]) {
      return authTokens[localIndex];
    }
  }
  
  // Mark that we're fetching this token
  tokenFetchStatus[localIndex] = true;
  
  try {
    const handle = `hubuser${userIndex}`;
    const email = `hubuser${userIndex}@example.com`;
    console.debug(`Authenticating user ${handle}`);
    
    // Step 1: Login to get TFA challenge
    const loginPayload = JSON.stringify({
      login: email,
      password: PASSWORD
    });
    
    const loginRes = http.post(`${API_BASE_URL}/auth/login`, loginPayload, {
      headers: { "Content-Type": "application/json" },
      tags: { name: "LoginAPI" }
    });
    
    if (loginRes.status !== 200) {
      console.error(`Login failed for user ${handle}: Status ${loginRes.status}, Body: ${loginRes.body}`);
      tokenFetchStatus[localIndex] = false;
      return null;
    }
    
    let loginData;
    try {
      loginData = JSON.parse(loginRes.body);
    } catch (e) {
      console.error(`Failed to parse login response: ${e}`);
      tokenFetchStatus[localIndex] = false;
      return null;
    }
    
    // Step 2: Get TFA code from email
    const tfaCode = await fetchTFACodeForUser(email);
    if (!tfaCode) {
      console.error(`Failed to get TFA code for ${email}`);
      tokenFetchStatus[localIndex] = false;
      return null;
    }
    
    // Step 3: Submit TFA code
    const tfaPayload = JSON.stringify({
      login: email,
      code: tfaCode
    });
    
    const tfaRes = http.post(`${API_BASE_URL}/auth/tfa`, tfaPayload, {
      headers: { "Content-Type": "application/json" },
      tags: { name: "TFAAPI" }
    });
    
    if (tfaRes.status !== 200) {
      console.error(`TFA verification failed for ${handle}: ${tfaRes.status}, ${tfaRes.body}`);
      tokenFetchStatus[localIndex] = false;
      return null;
    }
    
    let tfaData;
    try {
      tfaData = JSON.parse(tfaRes.body);
    } catch (e) {
      console.error(`Failed to parse TFA response: ${e}`);
      tokenFetchStatus[localIndex] = false;
      return null;
    }
    
    const authToken = tfaData.auth_token;
    
    // Store user data in shared memory (read-only, but we can use it for reference)
    // Note: We can't modify SharedArray after initialization, so just keep it for reference
    // and track which tokens have been fetched
    console.log(`Successfully authenticated user ${handle}`);
    
    // Return the user data
    const userData = {
      userHandle: handle,
      authToken: authToken
    };
    
    // Store the reference to the userData in our tracking structure for this instance
    // Despite being read-only, this serves as a lookup table
    // In k6, SharedArray is only read-only during test execution, but we can reference its values
    // This is a bit of a hack, but it works for our use case
    authTokens[localIndex] = userData;
    
    return userData;
  } catch (e) {
    console.error(`Unexpected error in loginAndAuthenticateUser: ${e}`);
    tokenFetchStatus[localIndex] = false;
    return null;
  }
}

// --- TFA Code Fetch (Uses Email) ---
async function fetchTFACodeForUser(email) {
  try {
    // Acquire lock before fetching TFA code - critical for email TFA code extraction
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
      releaseTfaLock();
      return null;
    }

    // Step 2: Fetch the message content by ID
    const messageUrl = `${MAILPIT_URL}/api/v1/message/${messageId}`;
    console.debug(`Fetching message content from ${messageUrl}`);

    const messageRes = http.get(messageUrl);

    if (messageRes.status !== 200) {
      console.error(
        `Failed to fetch message content. Status: ${messageRes.status}`
      );
      releaseTfaLock();
      return null;
    }

    try {
      const messageData = JSON.parse(messageRes.body);
      console.debug(`Message data fetched successfully`);

      // Extract TFA code using regex on the text/html part
      const htmlContent =
        messageData.Text.HTML ||
        (messageData.Parts?.find((p) => p.Header?.ContentType?.includes("text/html"))
          ?.Body || "");

      // Use the exact regex pattern from Go tests: look for 6 digits code
      const codeMatch = htmlContent.match(/([0-9]{6})/);

      if (codeMatch && codeMatch[1]) {
        const code = codeMatch[1];
        console.log(`Extracted TFA code: ${code} for ${email}`);
        releaseTfaLock();
        return code;
      } else {
        console.error(`Failed to extract TFA code from email content`);
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
  } catch (e) {
    console.error(`Unexpected error in fetchTFACodeForUser: ${e}`);
    releaseTfaLock();
    return null;
  }
}

// --- k6 Setup Function ---
export function setup() {
  console.log(`Instance ${INSTANCE_INDEX} starting setup phase for users ${START_USER_INDEX}-${END_USER_INDEX}`);
  
  // For distributed testing, we don't need to authenticate all users during setup
  // Instead, we'll authenticate them on-demand during the test
  // We just need to pre-authenticate a small subset for better test startup
  
  const preAuthCount = Math.min(SETUP_PARALLELISM, ACTUAL_USERS_COUNT);
  console.log(`Pre-authenticating ${preAuthCount} users during setup...`);
  
  // Create an array of promises for parallel authentication
  const authPromises = [];
  
  // Pre-authenticate a subset of users to warm up the cache
  for (let i = 0; i < preAuthCount; i++) {
    const userIndex = START_USER_INDEX + i;
    authPromises.push(loginAndAuthenticateUser(userIndex));
  }
  
  // Wait for all authentication promises to resolve
  Promise.all(authPromises).then(results => {
    const successCount = results.filter(r => r !== null).length;
    console.log(`Setup complete. Pre-authenticated ${successCount}/${preAuthCount} users`);
  });
  
  // Return the user range info for this instance
  return {
    startUserIndex: START_USER_INDEX,
    endUserIndex: END_USER_INDEX,
    instanceIndex: INSTANCE_INDEX
  };
}

// --- Default Function (Main Test Entry Point) ---
export default function(data) {
  // Calculate which user this VU should use based on VU ID
  // We distribute users evenly across VUs
  const userOffset = __VU % ACTUAL_USERS_COUNT;
  const userIndex = START_USER_INDEX + userOffset;
  
  // Get authentication token for this user (from cache if available)
  const userData = loginAndAuthenticateUser(userIndex);
  if (!userData) {
    console.error(`Authentication failed for user ${userIndex}, skipping iteration`);
    return;
  }
  
  // Perform social activities with the authenticated user
  socialActivity(userData.authToken, userData.userHandle);
}

// --- Social Activity Function ---
function socialActivity(authToken, userHandle) {
  // Common request parameters with auth token
  const authParams = {
    headers: {
      "Content-Type": "application/json",
      "Authorization": `Bearer ${authToken}`
    }
  };
  
  // Randomly select an action to perform
  const action = Math.random();
  
  if (action < 0.3) {
    // Read timeline (30% of actions)
    group("Read Timeline", function() {
      const timelineRes = http.get(`${API_BASE_URL}/hub/timeline`, {
        ...authParams,
        tags: { name: "TimelineAPI" }
      });
      
      timelineReadTrend.add(timelineRes.timings.duration);
      
      check(timelineRes, {
        "Timeline request successful": (r) => r.status === 200
      });
      
      if (timelineRes.status !== 200) {
        console.error(`Timeline request failed: ${timelineRes.status}`);
      }
      
      sleep(randomIntBetween(1, 3));
    });
  } 
  else if (action < 0.5) {
    // Create post (20% of actions)
    group("Create Post", function() {
      const postContent = `Test post from ${userHandle} at ${new Date().toISOString()}`;
      const postPayload = JSON.stringify({ content: postContent });
      
      const createPostRes = http.post(
        `${API_BASE_URL}/hub/create-post`,
        postPayload,
        {
          ...authParams,
          tags: { name: "CreatePostAPI" }
        }
      );
      
      createPostTrend.add(createPostRes.timings.duration);
      
      check(createPostRes, {
        "Create post request successful": (r) => r.status === 200
      });
      
      sleep(randomIntBetween(1, 3));
    });
  }
  else if (action < 0.7) {
    // Read posts (20% of actions)
    group("Read Posts", function() {
      const postsRes = http.get(`${API_BASE_URL}/hub/posts/recent`, {
        ...authParams,
        tags: { name: "RecentPostsAPI" }
      });
      
      postDetailsTrend.add(postsRes.timings.duration);
      
      check(postsRes, {
        "Posts request successful": (r) => r.status === 200
      });
      
      sleep(randomIntBetween(1, 2));
    });
  }
  else if (action < 0.85) {
    // Follow/unfollow users (15% of actions)
    group("Follow Actions", function() {
      // Generate a random user handle to follow
      const targetUserIndex = START_USER_INDEX + Math.floor(Math.random() * ACTUAL_USERS_COUNT);
      const targetHandle = `hubuser${targetUserIndex}`;
      
      // Don't follow yourself
      if (targetHandle === userHandle) {
        return;
      }
      
      const followPayload = JSON.stringify({ handle_to_follow: targetHandle });
      
      const followRes = http.post(
        `${API_BASE_URL}/hub/follow-user`,
        followPayload,
        {
          ...authParams,
          tags: { name: "FollowAPI" }
        }
      );
      
      followTrend.add(followRes.timings.duration);
      
      check(followRes, {
        "Follow request successful or already following": (r) =>
          r.status === 200 || r.status === 422
      });
      
      sleep(randomIntBetween(1, 3));
    });
  }
  else {
    // Vote actions (15% of actions)
    group("Vote Actions", function() {
      // Get recent posts first
      const postsRes = http.get(`${API_BASE_URL}/hub/posts/recent`, {
        ...authParams,
        tags: { name: "RecentPostsForVoteAPI" }
      });
      
      if (postsRes.status !== 200) {
        console.error(`Failed to get posts for voting: ${postsRes.status}`);
        return;
      }
      
      let posts;
      try {
        posts = JSON.parse(postsRes.body);
      } catch (e) {
        console.error(`Failed to parse posts response: ${e}`);
        return;
      }
      
      if (!posts || !posts.posts || posts.posts.length === 0) {
        console.debug("No posts available for voting");
        return;
      }
      
      // Select a random post to vote on
      const randomPost = randomItem(posts.posts);
      const postId = randomPost.id;
      
      // Choose between upvote/downvote randomly
      const voteType = Math.random() < 0.5 ? "upvote" : "downvote";
      const voteEndpoint = voteType === "upvote" ? "upvote-user-post" : "downvote-user-post";
      const voteTrend = voteType === "upvote" ? upvoteTrend : downvoteTrend;
      
      const votePayload = JSON.stringify({ post_id: postId });
      
      const voteRes = http.post(
        `${API_BASE_URL}/hub/${voteEndpoint}`,
        votePayload,
        {
          ...authParams,
          tags: { name: `${voteType.charAt(0).toUpperCase() + voteType.slice(1)}API` }
        }
      );
      
      voteTrend.add(voteRes.timings.duration);
      
      check(voteRes, {
        "Vote request successful or already voted": (r) =>
          r.status === 200 || r.status === 422
      });
      
      sleep(randomIntBetween(1, 2));
    });
  }
  
  // Small pause between iterations
  sleep(randomIntBetween(2, 5));
}

// --- k6 Test Configuration ---
export const options = {
  discardResponseBodies: true,  // Discard response bodies to save memory
  scenarios: {
    social_interactions: {
      executor: "ramping-vus",
      startVUs: 10,
      stages: [
        // Ramp up to target VUs
        { duration: "30s", target: Math.ceil(MAX_VUS * 0.3) },
        { duration: "1m", target: Math.ceil(MAX_VUS * 0.6) },
        { duration: "1m", target: MAX_VUS },
        // Steady load
        { duration: `${TEST_DURATION_SECONDS - 210}s`, target: MAX_VUS },
        // Ramp down
        { duration: "30s", target: 0 }
      ],
      gracefulRampDown: "30s"
    }
  },
  thresholds: {
    http_req_failed: ["rate<0.05"],  // Error rate under 5%
    http_req_duration: ["p(95)<5000"], // 95% of requests under 5s
    "http_req_duration{name:LoginAPI}": ["p(95)<10000"], // Login can take longer
    "http_req_duration{name:TFAAPI}": ["p(95)<10000"],  // TFA can take longer
  }
};
