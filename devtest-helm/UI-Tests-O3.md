# UI End-to-End Test Matrix

This document lists **exhaustive UI test cases** for the two public-facing Next.js front-ends that power Vetchium:

* **Ronweasly** – Hub user site (free & paid tiers)
* **Harrypotter** – Employer / organisation site

Use this list as the canonical checklist when creating Playwright / Cypress suites or when performing manual regression testing.  
Wherever sensible tests are grouped under the same logical workflow – the checkbox list is still granular enough so that each item can fail or pass independently.

---

## Ronweasly – Hub User Workflows

### 1. Authentication
#### 1.1 Hub Login
- [ ] Attempt to load the `/login` page while already authenticated – expect automatic redirect to home dashboard.
- [ ] Input email in invalid format (`foo@`) – **Sign-in button** must remain disabled.
- [ ] Input non-existent email + any password – server returns **Invalid credentials** error toast.
- [ ] Input email longer than max length ( > 320 chars ) – client-side validation error and button disabled.
- [ ] Input password shorter than 8 chars – client-side validation prevents submission.
- [ ] Input extremely long password ( > 128 chars ) – client still allows, server rejects with *Invalid credentials*.
- [ ] Input password missing required complexity – shows *Password does not meet requirements*.
- [ ] Valid credentials – redirected to `/tfa` page and TOTP sent.

#### 1.2 Two-Factor Authentication (TFA)
- [ ] Stay idle on `/tfa` for 2 min until code expires – submitting returns *Code expired* error.
- [ ] Enter code of incorrect length (<6 and >6) – inline validation blocks submit.
- [ ] Enter wrong 6-digit code – server responds *Invalid code*.
- [ ] Enter exact valid code within expiry – user lands on `/` home timeline.
- [ ] Trigger *Resend code* – new email/sms arrives and old code becomes invalid.

#### 1.3 Forgot & Reset Password
- [ ] From `/forgot-password`, request reset for unknown email – success toast but no email sent.
- [ ] Request reset for registered email – reset email delivered containing single-use token.
- [ ] Click expired / malformed token link – `/reset-password` shows *Link invalid or expired*.
- [ ] Open valid link – new password form loads with masked e-mail.
- [ ] Submit mismatching `password` and `confirm password` – client validation error.
- [ ] Submit weak password (no special char) – server rejects.
- [ ] Submit strong password – redirected to `/login` with *Password updated* toast.

#### 1.4 Change Password ( while logged-in )
- [ ] Visit `/settings` → *Change Password* card.
- [ ] Enter wrong current password – error *Current password incorrect*.
- [ ] New password identical to current – disabled save.
- [ ] New password weak – rejected.
- [ ] Successful change – confirmation toast & forced re-login on next privileged call.

### 2. On-Boarding & Account Setup
#### 2.1 Signup
- [ ] Access `/signup` (invited flows) while logged-in – redirect home.
- [ ] Submit email that already exists – inline *Email already registered*.
- [ ] Submit password not meeting policy – validation error.
- [ ] Happy path – account created, redirected to `/upgrade` when tier is pending payment.

#### 2.2 Handle & Tier
- [ ] For paid tier user without handle, visit `/settings` – *Set Handle* modal must pop.
- [ ] Attempt handle shorter than 3 chars, contains spaces, or uppercase – rejected.
- [ ] Attempt handle that is already taken – server error.
- [ ] Successful handle set – profile accessible via `/u/<handle>`.
- [ ] Downgrade paid → free tier – feature gated buttons (upload profile picture, paid posts) disappear.

### 3. Profile Management
#### 3.1 Bio
- [ ] Empty bio – *Save* disabled.
- [ ] Exceed max length (2000 chars) – client blocks.
- [ ] Valid markdown formatted bio – persists and renders markdown.

#### 3.2 Profile Picture
- [ ] Free tier attempts upload – CTA hidden.
- [ ] Paid tier uploads >5 MB JPEG – client rejects file.
- [ ] Upload square PNG ≤5 MB – picture displays in avatar across site.
- [ ] Remove picture – avatar falls back to initials.
- [ ] Access `/profile-picture/<uuid>` anonymously – public HTTP 200 image or 403 depending on privacy.

#### 3.3 Official Email Addresses
- [ ] Add address on free tier – allowed.
- [ ] Enter non-corporate domain when policy disallows – server rejects.
- [ ] Successful add sends verification email; unverified badge shows.
- [ ] Click verification link – status flips to verified.
- [ ] Delete verified email – removed from list.

