export interface ValidationErrors {
  errors: string[];
}

export type EmailAddress = string;
export type Password = string;

/**
 * Validates an email address according to common.tsp requirements
 * @param email The email address to validate
 * @returns True if the email is valid, false otherwise
 */
export function validateEmailAddress(email: string): boolean {
  // Check for min and max length as per common.tsp
  if (email.length < 3 || email.length > 256) {
    return false;
  }

  // Basic email format validation
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
}

/**
 * Validates a password according to common.tsp requirements
 * @param password The password to validate
 * @returns An object with validation result and error message if any
 */
export function validatePassword(password: string): {
  isValid: boolean;
  error?: string;
} {
  // Check for min and max length as per common.tsp
  if (password.length < 12) {
    return {
      isValid: false,
      error: "Password must be at least 12 characters long",
    };
  }

  if (password.length > 64) {
    return {
      isValid: false,
      error: "Password must be at most 64 characters long",
    };
  }

  // Additional password strength requirements
  const hasUpperCase = /[A-Z]/.test(password);
  const hasLowerCase = /[a-z]/.test(password);
  const hasNumbers = /[0-9]/.test(password);
  const hasSpecialChar = /[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]/.test(password);

  if (!hasUpperCase) {
    return {
      isValid: false,
      error: "Password must contain at least one uppercase letter",
    };
  }

  if (!hasLowerCase) {
    return {
      isValid: false,
      error: "Password must contain at least one lowercase letter",
    };
  }

  if (!hasNumbers) {
    return {
      isValid: false,
      error: "Password must contain at least one number",
    };
  }

  if (!hasSpecialChar) {
    return {
      isValid: false,
      error: "Password must contain at least one special character",
    };
  }

  return { isValid: true };
}
export type City = string;
export type Handle = string;

export type CountryCode = string;

export const GlobalCountryCode = "ZZG";

export type Currency = string;

export type TimeZone = string;

