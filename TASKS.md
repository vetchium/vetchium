# Done

- [x] Employer Domain Verification
- [x] Employer Signup
- [x] Employer Signin
- [x] Locations, CostCenters, etc.
- [x] Openings CRUD
- [x] Applications CRUD
- [x] Candidacy CRUD
- [x] Hub user Login
- [x] Find Openings
- [x] Communication on Candidacies (Realtime not done yet)
- [x] Work history and Official emails for Hub
- [x] Profile Photos on Hub
- [x] Colleague Connections and various states and flows
- [x] Interviews - Scheduling, RSVPs
- [x] Candidacy state changes after Interviews
- [x] Rudimentary scoring of Applicants
- [x] Achievements - Patents, Posts, Certifications
- [x] Posts
- [x] Timelines and Follows !?
- [x] Automated Scale testing
- [x] Upvotes and Downvotes

# Features needed before launch

- [] Better Email Templates
- [] Database Profiling and Indexes
- [] CandiateOffer RSVP and some signing (Docusign, Zohosign, etc.) integrations
- [] Find good matches for suitable Openings for HubUsers
- [] Scoring of match for Applicants that is beyond just text compare !?
- [] Academia support to add educational institutions, degrees, etc.
- [] Articles
- [] VibeCheck for posts
- [] Bulk load of employer (F500+ more) names & domains from a configmap !? Or should we seed in the db ?!
- [] Audit logs
- [] Introductory videos for Employer as well as the Product
- [] Tag specific fetching of posts
- [] Follow support for Orgs
- [] Post message support for Orgs
- [] GreatHall support (upvotes, downvotes, anonymous comments, etc.)
- [] Tags is a bit of a mess. Should we allow everyone to create tags ? Should we gatekeep ? Should we (maintainers) hardcode a seed list and keep it updated ?
- [] Cleaner RBAC error codes and UI alerts
- [] LoadTests should test for EmployerPosts also
- [] API for EmployerPost details getting
- [] HubUser account deletion and the triggered cleanup workflows for posts, comments, etc.
- [] Employer account deletion and the triggered cleanup workflows for posts, comments, openings, candidacies, applications, etc.

# Future Features

## 1. Advanced Search and Filtering for Job Seekers

The API schema shows basic search functionality (`/hub/find-openings`), but could be enhanced with:

- Skills-based search (currently not in schema)
- Salary range filtering (schema has salary structures)
- Remote work preferences (schema has `remote_country_codes` and `remote_timezones`)

The current schema in `FindHubOpeningsRequest` shows basic filters for location, company, and experience, but modern job seekers expect more granular search capabilities.

## 2. Application Tracking System (ATS) Enhancement

The schema has basic application states (`ApplicationState` enum) but could benefit from:

- Stage-based tracking (beyond current APPLIED/REJECTED/SHORTLISTED states). The `ApplicationState` enum in the schema shows basic states, but modern hiring processes typically have more granular stages.
- Automated status notifications (building on existing application state changes)

## 3. Analytics Dashboard for Employers

Based on the existing endpoints, there's no analytics functionality. Could add:

- Application funnel metrics
- Time-to-hire tracking
- Source of applicants
- Opening performance metrics

HRs love having and showing reports.

## 4. Enhanced Interview Management

While basic interview functionality exists (`InterviewType` enum and related endpoints), could add:

- Calendar integration
- Automated reminder system
- Interview feedback templates
- Interview scoring standardization
- Automated AI feedback suggestions !?
- Jitsi or some such integration to have meetings within ourselves
- Integrations with zoom, hangouts, calendly, etc. for meeting urls, recordings, freebusy, scheduling, etc.

The schema has interview management (`/employer/add-interview`, `/employer/put-assessment`) but lacks integration features.

## 5. Candidate Communication Hub

Current schema shows basic commenting (`/employer/add-candidacy-comment`), but could add:

- Structured messaging system
- Template-based communications
- Bulk candidate communications
- Email integration

The schema has basic commenting through `/employer/add-candidacy-comment` and `/hub/add-candidacy-comment` but lacks comprehensive communication features.

## 6. Resume Parsing and Management

The schema shows basic resume handling (`/employer/get-resume`) but could add:

- Automated skill extraction
- Resume scoring
- Candidate matching
- Resume formatting standardization

Current schema only shows basic resume storage/retrieval through `/employer/get-resume` endpoint.

## 7. Enhanced Employer Branding

Current schema lacks employer branding features. Could add:

- Company profile management
- Culture page customization
- Employee testimonials
- Benefits showcase

## 8. Domain Ownership Changes

When a domain say x.com was owned by Paypal first and we had some employees verified their email address. Then the domain got transferred to Twitter. Now what happens to the old email addresses that were verified when it was with Paypal ? Will those employees show up as Twitter employees now ? We need a cleaner solution to deboard, re-onboard domains which are onboarded. There are a lot of corner cases involved in this kind of domain cross-ownerships.

## 9. SSO and Directory sync for Employers

Keycloak, Ory Kratos etc. used to do these, may not be worth building it from scratch. Needs investigation

### 10. Better Prometheus integration

### 11. Markdown support

### 12. Payment support

Depends on incorporation to choose a payment vendor.

### 13. DB Sharding based on regions, to support Data Sovereignity across countries

### 14. DB monitoring

### 15. Remove storing of the official email id in the db and use just hashes, for better privacy. Will have impact in "Potential Team Mates" in Opening creation too

# Code Cleanups
* Better templating
* The number of files under the postgres package has increased. Should we consider breaking it into subdirectories, like we do for handlers ? Does it offer any readability improvements ? Or is the current method best for passing to AI IDEs etc. ?
* Better validation handling in the typespec/**/*.ts files with a builtin IsValid function; instead of adding validation in the nextjs files
