export interface AddWorkHistoryRequest {
  employer_domain: string;
  title: string;
  start_date: string;
  end_date?: string;
  description?: string;
}

export interface AddWorkHistoryResponse {
  work_history_id: string;
}

export interface WorkHistory {
  id: string;
  employer_domain: string;
  employer_name?: string;
  employer_logo_url?: string;
  title: string;
  start_date: string;
  end_date?: string;
  description?: string;
}

export interface UpdateWorkHistoryRequest {
  id: string;
  title: string;
  start_date: string;
  end_date?: string;
  description?: string;
}

export interface ListWorkHistoryRequest {
  user_handle?: string;
}

export interface DeleteWorkHistoryRequest {
  id: string;
}

function isValidDate(dateStr: string): boolean {
  const date = new Date(dateStr);
  return date instanceof Date && !isNaN(date.getTime());
}

export function isValidAddWorkHistoryRequest(req: AddWorkHistoryRequest): {
  valid: boolean;
  error?: string;
} {
  if (!req.employer_domain) {
    return { valid: false, error: "employer_domain is required" };
  }
  if (!req.title) {
    return { valid: false, error: "title is required" };
  }
  if (!req.start_date || !isValidDate(req.start_date)) {
    return {
      valid: false,
      error: "start_date is required and must be a valid date",
    };
  }
  if (req.end_date && !isValidDate(req.end_date)) {
    return { valid: false, error: "end_date must be a valid date if provided" };
  }
  if (req.description && req.description.length > 1024) {
    return {
      valid: false,
      error: "description must not exceed 1024 characters",
    };
  }
  return { valid: true };
}

export function isValidUpdateWorkHistoryRequest(
  req: UpdateWorkHistoryRequest
): { valid: boolean; error?: string } {
  if (!req.id) {
    return { valid: false, error: "id is required" };
  }
  if (!req.title) {
    return { valid: false, error: "title is required" };
  }
  if (!req.start_date || !isValidDate(req.start_date)) {
    return {
      valid: false,
      error: "start_date is required and must be a valid date",
    };
  }
  if (req.end_date && !isValidDate(req.end_date)) {
    return { valid: false, error: "end_date must be a valid date if provided" };
  }
  if (req.description && req.description.length > 1024) {
    return {
      valid: false,
      error: "description must not exceed 1024 characters",
    };
  }
  return { valid: true };
}

export function isValidListWorkHistoryRequest(req: ListWorkHistoryRequest): {
  valid: boolean;
  error?: string;
} {
  // No validation needed for user_handle as it's optional and any string is valid
  return { valid: true };
}
