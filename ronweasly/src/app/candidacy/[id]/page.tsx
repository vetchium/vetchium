"use client";

import { useParams, useRouter } from "next/navigation";
import { useState, useEffect } from "react";
import {
  Box,
  Typography,
  Button,
  TextField,
  Paper,
  Divider,
  Link,
  CircularProgress,
  IconButton,
  Avatar,
  Chip,
  Accordion,
  AccordionSummary,
  AccordionDetails,
  ButtonGroup,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
} from "@mui/material";
import {
  OpenInNew as OpenInNewIcon,
  ExpandMore as ExpandMoreIcon,
  CheckCircle as CheckCircleIcon,
  Cancel as CancelIcon,
  HelpOutline as HelpOutlineIcon,
  Person as PersonIcon,
} from "@mui/icons-material";
import { useTranslation } from "@/hooks/useTranslation";
import {
  GetCandidacyInfoRequest,
  GetCandidacyCommentsRequest,
  MyCandidacy,
  CandidacyComment,
  CandidacyState,
  AddHubCandidacyCommentRequest,
  GetHubInterviewsByCandidacyRequest,
  HubInterview,
  InterviewState,
  InterviewType,
  InterviewStates,
  InterviewTypes,
  CommenterTypes,
  RSVPStatus,
  RSVPInterviewRequest,
  RSVPStatuses,
} from "@psankar/vetchi-typespec";
import { config } from "@/config";
import Cookies from "js-cookie";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";

// Helper function for consistent date formatting
const formatDateTime = (
  date: string | Date,
  options?: Intl.DateTimeFormatOptions
) => {
  const defaultOptions: Intl.DateTimeFormatOptions = {
    dateStyle: "full",
    timeStyle: "short",
  };
  return new Date(date).toLocaleString(undefined, options || defaultOptions);
};

function CandidacyStateLabel({
  state,
  t,
}: {
  state: CandidacyState;
  t: (key: string) => string;
}) {
  let color:
    | "primary"
    | "secondary"
    | "error"
    | "info"
    | "success"
    | "warning" = "info";
  switch (state) {
    case "INTERVIEWING":
      color = "info";
      break;
    case "OFFERED":
      color = "warning";
      break;
    case "OFFER_ACCEPTED":
      color = "success";
      break;
    case "OFFER_DECLINED":
    case "CANDIDATE_UNSUITABLE":
    case "CANDIDATE_NOT_RESPONDING":
    case "CANDIDATE_WITHDREW":
    case "EMPLOYER_DEFUNCT":
      color = "error";
      break;
  }
  return (
    <Chip label={t(`candidacies.states.${state}`)} color={color} size="small" />
  );
}

const RSVPStatusIcon = ({ status }: { status: RSVPStatus }) => {
  switch (status) {
    case RSVPStatuses.YES:
      return <CheckCircleIcon color="success" />;
    case RSVPStatuses.NO:
      return <CancelIcon color="error" />;
    case RSVPStatuses.NOT_SET:
    default:
      return null;
  }
};