### 4. Education
- [ ] Autocomplete institute search throttles API and keyboard navigation works.
- [ ] Add education with end year before start year – client validation error.
- [ ] Add duplicate entry – deduplicated on list.
- [ ] Delete education – entry disappears.
- [ ] List entries sorted by end year desc.

### 5. Achievements
- [ ] Add achievement with overlong title – blocked.
- [ ] Attach public URL – link rendered clickable.
- [ ] Delete achievement – confirmation modal then list refreshes.

### 6. Work History
- [ ] Add history entry with future start date – blocked.
- [ ] Leaving *end date* empty marks as *Current*.
- [ ] Update entry successfully.
- [ ] Delete entry – reflect immediately.

### 7. Social Graph – Colleagues
#### 7.1 Connect / Approve
- [ ] Search colleague by email not on platform – informative toast.
- [ ] Send connect request to valid colleague – request appears in their *My Approvals*.
- [ ] Approve request – both users now show in *My Colleagues*.
- [ ] Reject request – status set to *Rejected* and cannot re-request for 7 days (cool off).

#### 7.2 Endorsements
- [ ] Attempt endorsement without colleague link – CTA disabled.
- [ ] Endorse application – status shows *Pending employer approval*.
- [ ] Withdraw endorsement before employer action – ability allowed.

### 8. Posts & Timeline
#### 8.1 Creating Posts
- [ ] Free tier: open `/posts/new` – only plain-text post supported; media upload disabled.
- [ ] Paid tier: can embed image; oversize image >10 MB blocked.
- [ ] Post containing more than 5000 chars – blocked.
- [ ] Successful post appears at top of `/posts` and followers' timelines.

#### 8.2 Follow / Unfollow User
- [ ] Follow same user twice – second click becomes *Unfollow* and only one DB entry.
- [ ] Unfollow removes their posts from home timeline after refresh.

#### 8.3 Voting
- [ ] Upvote own post – action disabled.
- [ ] Upvote then downvote toggles correctly & score updates.
- [ ] Unvote returns score to original.

#### 8.4 Comments
- [ ] Add comment > 1000 chars – blocked.
- [ ] Disable comments on own post – *Add Comment* box hidden for others.
- [ ] Delete own comment – removed and replies orphan properly.

### 9. Openings & Applications
#### 9.1 Discover & Apply
- [ ] Visit `/find-openings` – infinite scroll fetches next page.
- [ ] Filter by virtual tag; query parameter persists in URL.
- [ ] View opening details at `/org/<domain>/opening/<id>` – contact & JD visible.
- [ ] Apply to closed opening – *Apply* button disabled.
- [ ] First-time apply requires resume upload; missing file -> error.
- [ ] Successful apply – toast and entry appears in `/my-applications`.
- [ ] Withdraw application – status changes to *Withdrawn* and employer notified.

#### 9.2 Application & Candidacy Tracking
- [ ] `/my-applications` list paginates correctly.
- [ ] `/my-candidacies` shows employer feedback columns.
- [ ] Candidacy details `/candidacy/<id>` loads timeline, interview schedule.
- [ ] Add comment with @mention colleague – autocomplete suggestions appear.

### 10. Interviews
- [ ] Accept interview RSVP – status updates, Google/iCal button appears.
- [ ] Decline interview – requires mandatory reason textarea enforced.
- [ ] Attempt RSVP after deadline – disabled with tooltip.

### 11. Notifications & Real-Time
- [ ] After colleague approval, toast + badge counter increments without reload via SSE/WebSocket.
- [ ] Background tab receives title badge increments when new post appears on timeline.

---

## Harrypotter – Employer Workflows

### 1. Authentication & Org-User Management
#### 1.1 Employer Sign-in
- [ ] `/signin` loads while authenticated redirecting to dashboard.
- [ ] Invalid email format blocks submit.
- [ ] Valid creds redirect to `/tfa`.
- [ ] TFA flows identical to hub tests, including expiry.

#### 1.2 Forgot & Reset Password
- [ ] Same negative and positive flows as hub.

#### 1.3 Change Password
- [ ] Located in `/settings` – verify flows like hub; ensure all roles can access.

#### 1.4 Signup Org User
- [ ] Access `/signup-orguser` with invalid invite token – *Invalid invite* page.
- [ ] Valid token but already used – *Already accepted* message.
- [ ] Complete form ( first name, last name, password ) – redirects to `/signin`.

