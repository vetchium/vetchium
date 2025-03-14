"use client";

import AuthenticatedLayout from "@/components/AuthenticatedLayout";
import { config } from "@/config";
import { useAuth } from "@/contexts/AuthContext";
import { useTranslation } from "@/hooks/useTranslation";
import {
  Cancel as CancelIcon,
  CheckCircle as CheckCircleIcon,
  Edit as EditIcon,
  ExpandMore as ExpandMoreIcon,
  Person as PersonIcon,
  Public as PublicIcon,
} from "@mui/icons-material";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Alert,
  Box,
  Button,
  ButtonGroup,
  Card,
  CardContent,
  Checkbox,
  Chip,
  CircularProgress,
  Container,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Divider,
  FormControl,
  FormControlLabel,
  Grid,
  InputLabel,
  MenuItem,
  Paper,
  Select,
  TextField,
  Typography,
} from "@mui/material";
import {
  EmployerInterview,
  InterviewersDecision,
  InterviewersDecisions,
  PutAssessmentRequest,
  RSVPInterviewRequest,
  RSVPStatus,
  RSVPStatuses,
} from "@psankar/vetchi-typespec";
import Cookies from "js-cookie";
import { useParams, useRouter } from "next/navigation";
import { useCallback, useEffect, useState } from "react";
import { styled, alpha } from "@mui/material/styles";

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

// FeedbackWidget component for displaying feedback sections (read-only)
interface FeedbackWidgetProps {
  title: string;
  value: string;
  placeholder: string;
  icon?: React.ReactNode;
}

const StyledPaper = styled(Paper)(({ theme }) => ({
  padding: theme.spacing(2),
  marginBottom: theme.spacing(3),
  backgroundColor: alpha(theme.palette.warning.light, 0.1),
}));

// Update FeedbackWidget to use StyledFeedbackCard
function FeedbackWidget({
  title,
  value,
  placeholder,
  icon,
}: FeedbackWidgetProps) {
  return (
    <Card
      sx={{
        mb: 2,
        "& .MuiCardContent-root": {
          bgcolor: (theme) => alpha(theme.palette.warning.light, 0.1),
        },
      }}
    >
      <CardContent>
        <Box
          sx={{
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            mb: 2,
          }}
        >
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            {icon}
            <Typography variant="h6">{title}</Typography>
          </Box>
        </Box>
        <Typography
          variant="body1"
          sx={{
            whiteSpace: "pre-wrap",
            color: value ? "text.primary" : "text.secondary",
            fontStyle: value ? "normal" : "italic",
          }}
        >
          {value || placeholder}
        </Typography>
      </CardContent>
    </Card>
  );
}

