// distributed_hub_scenario.js - Distributed k6 load test for Vetchium API
import {
  randomIntBetween,
  randomItem,
} from "https://jslib.k6.io/k6-utils/1.2.0/index.js";
import { check, group, sleep } from "k6";
import http from "k6/http";
import { Trend } from "k6/metrics";
import exec from "k6/execution";

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

// --- Calculate user range for this instance ---
const startUserIndex = INSTANCE_INDEX * USERS_PER_INSTANCE + 1;
const endUserIndex = Math.min(
  (INSTANCE_INDEX + 1) * USERS_PER_INSTANCE,
  TOTAL_USERS
);

console.log(
  `Instance ${INSTANCE_INDEX}/${INSTANCE_COUNT} handling users ${startUserIndex} to ${endUserIndex}`
);

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
      releaseTfaLock();
      return null;
    }

    // Step 2: Get the message content using the message ID
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
      const htmlContent = messageData.HTML || "";

      // Step 3: Extract TFA code using regex pattern
      // The pattern looks for a 6-digit code that appears after certain text patterns
      const tfaCodePattern = /verification code is[^\d]*(\d{6})|code:[^\d]*(\d{6})|code is[^\d]*(\d{6})/i;
      const tfaMatch = htmlContent.match(tfaCodePattern);

      if (tfaMatch) {
        // The code could be in any of the capture groups, find the non-null one
        const tfaCode = tfaMatch[1] || tfaMatch[2] || tfaMatch[3];
        console.debug(`Successfully extracted TFA code: ${tfaCode}`);
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
    const loginRes = http.post(
      `${API_BASE_URL}/hub/login`,
      loginPayload,
      {
        headers: { "Content-Type": "application/json" },
        tags: { name: "HubLoginAPI" },
      }
    );

    // Check login response
    if (loginRes.status !== 200) {
      console.error(
        `Login failed for ${user.email}. Status: ${loginRes.status}, Body: ${loginRes.body}`
      );
      sleep(2);
      continue;
    }

    try {
      const loginData = JSON.parse(loginRes.body);

      // Extract TFA token from response
      const tfaToken = loginData.tfa_token;
      if (!tfaToken) {
        console.error(
          `TFA token not found in login response for ${user.email}. Response: ${loginRes.body}`
        );
        sleep(2);
        continue;
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
      const tfaPayload = JSON.stringify({
        tfa_token: tfaToken,
        tfa_code: tfaCode,
      });

      const tfaRes = http.post(
        `${API_BASE_URL}/hub/tfa`,
        tfaPayload,
        {
          headers: { "Content-Type": "application/json" },
          tags: { name: "HubTFAAPI" },
        }
      );

      // Check TFA response
      if (tfaRes.status !== 200) {
        console.error(
          `TFA verification failed for ${user.email}. Status: ${tfaRes.status}, Body: ${tfaRes.body}`
        );
        sleep(2);
        continue;
      }

      try {
        const tfaData = JSON.parse(tfaRes.body);
        authToken = tfaData.auth_token;

        if (!authToken) {
          console.error(
            `Auth token not found in TFA response for ${user.email}. Response: ${tfaRes.body}`
          );
          sleep(2);
          continue;
        }

        console.debug(`Successfully authenticated ${user.email}`);
        return authToken;
      } catch (e) {
        console.error(
          `Error parsing TFA response for ${user.email}: ${e}. Body: ${tfaRes.body}`
        );
        sleep(2);
      }
    } catch (e) {
      console.error(
        `Error parsing login response for ${user.email}: ${e}. Body: ${loginRes.body}`
      );
      sleep(2);
    }
  }

  return null;
}

// --- k6 Setup Function ---
async function setup() {
  console.log(`Starting setup for instance ${INSTANCE_INDEX} with users ${startUserIndex} to ${endUserIndex}`);
  
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
  const batchSize = SETUP_PARALLELISM;
  const batches = Math.ceil(users.length / batchSize);
  
  for (let batchIndex = 0; batchIndex < batches; batchIndex++) {
    const batchStart = batchIndex * batchSize;
    const batchEnd = Math.min((batchIndex + 1) * batchSize, users.length);
    const batchUsers = users.slice(batchStart, batchEnd);
    
    console.log(`Processing authentication batch ${batchIndex + 1}/${batches} (users ${batchStart + 1} to ${batchEnd})`);
    
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
    const successfulAuths = batchResults.filter(result => result !== null);
    authenticatedUsers.push(...successfulAuths);
    
    console.log(`Batch ${batchIndex + 1} complete: ${successfulAuths.length}/${batchUsers.length} users authenticated`);
  }
  
  console.log(`Setup complete: ${authenticatedUsers.length}/${users.length} users authenticated`);
  
  // Return the authenticated users and all handles for the test
  return {
    authenticatedUsers,
    allUserHandles
  };
}

