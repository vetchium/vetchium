# UI Test Cases for Vetchium Social Networking and Jobs Site

## Hub User Tests (ronweasly)

### Authentication & Account Management

#### Hub Login
- [ ] Try to login with an invalid format email address. Sign in should be disabled
- [ ] Try to login with an invalid, non-existent email address. Sign in should fail
- [ ] Try to login with a very long email address more than the allowed length, it should fail
- [ ] Try to login with a short password. It should fail
- [ ] Try to login with a very long password more than the allowed length, it should fail
- [ ] Try to login with an insecure password that does not have all the necessary requirements
- [ ] Try to login with a valid password and email address and the login should succeed and show the TFA page
- [ ] After the login page, in the TFA page, wait long before entering the TFA token and ensure that the TFA authentication fails
- [ ] After the login page, in the TFA page, enter a TFA code that is different from what was received, and ensure that the authentication fails
- [ ] After the login page, in the TFA page, enter a TFA code that is of different lengths, lesser than 6 characters and longer than 6 characters and ensure that the authentication fails
- [ ] After the login page, in the TFA page, enter the correct TFA code and ensure that the authentication succeeds
- [ ] Try to access protected pages without logging in and ensure redirection to login page
- [ ] After successful login, verify that the user is redirected to the intended page or dashboard

#### Hub Signup
- [ ] Try to signup with an already registered email address, it should fail
- [ ] Try to signup with invalid email format, validation should fail
- [ ] Try to signup with weak password not meeting security requirements
- [ ] Try to signup with mismatched password confirmation
- [ ] Try to signup with empty required fields
- [ ] Try to signup with special characters in name fields where not allowed
- [ ] Try to signup with extremely long values in each field
- [ ] Successfully signup with valid data and verify email verification is triggered
- [ ] After signup, verify that the user cannot login until email is verified
- [ ] Verify that the invitation code field works correctly if required

#### Hub Onboarding
- [ ] Access onboarding page without valid signup, should be redirected
- [ ] Try to skip required onboarding steps
- [ ] Upload invalid file formats for any document uploads during onboarding
- [ ] Upload files exceeding size limits during onboarding
- [ ] Complete onboarding with minimum required information
- [ ] Complete onboarding with all optional fields filled
- [ ] Verify onboarding progress is saved if user exits and returns

#### Forgot/Reset Password
- [ ] Try forgot password with non-existent email
- [ ] Try forgot password with invalid email format
- [ ] Successfully request password reset and verify email is sent
- [ ] Try to use expired reset password link
- [ ] Try to use already-used reset password link
- [ ] Try to reset with password not meeting requirements
- [ ] Try to reset with mismatched password confirmation
- [ ] Successfully reset password and verify can login with new password
- [ ] Verify old password no longer works after reset

#### Change Password
- [ ] Try to change password without entering current password
- [ ] Try to change password with incorrect current password
- [ ] Try to change to a password that doesn't meet requirements
- [ ] Try to change to the same password as current
- [ ] Try with mismatched new password confirmation
- [ ] Successfully change password and verify immediate re-login required
- [ ] Verify all other sessions are invalidated after password change

#### Change Email Address
- [ ] Try to change to an already registered email address
- [ ] Try to change to invalid email format
- [ ] Try to change to the same email address
- [ ] Successfully request email change and verify verification email sent to new address
- [ ] Verify cannot login with new email until verified
- [ ] Verify old email still works until new email is verified
- [ ] Complete email verification and verify old email no longer works

### Profile Management

#### Bio Management
- [ ] Try to save bio exceeding character limit
- [ ] Try to save bio with prohibited content/HTML injection
- [ ] Save bio with special characters and emojis
- [ ] Save bio with multiple paragraphs and formatting
- [ ] Clear bio and save empty bio
- [ ] Update existing bio and verify changes persist
- [ ] Verify bio displays correctly on public profile

#### Profile Picture
- [ ] Upload image exceeding size limit (should fail)
- [ ] Upload non-image file formats (should fail)
- [ ] Upload corrupted image file (should fail)
- [ ] Upload very small resolution image
- [ ] Upload very large resolution image (should be resized)
- [ ] Upload supported formats (JPG, PNG, GIF)
- [ ] Successfully upload and verify picture displays
- [ ] Remove profile picture and verify default avatar shows
- [ ] Replace existing picture with new one
- [ ] Verify profile picture appears correctly in all contexts (posts, comments, etc.)

