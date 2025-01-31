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