export const validTimezones = new Set([
  "ACDT Australian Central Daylight Time GMT+1030",
  "ACST Australian Central Standard Time GMT+0930",
  "AEDT Australian Eastern Daylight Time GMT+1100",
  "AEST Australian Eastern Standard Time GMT+1000",
  "AFT Afghanistan Time GMT+0430",
  "AKDT Alaska Daylight Time GMT-0800",
  "AKST Alaska Standard Time GMT-0900",
  "ALMT Alma-Ata Time GMT+0600",
  "AMST Amazon Summer Time (Brazil) GMT-0300",
  "AMT Amazon Time (Brazil) GMT-0400",
  "ANAST Anadyr Summer Time GMT+1200",
  "ANAT Anadyr Time GMT+1200",
  "AQTT Aqtobe Time GMT+0500",
  "ART Argentina Time GMT-0300",
  "AST Arabia Standard Time GMT+0300",
  "AST Atlantic Standard Time GMT-0400",
  "AWST Australian Western Standard Time GMT+0800",
  "AZOST Azores Summer Time GMT+0000",
  "AZOT Azores Standard Time GMT-0100",
  "AZT Azerbaijan Time GMT+0400",
  "BNT Brunei Darussalam Time GMT+0800",
  "BOT Bolivia Time GMT-0400",
  "BRST Brasilia Summer Time GMT-0200",
  "BRT Brasilia Time GMT-0300",
  "BST Bangladesh Standard Time GMT+0600",
  "BST Bougainville Standard Time GMT+1100",
  "BST British Summer Time GMT+0100",
  "BTT Bhutan Time GMT+0600",
  "CAT Central Africa Time GMT+0200",
  "CCT Cocos Islands Time GMT+0630",
  "CDT Central Daylight Time (North America) GMT-0500",
  "CEST Central European Summer Time GMT+0200",
  "CET Central European Time GMT+0100",
  "CHADT Chatham Island Daylight Time GMT+1345",
  "CHAST Chatham Island Standard Time GMT+1245",
  "CKT Cook Island Time GMT-1000",
  "CLST Chile Summer Time GMT-0300",
  "CLT Chile Standard Time GMT-0400",
  "COT Colombia Time GMT-0500",
  "CST Central Standard Time (North America) GMT-0600",
  "CST China Standard Time GMT+0800",
  "CST Cuba Standard Time GMT-0500",
  "CVT Cape Verde Time GMT-0100",
  "CXT Christmas Island Time GMT+0700",
  "DAVT Davis Time GMT+0700",
  "EASST Easter Island Summer Time GMT-0500",
  "EAST Easter Island Standard Time GMT-0600",
  "EAT East Africa Time GMT+0300",
  "ECT Ecuador Time GMT-0500",
  "EDT Eastern Daylight Time (North America) GMT-0400",
  "EEST Eastern European Summer Time GMT+0300",
  "EET Eastern European Time GMT+0200",
  "EGST Eastern Greenland Summer Time GMT+0000",
  "EGT Eastern Greenland Time GMT-0100",
  "FET Further-eastern European Time GMT+0300",
  "FJT Fiji Time GMT+1200",
  "FKST Falkland Islands Summer Time GMT-0300",
  "FKT Falkland Islands Time GMT-0400",
  "FNT Fernando de Noronha Time GMT-0200",
  "GALT Galapagos Time GMT-0600",
  "GAMT Gambier Time GMT-0900",
  "GET Georgia Standard Time GMT+0400",
  "GFT French Guiana Time GMT-0300",
  "GILT Gilbert Island Time GMT+1200",
  "GMT Greenwich Mean Time GMT+0000",
  "GST South Georgia Time GMT-0200",
  "GST Gulf Standard Time GMT+0400",
  "GYT Guyana Time GMT-0400",
  "HKT Hong Kong Time GMT+0800",
  "HOVT Hovd Time GMT+0700",
  "HST Hawaii-Aleutian Standard Time GMT-1000",
  "ICT Indochina Time GMT+0700",
  "IDT Israel Daylight Time GMT+0300",
  "IOT Indian Chagos Time GMT+0600",
  "IRDT Iran Daylight Time GMT+0430",
  "IRKT Irkutsk Time GMT+0800",
  "IRST Iran Standard Time GMT+0330",
  "IST Indian Standard Time GMT+0530",
  "IST Irish Standard Time GMT+0100",
  "JST Japan Standard Time GMT+0900",
  "KGT Kyrgyzstan Time GMT+0600",
  "KOST Kosrae Time GMT+1100",
  "KRAT Krasnoyarsk Time GMT+0700",
  "KST Korea Standard Time GMT+0900",
  "LHDT Lord Howe Daylight Time GMT+1100",
  "LHST Lord Howe Standard Time GMT+1030",
  "LINT Line Islands Time GMT+1400",
  "MAGT Magadan Time GMT+1100",
  "MART Marquesas Time GMT-0930",
  "MAWT Mawson Time GMT+0500",
  "MDT Mountain Daylight Time (North America) GMT-0600",
  "MET Middle European Time GMT+0100",
  "MEST Middle European Summer Time GMT+0200",
  "MHT Marshall Islands Time GMT+1200",
  "MIST Macquarie Island Station Time GMT+1100",
  "MIT Marquesas Islands Time GMT-0930",
  "MMT Myanmar Time GMT+0630",
  "MSK Moscow Time GMT+0300",
  "MST Malaysia Standard Time GMT+0800",
  "MST Mountain Standard Time (North America) GMT-0700",
  "MUT Mauritius Time GMT+0400",
  "MVT Maldives Time GMT+0500",
  "MYT Malaysia Time GMT+0800",
  "NCT New Caledonia Time GMT+1100",
  "NDT Newfoundland Daylight Time GMT-0230",
  "NFT Norfolk Island Time GMT+1130",
  "NOVT Novosibirsk Time GMT+0700",
  "NPT Nepal Time GMT+0545",
  "NST Newfoundland Standard Time GMT-0330",
  "NT Newfoundland Time GMT-0330",
  "NUT Niue Time GMT-1100",
  "NZDT New Zealand Daylight Time GMT+1300",
  "NZST New Zealand Standard Time GMT+1200",
  "OMST Omsk Time GMT+0600",
  "ORAT Oral Time GMT+0500",
  "PDT Pacific Daylight Time (North America) GMT-0700",
  "PET Peru Time GMT-0500",
  "PETT Kamchatka Time GMT+1200",
  "PGT Papua New Guinea Time GMT+1000",
  "PHOT Phoenix Island Time GMT+1300",
  "PKT Pakistan Standard Time GMT+0500",
  "PMDT Saint Pierre and Miquelon Daylight Time GMT-0200",
  "PMST Saint Pierre and Miquelon Standard Time GMT-0300",
  "PONT Pohnpei Standard Time GMT+1100",
  "PST Pacific Standard Time (North America) GMT-0800",
  "PYST Paraguay Summer Time GMT-0300",
  "PYT Paraguay Time GMT-0400",
  "RET Réunion Time GMT+0400",
  "ROTT Rothera Research Station Time GMT-0300",
  "SAKT Sakhalin Island Time GMT+1100",
  "SAMT Samara Time GMT+0400",
  "SAST South Africa Standard Time GMT+0200",
  "SBT Solomon Islands Time GMT+1100",
  "SCT Seychelles Time GMT+0400",
  "SGT Singapore Time GMT+0800",
  "SLST Sri Lanka Standard Time GMT+0530",
  "SRET Srednekolymsk Time GMT+1100",
  "SRT Suriname Time GMT-0300",
  "SST Samoa Standard Time GMT-1100",
  "SYOT Syowa Time GMT+0300",
  "TAHT Tahiti Time GMT-1000",
  "THA Thailand Standard Time GMT+0700",
  "TFT French Southern and Antarctic Time GMT+0500",
  "TJT Tajikistan Time GMT+0500",
  "TKT Tokelau Time GMT+1300",
  "TLT Timor Leste Time GMT+0900",
  "TMT Turkmenistan Time GMT+0500",
  "TRT Turkey Time GMT+0300",
  "TOT Tonga Time GMT+1300",
  "TVT Tuvalu Time GMT+1200",
  "ULAST Ulaanbaatar Summer Time GMT+0900",
  "ULAT Ulaanbaatar Time GMT+0800",
  "UTC Coordinated Universal Time GMT+0000",
  "UYST Uruguay Summer Time GMT-0200",
  "UYT Uruguay Time GMT-0300",
  "VET Venezuelan Standard Time GMT-0400",
  "VLAST Vladivostok Summer Time GMT+1100",
  "VLAT Vladivostok Time GMT+1000",
  "VOST Vostok Station Time GMT+0600",
  "VUT Vanuatu Time GMT+1100",
  "WAKT Wake Island Time GMT+1200",
  "WAT West Africa Time GMT+0100",
  "WEDT Western European Daylight Time GMT+0100",
  "WEST Western European Summer Time GMT+0100",
  "WET Western European Time GMT+0000",
  "WGST Western Greenland Summer Time GMT-0200",
  "WGT Western Greenland Time GMT-0300",
  "WIB Western Indonesia Time GMT+0700",
  "WIT Eastern Indonesia Time GMT+0900",
  "WITA Central Indonesia Time GMT+0800",
  "WST Western Standard Time (Australia) GMT+0800",
  "WT Western Sahara Standard Time GMT+0000",
  "YAKT Yakutsk Time GMT+0900",
  "YEKT Yekaterinburg Time GMT+0500",
]);

