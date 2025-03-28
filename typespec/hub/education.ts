import { Handle } from "../common/common";
import { Institute } from "../common/education";

export interface AddEducationRequest {
  institute_domain: string;
  degree?: string;
  start_date?: string;
  end_date?: string;
  description?: string;
}

export interface AddEducationResponse {
  education_id: string;
}

export interface FilterInstitutesRequest {
  prefix: string;
}

export interface DeleteEducationRequest {
  education_id: string;
}

export interface ListEducationRequest {
  user_handle?: Handle;
}

export interface FilterInstitutesResponse {
  // Maximum 10 institutes will be returned in random order.
  institutes: Institute[];
}

// Validation functions
export function isValidDomain(domain: string): boolean {
  // Match the regex from validations.go
  const domainRegex = /^([a-zA-Z0-9-]+\.)+[a-zA-Z0-9-]{2,}$/;
  return domainRegex.test(domain);
}

export function isValidDate(dateStr?: string): boolean {
  if (!dateStr) return true;

  // Check if date is in YYYY-MM-DD format
  const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
  if (!dateRegex.test(dateStr)) return false;

  // Check if date is valid
  const date = new Date(dateStr);
  return !isNaN(date.getTime());
}

export function isEndDateAfterStartDate(
  startDate?: string,
  endDate?: string
): boolean {
  if (!startDate || !endDate) return true;

  // Validate both dates first
  if (!isValidDate(startDate) || !isValidDate(endDate)) return false;

  // Parse dates
  const start = new Date(startDate);
  const end = new Date(endDate);

  return end >= start;
}

export function isNotFutureDate(dateStr?: string): boolean {
  if (!dateStr) return true;

  // Validate date first
  if (!isValidDate(dateStr)) return false;

  // Get today's date (strip time)
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const date = new Date(dateStr);
  return date <= today;
}

export function validateEducationForm(
  education: AddEducationRequest
): string[] {
  const errors: string[] = [];

  // Validate domain
  if (!education.institute_domain) {
    errors.push("Institute domain is required");
  } else if (!isValidDomain(education.institute_domain)) {
    errors.push("Invalid institute domain format");
  }

  // Validate degree
  if (
    education.degree &&
    (education.degree.length < 3 || education.degree.length > 64)
  ) {
    errors.push("Degree must be between 3 and 64 characters");
  }

  // Validate dates
  if (education.start_date && !isValidDate(education.start_date)) {
    errors.push("Invalid start date format");
  }

  if (education.end_date && !isValidDate(education.end_date)) {
    errors.push("Invalid end date format");
  }

  if (
    education.start_date &&
    education.end_date &&
    !isEndDateAfterStartDate(education.start_date, education.end_date)
  ) {
    errors.push("End date must be after start date");
  }

  // Validate description length
  if (education.description && education.description.length > 1024) {
    errors.push("Description cannot exceed 1024 characters");
  }

  return errors;
}