// --- Social Activity Function ---
function socialActivity(authToken, userHandle, allUserHandles, vuState) {
  // Common auth parameters for all requests
  const authParams = {
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${authToken}`
    }
  };

  // Randomly select a social action to perform
  const socialActions = [
    'follow_user',
    'unfollow_user',
    'create_post',
    'read_timeline',
    'view_post_details',
    'upvote',
    'downvote',
    'unvote'
  ];

  // Weight the actions to create a realistic distribution
  // Reading timeline and viewing posts should be more common than posting
  const actionWeights = [15, 5, 10, 30, 20, 10, 5, 5]; // Percentages
  
  // Calculate cumulative weights
  const cumulativeWeights = [];
  let sum = 0;
  for (const weight of actionWeights) {
    sum += weight;
    cumulativeWeights.push(sum);
  }
  
  // Select action based on weights
  const randomValue = Math.random() * 100;
  let selectedActionIndex = 0;
  
  for (let i = 0; i < cumulativeWeights.length; i++) {
    if (randomValue <= cumulativeWeights[i]) {
      selectedActionIndex = i;
      break;
    }
  }
  
  const selectedAction = socialActions[selectedActionIndex];
  
  // Execute the selected social action
  group(`Social Action: ${selectedAction}`, () => {
    switch (selectedAction) {
      case 'follow_user': {
        // Select a random user to follow (not self)
        const otherHandles = allUserHandles.filter(h => h !== userHandle);
        if (otherHandles.length === 0) {
          console.debug(`VU ${__VU} (${userHandle}): No other users to follow.`);
          break;
        }
        
        // Don't follow users we already follow
        const potentialHandlesToFollow = otherHandles.filter(
          h => !vuState.followedUsers.includes(h)
        );
        
        if (potentialHandlesToFollow.length === 0) {
          console.debug(`VU ${__VU} (${userHandle}): Already following all available users.`);
          break;
        }
        
        const handleToFollow = randomItem(potentialHandlesToFollow);
        const followPayload = JSON.stringify({ handle: handleToFollow });
        
        console.debug(`VU ${__VU} (${userHandle}): Attempting to follow user ${handleToFollow}`);
        
        const followRes = http.post(
          `${API_BASE_URL}/hub/follow-user`,
          followPayload,
          {
            ...authParams,
            tags: { name: 'HubFollowUserAPI' }
          }
        );
        
        followTrend.add(followRes.timings.duration);
        
        // Check for success
        check(followRes, {
          'Follow request successful': (r) => r.status === 200
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
      
      case 'unfollow_user': {
        // Check if we have any users to unfollow
        if (vuState.followedUsers.length === 0) {
          console.debug(`VU ${__VU} (${userHandle}): No users to unfollow.`);
          break;
        }
        
        const handleToUnfollow = randomItem(vuState.followedUsers);
        const unfollowPayload = JSON.stringify({ handle: handleToUnfollow });
        
        console.debug(`VU ${__VU} (${userHandle}): Attempting to unfollow user ${handleToUnfollow}`);
        
        const unfollowRes = http.post(
          `${API_BASE_URL}/hub/unfollow-user`,
          unfollowPayload,
          {
            ...authParams,
            tags: { name: 'HubUnfollowUserAPI' }
          }
        );
        
        unfollowTrend.add(unfollowRes.timings.duration);
        
        // Check for success
        check(unfollowRes, {
          'Unfollow request successful': (r) => r.status === 200
        });
        
        // If successful, remove from followed users list
        if (unfollowRes.status === 200) {
          vuState.followedUsers = vuState.followedUsers.filter(
            h => h !== handleToUnfollow
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
      
      case 'create_post': {
        // Generate random post content
        const postContent = `Test post from ${userHandle} at ${new Date().toISOString()} - ${randomString(20)}`;
        const postPayload = JSON.stringify({ content: postContent });
        
        console.debug(`VU ${__VU} (${userHandle}): Creating new post`);
        
        const createPostRes = http.post(
          `${API_BASE_URL}/hub/create-post`,
          postPayload,
          {
            ...authParams,
            tags: { name: 'HubCreatePostAPI' }
          }
        );
        
        createPostTrend.add(createPostRes.timings.duration);
        
        // Check for success
        check(createPostRes, {
          'Create post request successful': (r) => r.status === 200
        });
        
        // If successful, extract post ID and add to created posts
        if (createPostRes.status === 200) {
          try {
            const postData = JSON.parse(createPostRes.body);
            if (postData.post_id) {
              vuState.createdPostIds.push(postData.post_id);
              console.debug(`VU ${__VU} (${userHandle}): Created post with ID ${postData.post_id}`);
            }
          } catch (e) {
            console.error(`VU ${__VU} (${userHandle}): Error parsing create post response: ${e}`);
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
      
      case 'read_timeline': {
        // Prepare query parameters
        const limit = 20; // Number of posts to fetch
        let timelineUrl = `${API_BASE_URL}/hub/timeline?limit=${limit}`;
        
        // Add cursor for pagination if we have one
        if (vuState.timelineCursor) {
          timelineUrl += `&cursor=${vuState.timelineCursor}`;
        }
        
        console.debug(`VU ${__VU} (${userHandle}): Reading timeline${vuState.timelineCursor ? ' with cursor' : ''}`);
        
        const timelineRes = http.get(timelineUrl, {
          ...authParams,
          tags: { name: 'HubTimelineAPI' }
        });
        
        timelineReadTrend.add(timelineRes.timings.duration);
        
        // Check for success
        check(timelineRes, {
          'Timeline read successful': (r) => r.status === 200
        });
        
        // Process timeline response
        if (timelineRes.status === 200) {
          try {
            const timelineData = JSON.parse(timelineRes.body);
            
            // Update cursor for next pagination
            if (timelineData.next_cursor) {
              vuState.timelineCursor = timelineData.next_cursor;
            }
            
            // Store post IDs for potential interactions
            if (timelineData.posts && timelineData.posts.length > 0) {
              // Add new post IDs to the timeline posts collection
              timelineData.posts.forEach(post => {
                if (!vuState.timelinePostIds.includes(post.post_id)) {
                  vuState.timelinePostIds.push(post.post_id);
                }
              });
              
              console.debug(`VU ${__VU} (${userHandle}): Retrieved ${timelineData.posts.length} posts, total in memory: ${vuState.timelinePostIds.length}`);
            } else {
              console.debug(`VU ${__VU} (${userHandle}): No posts in timeline response`);
            }
          } catch (e) {
            console.error(`VU ${__VU} (${userHandle}): Error parsing timeline response: ${e}`);
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
      
      case 'view_post_details': {
        // Combine timeline posts and created posts for selection
        const availablePosts = [...vuState.timelinePostIds, ...vuState.createdPostIds];
        
        if (availablePosts.length === 0) {
          console.debug(`VU ${__VU} (${userHandle}): No posts available to view details.`);
          break;
        }
        
        const postIdToView = randomItem(availablePosts);
        const postDetailsUrl = `${API_BASE_URL}/hub/post/${postIdToView}`;
        
        console.debug(`VU ${__VU} (${userHandle}): Viewing details for post ${postIdToView}`);
        
        const postDetailsRes = http.get(postDetailsUrl, {
          ...authParams,
          tags: { name: 'HubPostDetailsAPI' }
        });
        
        postDetailsTrend.add(postDetailsRes.timings.duration);
        
        // Check for success
        check(postDetailsRes, {
          'Post details request successful': (r) => r.status === 200
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
      
      case 'upvote':
      case 'downvote': {
        // Combine timeline posts for voting (exclude own posts)
        const availablePosts = vuState.timelinePostIds.filter(
          id => !vuState.createdPostIds.includes(id)
        );
        
        // Also exclude posts we've already voted on
        const votablePosts = availablePosts.filter(
          id => !vuState.upvotedPostIds.includes(id) && !vuState.downvotedPostIds.includes(id)
        );
        
        if (votablePosts.length === 0) {
          console.debug(`VU ${__VU} (${userHandle}): No posts available to ${selectedAction}.`);
          break;
        }
        
        const postIdToVote = randomItem(votablePosts);
        const votePayload = JSON.stringify({ post_id: postIdToVote });
        const voteUrl = `${API_BASE_URL}/hub/${selectedAction === 'upvote' ? 'upvote' : 'downvote'}-user-post`;
        
        console.debug(`VU ${__VU} (${userHandle}): Attempting to ${selectedAction} post ${postIdToVote}`);
        
        // Set appropriate trend and tag based on vote type
        let voteTrend, voteTag;
        if (selectedAction === 'upvote') {
          voteTrend = upvoteTrend;
          voteTag = 'HubUpvoteAPI';
        } else {
          voteTrend = downvoteTrend;
          voteTag = 'HubDownvoteAPI';
        }
        
        const voteRes = http.post(voteUrl, votePayload, {
          ...authParams,
          tags: { name: voteTag }
        });
        
        voteTrend.add(voteRes.timings.duration);
        
        // Check for success or expected error
        check(voteRes, {
          [`${selectedAction} request successful or expected error`]: (r) => r.status === 200 || r.status === 422
        });
        
        // If successful, track the voted post
        if (voteRes.status === 200) {
          if (selectedAction === 'upvote') {
            vuState.upvotedPostIds.push(postIdToVote);
          } else {
            vuState.downvotedPostIds.push(postIdToVote);
          }
        }
        
        // Log unexpected errors
        if (voteRes.status !== 200 && voteRes.status !== 422) {
          console.error(
            `VU ${__VU} (${userHandle}): ${selectedAction} API Error! PostID: ${postIdToVote}, Status: ${voteRes.status}, Body: ${voteRes.body}`
          );
        }
        
        sleep(randomIntBetween(1, 4));
        break;
      }
      
      case 'unvote': {
        // Combine upvoted and downvoted posts to pick from
        const votedPosts = [...vuState.upvotedPostIds, ...vuState.downvotedPostIds];
        
        if (votedPosts.length === 0) {
          console.debug(`VU ${__VU} (${userHandle}): No voted posts available to unvote.`);
          break;
        }
        
        const postIdToUnvote = randomItem(votedPosts);
        const unvotePayload = JSON.stringify({ post_id: postIdToUnvote });
        
        console.debug(`VU ${__VU} (${userHandle}): Attempting to unvote post ${postIdToUnvote}`);
        
        const unvoteRes = http.post(
          `${API_BASE_URL}/hub/unvote-user-post`,
          unvotePayload,
          {
            ...authParams,
            tags: { name: 'HubUnvoteAPI' }
          }
        );
        
        unvoteTrend.add(unvoteRes.timings.duration);
        
        // Check for success or expected error
        check(unvoteRes, {
          'Unvote request successful or expected error': (r) => r.status === 200 || r.status === 422
        });
        
        // If successful, remove from voted lists
        if (unvoteRes.status === 200) {
          vuState.upvotedPostIds = vuState.upvotedPostIds.filter(id => id !== postIdToUnvote);
          vuState.downvotedPostIds = vuState.downvotedPostIds.filter(id => id !== postIdToUnvote);
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
    }
    
    // Think time between actions
    sleep(randomIntBetween(2, 6));
  });
}

// --- Main Test Logic (Accepts setup data) ---
export default function (data) {
  // Initialize VU state for tracking social interactions
  const vuState = {
    // Timeline pagination
    timelineCursor: null,

    // Posts tracking
    timelinePostIds: [],
    createdPostIds: [],
    upvotedPostIds: [],
    downvotedPostIds: [],

    // Social graph tracking
    followedUsers: []
  };

  // Check if we have authenticated users from setup
  if (!data || !data.authenticatedUsers || data.authenticatedUsers.length === 0) {
    console.error(`No authenticated users available for testing. Setup may have failed.`);
    return;
  }

  // Map VU number to a user from our authenticated pool
  const vuIndex = (__VU - 1) % data.authenticatedUsers.length;
  const currentUser = data.authenticatedUsers[vuIndex];

  if (!currentUser || !currentUser.authToken) {
    console.error(`No valid user found for VU ${__VU}. Using index ${vuIndex} in a pool of ${data.authenticatedUsers.length} users.`);
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
export const options = {
  // Calculate setup timeout based on number of users and parallelism
  setupTimeout: `${Math.max(300, Math.ceil((endUserIndex - startUserIndex + 1) / SETUP_PARALLELISM) * 60)}s`,

  scenarios: {
    social_interactions: {
      executor: "shared-iterations",
      vus: Math.min(endUserIndex - startUserIndex + 1, 100), // Cap at 100 VUs per instance for stability
      iterations: 1000, // This is per VU
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
