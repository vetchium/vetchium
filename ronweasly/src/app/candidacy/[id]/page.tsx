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
} from "@mui/material";
import { OpenInNew as OpenInNewIcon } from "@mui/icons-material";
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
} from "@psankar/vetchi-typespec";
import { config } from "@/config";
import Cookies from "js-cookie";
import AuthenticatedLayout from "@/components/AuthenticatedLayout";

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

export default function CandidacyDetailPage() {
  const params = useParams();
  const candidacyId = params.id as string;
  const { t } = useTranslation();
  const router = useRouter();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [candidacy, setCandidacy] = useState<MyCandidacy | null>(null);
  const [comments, setComments] = useState<CandidacyComment[]>([]);
  const [interviews, setInterviews] = useState<HubInterview[]>([]);
  const [newComment, setNewComment] = useState("");
  const [submitting, setSubmitting] = useState(false);

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

  // Fetch data on mount
  useEffect(() => {
    fetchData();
    fetchInterviews();
  }, []); // Empty dependency array means this runs once on mount

  const formatDateTime = (dateTime: string): string => {
    return new Date(dateTime).toLocaleString();
  };

  const getInterviewStateColor = (
    state: InterviewState
  ): "warning" | "success" | "error" | "default" => {
    switch (state) {
      case InterviewStates.SCHEDULED_INTERVIEW:
        return "warning";
      case InterviewStates.COMPLETED_INTERVIEW:
        return "success";
      case InterviewStates.CANCELLED_INTERVIEW:
        return "error";
      default:
        return "default";
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
            <Paper sx={{ p: 3, mb: 3 }}>
              <Typography variant="h6" gutterBottom>
                {t("interviews.title")}
              </Typography>
              {interviews.length > 0 ? (
                interviews.map((interview) => (
                  <Box
                    key={interview.interview_id}
                    sx={{
                      py: 2,
                      "&:not(:last-child)": {
                        borderBottom: 1,
                        borderColor: "divider",
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
                      <Typography variant="subtitle1">
                        {getInterviewTypeLabel(interview.interview_type)}
                      </Typography>
                      <Chip
                        label={t(
                          `interviews.states.${interview.interview_state}`
                        )}
                        color={getInterviewStateColor(
                          interview.interview_state
                        )}
                        size="small"
                      />
                    </Box>
                    <Typography variant="body2" color="text.secondary">
                      {t("interviews.timeRange", {
                        start: new Date(interview.start_time).toLocaleString(),
                        end: new Date(interview.end_time).toLocaleString(),
                      })}
                    </Typography>
                    {interview.description && (
                      <Typography variant="body2" sx={{ mt: 1 }}>
                        {interview.description}
                      </Typography>
                    )}
                    {interview.interviewers &&
                      interview.interviewers.length > 0 && (
                        <Box sx={{ mt: 1 }}>
                          <Typography variant="body2" color="text.secondary">
                            {t("interviews.interviewers")}:{" "}
                            {interview.interviewers.join(", ")}
                          </Typography>
                        </Box>
                      )}
                  </Box>
                ))
              ) : (
                <Typography color="text.secondary">
                  {t("interviews.noInterviews")}
                </Typography>
              )}
            </Paper>

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
                            comment.commenter_type === "ORG_USER"
                              ? "row"
                              : "row-reverse",
                        }}
                      >
                        <Avatar
                          sx={{
                            width: 40,
                            height: 40,
                            bgcolor: (theme) =>
                              comment.commenter_type === "ORG_USER"
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
                                ...(comment.commenter_type === "ORG_USER"
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
                                    comment.commenter_type === "ORG_USER"
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

  return <AuthenticatedLayout>{content}</AuthenticatedLayout>;
}