export default function InterviewDetailPage() {
  const params = useParams();
  const interviewId = params.id as string;
  const router = useRouter();
  const { t } = useTranslation();
  const { userEmail } = useAuth();

  const [interview, setInterview] = useState<EmployerInterview | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [rsvpLoading, setRsvpLoading] = useState(false);
  const [confirmDialog, setConfirmDialog] = useState<{
    open: boolean;
    status: RSVPStatus;
  }>({
    open: false,
    status: RSVPStatuses.NOT_SET,
  });
  const [editDialogOpen, setEditDialogOpen] = useState(false);
  const [editFormData, setEditFormData] = useState<{
    interviewers_decision: InterviewersDecision;
    positives: string;
    negatives: string;
    overall_assessment: string;
    feedback_to_candidate: string;
    mark_interview_completed: boolean;
  }>({
    interviewers_decision: InterviewersDecisions.NEUTRAL,
    positives: "",
    negatives: "",
    overall_assessment: "",
    feedback_to_candidate: "",
    mark_interview_completed: false,
  });

  const fetchInterview = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/signin");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/employer/get-interview/${interviewId}`,
        {
          headers: {
            Authorization: `Bearer ${token}`,
          },
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/signin");
        return;
      }

      if (!response.ok) {
        throw new Error(t("interviews.fetchError"));
      }

      const data = await response.json();
      setInterview(data);
    } catch {
      setError(t("interviews.fetchError"));
    } finally {
      setLoading(false);
    }
  }, [interviewId, router, t]);

  useEffect(() => {
    fetchInterview();
  }, [fetchInterview]);

  const handleRSVPClick = (status: RSVPStatus) => {
    setConfirmDialog({
      open: true,
      status,
    });
  };

  const handleConfirmRSVP = () => {
    handleRSVP(confirmDialog.status);
    setConfirmDialog((prev) => ({ ...prev, open: false }));
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

      if (response.status === 403) {
        throw new Error(t("interviews.assessment.forbiddenError"));
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
    } catch {
      setError(t("common.error"));
    } finally {
      setRsvpLoading(false);
    }
  };

  const isInterviewer = Boolean(
    interview?.interviewers?.some(
      (interviewer) => interviewer.email === userEmail
    )
  );

  const handleFeedbackSave = async (updates: Partial<EmployerInterview>) => {
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
        decision:
          updates.interviewers_decision ||
          interview?.interviewers_decision ||
          InterviewersDecisions.NEUTRAL,
        positives: updates.positives || interview?.positives || "",
        negatives: updates.negatives || interview?.negatives || "",
        overall_assessment:
          updates.overall_assessment || interview?.overall_assessment || "",
        feedback_to_candidate:
          updates.feedback_to_candidate ||
          interview?.feedback_to_candidate ||
          "",
        mark_interview_completed: editFormData.mark_interview_completed,
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

      if (!response.ok) {
        throw new Error(t("interviews.assessment.saveError"));
      }

      setSuccessMessage(t("interviews.assessment.saveSuccess"));
      await fetchInterview();
    } catch {
      setError(t("common.error"));
    } finally {
      setSaving(false);
    }
  };

  const isInterviewScheduled =
    interview?.interview_state === "SCHEDULED_INTERVIEW";

  // Update edit form data when interview data changes
  useEffect(() => {
    if (interview) {
      setEditFormData({
        interviewers_decision:
          interview.interviewers_decision || InterviewersDecisions.NEUTRAL,
        positives: interview.positives || "",
        negatives: interview.negatives || "",
        overall_assessment: interview.overall_assessment || "",
        feedback_to_candidate: interview.feedback_to_candidate || "",
        mark_interview_completed: false,
      });
    }
  }, [interview]);

  const handleEditDialogClose = () => {
    setEditDialogOpen(false);
  };

  const handleEditDialogSave = async () => {
    await handleFeedbackSave(editFormData);
    setEditDialogOpen(false);
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
                      <Box
                        sx={{
                          display: "flex",
                          flexDirection: "column",
                          alignItems: "center",
                        }}
                      >
                        <Box>
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
                              (
                              {Intl.DateTimeFormat().resolvedOptions().timeZone}
                              )
                            </Typography>
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            UTC:{" "}
                            {formatUTCDateTime(interview?.start_time || "")}
                          </Typography>
                        </Box>

                        <Typography color="text.secondary" sx={{ my: 1 }}>
                          -
                        </Typography>

                        <Box>
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
                              (
                              {Intl.DateTimeFormat().resolvedOptions().timeZone}
                              )
                            </Typography>
                          </Typography>
                          <Typography variant="caption" color="text.secondary">
                            UTC: {formatUTCDateTime(interview?.end_time || "")}
                          </Typography>
                        </Box>
                      </Box>
                    </Grid>

                    <Grid item xs={12}>
                      <Typography
                        variant="subtitle2"
                        color="text.secondary"
                        gutterBottom
                      >
                        {t("interviews.candidate")}
                      </Typography>
                      <Box
                        sx={{ display: "flex", alignItems: "center", gap: 1 }}
                      >
                        <Chip
                          icon={<PersonIcon />}
                          variant="outlined"
                          label={
                            <Box
                              sx={{
                                display: "flex",
                                alignItems: "center",
                                gap: 1,
                              }}
                            >
                              <Typography>
                                {interview?.candidate_name} (
                                {interview?.candidate_handle})
                              </Typography>
                              {interview?.candidate_rsvp_status !==
                                RSVPStatuses.NOT_SET && (
                                <Box
                                  component="span"
                                  sx={{
                                    display: "inline-flex",
                                    verticalAlign: "middle",
                                  }}
                                >
                                  {interview?.candidate_rsvp_status ===
                                  RSVPStatuses.YES ? (
                                    <CheckCircleIcon
                                      color="success"
                                      fontSize="small"
                                    />
                                  ) : (
                                    <CancelIcon
                                      color="error"
                                      fontSize="small"
                                    />
                                  )}
                                </Box>
                              )}
                            </Box>
                          }
                        />
                      </Box>
                    </Grid>

                    {interview?.description && (
                      <Grid item xs={12}>
                        <Typography variant="body1">
                          {interview.description}
                        </Typography>
                      </Grid>
                    )}

                    <Grid item xs={12}>
                      <Typography variant="caption" color="text.secondary">
                        {t("interviews.interviewers")}
                      </Typography>
                      {interview?.interviewers &&
                      interview.interviewers.length > 0 ? (
                        <Box
                          sx={{
                            display: "flex",
                            flexDirection: "column",
                            gap: 2,
                          }}
                        >
                          {/* Current user's interviewer card if they are an interviewer */}
                          {interview.interviewers
                            .filter(
                              (interviewer) => interviewer.email === userEmail
                            )
                            .map((interviewer, index) => (
                              <Box
                                key={index}
                                sx={{
                                  display: "flex",
                                  flexDirection: "column",
                                  alignItems: "flex-start",
                                  gap: 1,
                                  width: "100%",
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
                                      <Typography
                                        variant="caption"
                                        sx={{
                                          color: "primary.main",
                                          fontWeight: "bold",
                                        }}
                                      >
                                        ({t("interviews.you")})
                                      </Typography>
                                    </Box>
                                  }
                                  variant="outlined"
                                  sx={{
                                    bgcolor: "primary.light",
                                    mb: 1,
                                  }}
                                />
                                {interview.interview_state ===
                                  "SCHEDULED_INTERVIEW" && (
                                  <Box
                                    sx={{
                                      display: "flex",
                                      flexDirection: "column",
                                      gap: 1,
                                      width: "100%",
                                    }}
                                  >
                                    <Box
                                      sx={{
                                        display: "flex",
                                        alignItems: "center",
                                        gap: 1,
                                      }}
                                    >
                                      <Typography variant="subtitle2">
                                        {t("interviews.yourRSVP")}:
                                      </Typography>
                                      {interviewer.rsvp_status !==
                                        RSVPStatuses.NOT_SET && (
                                        <Box
                                          component="span"
                                          sx={{
                                            display: "flex",
                                            alignItems: "center",
                                          }}
                                        >
                                          {interviewer.rsvp_status ===
                                          RSVPStatuses.YES ? (
                                            <CheckCircleIcon
                                              color="success"
                                              fontSize="small"
                                            />
                                          ) : (
                                            <CancelIcon
                                              color="error"
                                              fontSize="small"
                                            />
                                          )}
                                        </Box>
                                      )}
                                    </Box>
                                    <ButtonGroup size="small">
                                      <Button
                                        variant={
                                          interviewer.rsvp_status ===
                                          RSVPStatuses.YES
                                            ? "contained"
                                            : "outlined"
                                        }
                                        onClick={() =>
                                          handleRSVPClick(RSVPStatuses.YES)
                                        }
                                        color="success"
                                        disabled={
                                          rsvpLoading ||
                                          interviewer.rsvp_status ===
                                            RSVPStatuses.YES ||
                                          !isInterviewScheduled
                                        }
                                        sx={{
                                          "&.Mui-disabled": {
                                            backgroundColor:
                                              interviewer.rsvp_status ===
                                              RSVPStatuses.YES
                                                ? "success.main"
                                                : "transparent",
                                            color:
                                              interviewer.rsvp_status ===
                                              RSVPStatuses.YES
                                                ? "white"
                                                : undefined,
                                            opacity: 0.7,
                                          },
                                        }}
                                      >
                                        {t("interviews.rsvp.yes")}
                                      </Button>
                                      <Button
                                        variant={
                                          interviewer.rsvp_status ===
                                          RSVPStatuses.NO
                                            ? "contained"
                                            : "outlined"
                                        }
                                        onClick={() =>
                                          handleRSVPClick(RSVPStatuses.NO)
                                        }
                                        color="error"
                                        disabled={
                                          rsvpLoading ||
                                          interviewer.rsvp_status ===
                                            RSVPStatuses.NO ||
                                          !isInterviewScheduled
                                        }
                                        sx={{
                                          "&.Mui-disabled": {
                                            backgroundColor:
                                              interviewer.rsvp_status ===
                                              RSVPStatuses.NO
                                                ? "error.main"
                                                : "transparent",
                                            color:
                                              interviewer.rsvp_status ===
                                              RSVPStatuses.NO
                                                ? "white"
                                                : undefined,
                                            opacity: 0.7,
                                          },
                                        }}
                                      >
                                        {t("interviews.rsvp.no")}
                                      </Button>
                                    </ButtonGroup>
                                  </Box>
                                )}
                              </Box>
                            ))}

                          {/* Divider and "Other Interviewers" title only if user is an interviewer */}
                          {interview.interviewers.some(
                            (interviewer) => interviewer.email === userEmail
                          ) &&
                            interview.interviewers.some(
                              (interviewer) => interviewer.email !== userEmail
                            ) && (
                              <>
                                <Divider sx={{ width: "100%" }} />
                                <Typography
                                  variant="subtitle2"
                                  color="text.secondary"
                                >
                                  {t("interviews.otherInterviewers")}
                                </Typography>
                              </>
                            )}

                          {/* Other interviewers */}
                          <Box
                            sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}
                          >
                            {(interview?.interviewers || [])
                              .filter(
                                (interviewer) =>
                                  interviewer.email !== userEmail ||
                                  !interview?.interviewers?.some(
                                    (i) => i.email === userEmail
                                  )
                              )
                              .map((interviewer, index) => (
                                <Chip
                                  key={index}
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
                                      {interviewer.rsvp_status !==
                                        RSVPStatuses.NOT_SET && (
                                        <Box
                                          component="span"
                                          sx={{
                                            display: "flex",
                                            alignItems: "center",
                                          }}
                                        >
                                          {interviewer.rsvp_status ===
                                          RSVPStatuses.YES ? (
                                            <CheckCircleIcon
                                              color="success"
                                              fontSize="small"
                                            />
                                          ) : (
                                            <CancelIcon
                                              color="error"
                                              fontSize="small"
                                            />
                                          )}
                                        </Box>
                                      )}
                                    </Box>
                                  }
                                  variant="outlined"
                                />
                              ))}
                          </Box>
                        </Box>
                      ) : (
                        <Typography color="text.secondary">
                          {t("interviews.noInterviewers")}
                        </Typography>
                      )}
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
                          `interviews.assessment.ratings.${Object.keys(
                            InterviewersDecisions
                          ).find(
                            (key) =>
                              InterviewersDecisions[
                                key as keyof typeof InterviewersDecisions
                              ] === interview.interviewers_decision
                          )}`
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
                  {isInterviewer && isInterviewScheduled && (
                    <Box
                      sx={{
                        display: "flex",
                        justifyContent: "flex-end",
                        mb: 2,
                      }}
                    >
                      <Button
                        variant="contained"
                        startIcon={<EditIcon />}
                        onClick={() => setEditDialogOpen(true)}
                      >
                        {t("interviews.assessment.editFeedback")}
                      </Button>
                    </Box>
                  )}

                  <FeedbackWidget
                    title={t("interviews.assessment.rating")}
                    value={
                      interview?.interviewers_decision
                        ? t(
                            `interviews.assessment.ratings.${Object.keys(
                              InterviewersDecisions
                            ).find(
                              (key) =>
                                InterviewersDecisions[
                                  key as keyof typeof InterviewersDecisions
                                ] === interview.interviewers_decision
                            )}`
                          )
                        : ""
                    }
                    placeholder={t("interviews.assessment.ratingPlaceholder")}
                  />

                  <FeedbackWidget
                    title={t("interviews.assessment.positives")}
                    value={interview?.positives || ""}
                    placeholder={t(
                      "interviews.assessment.positivesPlaceholder"
                    )}
                  />

                  <FeedbackWidget
                    title={t("interviews.assessment.negatives")}
                    value={interview?.negatives || ""}
                    placeholder={t(
                      "interviews.assessment.negativesPlaceholder"
                    )}
                  />

                  <FeedbackWidget
                    title={t("interviews.assessment.overallAssessment")}
                    value={interview?.overall_assessment || ""}
                    placeholder={t(
                      "interviews.assessment.overallAssessmentPlaceholder"
                    )}
                  />

                  <FeedbackWidget
                    title={t("interviews.assessment.feedback")}
                    value={interview?.feedback_to_candidate || ""}
                    placeholder={t("interviews.assessment.feedbackPlaceholder")}
                    icon={<PublicIcon color="action" fontSize="small" />}
                  />

                  {interview?.feedback_submitted_by && (
                    <Box sx={{ mt: 2 }}>
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

                  {/* Edit Dialog */}
                  <Dialog
                    open={editDialogOpen}
                    onClose={handleEditDialogClose}
                    maxWidth="md"
                    fullWidth
                  >
                    <DialogTitle>
                      {t("interviews.assessment.editFeedback")}
                    </DialogTitle>
                    <DialogContent>
                      <FormControl fullWidth sx={{ mt: 2, mb: 3 }}>
                        <InputLabel id="rating-label">
                          {t("interviews.assessment.rating")}
                        </InputLabel>
                        <Select
                          labelId="rating-label"
                          value={editFormData.interviewers_decision}
                          label={t("interviews.assessment.rating")}
                          onChange={(e) =>
                            setEditFormData((prev) => ({
                              ...prev,
                              interviewers_decision: e.target
                                .value as InterviewersDecision,
                            }))
                          }
                        >
                          {Object.keys(InterviewersDecisions).map((key) => (
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
                        value={editFormData.positives}
                        onChange={(e) =>
                          setEditFormData((prev) => ({
                            ...prev,
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
                        value={editFormData.negatives}
                        onChange={(e) =>
                          setEditFormData((prev) => ({
                            ...prev,
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
                        value={editFormData.overall_assessment}
                        onChange={(e) =>
                          setEditFormData((prev) => ({
                            ...prev,
                            overall_assessment: e.target.value,
                          }))
                        }
                        sx={{ mb: 3 }}
                      />

                      <StyledPaper elevation={0}>
                        <TextField
                          fullWidth
                          multiline
                          rows={4}
                          label={t("interviews.assessment.feedback")}
                          placeholder={t(
                            "interviews.assessment.feedbackPlaceholder"
                          )}
                          value={editFormData.feedback_to_candidate}
                          onChange={(e) =>
                            setEditFormData((prev) => ({
                              ...prev,
                              feedback_to_candidate: e.target.value,
                            }))
                          }
                        />
                      </StyledPaper>

                      <Box sx={{ mb: 3, mt: 2 }}>
                        <FormControlLabel
                          control={
                            <Checkbox
                              checked={editFormData.mark_interview_completed}
                              onChange={(e) =>
                                setEditFormData((prev) => ({
                                  ...prev,
                                  mark_interview_completed: e.target.checked,
                                }))
                              }
                              color="primary"
                            />
                          }
                          label={
                            <Typography variant="body1" color="text.primary">
                              {t("interviews.assessment.markAsCompleted")}
                            </Typography>
                          }
                        />
                      </Box>
                    </DialogContent>
                    <DialogActions>
                      <Button onClick={handleEditDialogClose}>
                        {t("common.cancel")}
                      </Button>
                      <Button
                        onClick={handleEditDialogSave}
                        variant="contained"
                        disabled={saving}
                      >
                        {saving ? (
                          <CircularProgress size={24} color="inherit" />
                        ) : (
                          t("common.save")
                        )}
                      </Button>
                    </DialogActions>
                  </Dialog>
                </AccordionDetails>
              </Accordion>
            </Grid>
          </Grid>
        </Box>
      </Container>
      <Dialog
        open={confirmDialog.open}
        onClose={() => setConfirmDialog((prev) => ({ ...prev, open: false }))}
        maxWidth="sm"
        fullWidth
      >
        <DialogTitle>
          {confirmDialog.status === RSVPStatuses.YES
            ? t("interviews.rsvp.confirmYes")
            : t("interviews.rsvp.confirmNo")}
        </DialogTitle>
        <DialogContent>
          <Typography>
            {interview?.interviewers?.find((i) => i.email === userEmail)
              ?.rsvp_status === RSVPStatuses.NOT_SET
              ? confirmDialog.status === RSVPStatuses.YES
                ? t("interviews.rsvp.confirmYesMessage")
                : t("interviews.rsvp.confirmNoMessage")
              : confirmDialog.status === RSVPStatuses.YES
              ? t("interviews.rsvp.confirmChangeYesMessage")
              : t("interviews.rsvp.confirmChangeNoMessage")}
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button
            onClick={() =>
              setConfirmDialog((prev) => ({ ...prev, open: false }))
            }
            color="inherit"
          >
            {t("common.cancel")}
          </Button>
          <Button
            onClick={handleConfirmRSVP}
            color={
              confirmDialog.status === RSVPStatuses.YES ? "success" : "error"
            }
            variant="contained"
            autoFocus
          >
            {confirmDialog.status === RSVPStatuses.YES
              ? t("interviews.rsvp.yes")
              : t("interviews.rsvp.no")}
          </Button>
        </DialogActions>
      </Dialog>
    </AuthenticatedLayout>
  );
}