#### Handle Management (Paid Tier Only)
- [ ] Try to set handle as free tier user (should fail)
- [ ] Check handle availability with already taken handle
- [ ] Check handle availability with invalid characters
- [ ] Check handle availability with too short handle
- [ ] Check handle availability with too long handle
- [ ] Check handle availability with reserved words
- [ ] Successfully set available handle
- [ ] Try to change handle after already setting one
- [ ] Verify handle appears in profile URL

#### Official Emails
- [ ] Add official email with personal email domain (should fail if restricted)
- [ ] Add official email with invalid format
- [ ] Add duplicate official email
- [ ] Add official email exceeding maximum allowed count
- [ ] Successfully add official email and verify verification email sent
- [ ] Trigger re-verification for unverified email
- [ ] Try to verify with incorrect verification code
- [ ] Try to verify with expired verification code
- [ ] Successfully verify official email
- [ ] Delete verified official email
- [ ] Delete unverified official email
- [ ] Verify at least one email must remain (cannot delete all)

### Education & Achievements

#### Education Management
- [ ] Add education with future graduation date
- [ ] Add education with graduation date before start date
- [ ] Add education without required fields
- [ ] Filter institutes with partial name match
- [ ] Filter institutes with no results
- [ ] Select institute from filtered results
- [ ] Add custom institute not in database
- [ ] Add multiple education entries
- [ ] Delete education entry with confirmation
- [ ] Edit existing education entry
- [ ] Verify education appears in chronological order

#### Achievements Management
- [ ] Add achievement without required fields
- [ ] Add achievement with very long description
- [ ] Add achievement with special characters
- [ ] Add achievement with past date
- [ ] Add achievement with future date
- [ ] Add multiple achievements
- [ ] Delete achievement with confirmation
- [ ] Verify achievements display in correct order
- [ ] Verify achievement count limits if any

### Work History

#### Work History Management
- [ ] Add work history with end date before start date
- [ ] Add work history without end date (current job)
- [ ] Add work history with all optional fields
- [ ] Add work history with minimum required fields
- [ ] Filter and select employer from database
- [ ] Add custom employer not in database
- [ ] Update existing work history entry
- [ ] Delete work history entry
- [ ] Add overlapping work history entries
- [ ] Verify work history displays in chronological order
- [ ] Add work history with very long job description

### Job Search & Applications

#### Opening Search
- [ ] Search openings with no filters (view all)
- [ ] Search with single keyword
- [ ] Search with multiple keywords
- [ ] Filter by single VTag
- [ ] Filter by multiple VTags
- [ ] Filter by location
- [ ] Filter by employer
- [ ] Filter by salary range
- [ ] Filter by experience level
- [ ] Combine multiple filters
- [ ] Clear individual filters
- [ ] Clear all filters at once
- [ ] Sort results by relevance
- [ ] Sort results by date posted
- [ ] Sort results by salary
- [ ] Paginate through large result sets
- [ ] Search with no results and verify appropriate message

#### Opening Details
- [ ] View opening details as logged-in user
- [ ] View expired opening (should show appropriate status)
- [ ] View opening from blocked/defunct employer
- [ ] Verify all opening information displays correctly
- [ ] Check employer information link works
- [ ] Verify salary range displays correctly
- [ ] Verify required qualifications display
- [ ] Check application deadline if present

#### Apply for Opening
- [ ] Apply for opening without uploading resume (if required)
- [ ] Apply with resume exceeding size limit
- [ ] Apply with invalid resume format
- [ ] Apply for same opening twice (should fail)
- [ ] Apply for opening after deadline
- [ ] Apply for opening that's been closed
- [ ] Successfully apply with all required information
- [ ] Apply with optional cover letter
- [ ] Verify application appears in "My Applications"
- [ ] Verify cannot apply if cool-off period not met

#### My Applications
- [ ] View all applications with no filter
- [ ] Filter applications by status (applied, shortlisted, rejected)
- [ ] Filter applications by date range
- [ ] Search applications by employer name
- [ ] Search applications by job title
- [ ] Sort applications by date applied
- [ ] Sort applications by last updated
- [ ] View application details
- [ ] Withdraw pending application
- [ ] Try to withdraw already processed application
- [ ] Verify pagination works correctly

### Candidacy Management

#### Candidacy Comments
- [ ] Add comment to candidacy
- [ ] Add very long comment exceeding limit
- [ ] Add comment with special characters
- [ ] Add empty comment (should fail)
- [ ] View all comments in chronological order
- [ ] Add comment and verify real-time update
- [ ] Verify only authorized users can view comments

#### Candidacy Info
- [ ] View candidacy details
- [ ] Verify interview schedule displays correctly
- [ ] Verify offer details display when applicable
- [ ] Check candidacy status updates
- [ ] Verify all candidacy milestones shown