export default function CandidacyDetailPage() {
  const params = useParams();
  const candidacyId = params?.id as string;
  if (!candidacyId) {
    return <div>Invalid candidacy ID</div>;
  }
  const { t } = useTranslation();
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [candidacy, setCandidacy] = useState<MyCandidacy | null>(null);
  const [comments, setComments] = useState<CandidacyComment[]>([]);
  const [interviews, setInterviews] = useState<HubInterview[]>([]);
  const [newComment, setNewComment] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [rsvpLoading, setRsvpLoading] = useState<string | null>(null);
  const [confirmDialog, setConfirmDialog] = useState<{
    open: boolean;
    interviewId: string;
    status: RSVPStatus;
  }>({
    open: false,
    interviewId: "",
    status: RSVPStatuses.NOT_SET,
  });

  // Fetch candidacy info and comments
  const fetchData = async () => {
    setLoading(true);
    setError(null);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      // Fetch candidacy info
      const infoResponse = await fetch(
        `${config.API_SERVER_PREFIX}/hub/get-candidacy-info`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            candidacy_id: candidacyId,
          } as GetCandidacyInfoRequest),
        }
      );

      if (infoResponse.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!infoResponse.ok) throw new Error(t("candidacies.fetchError"));
      const candidacyData = await infoResponse.json();
      if (!candidacyData) {
        throw new Error(t("candidacies.fetchError"));
      }
      setCandidacy(candidacyData);

      // Fetch comments
      const commentsResponse = await fetch(
        `${config.API_SERVER_PREFIX}/hub/get-candidacy-comments`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            candidacy_id: candidacyId,
          } as GetCandidacyCommentsRequest),
        }
      );

      if (!commentsResponse.ok) throw new Error(t("candidacies.fetchError"));
      const commentsData = await commentsResponse.json();
      setComments(commentsData || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
    } finally {
      setLoading(false);
    }
  };

  // Fetch interviews
  const fetchInterviews = async () => {
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/get-interviews-by-candidacy`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            candidacy_id: candidacyId,
          } as GetHubInterviewsByCandidacyRequest),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) throw new Error(t("interviews.fetchError"));
      const interviewsData = await response.json();
      setInterviews(interviewsData || []);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
    }
  };

  // Add new comment
  const handleAddComment = async () => {
    if (!newComment.trim()) return;
    setSubmitting(true);
    try {
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/add-candidacy-comment`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${token}`,
          },
          body: JSON.stringify({
            candidacy_id: candidacyId,
            comment: newComment.trim(),
          } as AddHubCandidacyCommentRequest),
        }
      );

      if (response.status === 401) {
        Cookies.remove("session_token");
        router.push("/login");
        return;
      }

      if (!response.ok) throw new Error(t("common.error.serverError"));
      setNewComment("");
      // Refresh comments
      await fetchData();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
    } finally {
      setSubmitting(false);
    }
  };

  const handleRSVPClick = (
    interviewId: string,
    status: RSVPStatus,
    currentStatus: RSVPStatus
  ) => {
    // Always show confirmation dialog when changing status
    setConfirmDialog({
      open: true,
      interviewId,
      status,
    });
  };

  const handleConfirmRSVP = () => {
    handleRSVP(confirmDialog.interviewId, confirmDialog.status);
    setConfirmDialog((prev) => ({ ...prev, open: false }));
  };

  // Add RSVP mutation function
  const handleRSVP = async (interviewId: string, status: RSVPStatus) => {
    try {
      setRsvpLoading(interviewId);
      const token = Cookies.get("session_token");
      if (!token) {
        router.push("/login");
        return;
      }

      const request: RSVPInterviewRequest = {
        interview_id: interviewId,
        rsvp_status: status,
      };

      const response = await fetch(
        `${config.API_SERVER_PREFIX}/hub/rsvp-interview`,
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

      if (!response.ok) throw new Error(t("interviews.rsvpError"));

      // Refresh interviews after successful RSVP
      await fetchInterviews();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : t("common.error.serverError")
      );
    } finally {
      setRsvpLoading(null);
    }
  };

  // Fetch data on mount
  useEffect(() => {
    fetchData();
    fetchInterviews();
  }, []); // Empty dependency array means this runs once on mount

  const getInterviewStateColor = (
    state: InterviewState
  ): "primary" | "success" | "error" => {
    switch (state) {
      case InterviewStates.SCHEDULED_INTERVIEW:
        return "primary";
      case InterviewStates.COMPLETED_INTERVIEW:
        return "success";
      case InterviewStates.CANCELLED_INTERVIEW:
      default:
        return "error";
    }
  };

  const getInterviewTypeLabel = (type: InterviewType) => {
    switch (type) {
      case InterviewTypes.IN_PERSON:
        return t("interviews.types.IN_PERSON");
      case InterviewTypes.VIDEO_CALL:
        return t("interviews.types.VIDEO_CALL");
      case InterviewTypes.TAKE_HOME:
        return t("interviews.types.TAKE_HOME");
      case InterviewTypes.OTHER_INTERVIEW:
        return t("interviews.types.OTHER_INTERVIEW");
      default:
        return type; // fallback
    }
  };

  const getInterviewStateLabel = (state: InterviewState): string => {
    switch (state) {
      case InterviewStates.SCHEDULED_INTERVIEW:
        return t("interviews.states.SCHEDULED_INTERVIEW");
      case InterviewStates.COMPLETED_INTERVIEW:
        return t("interviews.states.COMPLETED_INTERVIEW");
      case InterviewStates.CANCELLED_INTERVIEW:
        return t("interviews.states.CANCELLED_INTERVIEW");
      default:
        return state;
    }
  };

  const content = (
    <Box sx={{ p: 3 }}>
      <Box sx={{ display: "flex", justifyContent: "space-between", mb: 3 }}>
        <Typography variant="h4">{t("candidacies.viewCandidacy")}</Typography>
        <Button variant="outlined" onClick={() => router.back()}>
          {t("common.back")}
        </Button>
      </Box>

      {loading ? (
        <Box sx={{ display: "flex", justifyContent: "center", p: 3 }}>
          <CircularProgress />
        </Box>
      ) : error ? (
        <Box sx={{ p: 3 }}>
          <Typography color="error">{error}</Typography>
          <Button variant="contained" onClick={fetchData} sx={{ mt: 2 }}>
            {t("common.retry")}
          </Button>
        </Box>
      ) : (
        candidacy && (
          <>
            <Paper sx={{ p: 3, mb: 3 }}>
              <Box
                sx={{
                  display: "flex",
                  justifyContent: "space-between",
                  alignItems: "center",
                  mb: 1,
                }}
              >
                <Typography variant="h6">{candidacy.opening_title}</Typography>
                <CandidacyStateLabel state={candidacy.candidacy_state} t={t} />
              </Box>
              <Typography
                variant="subtitle1"
                gutterBottom
                sx={{ color: "text.secondary" }}
              >
                {candidacy.company_name}
              </Typography>
              {candidacy.company_domain && (
                <Link
                  href={`https://${candidacy.company_domain}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  sx={{
                    color: "primary.main",
                    textDecoration: "none",
                    display: "flex",
                    alignItems: "center",
                    gap: 0.5,
                    mb: 2,
                  }}
                >
                  <Typography variant="body2">
                    {candidacy.company_domain}
                  </Typography>
                  <OpenInNewIcon fontSize="small" />
                </Link>
              )}
              <Divider sx={{ my: 2 }} />
              <Typography variant="body2" color="text.secondary">
                {candidacy.opening_description}
              </Typography>
            </Paper>

            {/* Interviews Section */}
            <Accordion defaultExpanded sx={{ mb: 3 }}>
              <AccordionSummary expandIcon={<ExpandMoreIcon />}>
                <Typography variant="h6">
                  {t("interviews.title")} ({interviews.length})
                </Typography>
              </AccordionSummary>
              <AccordionDetails>
                {interviews.length > 0 ? (
                  interviews.map((interview) => (
                    <Accordion
                      key={interview.interview_id}
                      sx={{
                        "&:not(:last-child)": { mb: 1 },
                        borderTop: "1px solid",
                        borderColor: "divider",
                      }}
                    >
                      <AccordionSummary
                        expandIcon={<ExpandMoreIcon />}
                        sx={{
                          flexDirection: "row-reverse",
                          "& .MuiAccordionSummary-expandIconWrapper": {
                            mr: 1,
                          },
                        }}
                      >
                        <Box
                          sx={{
                            display: "flex",
                            justifyContent: "space-between",
                            alignItems: "center",
                            width: "100%",
                            ml: 1,
                          }}
                        >
                          <Box
                            sx={{
                              display: "flex",
                              alignItems: "baseline",
                              gap: 1,
                            }}
                          >
                            <Typography variant="body1" color="text.secondary">
                              {formatDateTime(interview.start_time, {
                                weekday: "short",
                                year: "numeric",
                                month: "short",
                                day: "numeric",
                              })}
                            </Typography>
                            <Typography
                              variant="subtitle1"
                              sx={{ fontWeight: 500, color: "primary.main" }}
                            >
                              {formatDateTime(interview.start_time, {
                                hour: "2-digit",
                                minute: "2-digit",
                              })}
                            </Typography>
                          </Box>
                          <Box
                            sx={{
                              display: "flex",
                              alignItems: "center",
                              gap: 2,
                            }}
                          >
                            <Chip
                              label={getInterviewTypeLabel(
                                interview.interview_type
                              )}
                              size="small"
                              color="primary"
                              variant="outlined"
                            />
                            <Chip
                              label={getInterviewStateLabel(
                                interview.interview_state
                              )}
                              color={getInterviewStateColor(
                                interview.interview_state
                              )}
                              size="small"
                            />
                          </Box>
                        </Box>
                      </AccordionSummary>
                      <AccordionDetails>
                        <Box
                          sx={{
                            display: "flex",
                            flexDirection: "column",
                            gap: 2,
                          }}
                        >
                          {/* Interview Details */}
                          <Box>
                            <Typography variant="subtitle1" gutterBottom>
                              {t("interviews.details")}
                            </Typography>
                            <Typography
                              variant="body2"
                              color="text.secondary"
                              paragraph
                            >
                              {interview.description}
                            </Typography>
                            <Typography variant="body2" color="text.secondary">
                              {t("interviews.timeRange", {
                                start: formatDateTime(interview.start_time),
                                end: formatDateTime(interview.end_time),
                              })}
                            </Typography>
                            <Typography
                              variant="caption"
                              color="text.secondary"
                              sx={{ mt: 1, display: "block" }}
                            >
                              {t("interviews.timezone", {
                                zone: Intl.DateTimeFormat().resolvedOptions()
                                  .timeZone,
                              })}
                            </Typography>
                          </Box>

                          {/* RSVP Section */}
                          <Box>
                            <Typography variant="subtitle1" gutterBottom>
                              {t("interviews.yourRSVP")}
                            </Typography>
                            <Box
                              sx={{
                                display: "flex",
                                alignItems: "center",
                                gap: 2,
                              }}
                            >
                              <RSVPStatusIcon
                                status={interview.candidate_rsvp_status}
                              />
                              {interview.interview_state ===
                                InterviewStates.SCHEDULED_INTERVIEW && (
                                <ButtonGroup size="small">
                                  <Button
                                    variant={
                                      interview.candidate_rsvp_status ===
                                      RSVPStatuses.YES
                                        ? "contained"
                                        : "outlined"
                                    }
                                    onClick={() =>
                                      handleRSVPClick(
                                        interview.interview_id,
                                        RSVPStatuses.YES,
                                        interview.candidate_rsvp_status
                                      )
                                    }
                                    color="success"
                                    disabled={
                                      rsvpLoading === interview.interview_id ||
                                      interview.candidate_rsvp_status ===
                                        RSVPStatuses.YES
                                    }
                                    startIcon={
                                      rsvpLoading ===
                                        interview.interview_id && (
                                        <CircularProgress size={16} />
                                      )
                                    }
                                    sx={{
                                      "&.Mui-disabled": {
                                        backgroundColor:
                                          interview.candidate_rsvp_status ===
                                          RSVPStatuses.YES
                                            ? "success.main"
                                            : "transparent",
                                        color:
                                          interview.candidate_rsvp_status ===
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
                                      interview.candidate_rsvp_status ===
                                      RSVPStatuses.NO
                                        ? "contained"
                                        : "outlined"
                                    }
                                    onClick={() =>
                                      handleRSVPClick(
                                        interview.interview_id,
                                        RSVPStatuses.NO,
                                        interview.candidate_rsvp_status
                                      )
                                    }
                                    color="error"
                                    disabled={
                                      rsvpLoading === interview.interview_id ||
                                      interview.candidate_rsvp_status ===
                                        RSVPStatuses.NO
                                    }
                                    startIcon={
                                      rsvpLoading ===
                                        interview.interview_id && (
                                        <CircularProgress size={16} />
                                      )
                                    }
                                    sx={{
                                      "&.Mui-disabled": {
                                        backgroundColor:
                                          interview.candidate_rsvp_status ===
                                          RSVPStatuses.NO
                                            ? "error.main"
                                            : "transparent",
                                        color:
                                          interview.candidate_rsvp_status ===
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
                              )}
                            </Box>
                          </Box>

                          {/* Interviewers Section */}
                          <Box>
                            <Typography variant="subtitle1" gutterBottom>
                              {t("interviews.interviewers")}
                            </Typography>
                            <Box
                              sx={{ display: "flex", flexWrap: "wrap", gap: 1 }}
                            >
                              {interview.interviewers &&
                              interview.interviewers.length > 0 ? (
                                interview.interviewers.map(
                                  (interviewer, index) => (
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
                                          <RSVPStatusIcon
                                            status={interviewer.rsvp_status}
                                          />
                                        </Box>
                                      }
                                      variant="outlined"
                                    />
                                  )
                                )
                              ) : (
                                <Typography color="text.secondary">
                                  {t("interviews.noInterviewers")}
                                </Typography>
                              )}
                            </Box>
                          </Box>
                        </Box>
                      </AccordionDetails>
                    </Accordion>
                  ))
                ) : (
                  <Typography color="text.secondary">
                    {t("interviews.noInterviews")}
                  </Typography>
                )}
              </AccordionDetails>
            </Accordion>

            <Paper sx={{ p: 3 }}>
              <Typography variant="h6" gutterBottom>
                {t("comments.title")}
              </Typography>

              {comments.length > 0 ? (
                <Box sx={{ mt: 3 }}>
                  {[...comments]
                    .sort(
                      (a, b) =>
                        new Date(a.created_at).getTime() -
                        new Date(b.created_at).getTime()
                    )
                    .map((comment) => (
                      <Box
                        key={comment.comment_id}
                        sx={{
                          display: "flex",
                          gap: 2,
                          mb: 3,
                          flexDirection:
                            comment.commenter_type === CommenterTypes.ORG_USER
                              ? "row"
                              : "row-reverse",
                        }}
                      >
                        <Avatar
                          sx={{
                            width: 40,
                            height: 40,
                            bgcolor: (theme) =>
                              comment.commenter_type === CommenterTypes.ORG_USER
                                ? theme.palette.primary.main
                                : theme.palette.grey[400],
                          }}
                        >
                          {comment.commenter_name.charAt(0).toUpperCase()}
                        </Avatar>
                        <Box sx={{ flexGrow: 1 }}>
                          <Paper
                            sx={{
                              p: 2,
                              borderRadius: 2,
                              border: "1px solid",
                              borderColor: "divider",
                              position: "relative",
                              "&::before": {
                                content: '""',
                                position: "absolute",
                                ...(comment.commenter_type ===
                                CommenterTypes.ORG_USER
                                  ? {
                                      left: -8,
                                      borderRight: (theme) =>
                                        `8px solid ${theme.palette.divider}`,
                                    }
                                  : {
                                      right: -8,
                                      borderLeft: (theme) =>
                                        `8px solid ${theme.palette.divider}`,
                                    }),
                                top: 16,
                                width: 0,
                                height: 0,
                                borderTop: "8px solid transparent",
                                borderBottom: "8px solid transparent",
                              },
                            }}
                          >
                            <Box
                              sx={{
                                display: "flex",
                                justifyContent: "space-between",
                                alignItems: "center",
                                mb: 1,
                              }}
                            >
                              <Typography
                                variant="subtitle2"
                                sx={{
                                  fontWeight: "bold",
                                  color: (theme) =>
                                    comment.commenter_type ===
                                    CommenterTypes.ORG_USER
                                      ? theme.palette.primary.main
                                      : theme.palette.text.primary,
                                }}
                              >
                                {comment.commenter_name}
                              </Typography>
                              <Typography
                                variant="caption"
                                color="text.secondary"
                              >
                                {new Date(
                                  comment.created_at
                                ).toLocaleDateString(undefined, {
                                  year: "numeric",
                                  month: "short",
                                  day: "2-digit",
                                })}{" "}
                                {new Date(
                                  comment.created_at
                                ).toLocaleTimeString(undefined, {
                                  hour: "2-digit",
                                  minute: "2-digit",
                                })}
                              </Typography>
                            </Box>
                            <Typography
                              sx={{
                                whiteSpace: "pre-wrap",
                                wordBreak: "break-word",
                              }}
                            >
                              {comment.content}
                            </Typography>
                          </Paper>
                        </Box>
                      </Box>
                    ))}
                </Box>
              ) : (
                <Typography color="text.secondary" sx={{ my: 2 }}>
                  {t("comments.noComments")}
                </Typography>
              )}

              <Divider sx={{ my: 3 }} />

              <Box>
                <TextField
                  fullWidth
                  multiline
                  rows={4}
                  value={newComment}
                  onChange={(e) => setNewComment(e.target.value)}
                  placeholder={t("comments.addPlaceholder")}
                  disabled={submitting}
                />
                <Box
                  sx={{ display: "flex", justifyContent: "flex-end", mt: 2 }}
                >
                  <Button
                    variant="contained"
                    onClick={handleAddComment}
                    disabled={!newComment.trim() || submitting}
                  >
                    {submitting ? t("common.loading") : t("comments.add")}
                  </Button>
                </Box>
              </Box>
            </Paper>
          </>
        )
      )}
    </Box>
  );

  return (
    <AuthenticatedLayout>
      {content}
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
            {confirmDialog.status === RSVPStatuses.YES
              ? interviews.find(
                  (i) => i.interview_id === confirmDialog.interviewId
                )?.candidate_rsvp_status === RSVPStatuses.NOT_SET
                ? t("interviews.rsvp.confirmYesMessage")
                : t("interviews.rsvp.confirmChangeYesMessage")
              : interviews.find(
                  (i) => i.interview_id === confirmDialog.interviewId
                )?.candidate_rsvp_status === RSVPStatuses.NOT_SET
              ? t("interviews.rsvp.confirmNoMessage")
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