export function isValidTimeZone(timezone: string): timezone is TimeZone {
  return validTimezones.has(timezone);
}

export type OrgUserRole = string;

export const OrgUserRoles = {
  ADMIN: "ADMIN" as OrgUserRole,
  ANY: "ANY" as OrgUserRole,
  APPLICATIONS_CRUD: "APPLICATIONS_CRUD" as OrgUserRole,
  APPLICATIONS_VIEWER: "APPLICATIONS_VIEWER" as OrgUserRole,
  COST_CENTERS_CRUD: "COST_CENTERS_CRUD" as OrgUserRole,
  COST_CENTERS_VIEWER: "COST_CENTERS_VIEWER" as OrgUserRole,
  LOCATIONS_CRUD: "LOCATIONS_CRUD" as OrgUserRole,
  LOCATIONS_VIEWER: "LOCATIONS_VIEWER" as OrgUserRole,
  OPENINGS_CRUD: "OPENINGS_CRUD" as OrgUserRole,
  OPENINGS_VIEWER: "OPENINGS_VIEWER" as OrgUserRole,
  ORG_USERS_CRUD: "ORG_USERS_CRUD" as OrgUserRole,
  ORG_USERS_VIEWER: "ORG_USERS_VIEWER" as OrgUserRole,
} as const;

export function isValidOrgUserRole(role: string): role is OrgUserRole {
  return Object.values(OrgUserRoles).includes(role as OrgUserRole);
}

export interface HubAuth {
  type: "http";
  scheme: "bearer";
}

export interface EmployerAuth {
  type: "http";
  scheme: "bearer";
}

export type TimelineItemType = "USER_POST" | "EMPLOYER_POST";

export const TimelineItemUserPost: TimelineItemType = "USER_POST";
export const TimelineItemEmployerPost: TimelineItemType = "EMPLOYER_POST";