#### Interview Management
- [ ] View scheduled interviews
- [ ] RSVP Yes to interview
- [ ] RSVP No to interview
- [ ] Change RSVP response
- [ ] Try to RSVP after deadline
- [ ] View interview details (time, location, interviewers)
- [ ] Verify calendar integration if available
- [ ] View past interviews

### Social Features

#### Colleague Connections
- [ ] Search for colleagues by name
- [ ] Search for colleagues by email
- [ ] Send colleague connection request
- [ ] Send duplicate connection request (should fail)
- [ ] Cancel pending connection request
- [ ] View pending connection requests sent
- [ ] View pending connection requests received
- [ ] Approve colleague connection
- [ ] Reject colleague connection
- [ ] Unlink existing colleague
- [ ] Filter colleagues by various criteria
- [ ] Verify colleague count limits if any

#### Endorsements
- [ ] Request endorsement for application
- [ ] Try to request endorsement from non-colleague
- [ ] View pending endorsement requests
- [ ] Approve endorsement request
- [ ] Reject endorsement request
- [ ] Withdraw endorsement request
- [ ] Verify endorsed applications show endorser details

#### Posts (User Posts)
- [ ] Create text-only post (paid tier)
- [ ] Create post with images (paid tier)
- [ ] Create post with maximum character limit
- [ ] Create post exceeding character limit
- [ ] Try to create post as free tier (limited)
- [ ] Create free tier post with restrictions
- [ ] Edit own post
- [ ] Delete own post
- [ ] View post details
- [ ] Share post
- [ ] Report inappropriate post

#### Timeline & Feed
- [ ] View home timeline
- [ ] Verify posts from followed users appear
- [ ] Verify posts from followed orgs appear
- [ ] Refresh timeline for new posts
- [ ] Infinite scroll on timeline
- [ ] Filter timeline by post type
- [ ] Verify blocked users' posts don't appear

#### Following System
- [ ] Follow user
- [ ] Unfollow user
- [ ] Try to follow already-followed user
- [ ] Follow organization
- [ ] Unfollow organization
- [ ] View following list
- [ ] View followers list
- [ ] Verify follow count updates

#### Voting System
- [ ] Upvote post
- [ ] Downvote post
- [ ] Remove vote from post
- [ ] Change vote from up to down
- [ ] Change vote from down to up
- [ ] Verify vote counts update in real-time
- [ ] Try to vote on own post

#### Comments System
- [ ] Add comment to post
- [ ] Add reply to comment
- [ ] Add comment exceeding character limit
- [ ] Add empty comment (should fail)
- [ ] Edit own comment (if supported)
- [ ] Delete own comment
- [ ] Delete comment on own post
- [ ] View all comments
- [ ] Load more comments (pagination)
- [ ] Report inappropriate comment
- [ ] Disable comments on own post
- [ ] Enable comments on own post

### Account Settings

#### Tier Management
- [ ] View current tier details
- [ ] View tier benefits comparison
- [ ] Upgrade from free to paid tier
- [ ] Verify paid features become available
- [ ] Handle payment failure scenarios
- [ ] Cancel paid subscription
- [ ] Verify downgrade to free tier

#### Privacy Settings
- [ ] Set profile visibility options
- [ ] Control who can see contact information
- [ ] Manage blocked users list
- [ ] Control who can send messages
- [ ] Set job seeking status visibility

## Employer Tests (harrypotter)

### Authentication & Onboarding

#### Employer Onboarding
- [ ] Check onboard status with invalid token
- [ ] Check onboard status with expired token
- [ ] Check onboard status with already-used token
- [ ] Set password not meeting requirements
- [ ] Set password with mismatched confirmation
- [ ] Successfully set password for new employer account
- [ ] Try to set password again after already set

#### Employer Sign In
- [ ] Sign in with invalid email format
- [ ] Sign in with non-existent employer email
- [ ] Sign in with incorrect password
- [ ] Sign in with correct credentials
- [ ] Verify TFA flow same as hub users
- [ ] Sign in and verify correct role permissions loaded
- [ ] Verify session timeout works correctly

#### Employer Password Management
- [ ] Change password flow (same validations as hub)
- [ ] Forgot password flow
- [ ] Reset password with invalid token
- [ ] Successfully reset password

### Organization Management

