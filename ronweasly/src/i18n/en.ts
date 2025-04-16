export const en = {
  common: {
    login: "Sign In",
    logout: "Logout",
    email: "Email Address",
    password: "Password",
    verify: "Verify",
    search: "Search",
    loading: "Loading...",
    cancel: "Cancel",
    approve: "Approve",
    reject: "Reject",
    warning: "Warning",
    proceed: "Proceed",
    external_url_warning:
      "You are about to leave Vetchium and visit an external website. This link has not been verified by Vetchi. Please proceed with caution.",
    externalLink: {
      message: "This link will open in a new tab",
      warning: "You are about to open this link in a new tab",
    },
    error: {
      notAuthenticated: "Not authenticated. Please log in again.",
      sessionExpired: "Session expired. Please log in again.",
      serverError: "The server is experiencing issues. Please try again later",
      requiredField: "This field is required",
    },
    retry: "Retry",
    back: "Back",
    actions: "Actions",
  },
  auth: {
    tfa: {
      title: "Two-Factor Authentication",
      description: "Please enter the verification code sent to your email",
      codeLabel: "Verification Code",
    },
    errors: {
      invalidCredentials: "Invalid credentials",
      accountDisabled: "Your account has been disabled",
      serverError: "The server is experiencing issues. Please try again later",
      invalidTfaCode: "Invalid verification code",
    },
    loginFailed: "Login failed",
    tfaFailed: "TFA verification failed",
    backToLogin: "Back to Login",
    rememberMe: "Remember me",
  },
  officialEmails: {
    title: "Official Emails",
    addEmail: "Add Official Email",
    verifyEmail: "Verify Email",
    verificationCode: "Verification Code",
    verificationPending: "Verification Pending",
    verifyButton: "Verify",
    enterCode: "Enter Code",
    sendCode: "Send Code",
    enterVerificationCode: "Please enter the verification code sent to {email}",
    noEmails: "No official emails added yet",
    deleteEmail: "Delete email",
    addEmailSubmit: "Add email",
    verifyEmailSubmit: "Verify email",
    verificationExpired: "Verification expired",
    verifiedOn: "Verified on {date}",
    errors: {
      loadFailed: "Failed to load official emails",
      addFailed: "Failed to add email",
      emailExists: "This email is already registered",
      domainNotEmployer:
        "No employer with this domain found. Please add a work experience with this domain before adding it as an official email.",
      deleteFailed: "Failed to delete email",
      triggerFailed: "Failed to send verification code",
      invalidCode: "Invalid verification code",
      invalidCodeLength: "Code must be 4 characters",
    },
  },
  dashboard: {
    title: "Dashboard",
    recentApplications: "Recent Applications",
    activeOpenings: "Active Openings",
    upcomingInterviews: "Upcoming Interviews",
  },
  findOpenings: {
    title: "Find Openings",
    description: "Search for job openings across all locations",
    searchPlaceholder: "Search for job titles, skills, or keywords",
    noOpeningsFound: "No openings found",
  },
  myApplications: {
    title: "My Applications",
    noApplications: "You haven't applied to any openings yet",
    error: {
      loadFailed: "Failed to load your applications. Please try again later.",
    },
    applicationState: {
      applied: "Applied",
      rejected: "Rejected",
      shortlisted: "Shortlisted",
      withdrawn: "Withdrawn",
      expired: "Expired",
    },
    appliedOn: "Applied on {date}",
    viewOpening: "View Opening",
    withdrawApplication: "Withdraw Application",
    withdrawConfirmation: "Are you sure you want to withdraw this application?",
    withdrawSuccess: "Application withdrawn successfully",
    withdrawError: "Failed to withdraw application. Please try again later.",
  },
  openingDetails: {
    notFound: "Opening not found",
    hiringManager: "Hiring Manager",
    yearsExperience: "{min}-{max} years experience",
    apply: "Apply for this Opening",
    cannotApply: "You cannot apply to this Opening",
    educationLevel: {
      bachelor: "Bachelor's Degree",
      master: "Master's Degree",
      doctorate: "Doctorate",
      notMatters: "Any Education Level",
      unspecified: "Not Specified",
    },
    openingType: {
      fullTime: "Full Time",
      partTime: "Part Time",
      contract: "Contract",
      internship: "Internship",
      unspecified: "Not Specified",
    },
    state: {
      draft: "This opening is not yet active",
      suspended: "Applications are temporarily suspended",
      closed: "This opening is no longer accepting applications",
    },
    error: {
      loadFailed: "Failed to load opening details. Please try again later.",
      pdfOnly: "Please upload a PDF file only",
      fileTooLarge: "File size should be less than 5MB",
      noResume: "Please select a resume to upload",
      applyFailed: "Failed to apply for the opening. Please try again later.",
    },
    selectResume: "Select Resume (PDF)",
    resumeSelected: "Selected: {name}",
    endorsers: {
      title: "Endorsers",
      description:
        "Add up to 5 verified colleagues who can endorse your application. This can increase your chances of getting selected.",
      search: "Search for colleagues",
      maxReached: "Maximum 5 endorsers reached",
      remaining: "{count} more endorsers can be added",
      noColleagues:
        "No verified colleagues found. Connect with colleagues to add endorsers.",
    },
  },
  navigation: {
    home: "Home",
    findOpenings: "Find Openings",
    posts: "Posts",
    myApplications: "My Applications",
    myCandidacies: "My Candidacies",
    myProfile: "My Profile",
    myApprovals: "My Approvals",
    myRequisitions: "My Requisitions",
    settings: "Settings",
  },
  candidacies: {
    viewCandidacy: "View Candidacy",
    viewDetails: "View candidacy details",
    fetchError: "Failed to load candidacy details",
    noCandidacies: "You don't have any candidacies yet",
    states: {
      INTERVIEWING: "Interviewing",
      OFFERED: "Offered",
      OFFER_ACCEPTED: "Offer Accepted",
      OFFER_DECLINED: "Offer Declined",
      CANDIDATE_UNSUITABLE: "Not Selected",
      CANDIDATE_NOT_RESPONDING: "Not Responding",
      CANDIDATE_WITHDREW: "Withdrawn",
      EMPLOYER_DEFUNCT: "Position Closed",
    },
  },
  comments: {
    title: "Comments",
    noComments: "No comments yet",
    addPlaceholder: "Add a comment...",
    add: "Add Comment",
  },
  interviews: {
    title: "Interviews",
    noInterviews: "No interviews scheduled",
    fetchError: "Failed to load interviews",
    timeRange: "{start} - {end}",
    timezone: "Timezone: {zone}",
    interviewers: "Interviewers",
    noInterviewers: "No interviewers assigned",
    endTime: "End Time",
    details: "Interview Details",
    yourRSVP: "Your RSVP Status",
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
    rsvpError: "Failed to update RSVP status",
    states: {
      SCHEDULED_INTERVIEW: "Scheduled",
      COMPLETED_INTERVIEW: "Completed",
      CANCELLED_INTERVIEW: "Cancelled",
    },
    types: {
      IN_PERSON: "In Person",
      VIDEO_CALL: "Video Call",
      TAKE_HOME: "Take Home Assignment",
      OTHER_INTERVIEW: "Other",
    },
  },
  workHistory: {
    title: "Work Experience",
    addExperience: "Add Experience",
    updateExperience: "Update Experience",
    loading: "Loading work history...",
    noEntries: "No work experience entries yet",
    companyDomain: "Company Domain",
    jobTitle: "Title",
    startDate: "Start Date",
    endDate: "End Date",
    description: "Description",
    present: "Present",
    deleteConfirm: "Are you sure you want to delete this work history entry?",
    error: {
      fetchFailed: "Failed to fetch work history",
      saveFailed: "Failed to save work history",
      deleteFailed: "Failed to delete work history",
    },
    actions: {
      edit: "Edit",
      delete: "Delete",
      save: "Save",
      cancel: "Cancel",
    },
  },
  profile: {
    myProfile: "My Profile",
    editMyProfile: "Edit My Profile",
    actions: "Actions",
    addAsColleague: "Add as Colleague",
    cannotAddAsColleague:
      "You need a verified email in a common domain to connect with this person",
    requestPending: "Connection request pending",
    receivedColleagueRequest:
      "This person wants to connect with you as a colleague",
    approveRequest: "Approve Request",
    declineRequest: "Decline Request",
    mutuallyVerifiedColleague: "Mutually verified Colleague",
    unlinkConnection: "Unlink Connection",
    youRejectedTheirRequest: "You previously rejected their connection request",
    theyRejectedYourRequest: "They previously rejected your connection request",
    youUnlinkedConnection: "You previously unlinked this connection",
    theyUnlinkedConnection: "They previously unlinked this connection",
    error: {
      userNotFound: "User not found",
      cannotConnect: "Cannot connect with this user at this time",
      connectionFailed: "Failed to send colleague request",
      approvalFailed: "Failed to approve colleague request",
      rejectFailed: "Failed to reject colleague request",
      noRequestFound: "No pending request found",
      unlinkFailed: "Failed to unlink colleague connection",
      noConnectionFound: "No active connection found with this colleague",
      handleMismatch: "The handle you entered does not match",
    },
    bio: {
      error: {
        fetchFailed: "Failed to load profile information",
        updateFailed: "Failed to update profile information",
        uploadFailed: "Failed to upload profile picture",
      },
      title: "Bio",
      fullName: "Full Name",
      handle: "Handle",
      shortBio: "Short Bio",
      longBio: "Long Bio",
      save: "Save",
      cancel: "Cancel",
      verifiedDomains: "Verified Domains",
      verifiedDomainsInfo:
        "These are domains where this user has verified their email address ownership. This verification helps establish the user's professional affiliations and work history authenticity.",
    },
    picture: {
      change: "Change Profile Picture",
      upload: "Upload Profile Picture",
      remove: "Remove Profile Picture",
      removeFailed: "Failed to remove profile picture",
      removeConfirmTitle: "Remove Profile Picture?",
      removeConfirmMessage:
        "Are you sure you want to remove your profile picture? This action cannot be undone.",
      removeConfirm: "Yes, Remove Picture",
      fullSize: "View Full Size",
      upgradePrompt: "Profile picture upload is available for paid users.",
      upgradeLink: "Upgrade your account.",
    },
    unlinkConfirmTitle: "Unlink Colleague Connection",
    unlinkConfirmMessage:
      "To unlink your connection with {handle}, please type their handle below to confirm.",
    unlinkConfirmHandleLabel: "Type handle to confirm",
    unlinkConfirm: "Unlink Connection",
  },
  approvals: {
    title: "My Approvals",
    colleagueApprovals: "Colleague Approvals",
    endorsementApprovals: "Endorsement Approvals",
    noApprovals: "No pending approvals",
    noEndorsements: "No pending endorsement requests",
    error: {
      fetchFailed: "Failed to fetch approvals",
      endorsementActionFailed: "Failed to process endorsement request",
    },
    endorsement: {
      from: "From",
      for: "For",
      at: "at",
      appliedOn: "Applied on {date}",
      viewOpening: "View Opening",
    },
  },
  requisitions: {
    title: "My Requisitions",
    colleagueSeeks: "Colleague Connection Requests",
    noSeeks: "No pending connection requests",
    error: {
      fetchFailed: "Failed to fetch connection requests",
    },
  },
  settings: {
    title: "Settings",
    inviteUser: {
      title: "Invite User",
      description: "Invite a new user to join the platform",
      emailPlaceholder: "Enter email address",
      inviteButton: "Send Invite",
      success: "Invitation sent successfully",
      error: {
        failed: "Failed to send invitation",
        invalidEmail: "Please enter a valid email address",
      },
    },
    changeHandle: {
      title: "Change Handle",
      currentHandle: "Current Handle",
      newHandleLabel: "New Handle",
      newHandlePlaceholder: "Enter desired handle",
      formatHelp: "3-32 characters, letters, numbers, and underscores only.",
      checkAvailabilityButton: "Check Availability",
      setHandleButton: "Set New Handle",
      available: "Handle is available!",
      notAvailable: "Handle is not available.",
      suggestions: "Suggestions",
      success: "Handle updated successfully!",
      upgradePrompt: "Changing your handle is available for paid users.",
      upgradeLink: "Upgrade your account.",
      error: {
        invalidFormat:
          "Invalid handle format. Use 3-32 letters, numbers, or underscores.",
        checkFailed: "Failed to check handle availability. Please try again.",
        setFailed: "Failed to set new handle. Please try again.",
        conflict: "This handle has just been taken. Please try another.",
        notAvailableOrInvalid: "Handle is invalid or not available.",
      },
    },
  },
  hubUserOnboarding: {
    title: "Welcome to Vetchi",
    subtitle: "Complete your profile to get started",
    form: {
      fullName: "Full Name",
      fullNamePlaceholder: "Enter your full name",
      password: "Password",
      passwordPlaceholder: "Choose a secure password",
      confirmPassword: "Confirm Password",
      confirmPasswordPlaceholder: "Re-enter your password",
      countryCode: "Country of Residence",
      countryCodePlaceholder: "Select your country",
      tier: {
        label: "Select Your Plan",
        free: "Free Tier",
        paid: "Paid Tier",
        freeDescription: "Apply for Jobs, Look at Posts, See Ads",
        paidDescription: "Support Open Source software, No Ads",
      },
      preferredLanguage: "Preferred Language",
      preferredLanguagePlaceholder: "Select your preferred language",
      shortBio: "Short Bio",
      shortBioPlaceholder:
        "Brief introduction about yourself (max 64 characters)",
      longBio: "Long Bio",
      longBioPlaceholder:
        "Detailed description about your background and expertise (max 1024 characters)",
      submit: "Complete Registration",
    },
    error: {
      invalidToken:
        "The invitation token is invalid or has expired. Please request a new invitation.",
      onboardingFailed: "Failed to complete registration. Please try again.",
      requiredField: "This field is required",
      passwordLength: "Password must be between 12 and 64 characters",
      passwordMismatch: "Passwords do not match",
      validationError: "Please correct the following errors: {details}",
    },
    success: {
      title: "Registration Complete!",
      description: "Your account has been created successfully.",
      handle: "Your generated handle is: {handle}",
      redirecting: "Redirecting to dashboard...",
    },
  },
  education: {
    title: "Education",
    addEducation: "Add Education",
    updateEducation: "Update Education",
    loading: "Loading education...",
    noEntries: "No education entries yet",
    instituteDomain: "Institute Domain",
    degree: "Degree",
    startDate: "Start Date",
    endDate: "End Date",
    description: "Description",
    present: "Present",
    searchInstitute: "Search for institute",
    searchMinChars: "Type at least 3 characters to search",
    deleteConfirm: "Are you sure you want to delete this education entry?",
    charactersLimit: "characters",
    error: {
      fetchFailed: "Failed to fetch education",
      saveFailed: "Failed to save education",
      deleteFailed: "Failed to delete education",
      searchFailed: "Failed to search institutes",
      invalidDomain:
        "Please enter a valid domain name (e.g., harvard.edu, stanford.example)",
      invalidDate: "Please enter a valid date in YYYY-MM-DD format",
      endDateBeforeStart: "End date must be after or equal to start date",
      futureDate: "Date cannot be in the future",
      degreeLength: "Degree must be between 3 and 64 characters",
      descriptionTooLong: "Description cannot exceed 1024 characters",
    },
    actions: {
      edit: "Edit",
      delete: "Delete",
      save: "Save",
      cancel: "Cancel",
    },
  },
  achievements: {
    patents: {
      title: "Patents",
      addPatent: "Add Patent",
      updatePatent: "Update Patent",
      loading: "Loading patents...",
      noEntries: "No patents entries yet",
      title_field: "Title",
      description: "Description",
      url: "URL",
      date: "Date",
      deleteConfirm: "Are you sure you want to delete this patent?",
      error: {
        fetchFailed: "Failed to fetch patents",
        saveFailed: "Failed to save patent",
        deleteFailed: "Failed to delete patent",
        titleLength: "Title must be between 3 and 128 characters",
        descriptionTooLong: "Description cannot exceed 1024 characters",
        urlTooLong: "URL cannot exceed 1024 characters",
        invalidUrl: "Please enter a valid URL",
        invalidDate: "Please enter a valid date",
        futureDate: "Date cannot be in the future",
      },
    },
    publications: {
      title: "Publications",
      addPublication: "Add Publication",
      updatePublication: "Update Publication",
      loading: "Loading publications...",
      noEntries: "No publications entries yet",
      title_field: "Title",
      description: "Description",
      url: "URL",
      date: "Date",
      deleteConfirm: "Are you sure you want to delete this publication?",
      error: {
        fetchFailed: "Failed to fetch publications",
        saveFailed: "Failed to save publication",
        deleteFailed: "Failed to delete publication",
        titleLength: "Title must be between 3 and 128 characters",
        descriptionTooLong: "Description cannot exceed 1024 characters",
        urlTooLong: "URL cannot exceed 1024 characters",
        invalidUrl: "Please enter a valid URL",
        invalidDate: "Please enter a valid date",
        futureDate: "Date cannot be in the future",
      },
    },
    certifications: {
      title: "Certifications",
      addCertification: "Add Certification",
      updateCertification: "Update Certification",
      loading: "Loading certifications...",
      noEntries: "No certifications entries yet",
      title_field: "Title",
      description: "Description",
      url: "URL",
      date: "Date Obtained",
      deleteConfirm: "Are you sure you want to delete this certification?",
      error: {
        fetchFailed: "Failed to fetch certifications",
        saveFailed: "Failed to save certification",
        deleteFailed: "Failed to delete certification",
        titleLength: "Title must be between 3 and 128 characters",
        descriptionTooLong: "Description cannot exceed 1024 characters",
        urlTooLong: "URL cannot exceed 1024 characters",
        invalidUrl: "Please enter a valid URL",
        invalidDate: "Please enter a valid date",
        futureDate: "Date cannot be in the future",
      },
    },
    actions: {
      edit: "Edit",
      delete: "Delete",
      save: "Save",
      cancel: "Cancel",
    },
  },
  posts: {
    title: "Posts",
    compose: "Compose",
    placeholder: "What's on your mind?",
    publish: "Publish",
    addTag: "Add Tag",
    removeTag: "Remove",
    maxTags: "Maximum of 3 tags allowed",
    newTag: "Create new tag: {name}",
    searchTags: "Search for tags",
    following: "Following",
    trending: "Trending",
    error: {
      fetchFailed: "Failed to fetch posts",
      createFailed: "Failed to create post",
      tagsFailed: "Failed to fetch tags",
      contentRequired: "Post content is required",
    },
    success: "Post published successfully",
    noTimelinePosts:
      "No posts from people you follow yet. Start following users to see their posts here!",
    trendingComingSoon: "Trending posts coming soon",
    viewPost: "Post Details",
    postId: "Post ID",
    content: "Content",
    contentPlaceholder: "This content will be loaded in a future update",
    detailsComingSoon:
      "Detailed post view with comments and interactions coming soon",
  },
};
