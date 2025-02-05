export const en = {
  common: {
    home: "Home",
    openings: "Openings",
    logout: "Logout",
    loading: "Loading...",
    generalError: "Error",
    retry: "Retry",
    costCenters: "Cost Centers",
    locations: "Locations",
    actions: "Actions",
    add: "Add",
    save: "Save",
    cancel: "Cancel",
    loadMore: "Load More",
    serverError: "Please try again after some time.",
    none: "None",
    back: "Back",
    headerTitle: "Vetchi for Employers",
    success: "Success",
  },
  auth: {
    signin: "Sign In",
    domain: "Domain",
    email: "Email",
    password: "Password",
    tfa: "Two Factor Authentication",
    tfaCode: "Enter TFA Code",
    verify: "Verify",
    submit: "Submit",
    rememberMe: "Remember me",
    domainNotVerified:
      "Please add CNAME records to verify your domain with Vetchi.",
    domainVerifyPending:
      "Please ask your domain admin to check their email and complete the onboarding process.",
    accountDisabled: "Your account has been disabled.",
    invalidCredentials: "Invalid credentials.",
    unauthorized: "Unauthorized access",
  },
  dashboard: {
    welcome: "Welcome to your dashboard",
  },
  openings: {
    title: "Job Openings",
    noOpenings: "No openings found",
    create: "Create Opening",
    createTitle: "Create New Opening",
    openingTitle: "Title",
    positions: "Number of Positions",
    jobDescription: "Job Description",
    recruiter: "Recruiter",
    hiringManager: "Hiring Manager",
    costCenter: "Cost Center",
    type: "Opening Type",
    stateLabel: "State",
    minYoe: "Minimum Years of Experience",
    maxYoe: "Maximum Years of Experience",
    minEducation: "Minimum Education Level",
    employerNotes: "Employer Notes",
    locations: "Locations",
    physicalLocations: "Office Locations",
    remoteWork: "Remote Work",
    remoteTimezones: "Remote Timezones",
    remoteCountries: "Remote Countries",
    noLocationsError: "No locations specified",
    remoteTimezonesHelp:
      "Select the timezones where remote work is allowed. Leave empty if remote work is not allowed.",
    officeLocations: "Office Locations",
    officeLocationsHelp:
      "Select the office locations where this position is available. Leave empty if the position is fully remote.",
    globallyRemote: "Globally Remote (Available worldwide)",
    remoteCountriesHelp:
      "Select the countries where remote work is allowed. Leave empty if remote work is not allowed.",
    locationRequiredError:
      "Please select at least one location option (office locations, remote timezones, remote countries) or mark the position as globally remote.",
    showClosed: "Show Closed Openings",
    types: {
      FULL_TIME_OPENING: "Full Time",
      PART_TIME_OPENING: "Part Time",
      CONTRACT_OPENING: "Contract",
      INTERNSHIP_OPENING: "Internship",
      UNSPECIFIED_OPENING: "Unspecified",
    },
    education: {
      BACHELOR_EDUCATION: "Bachelor's Degree",
      MASTER_EDUCATION: "Master's Degree",
      DOCTORATE_EDUCATION: "Doctorate",
      NOT_MATTERS_EDUCATION: "Not Required",
      UNSPECIFIED_EDUCATION: "Unspecified",
    },
    state: {
      DRAFT_OPENING_STATE: "Draft",
      ACTIVE_OPENING_STATE: "Active",
      SUSPENDED_OPENING_STATE: "Suspended",
      CLOSED_OPENING_STATE: "Closed",
    },
    details: "Opening Details",
    id: "Opening ID",
    filledPositions: "Filled Positions",
    description: "Description",
    contacts: "Contact Information",
    actions: "Actions",
    publish: "Publish Opening",
    suspend: "Suspend Opening",
    reactivate: "Reactivate Opening",
    viewCandidacies: "View Candidacies",
    viewInterviews: "View Interviews",
    invalidStateTransition: "Invalid state transition",
    notFound: "Opening not found",
    fetchError: "Failed to fetch openings",
    createError: "Failed to create opening",
    fetchCostCentersError: "Failed to fetch cost centers",
    fetchLocationsError: "Failed to fetch locations",
    missingUserError: "Please select both recruiter and hiring manager",
    close: "Close Opening",
    closeConfirmTitle: "Close Opening",
    closeConfirmMessage:
      "Are you sure you want to close this opening? This action cannot be undone.",
    confirmClose: "Yes, Close Opening",
    stateChangeSuccess: "Opening state updated successfully",
    viewApplications: "View Applications",
    tags: "Opening Tags",
    selectTags: "Select Existing Tags",
    selectTagsPlaceholder: "Type to search or select tags...",
    tagsHelp: "Select up to 3 tags that best describe this opening",
    addNewTag: "Add New Tag",
    addNewTagPlaceholder: "Type a new tag and click Add",
    newTagHelp: "Can't find what you need? Add a new tag (max 3 tags total)",
    maxTagsError: "Maximum of 3 tags allowed (existing + new tags combined)",
    maxTagsReached: "Maximum tags reached (3)",
    noTagsFound: "No matching tags found",
    tagsRequiredError: "Please select at least one tag or add a new tag",
    fetchTagsError: "Failed to fetch opening tags",
    tagLengthError: "Tag must not exceed 32 characters",
  },
  costCenters: {
    title: "Cost Centers",
    addTitle: "Add Cost Center",
    editTitle: "Edit Cost Center",
    name: "Name",
    notes: "Notes",
    state: "State",
    add: "Add Cost Center",
    active: "Active",
    defunct: "Defunct",
    includeDefunct: "Show Defunct Cost Centers",
    noCostCenters:
      "No cost centers found. Click 'Add Cost Center' to create one.",
    fetchError: "Failed to fetch cost centers",
    addError: "Failed to add cost center",
    updateError: "Failed to update cost center",
    defunctError: "Failed to defunct cost center",
  },
  locations: {
    title: "Locations",
    addTitle: "Add Location",
    editTitle: "Edit Location",
    locationTitle: "Title",
    countryCode: "Country Code",
    countryCodeHelp: "Enter 3-letter ISO country code (e.g., USA, IND, GBR)",
    postalAddress: "Postal Address",
    postalCode: "Postal Code",
    mapUrl: "OpenStreetMap URL",
    cityAka: "Alternative City Names",
    cityAkaPlaceholder: "Enter alternative city name",
    state: "State",
    active: "Active",
    defunct: "Defunct",
    add: "Add Location",
    fetchError: "Failed to fetch locations",
    addError: "Failed to add location",
    updateError: "Failed to update location",
    defunctError: "Failed to defunct location",
    noLocations: "No locations found. Click 'Add Location' to create one.",
    includeDefunct: "Include Defunct Locations",
    viewMap: "View on Map",
  },
  validation: {
    title: {
      lengthError: "Title must be between 3 and 32 characters",
    },
    name: {
      length: {
        "2.64": "Name must be between 2 and 64 characters",
      },
      required: "Name is required",
    },
    email: {
      invalid: "Please enter a valid email address",
      required: "Email is required",
    },
    positions: {
      range: {
        "1.20": "Number of positions must be between 1 and 20",
      },
    },
    jobDescription: {
      lengthError: "Job description must be between 10 and 1024 characters",
    },
    employerNotes: {
      maxLength: {
        "1024": "Employer notes must not exceed 1024 characters",
      },
    },
    roles: {
      required: "At least one role must be selected",
    },
  },
  applications: {
    title: "Applications",
    filterByColor: "Filter by Color",
    allColors: "All Colors",
    colorGreen: "Green",
    colorYellow: "Yellow",
    colorRed: "Red",
    removeColor: "Remove Color",
    shortlist: "Shortlist",
    reject: "Reject",
    resumePreview: "Resume Preview",
    noApplications: "No applications found",
    setColor: "Set Color",
    noColor: "No color tag",
    clickToPreview: "Click to preview resume",
    pdfPreviewNotAvailable: "PDF preview not available",
  },
  candidacies: {
    title: "Candidacies",
    view: "View Candidacies",
    viewCandidacy: "View Candidacy",
    candidacyDetails: "Candidacy Details",
    noCandidacies: "No candidacies found for this opening.",
    fetchError: "Failed to fetch candidacies",
    applicantName: "Applicant Name",
    handle: "Handle",
    state: "State",
    filterPlaceholder: "Filter by",
    stateChanges: "Change State",
    makeOffer: {
      title: "Make Offer",
      description:
        "Upload an offer letter (PDF) and change the candidacy state to OFFERED. All pending interviews will be marked as cancelled.",
      button: "Make Offer",
      confirmTitle: "Confirm Make Offer",
      confirmDescription:
        "Are you sure you want to make an offer to this candidate?",
      uploadButton: "Upload Offer Letter (PDF)",
      selectedFile: "Selected file:",
      error: "Failed to make offer to candidate",
      success: "Offer has been successfully made to the candidate",
    },
    reject: {
      title: "Reject Candidacy",
      description:
        "Mark the candidate as unsuitable for this position. All pending interviews will be marked as cancelled.",
      button: "Reject Candidate",
      confirmTitle: "Confirm Reject Candidacy",
      confirmDescription: "Are you sure you want to reject this candidate?",
    },
    markUnresponsive: {
      title: "Mark as Unresponsive",
      description:
        "Mark the candidate as unresponsive if they have not been responding to communications. All pending interviews will be marked as cancelled.",
      button: "Mark Unresponsive",
      confirmTitle: "Confirm Mark as Unresponsive",
      confirmDescription:
        "Are you sure you want to mark this candidate as unresponsive?",
    },
    dialogActions: {
      confirm: "Confirm",
      cancel: "Cancel",
      dialogCancel: "Cancel",
    },
    dialogWarning: "This action cannot be undone.",
    dialogEffects: {
      title: "This action will:",
      cancelInterviews: "Mark all pending interviews as cancelled",
      stateChange: "Change the candidacy state to",
      uploadOffer: "Upload the offer letter",
    },
    states: {
      INTERVIEWING: "Interviewing",
      OFFERED: "Offered",
      OFFER_ACCEPTED: "Offer Accepted",
      OFFER_DECLINED: "Offer Declined",
      CANDIDATE_UNSUITABLE: "Candidate Unsuitable",
      CANDIDATE_NOT_RESPONDING: "Not Responding",
      CANDIDATE_WITHDREW: "Candidate Withdrew",
      EMPLOYER_DEFUNCT: "Employer Defunct",
    },
  },
  interviews: {
    title: "Interviews",
    addNew: "Add Interview",
    type: "Type",
    startTime: "Start Time",
    endTime: "End Time",
    state: "State",
    description: "Description",
    timezone: "Timezone",
    interviewers: "Interviewers",
    otherInterviewers: "Other Interviewers",
    yourRSVP: "Your RSVP status for Interview",
    noInterviewers: "No interviewers assigned",
    manage: "Manage Interview",
    manageInterview: "Manage Interview",
    details: "Interview Details",
    candidate: "Candidate",
    placeholder: "Interview management page is under construction",
    noInterviews: "No interviews scheduled",
    allowPastDates: "Allow setting dates in the past",
    endTimeBeforeStart: "End time cannot be before or equal to start time",
    use24HourFormat: "Use 24-hour time format",
    you: "You",
    types: {
      VIDEO_CALL: "Video Call",
      IN_PERSON: "In Person",
      TAKE_HOME: "Take Home",
      OTHER_INTERVIEW: "Other",
    },
    states: {
      SCHEDULED_INTERVIEW: "Scheduled",
      COMPLETED_INTERVIEW: "Completed",
      CANCELLED_INTERVIEW: "Cancelled",
    },
    addError: "Failed to add interview",
    fetchError: "Failed to fetch interviews",
    rsvp: {
      yes: "Accept",
      no: "Decline",
      confirmYes: "Accept Interview",
      confirmNo: "Decline Interview",
      confirmYesMessage:
        "Are you sure you want to accept this interview? The employer will be notified of your response.",
      confirmNoMessage:
        "Are you sure you want to decline this interview? The employer will be notified of your response.",
      confirmChangeYesMessage:
        "Are you sure you want to change your response to accept? The employer will be notified of this change.",
      confirmChangeNoMessage:
        "Are you sure you want to change your response to decline? The employer will be notified of this change.",
    },
    assessment: {
      title: "Interview Assessment",
      rating: "Rating",
      ratingPlaceholder: "Select your rating for the candidate",
      editFeedback: "Edit Feedback",
      feedback: "Feedback to Candidate (Public)",
      feedbackPlaceholder:
        "Enter feedback that will be shared with the candidate. This feedback will be visible to the candidate.",
      positives: "Positives",
      positivesPlaceholder:
        "Enter positive aspects of the candidate's performance",
      negatives: "Negatives",
      negativesPlaceholder:
        "Enter negative aspects of the candidate's performance",
      overallAssessment: "Overall Assessment",
      overallAssessmentPlaceholder:
        "Enter your overall assessment of the candidate",
      save: "Save Assessment",
      saveSuccess: "Assessment saved successfully",
      saveError: "Failed to save assessment",
      notFoundError: "Interview not found",
      validationError: "Please check your input and try again",
      forbiddenError: "You do not have permission to update this assessment",
      fetchError: "Failed to fetch assessment",
      lastUpdated: "Last updated by {{name}} on {{date}}",
      rsvpError: "Failed to update RSVP status",
      rsvpSuccess: "RSVP status updated successfully",
      invalidStateError: "Interview is not in a valid state for RSVP",
      markAsCompleted: "Mark Interview as Completed",
      ratings: {
        STRONG_YES: "Strong Yes",
        YES: "Yes",
        NEUTRAL: "Neutral",
        NO: "No",
        STRONG_NO: "Strong No",
      },
    },
  },
  comments: {
    title: "Comments",
    add: "Add Comment",
    addPlaceholder: "Write your comment here...",
    noComments: "No comments yet",
  },
  orgUsers: {
    title: "Users",
    addTitle: "Add Organization User",
    email: "Email",
    name: "Name",
    rolesList: "Roles",
    state: "State",
    add: "Add User",
    disable: "Disable User",
    enable: "Enable User",
    noUsers: "No users found. Click 'Add User' to create one.",
    fetchError: "Failed to fetch organization users",
    addError: "Failed to add user",
    disableError: "Failed to disable user",
    enableError: "Failed to enable user",
    searchPlaceholder: "Search by email or name",
    addSuccess: "User added successfully",
    disableSuccess: "User disabled successfully",
    enableSuccess: "User enabled successfully",
    includeDisabled: "Show Disabled Users",
    states: {
      ACTIVE_ORG_USER: "Active",
      ADDED_ORG_USER: "Added",
      DISABLED_ORG_USER: "Disabled",
    },
    roles: {
      ADMIN: "Admin",
      ORG_USERS_CRUD: "Org Users Manager",
      ORG_USERS_VIEWER: "Org Users Viewer",
      COST_CENTERS_CRUD: "Cost Centers Manager",
      COST_CENTERS_VIEWER: "Cost Centers Viewer",
      LOCATIONS_CRUD: "Locations Manager",
      LOCATIONS_VIEWER: "Locations Viewer",
      OPENINGS_CRUD: "Openings Manager",
      OPENINGS_VIEWER: "Openings Viewer",
      APPLICATIONS_CRUD: "Applications Manager",
      APPLICATIONS_VIEWER: "Applications Viewer",
    },
    confirmDisable: {
      modalTitle: "Disable User",
      message:
        "Are you sure you want to disable this user? They will no longer be able to access the system.",
      confirmButton: "Yes, Disable User",
      cancelButton: "Cancel",
    },
  },
};