#### Cost Centers
- [ ] Add cost center with duplicate name
- [ ] Add cost center with empty name
- [ ] Add cost center with special characters
- [ ] Add cost center with very long name
- [ ] Successfully add cost center
- [ ] View all cost centers
- [ ] Filter cost centers by status
- [ ] Search cost centers by name
- [ ] Rename cost center to existing name
- [ ] Successfully rename cost center
- [ ] Update cost center details
- [ ] Defunct cost center with active openings
- [ ] Successfully defunct cost center
- [ ] Try to use defunct cost center in new opening
- [ ] Verify only authorized roles can manage cost centers

#### Locations
- [ ] Add location with incomplete address
- [ ] Add location with invalid postal code
- [ ] Add duplicate location
- [ ] Successfully add location with all details
- [ ] View all locations
- [ ] Filter locations by city/state
- [ ] Search locations
- [ ] Get specific location details
- [ ] Rename location
- [ ] Update location address
- [ ] Defunct location with active openings
- [ ] Successfully defunct location
- [ ] Verify location appears in opening creation

#### Organization Users
- [ ] Add org user with existing email
- [ ] Add org user with invalid email
- [ ] Add org user without required fields
- [ ] Add org user with invalid role
- [ ] Successfully add org user with Admin role
- [ ] Successfully add org user with limited role
- [ ] Update org user role
- [ ] Update org user details
- [ ] Disable active org user
- [ ] Enable disabled org user
- [ ] Try to disable last admin user
- [ ] Filter org users by role
- [ ] Filter org users by status
- [ ] Search org users by name/email
- [ ] Verify org user receives invitation email
- [ ] Complete org user signup from invitation

### Opening Management

#### Create Opening
- [ ] Create opening without required fields
- [ ] Create opening with invalid salary range (min > max)
- [ ] Create opening with past deadline
- [ ] Create opening with invalid cost center
- [ ] Create opening with defunct location
- [ ] Create opening with all required fields
- [ ] Create opening with all optional fields
- [ ] Add multiple VTags to opening
- [ ] Set experience requirements
- [ ] Set education requirements
- [ ] Preview opening before publishing
- [ ] Save opening as draft
- [ ] Publish opening immediately

#### Filter/Search Openings
- [ ] View all openings with no filter
- [ ] Filter by opening state (draft, active, closed)
- [ ] Filter by cost center
- [ ] Filter by location
- [ ] Filter by date range
- [ ] Search by job title
- [ ] Search by description keywords
- [ ] Sort by creation date
- [ ] Sort by application count
- [ ] Sort by deadline
- [ ] Paginate through results

#### Update Opening
- [ ] Update draft opening
- [ ] Update active opening
- [ ] Try to update closed opening
- [ ] Change salary range
- [ ] Add/remove VTags
- [ ] Extend deadline
- [ ] Try to set deadline in past
- [ ] Update job description
- [ ] Verify changes reflect immediately

#### Opening State Management
- [ ] Change from draft to active
- [ ] Change from active to paused
- [ ] Change from paused to active
- [ ] Close opening
- [ ] Try to reopen closed opening
- [ ] Verify state changes affect visibility

#### Opening Watchers
- [ ] View current watchers
- [ ] Add single watcher
- [ ] Add multiple watchers at once
- [ ] Add non-existent user as watcher
- [ ] Add user without permission as watcher
- [ ] Remove watcher
- [ ] Try to remove last watcher
- [ ] Verify watchers receive notifications

### Application Management

#### View Applications
- [ ] View all applications for opening
- [ ] Filter by application status
- [ ] Filter by color tags
- [ ] Filter by date range
- [ ] Search by candidate name
- [ ] Sort by application date
- [ ] Sort by match score
- [ ] View application details
- [ ] Download/view resume
- [ ] View candidate profile
- [ ] Bulk select applications
- [ ] Export applications list

#### Application Color Tags
- [ ] Set color tag on application
- [ ] Change existing color tag
- [ ] Remove color tag
- [ ] Filter by specific color tag
- [ ] Bulk tag applications
- [ ] Verify tag history tracked

#### Application Actions
- [ ] Shortlist application
- [ ] Reject application
- [ ] Move application between stages
- [ ] Add internal notes
- [ ] Bulk shortlist applications
- [ ] Bulk reject applications
- [ ] Undo rejection (if supported)
- [ ] Send message to candidate

### Interview Management

#### Schedule Interviews
- [ ] Add interview without required fields
- [ ] Add interview with past date/time
- [ ] Add interview with conflicting time
- [ ] Successfully schedule interview
- [ ] Add multiple interview rounds
- [ ] Set interview type (phone, video, in-person)
- [ ] Add interview location
- [ ] Add interview agenda/details

#### Interviewer Management
- [ ] Add interviewer to interview
- [ ] Add non-employee as interviewer
- [ ] Add multiple interviewers
- [ ] Remove interviewer
- [ ] Try to remove last interviewer
- [ ] Change primary interviewer
- [ ] Verify interviewers notified

