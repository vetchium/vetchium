"use client";

import { useParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
  TextField,
  Alert,
  CircularProgress,
  Container,
  Paper,
  Grid,
  Card,
  CardContent,
  Chip,
  Avatar,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  ButtonGroup,
} from "@mui/material";
import { useTranslation } from "@/hooks/useTranslation";
import { config } from "@/config";
import Cookies from "js-cookie";
import {
  EmployerInterview,
  GetInterviewDetailsRequest,
  InterviewersDecision,
  InterviewersDecisions,
  PutAssessmentRequest,
  RSVPInterviewRequest,
  RSVPStatus,
  RSVPStatuses,
} from "@psankar/vetchi-typespec";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import {
  Person as PersonIcon,
  ExpandMore as ExpandMoreIcon,
  Public as PublicIcon,
} from "@mui/icons-material";
import { useAuth } from "@/contexts/AuthContext";

// Helper function for consistent date formatting
const formatDateTime = (
  isoString: string | Date,
  options?: Intl.DateTimeFormatOptions,
  timeZone?: string
) => {
  const defaultOptions: Intl.DateTimeFormatOptions = {
    dateStyle: "full",
    timeStyle: "short",
    timeZone: timeZone || Intl.DateTimeFormat().resolvedOptions().timeZone,
  };

  // If it's already a Date object, convert it to ISO string
  const dateStr =
    isoString instanceof Date ? isoString.toISOString() : isoString;

  return new Intl.DateTimeFormat(undefined, options || defaultOptions).format(
    new Date(dateStr)
  );
};

// Helper to format UTC time
const formatUTCDateTime = (isoString: string | Date) => {
  // If it's already a Date object, convert it to ISO string
  const dateStr =
    isoString instanceof Date ? isoString.toISOString() : isoString;

  return new Intl.DateTimeFormat(undefined, {
    weekday: "short",
    year: "numeric",
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    timeZone: "UTC",
  }).format(new Date(dateStr));
};

