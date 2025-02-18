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
    error: {
      notAuthenticated: "Not authenticated. Please log in again.",
      sessionExpired: "Session expired. Please log in again.",
      serverError: "The server is experiencing issues. Please try again later",
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
  },
  navigation: {
    home: "Home",
    findOpenings: "Find Openings",
    myApplications: "My Applications",
    myCandidacies: "My Candidacies",
    myProfile: "My Profile",
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
    bio: {
      error: {
        fetchFailed: "Failed to load profile information",
        updateFailed: "Failed to update profile information",
        uploadFailed: "Failed to upload profile picture",
      },
      title: "Edit Bio",
      fullName: "Full Name",
      handle: "Handle",
      shortBio: "Short Bio",
      longBio: "Long Bio",
      save: "Save Changes",
      cancel: "Cancel",
    },
    picture: {
      change: "Change Profile Picture",
      upload: "Upload Profile Picture",
    },
  },
};
