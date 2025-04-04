import { config } from "@/config";
import { useTranslation } from "@/hooks/useTranslation";
import DeleteIcon from "@mui/icons-material/Delete";
import Alert from "@mui/material/Alert";
import Autocomplete from "@mui/material/Autocomplete";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import CircularProgress from "@mui/material/CircularProgress";
import IconButton from "@mui/material/IconButton";
import Paper from "@mui/material/Paper";
import Stack from "@mui/material/Stack";
import TextField from "@mui/material/TextField";
import Typography from "@mui/material/Typography";
import {
  AddEducationRequest,
  Education as EducationType,
  Institute,
} from "@vetchium/typespec";
import Cookies from "js-cookie";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

interface EducationProps {
  userHandle: string;
  canEdit: boolean;
}

type EducationFormData = AddEducationRequest;

function formatDate(dateString: string): string {
  return new Intl.DateTimeFormat(undefined, {
    year: "numeric",
    month: "long",
    day: "numeric",
  }).format(new Date(dateString));
}

function formatDateForInput(dateString: string): string {
  const date = new Date(dateString);
  return date.toISOString().split("T")[0];
}

// Validation functions
function isValidDomain(domain: string): boolean {
  // Match the regex from validations.go
  const domainRegex = /^([a-zA-Z0-9-]+\.)+[a-zA-Z0-9-]{2,}$/;
  return domainRegex.test(domain);
}

function isValidDate(dateStr?: string): boolean {
  if (!dateStr) return true;

  // Check if date is in YYYY-MM-DD format
  const dateRegex = /^\d{4}-\d{2}-\d{2}$/;
  if (!dateRegex.test(dateStr)) return false;

  // Check if date is valid
  const date = new Date(dateStr);
  return !isNaN(date.getTime());
}