export default function InterviewDetailPage() {
  const params = useParams();
  const interviewId = params.id as string;
  const { t } = useTranslation();
  const router = useRouter();
  const { userEmail } = useAuth();

  const [interview, setInterview] = useState<EmployerInterview | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [rsvpLoading, setRsvpLoading] = useState(false);

  useEffect(() => {
    fetchInterview();
  }, [interviewId]);

  const fetchInterview = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-interview-details`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            interview_id: interviewId,
          } satisfies GetInterviewDetailsRequest),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("interviews.assessment.fetchError"));
      }

      const data = await response.json();
      setInterview(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : t("common.error"));
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    if (!interview) return;

    try {
      setSaving(true);
      setError(null);
      setSuccessMessage(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: PutAssessmentRequest = {
        interview_id: interviewId,
        decision: interview.interviewers_decision,
        positives: interview.positives,
        negatives: interview.negatives,
        overall_assessment: interview.overall_assessment,
        feedback_to_candidate: interview.feedback_to_candidate,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/put-assessment`,
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
        router.push("/signin");
        return;
      }

      if (response.status === 403) {
        throw new Error(t("interviews.assessment.forbiddenError"));
      }

      if (response.status === 404) {
        throw new Error(t("interviews.assessment.notFoundError"));
      }

      if (response.status === 422) {
        throw new Error(t("interviews.assessment.validationError"));
      }

      if (!response.ok) {
        throw new Error(t("interviews.assessment.saveError"));
      }

      setSuccessMessage(t("interviews.assessment.saveSuccess"));
      await fetchInterview();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("common.error"));
    } finally {
      setSaving(false);
    }
  };

  const handleRSVP = async (status: RSVPStatus) => {
    try {
      setRsvpLoading(true);
      setError(null);
      setSuccessMessage(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const request: RSVPInterviewRequest = {
        interview_id: interviewId,
        rsvp_status: status,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/rsvp-interview`,
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
        router.push("/signin");
        return;
      }

      if (response.status === 404) {
        throw new Error(t("interviews.assessment.notFoundError"));
      }

      if (response.status === 422) {
        throw new Error(t("interviews.assessment.invalidStateError"));
      }

      if (!response.ok) {
        throw new Error(t("interviews.assessment.rsvpError"));
      }

      setSuccessMessage(t("interviews.assessment.rsvpSuccess"));
      await fetchInterview();
    } catch (err) {
      setError(err instanceof Error ? err.message : t("common.error"));
    } finally {
      setRsvpLoading(false);
    }
  };

  if (loading) {
    return (
      <AuthenticatedLayout>
        <Container maxWidth="lg">
          <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
            <CircularProgress />
          </Box>
        </Container>
      </AuthenticatedLayout>
    );
  }

  return (
    <AuthenticatedLayout>
      <Container maxWidth="lg">
        <Box sx={{ mb: 4 }}>
          <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
            <Typography variant="h4">
              {t("interviews.manageInterview")}
            </Typography>
            <Button variant="outlined" onClick={() => router.back()}>
              {t("common.back")}
            </Button>
          </Box>

          {error && (
            <Alert
              severity="error"
              sx={{ mb: 2 }}
              onClose={() => setError(null)}
            >
              {error}
            </Alert>
          )}

          {successMessage && (
            <Alert
              severity="success"
              sx={{ mb: 2 }}
              onClose={() => setSuccessMessage(null)}
            >
              {successMessage}
            </Alert>
          )}

          <Grid container spacing={3}>
            <Grid item xs={12}>
              <Accordion defaultExpanded>
                <AccordionSummary
                  expandIcon={<ExpandMoreIcon />}
                  sx={{
                    "& .MuiAccordionSummary-content": {
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "center",
                    },
                  }}
                >
                  <Typography variant="h6">
                    {t("interviews.details")}
                  </Typography>
                  <Box sx={{ display: "flex", gap: 1 }}>
                    <Chip
                      label={t(`interviews.types.${interview?.interview_type}`)}
                      size="small"
                      color="primary"
                      variant="outlined"
                    />
                    <Chip
                      label={t(
                        `interviews.states.${interview?.interview_state}`
                      )}
                      size="small"
                      color={
                        interview?.interview_state === "SCHEDULED_INTERVIEW"
                          ? "primary"
                          : interview?.interview_state === "COMPLETED_INTERVIEW"
                          ? "success"
                          : "error"
                      }
                    />
                  </Box>
                </AccordionSummary>
                <AccordionDetails>
                  <Grid container spacing={2}>
                    <Grid item xs={12}>
                      <Typography variant="subtitle2">
                        {t("interviews.type")}
                      </Typography>
                      <Chip
                        label={t(
                          `interviews.types.${interview?.interview_type}`
                        )}
                        size="small"
                        color="primary"
                        variant="outlined"
                      />
                    </Grid>
                    <Grid item xs={12}>
                      <Typography variant="subtitle2">
                        {t("interviews.startTime")}
                      </Typography>
                      <Typography>
                        {formatDateTime(interview?.start_time || "", {
                          weekday: "short",
                          year: "numeric",
                          month: "short",
                          day: "numeric",
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                        <Typography
                          component="span"
                          variant="caption"
                          color="text.secondary"
                          sx={{ ml: 1 }}
                        >
                          ({Intl.DateTimeFormat().resolvedOptions().timeZone})
                        </Typography>
                        <Typography
                          component="span"
                          variant="caption"
                          color="text.secondary"
                          sx={{ ml: 1 }}
                        >
                          (UTC: {formatUTCDateTime(interview?.start_time || "")}
                          )
                        </Typography>
                      </Typography>
                    </Grid>
                    <Grid item xs={12}>
                      <Typography variant="subtitle2">
                        {t("interviews.endTime")}
                      </Typography>
                      <Typography>
                        {formatDateTime(interview?.end_time || "", {
                          weekday: "short",
                          year: "numeric",
                          month: "short",
                          day: "numeric",
                          hour: "2-digit",
                          minute: "2-digit",
                        })}
                        <Typography
                          component="span"
                          variant="caption"
                          color="text.secondary"
                          sx={{ ml: 1 }}
                        >
                          ({Intl.DateTimeFormat().resolvedOptions().timeZone})
                        </Typography>
                        <Typography
                          component="span"
                          variant="caption"
                          color="text.secondary"
                          sx={{ ml: 1 }}
                        >
                          (UTC: {formatUTCDateTime(interview?.end_time || "")})
                        </Typography>
                      </Typography>
                    </Grid>
                    <Grid item xs={12}>
                      <Typography variant="subtitle2">
                        {t("interviews.state")}
                      </Typography>
                      <Chip
                        label={t(
                          `interviews.states.${interview?.interview_state}`
                        )}
                        size="small"
                        color={
                          interview?.interview_state === "SCHEDULED_INTERVIEW"
                            ? "primary"
                            : interview?.interview_state ===
                              "COMPLETED_INTERVIEW"
                            ? "success"
                            : "error"
                        }
                      />
                    </Grid>
                    <Grid item xs={12}>
                      <Typography variant="subtitle2" gutterBottom>
                        {t("interviews.interviewers")}
                      </Typography>
                      <Box sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}>
                        {interview?.interviewers &&
                        interview.interviewers.length > 0 ? (
                          interview.interviewers.map((interviewer, index) => {
                            const isCurrentUser =
                              interviewer.email === userEmail;
                            return (
                              <Box
                                key={index}
                                sx={{
                                  display: "flex",
                                  flexDirection: "column",
                                  alignItems: "center",
                                  gap: 1,
                                }}
                              >
                                <Chip
                                  icon={<PersonIcon />}
                                  label={
                                    <Box
                                      sx={{
                                        display: "flex",
                                        alignItems: "center",
                                        gap: 1,
                                      }}
                                    >
                                      <span>{interviewer.name}</span>
                                      {isCurrentUser && (
                                        <Typography
                                          variant="caption"
                                          sx={{
                                            color: "primary.main",
                                            fontWeight: "bold",
                                          }}
                                        >
                                          ({t("interviews.you")})
                                        </Typography>
                                      )}
                                    </Box>
                                  }
                                  variant="outlined"
                                  sx={{
                                    bgcolor: isCurrentUser
                                      ? "primary.light"
                                      : undefined,
                                  }}
                                />
                                {isCurrentUser &&
                                  interview.interview_state ===
                                    "SCHEDULED_INTERVIEW" && (
                                    <ButtonGroup size="small">
                                      <Button
                                        variant="contained"
                                        color="success"
                                        onClick={() =>
                                          handleRSVP(RSVPStatuses.YES)
                                        }
                                        disabled={rsvpLoading}
                                      >
                                        {t("interviews.rsvp.yes")}
                                      </Button>
                                      <Button
                                        variant="contained"
                                        color="error"
                                        onClick={() =>
                                          handleRSVP(RSVPStatuses.NO)
                                        }
                                        disabled={rsvpLoading}
                                      >
                                        {t("interviews.rsvp.no")}
                                      </Button>
                                    </ButtonGroup>
                                  )}
                              </Box>
                            );
                          })
                        ) : (
                          <Typography color="text.secondary">
                            {t("interviews.noInterviewers")}
                          </Typography>
                        )}
                      </Box>
                    </Grid>
                  </Grid>
                </AccordionDetails>
              </Accordion>
            </Grid>

            <Grid item xs={12}>
              <Accordion defaultExpanded>
                <AccordionSummary
                  expandIcon={<ExpandMoreIcon />}
                  sx={{
                    "& .MuiAccordionSummary-content": {
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "center",
                    },
                  }}
                >
                  <Typography variant="h6">
                    {t("interviews.assessment.title")}
                  </Typography>
                  <Box sx={{ display: "flex", gap: 1 }}>
                    {interview?.interviewers_decision && (
                      <Chip
                        label={t(
                          `interviews.assessment.ratings.${
                            Object.entries(InterviewersDecisions).find(
                              ([_, value]) =>
                                value === interview.interviewers_decision
                            )?.[0]
                          }`
                        )}
                        size="small"
                        color={
                          interview.interviewers_decision ===
                            InterviewersDecisions.STRONG_YES ||
                          interview.interviewers_decision ===
                            InterviewersDecisions.YES
                            ? "success"
                            : interview.interviewers_decision ===
                              InterviewersDecisions.NEUTRAL
                            ? "primary"
                            : "error"
                        }
                      />
                    )}
                  </Box>
                </AccordionSummary>
                <AccordionDetails>
                  <FormControl fullWidth sx={{ mb: 3 }}>
                    <InputLabel id="rating-label">
                      {t("interviews.assessment.rating")}
                    </InputLabel>
                    <Select
                      labelId="rating-label"
                      value={interview?.interviewers_decision || ""}
                      label={t("interviews.assessment.rating")}
                      onChange={(e) =>
                        setInterview((prev) => ({
                          ...prev!,
                          interviewers_decision: e.target
                            .value as InterviewersDecision,
                        }))
                      }
                    >
                      {Object.entries(InterviewersDecisions).map(([key]) => (
                        <MenuItem
                          key={key}
                          value={
                            InterviewersDecisions[
                              key as keyof typeof InterviewersDecisions
                            ]
                          }
                        >
                          {t(`interviews.assessment.ratings.${key}`)}
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>

                  <TextField
                    fullWidth
                    multiline
                    rows={4}
                    label={t("interviews.assessment.positives")}
                    placeholder={t(
                      "interviews.assessment.positivesPlaceholder"
                    )}
                    value={interview?.positives || ""}
                    onChange={(e) =>
                      setInterview((prev) => ({
                        ...prev!,
                        positives: e.target.value,
                      }))
                    }
                    sx={{ mb: 3 }}
                  />

                  <TextField
                    fullWidth
                    multiline
                    rows={4}
                    label={t("interviews.assessment.negatives")}
                    placeholder={t(
                      "interviews.assessment.negativesPlaceholder"
                    )}
                    value={interview?.negatives || ""}
                    onChange={(e) =>
                      setInterview((prev) => ({
                        ...prev!,
                        negatives: e.target.value,
                      }))
                    }
                    sx={{ mb: 3 }}
                  />

                  <TextField
                    fullWidth
                    multiline
                    rows={4}
                    label={t("interviews.assessment.overallAssessment")}
                    placeholder={t(
                      "interviews.assessment.overallAssessmentPlaceholder"
                    )}
                    value={interview?.overall_assessment || ""}
                    onChange={(e) =>
                      setInterview((prev) => ({
                        ...prev!,
                        overall_assessment: e.target.value,
                      }))
                    }
                    sx={{ mb: 3 }}
                  />

                  <TextField
                    fullWidth
                    multiline
                    rows={4}
                    label={
                      <Box
                        sx={{ display: "flex", alignItems: "center", gap: 1 }}
                      >
                        <PublicIcon color="action" fontSize="small" />
                        {t("interviews.assessment.feedback")}
                      </Box>
                    }
                    placeholder={t("interviews.assessment.feedbackPlaceholder")}
                    value={interview?.feedback_to_candidate || ""}
                    onChange={(e) =>
                      setInterview((prev) => ({
                        ...prev!,
                        feedback_to_candidate: e.target.value,
                      }))
                    }
                    sx={{
                      mb: 3,
                      "& .MuiOutlinedInput-root": {
                        bgcolor: (theme) => theme.palette.warning.light + "10",
                      },
                    }}
                  />

                  {interview?.feedback_submitted_by && (
                    <Box sx={{ mb: 3 }}>
                      <Typography variant="body2" color="text.secondary">
                        {`${t("interviews.assessment.lastUpdated")
                          .replace(
                            "{{name}}",
                            interview.feedback_submitted_by.name
                          )
                          .replace(
                            "{{date}}",
                            formatDateTime(
                              interview.feedback_submitted_at || "",
                              {
                                dateStyle: "medium",
                                timeStyle: "short",
                              }
                            )
                          )}`}
                      </Typography>
                    </Box>
                  )}

                  <Button
                    variant="contained"
                    onClick={handleSave}
                    disabled={saving}
                    sx={{ minWidth: 120 }}
                  >
                    {saving ? (
                      <CircularProgress size={24} color="inherit" />
                    ) : (
                      t("interviews.assessment.save")
                    )}
                  </Button>
                </AccordionDetails>
              </Accordion>
            </Grid>
          </Grid>
        </Box>
      </Container>
    </AuthenticatedLayout>
  );
}