### 2. Dashboard Navigation
- [ ] Side nav highlights current section.
- [ ] Collapse/expand persists in `localStorage`.
- [ ] 404 route inside dashboard shows branded *Page Not Found*.

### 3. Cost Centres
- [ ] `/cost-centers` shows list sorted alphabetically.
- [ ] Add duplicate code – server rejects.
- [ ] Rename with empty string – blocked.
- [ ] Mark cost centre defunct – item visually greyed & hidden from *Create Opening* dropdown.
- [ ] Update budget field – persisted.

### 4. Locations
- [ ] Similar CRUD tests as cost centres plus geo-autocomplete for city / country.
- [ ] Defunct location removed from filters.

### 5. Org Users
- [ ] Invite new org user with role *OpeningsViewer* – email invitation issued.
- [ ] Disable user – login attempt by that user blocked.
- [ ] Re-enable restores access.
- [ ] Change role set; verify RBAC – user without `CostCentersCRUD` cannot access `/cost-centers` create.

### 6. Openings Lifecycle
#### 6.1 Create Opening
- [ ] Required fields: title, location, cost centre, description – blank blocks submit.
- [ ] Title longer than 256 chars – blocked.
- [ ] Add virtual tags via autocomplete; new tag allowed if `Enter` pressed.
- [ ] Save draft keeps *Draft* state; publish sets *Open* and email watchers.

#### 6.2 Update & Watchers
- [ ] Add/remove watchers – user list filters by active users only.
- [ ] Changing owner to disabled org user – blocked.

#### 6.3 State Transitions
- [ ] Transition *Open* → *On Hold* requires reason modal – saved in audit log.
- [ ] *On Hold* → *Closed* – verify applications cannot be added post close.

#### 6.4 Filter / Search Openings
- [ ] Filter by cost centre combines with search query correctly.
- [ ] Pagination – page 2 retains filters when reloading.

### 7. Applications Management
- [ ] `/openings/<id>/applications` list default sorts by created desc.
- [ ] Download resume as PDF link functions.
- [ ] Set colour tag – chip appears in list.
- [ ] Remove colour tag – removed.
- [ ] Shortlist application moves row to *Shortlisted* tab; duplicate shortlist action disabled.
- [ ] Reject application with comment – status set *Rejected* and candidate email sent.

### 8. Candidacy & Interviews
- [ ] `/candidacy/<id>` loads timeline and comments.
- [ ] Add internal comment > 2000 chars – blocked.
- [ ] Offer to candidate button enabled only when state == *Shortlisted*.
- [ ] Schedule interview – date picker cannot select past date.
- [ ] Add interviewer: autocomplete searches org users; cannot add same interviewer twice.
- [ ] Remove interviewer confirmation.
- [ ] Assessment form: required rating fields, supports draft save.
- [ ] Employer RSVP flows mirror hub RSVP.

### 9. Employer Posts
- [ ] `/posts` list shows org posts with pagination.
- [ ] Add post with image >10 MB – blocked.
- [ ] Editing post updates timestamp labelled *Edited*.
- [ ] Delete post – confirmation modal; hub user attempting to fetch deleted post gets 404.

### 10. Hub User Insights
- [ ] From application row, click candidate handle – opens `/u/<handle>` within employer dashboard.
- [ ] Profile picture fallback identical to hub.
- [ ] Education & achievements accordions collapse/expand.

### 11. Settings
- [ ] Change *Cool-Off Period* accepts days 1-90 only.
- [ ] Setting persists across browser refresh and affects colleague re-invite flows on hub side.

### 12. Access Control Smoke Tests
- [ ] Org user with `OpeningsViewer` tries to POST `/employer/create-opening` – 403 shown UI.
- [ ] Org user with `Admin` can access all sections.

---

### Cross-Cutting Concerns
- [ ] All pages pass axe-core accessibility scan with no critical issues.
- [ ] All forms prevent double submission (button disabled while busy).
- [ ] Network failures show retry banners not silent failure.
- [ ] Mobile viewport (375×667) layouts have no horizontal scrolling.
- [ ] Dark mode toggle (if enabled) persists and colours contrast AA.

---

> ⚠️ **Maintenance note**: Keep this document in sync whenever a new front-end page or back-end route is added.  
> Add both *happy path* and *edge / failure* cases.  
> When automating tests, reference the exact API response payloads defined in `typespec` to avoid brittle assertions.
