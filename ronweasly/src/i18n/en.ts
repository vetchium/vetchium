export const en = {
  common: {
    login: "Sign In",
    logout: "Logout",
    email: "Email Address",
    password: "Password",
    verify: "Verify",
    search: "Search",
    loading: "Loading...",
    error: {
      notAuthenticated: "Not authenticated. Please log in again.",
      sessionExpired: "Session expired. Please log in again.",
      serverError: "The server is experiencing issues. Please try again later",
    },
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
    },
  },
  navigation: {
    home: "Home",
    findOpenings: "Find Openings",
  },
};