function isEndDateAfterStartDate(
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

function isNotFutureDate(dateStr?: string): boolean {
  if (!dateStr) return true;

  // Validate date first
  if (!isValidDate(dateStr)) return false;

  // Get today's date (strip time)
  const today = new Date();
  today.setHours(0, 0, 0, 0);

  const date = new Date(dateStr);
  return date <= today;
}

export function Education({ userHandle, canEdit }: EducationProps) {
  const router = useRouter();
  const { t } = useTranslation();
  const [education, setEducation] = useState<EducationType[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isAddingNew, setIsAddingNew] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [isDeleting, setIsDeleting] = useState<string | null>(null);
  const [formData, setFormData] = useState<EducationFormData>({
    institute_domain: "",
  });

  // Validation errors
  const [domainError, setDomainError] = useState<string | null>(null);
  const [startDateError, setStartDateError] = useState<string | null>(null);
  const [endDateError, setEndDateError] = useState<string | null>(null);
  const [degreeError, setDegreeError] = useState<string | null>(null);
  const [descriptionError, setDescriptionError] = useState<string | null>(null);

  // For institute search
  const [searchQuery, setSearchQuery] = useState("");
  const [institutes, setInstitutes] = useState<Institute[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [searchError, setSearchError] = useState<string | null>(null);

  useEffect(() => {
    fetchEducation();
  }, [userHandle]);

  // Search institutes when query is at least 3 characters
  useEffect(() => {
    if (searchQuery.length >= 3) {
      searchInstitutes(searchQuery);
    } else {
      setInstitutes([]);
    }
  }, [searchQuery]);

  // Validate form fields when they change
  useEffect(() => {
    validateForm();
  }, [formData]);

  // Validate form fields and set appropriate error states
  const validateForm = () => {
    // Reset all errors
    setDomainError(null);
    setStartDateError(null);
    setEndDateError(null);
    setDegreeError(null);
    setDescriptionError(null);

    // Validate domain
    if (!formData.institute_domain) {
      setDomainError(t("common.error.requiredField"));
    } else if (
      formData.institute_domain &&
      !isValidDomain(formData.institute_domain)
    ) {
      setDomainError(t("education.error.invalidDomain"));
    }

    // Validate degree (now mandatory)
    if (!formData.degree) {
      setDegreeError(t("common.error.requiredField"));
    } else if (
      formData.degree &&
      (formData.degree.length < 3 || formData.degree.length > 64)
    ) {
      setDegreeError(t("education.error.degreeLength"));
    }

    // Validate dates
    if (formData.start_date && !isValidDate(formData.start_date)) {
      setStartDateError(t("education.error.invalidDate"));
    } else if (formData.start_date && !isNotFutureDate(formData.start_date)) {
      setStartDateError(t("education.error.futureDate"));
    }

    if (formData.end_date && !isValidDate(formData.end_date)) {
      setEndDateError(t("education.error.invalidDate"));
    } else if (formData.end_date && !isNotFutureDate(formData.end_date)) {
      setEndDateError(t("education.error.futureDate"));
    }

    // Validate end date is after start date
    if (
      formData.start_date &&
      formData.end_date &&
      isValidDate(formData.start_date) &&
      isValidDate(formData.end_date) &&
      !isEndDateAfterStartDate(formData.start_date, formData.end_date)
    ) {
      setEndDateError(t("education.error.endDateBeforeStart"));
    }

    // Validate description length
    if (formData.description && formData.description.length > 1024) {
      setDescriptionError(t("education.error.descriptionTooLong"));
    }
  };

  // Check if the form has any validation errors
  const hasValidationErrors = (): boolean => {
    return !!(
      domainError ||
      startDateError ||
      endDateError ||
      degreeError ||
      descriptionError
    );
  };

  async function fetchEducation() {
    try {
      setIsLoading(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const request = { user_handle: userHandle };
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/list-education`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(request),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("education.error.fetchFailed"));
      }

      const data = await response.json();
      setEducation(data || []);
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
      setEducation([]);
    } finally {
      setIsLoading(false);
    }
  }

  async function searchInstitutes(prefix: string) {
    if (prefix.length < 3) return;

    try {
      setIsSearching(true);
      setSearchError(null);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/filter-institutes`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ prefix }),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("education.error.searchFailed"));
      }

      const data = await response.json();
      setInstitutes(data || []);
    } catch (err) {
      setSearchError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
      setInstitutes([]);
    } finally {
      setIsSearching(false);
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();

    // Validate form before submission
    validateForm();
    if (hasValidationErrors()) {
      return;
    }

    // Check mandatory fields
    if (!formData.institute_domain || !formData.degree) {
      return;
    }

    try {
      setIsSaving(true);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      // For add education, we use /hub/add-education endpoint
      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/add-education`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify(formData),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("education.error.saveFailed"));
      }

      await fetchEducation();
      resetForm();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
    } finally {
      setIsSaving(false);
    }
  }

  async function handleDelete(id: string) {
    if (!confirm(t("education.deleteConfirm"))) {
      return;
    }

    try {
      setIsDeleting(id);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/delete-education`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({ education_id: id }),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) {
        throw new Error(t("education.error.deleteFailed"));
      }

      await fetchEducation();
      setError(null);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
    } finally {
      setIsDeleting(null);
    }
  }

  function resetForm() {
    setFormData({
      institute_domain: "",
    });
    setSearchQuery("");
    setIsAddingNew(false);
    setDomainError(null);
    setStartDateError(null);
    setEndDateError(null);
    setDegreeError(null);
    setDescriptionError(null);
  }

  if (isLoading) {
    return (
      <Box sx={{ display: "flex", justifyContent: "center", p: 4 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      {error && (
        <Alert severity="error" sx={{ mb: 4 }}>
          {error}
        </Alert>
      )}

      <Box
        sx={{
          display: "flex",
          justifyContent: "space-between",
          alignItems: "center",
          mb: 3,
        }}
      >
        <Typography variant="h5">{t("education.title")}</Typography>
        {canEdit && !isAddingNew && (
          <Button
            variant="contained"
            color="primary"
            onClick={() => setIsAddingNew(true)}
          >
            {t("education.addEducation")}
          </Button>
        )}
      </Box>

      {isAddingNew && canEdit && (
        <Paper sx={{ p: 3, mb: 4 }}>
          {searchError && (
            <Alert severity="error" sx={{ mb: 3 }}>
              {searchError}
            </Alert>
          )}
          <form onSubmit={handleSubmit}>
            <Stack spacing={3}>
              <Autocomplete
                freeSolo
                options={institutes}
                getOptionLabel={(option) =>
                  typeof option === "string"
                    ? option
                    : `${option.name} (${option.domain})`
                }
                inputValue={searchQuery}
                onInputChange={(_, newValue) => {
                  setSearchQuery(newValue);
                  if (newValue.length >= 3) {
                    setFormData({
                      ...formData,
                      institute_domain: newValue,
                    });
                  }
                }}
                onChange={(_, newValue) => {
                  if (newValue && typeof newValue !== "string") {
                    setFormData({
                      ...formData,
                      institute_domain: newValue.domain,
                    });
                  } else if (typeof newValue === "string") {
                    setFormData({
                      ...formData,
                      institute_domain: newValue,
                    });
                  }
                }}
                loading={isSearching}
                renderInput={(params) => (
                  <TextField
                    {...params}
                    label={t("education.instituteDomain")}
                    required
                    fullWidth
                    error={!!domainError}
                    helperText={
                      domainError
                        ? domainError
                        : searchQuery.length > 0 && searchQuery.length < 3
                        ? t("education.searchMinChars")
                        : ""
                    }
                    InputProps={{
                      ...params.InputProps,
                      endAdornment: (
                        <>
                          {isSearching ? (
                            <CircularProgress color="inherit" size={20} />
                          ) : null}
                          {params.InputProps.endAdornment}
                        </>
                      ),
                    }}
                  />
                )}
              />
              <TextField
                label={t("education.degree")}
                value={formData.degree || ""}
                onChange={(e) =>
                  setFormData({ ...formData, degree: e.target.value })
                }
                fullWidth
                required
                error={!!degreeError}
                helperText={
                  degreeError ||
                  `${formData.degree ? formData.degree.length : 0}/64 ${t(
                    "education.charactersLimit"
                  )}`
                }
                inputProps={{ maxLength: 64 }}
              />
              <TextField
                label={t("education.startDate")}
                type="date"
                value={formData.start_date ? formData.start_date : ""}
                onChange={(e) =>
                  setFormData({ ...formData, start_date: e.target.value })
                }
                fullWidth
                InputLabelProps={{ shrink: true }}
                error={!!startDateError}
                helperText={startDateError || ""}
              />
              <TextField
                label={t("education.endDate")}
                type="date"
                value={formData.end_date ? formData.end_date : ""}
                onChange={(e) =>
                  setFormData({ ...formData, end_date: e.target.value })
                }
                fullWidth
                InputLabelProps={{ shrink: true }}
                error={!!endDateError}
                helperText={endDateError || ""}
              />
              <TextField
                label={t("education.description")}
                value={formData.description || ""}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
                multiline
                rows={4}
                fullWidth
                error={!!descriptionError}
                helperText={
                  descriptionError ||
                  `${
                    formData.description ? formData.description.length : 0
                  }/1024 ${t("education.charactersLimit")}`
                }
                inputProps={{ maxLength: 1024 }}
              />
              <Box sx={{ display: "flex", gap: 2 }}>
                <Button
                  type="submit"
                  variant="contained"
                  color="primary"
                  disabled={
                    isSaving ||
                    !formData.institute_domain ||
                    !formData.degree ||
                    hasValidationErrors()
                  }
                >
                  {isSaving ? (
                    <CircularProgress size={24} color="inherit" />
                  ) : (
                    t("education.actions.save")
                  )}
                </Button>
                <Button variant="outlined" color="inherit" onClick={resetForm}>
                  {t("education.actions.cancel")}
                </Button>
              </Box>
            </Stack>
          </form>
        </Paper>
      )}

      {!isAddingNew && (!education || education.length === 0) ? (
        <Paper
          elevation={1}
          sx={{
            p: { xs: 3, sm: 4 },
            textAlign: "center",
            bgcolor: "background.paper",
            borderRadius: 2,
          }}
        >
          <Typography color="text.secondary">
            {t("education.noEntries")}
          </Typography>
        </Paper>
      ) : (
        !isAddingNew && (
          <Stack spacing={2}>
            {education.map((entry) => (
              <Paper
                key={entry.id}
                elevation={1}
                sx={{
                  p: { xs: 3, sm: 4 },
                  bgcolor: (theme) =>
                    theme.palette.mode === "light" ? "grey.50" : "grey.900",
                  borderRadius: 2,
                  transition: "all 0.2s ease-in-out",
                  border: "1px solid",
                  borderColor: "divider",
                  "&:hover": {
                    boxShadow: (theme) => theme.shadows[2],
                    transform: canEdit ? "translateY(-2px)" : "none",
                    bgcolor: (theme) =>
                      theme.palette.mode === "light" ? "#ffffff" : "grey.800",
                  },
                }}
              >
                <Box
                  sx={{
                    display: "flex",
                    justifyContent: "space-between",
                    gap: 2,
                  }}
                >
                  <Box sx={{ flex: 1, minWidth: 0 }}>
                    <Typography
                      variant="h6"
                      gutterBottom
                      sx={{
                        color: "primary.main",
                        fontWeight: 600,
                      }}
                    >
                      {entry.degree || t("education.degree")}
                    </Typography>
                    <Typography
                      variant="subtitle1"
                      sx={{
                        color: "text.primary",
                        mb: 1,
                      }}
                    >
                      {entry.institute_domain}
                    </Typography>
                    <Typography
                      variant="body2"
                      sx={{
                        color: "text.secondary",
                        mb: entry.description ? 2 : 0,
                      }}
                    >
                      {entry.start_date ? formatDate(entry.start_date) : ""}
                      {entry.start_date && entry.end_date ? " - " : ""}
                      {entry.end_date
                        ? formatDate(entry.end_date)
                        : entry.start_date
                        ? t("education.present")
                        : ""}
                    </Typography>
                    {entry.description && (
                      <Typography
                        variant="body2"
                        sx={{
                          color: "text.primary",
                          whiteSpace: "pre-wrap",
                          lineHeight: 1.6,
                        }}
                      >
                        {entry.description}
                      </Typography>
                    )}
                  </Box>
                  {canEdit && (
                    <Box sx={{ display: "flex", gap: 1 }}>
                      <IconButton
                        onClick={() => handleDelete(entry.id!)}
                        color="error"
                        size="small"
                        sx={{
                          "&:hover": {
                            bgcolor: "error.lighter",
                          },
                        }}
                      >
                        <DeleteIcon />
                      </IconButton>
                    </Box>
                  )}
                </Box>
              </Paper>
            ))}
          </Stack>
        )
      )}
    </Box>
  );
}