#### Interview Feedback
- [ ] RSVP to interview as interviewer
- [ ] View interview details
- [ ] Submit assessment/feedback
- [ ] Update submitted assessment
- [ ] View other interviewers' assessments
- [ ] Rate candidate on various parameters
- [ ] Add detailed comments
- [ ] Recommend next action

### Candidacy Management

#### Candidacy Tracking
- [ ] View all candidacies for opening
- [ ] Filter candidacies by stage
- [ ] Filter by interview status
- [ ] Search candidacies
- [ ] View candidacy timeline
- [ ] Add comments to candidacy
- [ ] View candidacy history
- [ ] Track candidacy progress

#### Make Offer
- [ ] Create offer without required details
- [ ] Create offer with invalid salary
- [ ] Create offer with past joining date
- [ ] Successfully create offer
- [ ] Set offer expiry date
- [ ] Add offer terms
- [ ] Preview offer letter
- [ ] Send offer to candidate
- [ ] Revoke sent offer
- [ ] Extend offer deadline

### Employer Settings

#### Cool-off Period
- [ ] View current cool-off period
- [ ] Change to invalid period (negative)
- [ ] Change to very large period
- [ ] Successfully change cool-off period
- [ ] Verify changes apply to new applications

### Employer Posts

#### Create Posts
- [ ] Create post without content
- [ ] Create post exceeding length limit
- [ ] Create post with images
- [ ] Create post with links
- [ ] Preview post before publishing
- [ ] Save post as draft
- [ ] Schedule post for later
- [ ] Publish post immediately

#### Manage Posts
- [ ] View all posts
- [ ] Filter posts by status
- [ ] Edit published post
- [ ] Edit draft post
- [ ] Delete post
- [ ] View post analytics
- [ ] Share post
- [ ] Verify post appears on company page

### Role-Based Access

#### Admin Role
- [ ] Verify access to all features
- [ ] Test all CRUD operations
- [ ] Verify can manage other users

#### Limited Roles
- [ ] Verify CostCentersViewer can only view
- [ ] Verify LocationsCRUD cannot manage users
- [ ] Verify OpeningsViewer cannot edit
- [ ] Test each role's specific permissions
- [ ] Verify proper error messages for unauthorized actions

### Cross-Platform Features

#### Employer-Hub User Interaction
- [ ] View hub user profiles from applications
- [ ] View education history
- [ ] View achievements
- [ ] Download resume
- [ ] View profile picture
- [ ] Send messages (if supported)
- [ ] View endorsements

### Error Handling & Edge Cases

#### Network Issues
- [ ] Handle slow network connections
- [ ] Handle network timeouts
- [ ] Handle connection drops during form submission
- [ ] Verify proper retry mechanisms

#### Data Validation
- [ ] Test all form validations
- [ ] Test server-side validations
- [ ] Verify XSS prevention
- [ ] Test SQL injection prevention
- [ ] Verify CSRF protection

#### Browser Compatibility
- [ ] Test on Chrome
- [ ] Test on Firefox
- [ ] Test on Safari
- [ ] Test on Edge
- [ ] Test mobile browsers

#### Responsive Design
- [ ] Test on mobile devices
- [ ] Test on tablets
- [ ] Test on different desktop resolutions
- [ ] Verify touch interactions work
- [ ] Test landscape/portrait orientations

### Performance Tests

#### Page Load Times
- [ ] Test initial page loads
- [ ] Test navigation between pages
- [ ] Test with slow internet
- [ ] Verify lazy loading works

#### Search Performance
- [ ] Test with large result sets
- [ ] Test complex filter combinations
- [ ] Verify pagination performance
- [ ] Test real-time search

#### File Uploads
- [ ] Test large file uploads
- [ ] Test multiple file uploads
- [ ] Verify progress indicators
- [ ] Test upload cancellation

### Accessibility Tests

#### Screen Reader Compatibility
- [ ] Test all forms with screen reader
- [ ] Verify aria labels present
- [ ] Test keyboard navigation
- [ ] Verify focus indicators

#### Color Contrast
- [ ] Verify text readability
- [ ] Test color blind modes
- [ ] Verify error states visible

### Security Tests

#### Authentication Security
- [ ] Test session hijacking prevention
- [ ] Verify secure password storage
- [ ] Test brute force protection
- [ ] Verify TFA implementation

#### Data Privacy
- [ ] Verify PII protection
- [ ] Test data encryption
- [ ] Verify secure API calls
- [ ] Test authorization on all endpoints
